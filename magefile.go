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
	args := []string{"test", "./...", "--count=1"}
	if mg.Verbose() {
		args = append(args, "-v")
	}
	return sh.RunV("go", args...)
}

// Lint will check the Dockerfile and Go files for errors.
func Lint() {
	mg.Deps(LintDocker, LintGo)
}

// LintDocker lints the Dockerfile.
func LintDocker() error {
	return sh.Run("hadolint", "build/Dockerfile")
}

// LintGo lints all Go files.
func LintGo() error {
	return sh.Run("golangci-lint", "run", "--enable-all")
}

// Build will compile the REST API binary locally.
func Build() error {
	mg.Deps(Clean)
	return sh.Run("go", "build", "-a", "-installsuffix", "cgo", "-ldflags", `-extldflags "-static"`, "-o", "rest-api", "./cmd/api")
}

// DockerBuild will build the Docker image, which executes a layered build of the REST API.
func DockerBuild() error {
	//     --label list              Set metadata for an image // todo?
	return sh.Run("docker", "build", "--file", "build/Dockerfile", "--tag", "jlucktay/rest-api", ".")
}

// DockerRun will run the Docker image.
func DockerRun() error {
	mg.Deps(DockerBuild)
	return sh.RunV("docker", "run", "--publish", "8080:8080", "jlucktay/rest-api")
}

// Clean will removed compiled binaries.
func Clean() error {
	return sh.Run("rm", "-fv", "rest-api")
}

// Full runs all targets; linting and testing in parallel, then the Docker build.
func Full() {
	mg.Deps(Lint, Test)
	mg.Deps(Build, DockerBuild)
	mg.Deps(DockerRun)
}
