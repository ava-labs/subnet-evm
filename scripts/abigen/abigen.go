// The abigen command runs `solc | abigen` and writes the generated bindings to
// stdout.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	var c config

	flag.StringVar(&c.solc.version, "solc.version", "0.8.24", "Version of solc expected; the version used will be sourced from $PATH")
	flag.StringVar(&c.solc.evmVersion, "solc.evm-version", "paris", "solc --evm-version flag")
	flag.StringVar(&c.solc.basePath, "solc.base-path", "./", "solc --base-path flag")
	flag.StringVar(&c.solc.includePath, "solc.include-path", "", "solc --include-path flag; only propagated if not empty")
	flag.StringVar(&c.solc.output, "solc.output", "abi,bin", "solc --combined-json flag")
	flag.StringVar(&c.abigen.pkg, "abigen.pkg", "", "abigen --pkg flag")

	help := flag.Bool("help", false, "Print usage message")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if err := c.run(context.Background(), os.Stdout, os.Stderr); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

type config struct {
	solc struct {
		version, evmVersion, basePath, includePath, output string
	}
	abigen struct {
		pkg string
	}
}

func (cfg *config) run(ctx context.Context, stdout, stderr io.Writer) error {
	solcV := exec.CommandContext(ctx, "solc", "--version")
	buf, err := solcV.CombinedOutput()
	if err != nil {
		return fmt.Errorf("solc --version: %v", err)
	}
	if !bytes.Contains(buf, []byte(cfg.solc.version)) {
		fmt.Fprintf(stderr, "solc --version:\n%s", buf)
		return fmt.Errorf("solc version mismatch; not %q", cfg.solc.version)
	}

	args := append(nonEmptyArgs(map[string]string{
		"--evm-version":   cfg.solc.evmVersion,
		"--base-path":     cfg.solc.basePath,
		"--include-path":  cfg.solc.includePath,
		"--combined-json": cfg.solc.output,
	}), flag.Args()...)

	solc := exec.CommandContext(ctx, "solc", args...)
	// Although we could use io.Pipe(), it's much easier to reason about
	// non-concurrent processes, and solc doesn't create huge outputs.
	var solcOut bytes.Buffer
	solc.Stdout = &solcOut
	solc.Stderr = stderr
	if err := solc.Run(); err != nil {
		return fmt.Errorf("solc: %w", err)
	}

	abigen := exec.CommandContext(ctx, "abigen", nonEmptyArgs(map[string]string{
		"--combined-json": "-", // stdin
		"--pkg":           cfg.abigen.pkg,
	})...)
	abigen.Stdin = &solcOut
	var abigenOut bytes.Buffer
	abigen.Stdout = &abigenOut
	abigen.Stderr = stderr
	if err := abigen.Run(); err != nil {
		return fmt.Errorf("abigen: %w", err)
	}

	re := regexp.MustCompile(`"github\.com/ethereum/go-ethereum/(accounts|core)/`)
	_, err = stdout.Write(re.ReplaceAll(
		abigenOut.Bytes(),
		[]byte(`"github.com/ava-labs/subnet-evm/${1}/`)),
	)
	return err
}

// nonEmptyArgs returns a slice of arguments suitable for use with exec.Command.
// Any empty values are skipped.
func nonEmptyArgs(from map[string]string) []string {
	var args []string
	for k, v := range from {
		if v != "" {
			args = append(args, k, v)
		}
	}
	return args
}
