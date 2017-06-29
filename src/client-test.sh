#!/bin/bash
set -euo pipefail
source "$(dirname "$0")/vt6.inc"
vt6_connect "$0" "$@"

echo hallo >&3
