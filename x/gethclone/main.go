// The gethclone binary clones ethereum/go-ethereum Go packages, applying
// semantic patches.
package main

import (
	"context"
	"log"
	"os"

	"github.com/ava-labs/subnet-evm/x/gethclone/astpatch"
	"github.com/spf13/pflag"
)

func main() {
	c := config{
		astPatches: make(astpatch.PatchRegistry),
	}

	pflag.StringSliceVar(&c.packages, "packages", []string{"core/vm"}, `Geth packages to clone, with or without "github.com/ethereum/go-ethereum" prefix.`)
	pflag.StringVar(&c.outputGoMod, "output_go_mod", "", "go.mod file of the destination to which geth will be cloned.")
	pflag.Parse()

	log.SetOutput(os.Stderr)
	log.Print("START")
	if err := c.run(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Print("DONE")
}
