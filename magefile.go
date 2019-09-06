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
	"github.com/magefile/mage/sh"
	"github.com/mysteriumnetwork/go-ci/commands"
)

// Build builds the service
func Build() error {
	return sh.RunV("go", "build", "-o", "./build/feedback", "./cmd/main.go")
}

// Test runs tests
func Test() error {
	return commands.Test("./...")
}

// CheckGoImports checks for issues with go imports
func CheckGoImports() error {
	return commands.GoImports("./...", "docs")
}

// CheckGoLint reports linting errors
func CheckGoLint() error {
	return commands.GoLint("./...", "docs")
}

// CheckGoVet checks that the source is compliant with go vet
func CheckGoVet() error {
	return commands.GoVet("./...")
}

// Check runs all checks
func Check() error {
	return commands.Check("./...", "docs")
}

// Run runs the service
func Run() error {
	return sh.RunV("go", "run", "./cmd/main.go")
}
