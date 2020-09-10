// +build tools

package main

// Empty file for keeping CI tools in go.mod

import (
	_ "github.com/go-bindata/go-bindata"
	_ "github.com/go-openapi/jsonreference"
	_ "github.com/go-openapi/runtime"
	_ "github.com/go-swagger/go-swagger"
	_ "github.com/mailru/easyjson"
	_ "golang.org/x/tools/imports"
)
