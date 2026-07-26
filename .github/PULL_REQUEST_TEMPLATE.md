## What and why

<!-- What changes, and what problem it solves. -->

## Checks

- [ ] `go build ./... && go vet ./... && go test ./...` all pass
- [ ] No credentials, tokens, or upstream source URLs added to files that ship publicly
- [ ] If this touches published output: node names still do not encode rank or
      position (unstable names rewrite the whole subscription every hour and
      throw away clients' latency history)
- [ ] If this touches `README*.md` or `docs/`: the change is in the generators
      under `internal/readme/` or `internal/pages/`, not in the generated files
      (those are overwritten by the next hourly publish)
