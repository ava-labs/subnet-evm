# `gethclone`

This is an experimental module for tracking upstream `go-ethereum` changes.
The approach of `git merge`ing the upstream branch into `subnet-evm` is brittle as it relies on purely syntactic patching.
`gethclone` is intended to follow a rebase-like pattern, applying a set of semantic patches (e.g. AST modification, [Uber's `gopatch`](https://pkg.go.dev/github.com/uber-go/gopatch), etc.) that (a) should be more robust to refactoring; and (b) act as explicit documentation of the diffs.