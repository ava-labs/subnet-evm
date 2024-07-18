// The getsolc binary downloads the Solidity compiler from official sources. If
// a `solc` binary in the PATH is of the requested version then a symlink is
// created instead of downloading a new copy.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	var c config

	flag.StringVar(&c.version, "version", latestVersion, fmt.Sprintf("Version of solc; {major}.{minor}.{patch} or %q", latestVersion))
	flag.StringVar(&c.outputFile, "out", "./solc", "Path to which the `solc` binary will be saved")
	flag.BoolVar(&c.ignoreGOARCH, "ignore_goarch", false, "Download amd64 binary even if on another architecture")

	flag.Parse()

	if err := c.run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const latestVersion = "latest"

type config struct {
	version, outputFile string
	ignoreGOARCH        bool
}

func (c *config) run(ctx context.Context) error {
	goos := runtime.GOOS
	switch goos {
	case "darwin":
		goos = "macosx"
	case "linux":
	default:
		return fmt.Errorf("unsupported OS %q", goos)
	}

	// solc only provides amd64 binaries, but this can be ignored if there is a
	// translator (e.g. Rosetta on MacOS)
	if !c.ignoreGOARCH && runtime.GOARCH != "amd64" {
		return fmt.Errorf("unsupported GOARCH %q", runtime.GOARCH)
	}

	jsonList, err := httpGetSolFile(ctx, goos, "list.json")
	if err != nil {
		return err
	}
	defer jsonList.Body.Close()
	var list solcList
	if err := json.NewDecoder(jsonList.Body).Decode(&list); err != nil {
		return err
	}

	if c.version == latestVersion {
		c.version = list.LatestRelease
	}
	if p, ok := c.bestEffortFindInPATH(ctx); ok { // NOTE: this is not an error path
		fmt.Fprintf(os.Stderr, "Creating symlink from %q to %q\n", p, c.outputFile)
		return os.Link(p, c.outputFile)
	}
	fmt.Fprintln(os.Stderr, "Downloading solc...")

	solc, err := httpGetSolFile(ctx, goos, list.Releases[c.version])
	if err != nil {
		return err
	}
	defer solc.Body.Close()

	hash := crypto.NewKeccakState()
	tee := io.TeeReader(solc.Body, hash)

	out, err := os.OpenFile(c.outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %v", err)
	}
	if _, err := io.Copy(out, tee); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := out.Close(); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Keccak256 of download: %#x\n", hash.Sum(nil))
	return nil
}

// httpGetSolFile issues an HTTP(s) GET to download the specified file from the
// official binaries.soliditylang.org source. `goos` can be either linux or
// macosx, and file can be "list.json" or any of the paths in [solcList].
func httpGetSolFile(ctx context.Context, goos, file string) (*http.Response, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "binaries.soliditylang.org",
		Path:   path.Join(fmt.Sprintf("%s-amd64", goos), file),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP GET %q; Status = %d; err = %v", u.String(), resp.StatusCode, err)
	}
	return resp, nil
}

// solcList mirrors the JSON lists published by the solc team; e.g:
// https://binaries.soliditylang.org/linux-amd64/list.json
type solcList struct {
	Releases      map[string]string `json:"releases"`
	LatestRelease string            `json:"latestRelease"`
	Builds        []struct {
		Path, Version, Build, Keccak256 string
		LongVersion                     string `json:"longVersion"`
		SHA256                          string `json:"sha256"`
	} `json:"builds"`
}

// bestEffortFindInPATH attempts to find `solc` in the PATH and, if it has the
// required version, the path to said binary is returned. The boolean indicates
// succesful location of a matching binary.
func (c *config) bestEffortFindInPATH(ctx context.Context) (string, bool) {
	solc := exec.CommandContext(ctx, "solc", "--version")
	out, err := solc.CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte(c.version)) {
		return "", false
	}

	which := exec.CommandContext(ctx, "which", "solc")
	solcPath, err := which.CombinedOutput()
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(solcPath)), true
}
