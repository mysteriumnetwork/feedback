// +build ignore

package main

import (
	"log"

	"github.com/go-bindata/go-bindata"
)

func main() {
	cfg := &bindata.Config{
		Package: "docs",
		Input: []bindata.InputConfig{
			{Path: "./swagger.json"},
		},
		Output: "./swagger_bindata.go",
	}
	err := bindata.Translate(cfg)
	if err != nil {
		log.Fatalln(err)
	}
}
