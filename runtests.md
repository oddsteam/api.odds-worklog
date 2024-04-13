#!/bin/bash

#### Jump to the folder in which the `runtests.md` script is in. Return if we fail to jump to the expected folder so we do not run tests in the wrong folder.

cd "$(dirname "$0")" || return

#### Run the tests from the current package and its subfolder

go test ./...
