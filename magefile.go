// +build mage

package main

import (
	"github.com/fatih/color"
	"github.com/magefile/mage/sh"
	"github.com/mysteriumnetwork/go-ci/commands"
	"github.com/mysteriumnetwork/go-ci/shell"
)

// Build builds the service
func Build() error {
	return shell.NewCmd("go build -o ./build/feedback ./cmd/main.go").Run()
}

// Regen re-generates API schema (swagger.json) and related bindata files
func Regen() error {
	color.Cyan("Installing stuff")
	err := shell.NewCmd("go get -u github.com/go-swagger/go-swagger/...").Run()
	if err != nil {
		return err
	}
	color.Cyan("Generating swagger.json")
	err = shell.NewCmd("swagger generate spec --scan-models --output=./docs/swagger.json").Run()
	if err != nil {
		return err
	}
	color.Cyan("Generating assets for serving swagger.json")
	return shell.NewCmd("go generate ./...").Run()
}

// Validate validates API schema (swagger.json)
func Validate() error {
	color.Cyan("Installing stuff")
	err := shell.NewCmd("go get -u github.com/go-swagger/go-swagger/...").Run()
	if err != nil {
		return err
	}
	color.Cyan("Validating swagger.json")
	return shell.NewCmd("swagger validate ./docs/swagger.json").Run()
}

// Test runs tests
func Test() error {
	return commands.Test("./...")
}

// CheckCopyright checks for issues with go imports
func CheckCopyright() error {
	return commands.CopyrightD(".", "docs")
}

// CheckGo checks for issues with go imports
func CheckGoImports() error {
	return commands.GoImportsD(".", "docs")
}

// CheckGoLint reports linting errors
func CheckGoLint() error {
	return commands.GoLintD(".", "docs")
}

// CheckGoVet checks that the source is compliant with go vet
func CheckGoVet() error {
	return commands.GoVet(".")
}

// Check runs all checks
func Check() error {
	return commands.CheckD(".", "docs")
}

// Run runs the service
func Run() error {
	return sh.RunV("go", "run", "./cmd/main.go")
}
