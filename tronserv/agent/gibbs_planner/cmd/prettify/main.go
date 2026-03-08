// One-off: go run ./cmd/prettify testdata/no_draw_btbrjx.json
// Reads JSON (strips BOM), unmarshals into game.Game, writes back with tab indent.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/skelterjohn/tronimoes/tronserv/game"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: go run ./cmd/prettify <file.json>\n")
		os.Exit(1)
	}
	path := os.Args[1]
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		os.Exit(1)
	}
	// Strip UTF-8 BOM
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}
	var g game.Game
	if err := json.Unmarshal(data, &g); err != nil {
		fmt.Fprintf(os.Stderr, "unmarshal: %v\n", err)
		os.Exit(1)
	}
	out, err := json.MarshalIndent(&g, "", "\t")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(path, out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("ok")
}
