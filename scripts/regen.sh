#!/usr/bin/env bash
# Regenerate the platform-go SDK from specs/production/openapi.json.
#
# fern-go-sdk's --local generation WIPES every non-generated file in its output
# directory (verified: it deletes drivers/, .github/, VERSION even when they are
# .fernignore-listed). So we generate into a throwaway ./.gen and rsync the
# result up to the repo root, preserving hand-written paths: the Go MobileDriver
# (drivers/, AXI-1253), CI (.github/), the release VERSION, the fern config, and
# the spec. Run from the repo root.
set -euo pipefail

command -v fern >/dev/null || { echo "fern CLI not found: npm i -g fern-api" >&2; exit 1; }

rm -rf .gen
fern generate --local --group go-sdk --api backend --force --log-level warn

# --delete prunes generated files that no longer exist upstream, while the
# excludes keep every hand-written / scaffolding path off-limits. THIS LIST IS
# THE ONLY THING PROTECTING HAND-WRITTEN CODE: anything at the root that is not
# generated and not excluded here is deleted. Adding a path to .fernignore does
# nothing (see the header above) — add it here instead.
#
# CONTRIBUTING.md is excluded because Fern's generated version documents the
# .fernignore mechanism that this generator ignores in --local mode, i.e. it
# actively instructs contributors into silent data loss. The repo owns it.
# README.md is deliberately NOT excluded: it is mostly generated API reference
# that should track the spec.
rsync -a --delete \
  --exclude='.git' \
  --exclude='.gen' \
  --exclude='fern' \
  --exclude='specs' \
  --exclude='scripts' \
  --exclude='drivers' \
  --exclude='.github' \
  --exclude='VERSION' \
  --exclude='.gitignore' \
  --exclude='CONTRIBUTING.md' \
  .gen/ ./

rm -rf .gen

# fern-go-sdk regenerates go.mod/go.sum from the *generated* code's imports only,
# which drops the hand-written driver's deps (drivers/mobile pulls in
# coder/websocket). Re-tidy against the full tree (generated code + drivers/) so
# those deps survive every regen. Needs network (CI has it).
go mod tidy

echo "platform-go regenerated from specs/production/openapi.json"
