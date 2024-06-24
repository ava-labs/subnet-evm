package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"go.uber.org/zap"
)

// A goRunner runs the `go` binary.
type goRunner struct {
	bin string // if empty, `go` will be found in $PATH
	log *zap.SugaredLogger
}

// list runs `go list -json [patterns...]` and returns the parsed output. It
// returns an error if `PackagePublic.Err` is non-nil, but ignores
// `PackagePublic.DepsErrors` as they may pertain to a non-geth package that
// `gethclone` doesn't have available. Any dependency error pertaining to a geth
// dep will eventually be reached when traversing the import tree.
func (r *goRunner) list(ctx context.Context, patterns ...string) ([]*PackagePublic, error) {
	if r.bin == "" {
		// Although we could just use "go" as the command name later, this
		// allows logging for debugging.
		p, err := exec.LookPath("go")
		if err != nil {
			return nil, fmt.Errorf("finding `go` binary: %v", err)
		}
		r.log.Infof("Found `go` in PATH: %q", p)
		r.bin = p
	}

	var result []*PackagePublic

	// When `go list` finds more than one package, its output is not a JSON list
	// but a new-line concatenation of each JSON object. It's simpler to run `go
	// list` multiple times than it is to convert the output to a JSON list.
	for _, p := range patterns {
		// -find stops `go list` from traversing the dependency tree
		cmd := exec.CommandContext(ctx, r.bin, "list", "-find", "-json", p)

		start := time.Now()
		buf, err := cmd.Output()
		end := time.Now()
		if ee, ok := err.(*exec.ExitError); ok {
			r.log.Errorf("stderr of `go list`:\n%s", ee.Stderr)
		}
		if err != nil {
			return nil, fmt.Errorf("running `go list`: %v", err)
		}
		r.log.Debugf("`go list ... %q` ran in %s", p, end.Sub(start))

		pkg := new(PackagePublic)
		if err := json.Unmarshal(buf, pkg); err != nil {
			return nil, fmt.Errorf("unmarshal JSON output of `go list`: %v", err)
		}
		if pkg.Error != nil {
			return nil, fmt.Errorf("package error encountered by `go list`: %v", pkg.Error)
		}
		result = append(result, pkg)
	}

	return result, nil
}

/*
Types included below are copied (almost) verbatim from the Go source code, under
the included license (ny modifications are documented before the respective
code). This is necessary because they're in `internal` packages that we can't
otherwise access. Despite the internal designation, their output is considered
part of the public API of the `go list` command so it is safe to assume
stability due to the Go 1 compatibility promise.
*/

