/*
Copyright 2017 Google Inc. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/binary"
	"log"
	"sync"

	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

type Stack struct {
	mu     sync.Mutex // protexts layers
	layers []Layer
}

func (s *Stack) Push(l Layer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.layers = append(s.layers, l)
}

func (s *Stack) Pop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.layers) == 0 {
		return
	}
	s.layers = s.layers[:len(s.layers)-1]
}

func (s *Stack) Peek() Layer {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.layers) == 0 {
		return nil
	}
	return s.layers[len(s.layers)-1]
}

type Layer interface {
	Paint(glctx gl.Context, sz size.Event)
	Event(e interface{})
	Release(glctx gl.Context)
}

type sillyPainter struct {
	touchX   float32
	touchY   float32
	buf      gl.Buffer
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	size     gl.Uniform
	color    gl.Uniform
}

var triangleData = f32.Bytes(binary.LittleEndian,
	-1.0, 1.0, 0.0, // top left
	-1.0, -1.0, 0.0, // bottom left
	1.0, 1.0, 0.0, // top right
	1.0, -1.0, 0.0, // bottom right
)

const (
	coordsPerVertex = 3
	vertexCount     = 4
)

func (sp *sillyPainter) init(glctx gl.Context) {
	sp.buf = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, sp.buf)
	glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	const vertexShader = `#version 100
uniform vec2 offset;
uniform vec2 size;

attribute vec4 position;
void main() {
	// offset comes in with x/y values between 0 and w/h.
	// position bounds are -1 to 1.
	vec4 offset4 = vec4(2.0*offset.x/size.x-1.0, 1.0-2.0*offset.y/size.y, 0, 0);
	gl_Position = position + offset4;
}`

	const fragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`

	var err error
	sp.program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	sp.position = glctx.GetAttribLocation(sp.program, "position")
	sp.color = glctx.GetUniformLocation(sp.program, "color")
	sp.offset = glctx.GetUniformLocation(sp.program, "offset")
	sp.size = glctx.GetUniformLocation(sp.program, "size")

}

func (sp *sillyPainter) Release(glctx gl.Context) {
	glctx.DeleteBuffer(sp.buf)
	glctx.DeleteProgram(sp.program)
}

func (sp *sillyPainter) Paint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(1, 0, 0, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	glctx.UseProgram(sp.program)

	glctx.Uniform4f(sp.color, 0, 1, 0, 1)
	glctx.Uniform2f(sp.size, float32(sz.WidthPx), float32(sz.HeightPx))
	glctx.Uniform2f(sp.offset, sp.touchX, sp.touchY)

	glctx.BindBuffer(gl.ARRAY_BUFFER, sp.buf)
	glctx.EnableVertexAttribArray(sp.position)
	glctx.VertexAttribPointer(sp.position, coordsPerVertex, gl.FLOAT, false, 0, 0)
	glctx.DrawArrays(gl.TRIANGLE_STRIP, 0, vertexCount)
	glctx.DisableVertexAttribArray(sp.position)

}

func (sp *sillyPainter) Event(e interface{}) {
	switch e := e.(type) {
	case touch.Event:
		sp.touchX = e.X
		sp.touchY = e.Y
	case size.Event:
		sp.touchX = float32(e.WidthPx / 2)
		sp.touchY = float32(e.HeightPx / 2)
	}
}
