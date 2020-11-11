package standard

import (
	"context"
	"os"
	"os/exec"

	"github.com/paganotoni/x/internal/info"
)

// Build runs the Go compiler to generate the desired binary. Assuming the
// Go executable installed and can be invoked with `go`.
//
// IMPORTANT: it uses the static build flags.
func (g *Plugin) Build(ctx context.Context, root string, args []string) error {
	name, err := info.BuildName()
	if err != nil {
		return err
	}

	buildArgs := []string{
		"build",

		//--static
		"--ldflags",
		"-linkmode external",
		"--ldflags",
		`-extldflags "-static"`,

		//-o
		"-o",
		g.binaryOutput(name),

		// The main we're going to build
		"./cmd/" + name,
	}

	cmd := exec.CommandContext(ctx, "go", buildArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// binaryOutput considers the output passed to
// use it or default to bin/name.
func (g *Plugin) binaryOutput(name string) string {
	output := "bin/" + name
	if g.output != "" {
		output = g.output
	}

	return output
}