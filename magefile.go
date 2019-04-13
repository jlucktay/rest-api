// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified.
var Default = Test

// Test will run all tests across all sub-directories once.
func Test() error {
	return sh.Run("go", "test", "./...", "--count=1")
}

// Lint will check the Dockerfile and Go files for errors.
func Lint() error {
	if errDockerLint := sh.RunV("hadolint", "build/Dockerfile"); errDockerLint != nil {
		return errDockerLint
	}
	return sh.Run("golangci-lint", "run", "--enable-all")
}

// Dock will build the Docker image containing the REST API.
func Dock() error {
	//     --label list              Set metadata for an image // todo?
	return sh.Run("docker", "build", "--file", "build/Dockerfile", "--tag", "jlucktay/rest-api", ".")
}

// Build will compile the REST API binary.
func Build() error {
	return sh.Run("go", "build", "-a", "-installsuffix", "cgo", "-ldflags", `-extldflags "-static"`, "-o", "main", "./cmd/api")
}

// Clean will removed compiled binaries.
func Clean() error {
	return sh.Run("rm", "-fv", "main")
}

// Full runs all targets; linting and testing in parallel, then the Docker build.
func Full() {
	mg.Deps(Clean)
	mg.Deps(Lint, Test)
	mg.Deps(Build, Dock)
}
