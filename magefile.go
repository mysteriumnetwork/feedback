// +build mage

/*
 * Copyright (C) 2019 The "MysteriumNetwork/feedback" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