/*
https://go.googlesource.com/go/+/refs/tags/go1.22.4/LICENSE

Copyright (c) 2009 The Go Authors. All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// https://go.googlesource.com/go/+/refs/tags/go1.22.4/src/cmd/go/internal/load/pkg.go#59

// The `Module` field has been modified to use a package-local type.

type PackagePublic struct {
	// Note: These fields are part of the go command's public API.
	// See list.go. It is okay to add fields, but not to change or
	// remove existing ones. Keep in sync with ../list/list.go
	Dir            string        `json:",omitempty"` // directory containing package sources
	ImportPath     string        `json:",omitempty"` // import path of package in dir
	ImportComment  string        `json:",omitempty"` // path in import comment on package statement
	Name           string        `json:",omitempty"` // package name
	Doc            string        `json:",omitempty"` // package documentation string
	Target         string        `json:",omitempty"` // installed target for this package (may be executable)
	Shlib          string        `json:",omitempty"` // the shared library that contains this package (only set when -linkshared)
	Root           string        `json:",omitempty"` // Go root, Go path dir, or module root dir containing this package
	ConflictDir    string        `json:",omitempty"` // Dir is hidden by this other directory
	ForTest        string        `json:",omitempty"` // package is only for use in named test
	Export         string        `json:",omitempty"` // file containing export data (set by go list -export)
	BuildID        string        `json:",omitempty"` // build ID of the compiled package (set by go list -export)
	Module         *ModulePublic `json:",omitempty"` // info about package's module, if any
	Match          []string      `json:",omitempty"` // command-line patterns matching this package
	Goroot         bool          `json:",omitempty"` // is this package found in the Go root?
	Standard       bool          `json:",omitempty"` // is this package part of the standard Go library?
	DepOnly        bool          `json:",omitempty"` // package is only as a dependency, not explicitly listed
	BinaryOnly     bool          `json:",omitempty"` // package cannot be recompiled
	Incomplete     bool          `json:",omitempty"` // was there an error loading this package or dependencies?
	DefaultGODEBUG string        `json:",omitempty"` // default GODEBUG setting (only for Name=="main")
	// Stale and StaleReason remain here *only* for the list command.
	// They are only initialized in preparation for list execution.
	// The regular build determines staleness on the fly during action execution.
	Stale       bool   `json:",omitempty"` // would 'go install' do anything for this package?
	StaleReason string `json:",omitempty"` // why is Stale true?
	// Source files
	// If you add to this list you MUST add to p.AllFiles (below) too.
	// Otherwise file name security lists will not apply to any new additions.
	GoFiles           []string `json:",omitempty"` // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles          []string `json:",omitempty"` // .go source files that import "C"
	CompiledGoFiles   []string `json:",omitempty"` // .go output from running cgo on CgoFiles
	IgnoredGoFiles    []string `json:",omitempty"` // .go source files ignored due to build constraints
	InvalidGoFiles    []string `json:",omitempty"` // .go source files with detected problems (parse error, wrong package name, and so on)
	IgnoredOtherFiles []string `json:",omitempty"` // non-.go source files ignored due to build constraints
	CFiles            []string `json:",omitempty"` // .c source files
	CXXFiles          []string `json:",omitempty"` // .cc, .cpp and .cxx source files
	MFiles            []string `json:",omitempty"` // .m source files
	HFiles            []string `json:",omitempty"` // .h, .hh, .hpp and .hxx source files
	FFiles            []string `json:",omitempty"` // .f, .F, .for and .f90 Fortran source files
	SFiles            []string `json:",omitempty"` // .s source files
	SwigFiles         []string `json:",omitempty"` // .swig files
	SwigCXXFiles      []string `json:",omitempty"` // .swigcxx files
	SysoFiles         []string `json:",omitempty"` // .syso system object files added to package
	// Embedded files
	EmbedPatterns []string `json:",omitempty"` // //go:embed patterns
	EmbedFiles    []string `json:",omitempty"` // files matched by EmbedPatterns
	// Cgo directives
	CgoCFLAGS    []string `json:",omitempty"` // cgo: flags for C compiler
	CgoCPPFLAGS  []string `json:",omitempty"` // cgo: flags for C preprocessor
	CgoCXXFLAGS  []string `json:",omitempty"` // cgo: flags for C++ compiler
	CgoFFLAGS    []string `json:",omitempty"` // cgo: flags for Fortran compiler
	CgoLDFLAGS   []string `json:",omitempty"` // cgo: flags for linker
	CgoPkgConfig []string `json:",omitempty"` // cgo: pkg-config names
	// Dependency information
	Imports   []string          `json:",omitempty"` // import paths used by this package
	ImportMap map[string]string `json:",omitempty"` // map from source import to ImportPath (identity entries omitted)
	Deps      []string          `json:",omitempty"` // all (recursively) imported dependencies
	// Error information
	// Incomplete is above, packed into the other bools
	Error      *PackageError   `json:",omitempty"` // error loading this package (not dependencies)
	DepsErrors []*PackageError `json:",omitempty"` // errors loading dependencies, collected by go list before output
	// Test information
	// If you add to this list you MUST add to p.AllFiles (below) too.
	// Otherwise file name security lists will not apply to any new additions.
	TestGoFiles        []string `json:",omitempty"` // _test.go files in package
	TestImports        []string `json:",omitempty"` // imports from TestGoFiles
	TestEmbedPatterns  []string `json:",omitempty"` // //go:embed patterns
	TestEmbedFiles     []string `json:",omitempty"` // files matched by TestEmbedPatterns
	XTestGoFiles       []string `json:",omitempty"` // _test.go files outside package
	XTestImports       []string `json:",omitempty"` // imports from XTestGoFiles
	XTestEmbedPatterns []string `json:",omitempty"` // //go:embed patterns
	XTestEmbedFiles    []string `json:",omitempty"` // files matched by XTestEmbedPatterns
}

// https://go.googlesource.com/go/+/refs/tags/go1.22.4/src/cmd/go/internal/load/pkg.go#455

// The `Err` field has been changed to a `jsonErr` to allow for JSON unmarshalling.
type jsonErr string

func (e jsonErr) Error() string { return string(e) }

// A PackageError describes an error loading information about a package.
type PackageError struct {
	ImportStack      []string // shortest path from package named on command line to this one
	Pos              string   // position of error
	Err              jsonErr  // the error itself
	IsImportCycle    bool     // the error is an import cycle
	Hard             bool     // whether the error is soft or hard; soft errors are ignored in some places
	alwaysPrintStack bool     // whether to always print the ImportStack
}

func (p *PackageError) Error() string {
	// TODO(#43696): decide when to print the stack or the position based on
	// the error type and whether the package is in the main module.
	// Document the rationale.
	if p.Pos != "" && (len(p.ImportStack) == 0 || !p.alwaysPrintStack) {
		// Omit import stack. The full path to the file where the error
		// is the most important thing.
		return p.Pos + ": " + p.Err.Error()
	}
	// If the error is an ImportPathError, and the last path on the stack appears
	// in the error message, omit that path from the stack to avoid repetition.
	// If an ImportPathError wraps another ImportPathError that matches the
	// last path on the stack, we don't omit the path. An error like
	// "package A imports B: error loading C caused by B" would not be clearer
	// if "imports B" were omitted.
	if len(p.ImportStack) == 0 {
		return p.Err.Error()
	}
	var optpos string
	if p.Pos != "" {
		optpos = "\n\t" + p.Pos
	}
	return "package " + strings.Join(p.ImportStack, "\n\timports ") + optpos + ": " + p.Err.Error()
}
func (p *PackageError) Unwrap() error { return p.Err }

// https://go.googlesource.com/go/+/refs/tags/go1.22.4/src/cmd/go/internal/modinfo/info.go#16

// The `Origin` field has been modified to use a package-local type.

type ModulePublic struct {
	Path       string        `json:",omitempty"` // module path
	Version    string        `json:",omitempty"` // module version
	Query      string        `json:",omitempty"` // version query corresponding to this version
	Versions   []string      `json:",omitempty"` // available module versions
	Replace    *ModulePublic `json:",omitempty"` // replaced by this module
	Time       *time.Time    `json:",omitempty"` // time version was created
	Update     *ModulePublic `json:",omitempty"` // available update (with -u)
	Main       bool          `json:",omitempty"` // is this the main module?
	Indirect   bool          `json:",omitempty"` // module is only indirectly needed by main module
	Dir        string        `json:",omitempty"` // directory holding local copy of files, if any
	GoMod      string        `json:",omitempty"` // path to go.mod file describing module, if any
	GoVersion  string        `json:",omitempty"` // go version used in module
	Retracted  []string      `json:",omitempty"` // retraction information, if any (with -retracted or -u)
	Deprecated string        `json:",omitempty"` // deprecation message, if any (with -u)
	Error      *ModuleError  `json:",omitempty"` // error loading module
	Origin     *Origin       `json:",omitempty"` // provenance of module
	Reuse      bool          `json:",omitempty"` // reuse of old module info is safe
}
type ModuleError struct {
	Err string // error text
}
type moduleErrorNoMethods ModuleError

// UnmarshalJSON accepts both {"Err":"text"} and "text",
// so that the output of go mod download -json can still
// be unmarshalled into a ModulePublic during -reuse processing.
func (e *ModuleError) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		return json.Unmarshal(data, &e.Err)
	}
	return json.Unmarshal(data, (*moduleErrorNoMethods)(e))
}

// https://go.googlesource.com/go/+/refs/tags/go1.22.4/src/cmd/go/internal/modfetch/codehost/codehost.go#90

// An Origin describes the provenance of a given repo method result.
// It can be passed to CheckReuse (usually in a different go command invocation)
// to see whether the result remains up-to-date.
type Origin struct {
	VCS    string `json:",omitempty"` // "git" etc
	URL    string `json:",omitempty"` // URL of repository
	Subdir string `json:",omitempty"` // subdirectory in repo
	Hash   string `json:",omitempty"` // commit hash or ID
	// If TagSum is non-empty, then the resolution of this module version
	// depends on the set of tags present in the repo, specifically the tags
	// of the form TagPrefix + a valid semver version.
	// If the matching repo tags and their commit hashes still hash to TagSum,
	// the Origin is still valid (at least as far as the tags are concerned).
	// The exact checksum is up to the Repo implementation; see (*gitRepo).Tags.
	TagPrefix string `json:",omitempty"`
	TagSum    string `json:",omitempty"`
	// If Ref is non-empty, then the resolution of this module version
	// depends on Ref resolving to the revision identified by Hash.
	// If Ref still resolves to Hash, the Origin is still valid (at least as far as Ref is concerned).
	// For Git, the Ref is a full ref like "refs/heads/main" or "refs/tags/v1.2.3",
	// and the Hash is the Git object hash the ref maps to.
	// Other VCS might choose differently, but the idea is that Ref is the name
	// with a mutable meaning while Hash is a name with an immutable meaning.
	Ref string `json:",omitempty"`
	// If RepoSum is non-empty, then the resolution of this module version
	// failed due to the repo being available but the version not being present.
	// This depends on the entire state of the repo, which RepoSum summarizes.
	// For Git, this is a hash of all the refs and their hashes.
	RepoSum string `json:",omitempty"`
}
