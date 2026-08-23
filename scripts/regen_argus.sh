#!/usr/bin/env bash
# Regenerate the argus (vision inference) client from
# specs/production/argus-openapi.json into ./argus.
#
# Same shape as regen.sh, for the same reason: fern-go-sdk's --local generation
# WIPES every non-generated file in its output directory, so we generate into a
# throwaway ./.gen-argus and rsync the result into ./argus. Run from the repo
# root.
#
# The synced spec artifact carries no `servers` array (argus's FastAPI app does
# not declare one), but the generated raw clients bake the first server URL in
# as the base-URL fallback. We want that fallback to be the production argus
# host — the same default platform-python hardcodes in its Client wrapper — so
# this script injects it into a working copy of the spec at generation time and
# restores the pristine artifact afterwards. If argus ever declares its server
# upstream, delete the injection.
set -euo pipefail

ARGUS_BASE_URL="https://argus.axilio.ai"
SPEC="specs/production/argus-openapi.json"

command -v fern >/dev/null || { echo "fern CLI not found: npm i -g fern-api" >&2; exit 1; }

# Inject the production server URL; put the pristine spec back no matter how
# generation exits so the committed artifact stays byte-identical to the sync.
cp "$SPEC" "$SPEC.orig"
trap 'mv "$SPEC.orig" "$SPEC"' EXIT
python3 - "$SPEC" "$ARGUS_BASE_URL" <<'EOF'
import json, sys
path, url = sys.argv[1], sys.argv[2]
spec = json.load(open(path))
if "servers" in spec:
    sys.exit(f"{path} now declares servers; drop the injection from regen_argus.sh")
spec["servers"] = [{"url": url}]
json.dump(spec, open(path, "w"), indent=2)
EOF

rm -rf .gen-argus
fern generate --local --group go-sdk --api argus --force --log-level warn

# --delete keeps ./argus an exact mirror of the generated tree. The excludes
# are the ONLY protection for non-generated files under ./argus (same mechanism
# as regen.sh): go.mod/go.sum go because argus is a package of the root module,
# not a nested module; CONTRIBUTING.md documents the .fernignore mechanism the
# local generator ignores (see regen.sh); baseurl_test.go is the hand-written
# guard pinning the servers injection above.
mkdir -p argus
rsync -a --delete \
  --exclude='go.mod' \
  --exclude='go.sum' \
  --exclude='CONTRIBUTING.md' \
  --exclude='baseurl_test.go' \
  .gen-argus/ argus/

rm -rf .gen-argus

# Same dead-surface strip as the backend client (see the script's header):
# argus has no streaming endpoints, so Fern's SSE reconnect options are
# doubly dead here.
(cd argus && go run ../scripts/strip_dead_stream_options.go)

# The strip leaves struct-field alignment behind; keep the tree gofmt-clean.
gofmt -w argus

# Fold the generated code's deps into the root module.
go mod tidy

echo "argus client regenerated from $SPEC"
