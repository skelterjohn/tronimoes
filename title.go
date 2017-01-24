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
	"image"
	"image/draw"
	_ "image/png"
	"log"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

type titleLayer struct {
	s *Stack

	tronimoesImage *glutil.Image

	buf      gl.Buffer
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	size     gl.Uniform
	color    gl.Uniform
}

func (tl *titleLayer) init(glctx gl.Context) {
	pngReader, err := asset.Open("tronimoes.png")
	if err != nil {
		log.Fatalf("Error loading image: %v", err)
	}
	defer pngReader.Close()

	png, _, err := image.Decode(pngReader)
	if err != nil {
		log.Fatalf("Error decoding image: %v", err)
	}
	b := png.Bounds()

	tl.tronimoesImage = images.NewImage(b.Dx(), b.Dy())

	log.Print(b, b.Min)
	//raw(dst Image, r image.Rectangle, src image.Image, sp image.Point, op Op)
	draw.Draw(tl.tronimoesImage.RGBA, b, png, image.Point{0, 0}, draw.Src)

}

func (tl *titleLayer) Paint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(1, 0, 0, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	b := tl.tronimoesImage.RGBA.Bounds()
	tl.tronimoesImage.Draw(
		sz,
		geom.Point{0, 0},
		geom.Point{geom.Pt(b.Max.X), 0},
		geom.Point{0, geom.Pt(b.Max.Y)},
		b,
	)
}

func (tl *titleLayer) Event(e interface{}) {

}

func (tl *titleLayer) Release(glctx gl.Context) {
	tl.tronimoesImage.Release()
	glctx.DeleteProgram(tl.program)
}
