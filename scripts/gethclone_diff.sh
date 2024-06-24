#!/bin/bash

#
# Usage: scripts/gethclone_diff.sh {filepath}
#
# Example: `scripts/gethclone_diff.sh core/types/block.go | less`
#
# Convenience script for performing side-by-side diff of a file with the output
# of the `x/gethclone` command.
#

set -eu;

ROOT=$(git rev-parse --show-toplevel);
diff -y "${ROOT}/${1}" "${ROOT}/${1}.gethclone";
