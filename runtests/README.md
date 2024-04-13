#!/bin/bash

#### Jump to the folder in which the this script is in. Return if we fail to jump to the expected folder so we do not run tests in the wrong folder.

cd "$(dirname "$0")" || return

#### Jump out of runtests folder to root folder of the project

cd .. || return

#### Run the tests from the current package and its subfolder

go test ./...
