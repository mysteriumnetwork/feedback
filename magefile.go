//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mysteriumnetwork/feedback/ci"
	"github.com/mysteriumnetwork/go-ci/commands"
	"github.com/mysteriumnetwork/go-ci/shell"
	"github.com/mysteriumnetwork/go-ci/util"
)

// Build builds the service
func Build() error {
	if os.Getenv("GITHUB_CI") == "" {
		mg.Deps(ci.Swag)
		mg.Deps(Generate)
	} else {
		fmt.Println("Skipping swagger generation in CI environment")
	}

	return shell.NewCmd("go build -o ./build/feedback ./cmd/main.go").Run()
}

// Installs the swag generation tool
func swagInstall() error {
	return sh.RunV("go", "install", "github.com/swaggo/swag/cmd/swag@v1.16.2")
}

// Swag generates the swagger documents
func Swag() error {
	mg.Deps(swagInstall)
	swag, err := util.GetGoBinaryPath("swag")
	if err != nil {
		return err
	}

	return sh.RunV(
		swag, "init",
		"--generalInfo", "./cmd/main.go",
		"--output", "./docs",
		"--parseDependency",
		"--parseDepth", "1",
	)
}

// Test runs unit tests
func Test() error {
	packages, err := unitTestPackages()
	if err != nil {
		return err
	}
	args := append([]string{"test", "-race", "-count=1", "-timeout", "5m"}, packages...)
	return sh.RunV("go", args...)
}

func unitTestPackages() ([]string, error) {
	allPackages, err := listPackages()
	if err != nil {
		return nil, err
	}
	packages := make([]string, 0)
	for _, p := range allPackages {
		if !strings.Contains(p, "e2e") {
			packages = append(packages, p)
		}
	}
	return packages, nil
}

func listPackages() ([]string, error) {
	output, err := sh.Output("go", "list", "./...")
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.Replace(output, "\r\n", "\n", -1), "\n"), nil
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

// E2E runs the e2e test suite
func E2E() error {
	return ci.E2E()
}

// Generate runs go generate
func Generate() error {
	return sh.RunV("go", "generate", "./...")
}
