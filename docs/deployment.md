# Deployment And Operations

## End-to-end publishing path

The public repository does not generate node data by itself. The publishing chain is:

1. The private control repository `Au1rxx/vpn-lab` monitors nodes and generates public outputs.
2. The private repo synchronizes only `output/` into this public repo.
3. A push to this public repo triggers Pages deployment, release creation, update-feed refresh, and validation.

That split keeps public delivery and public docs in this repo, while secrets and control-plane state stay upstream.

## Workflow map

### Upstream private repository: `Au1rxx/vpn-lab`

#### `Monitor Nodes`

- File: `vpn-lab/.github/workflows/monitor.yml`
- Runner: self-hosted control VM
- Trigger: hourly cron plus manual dispatch
- Job:
  runs health checks, updates private state, and commits sanitized `output/status.json`

#### `Publish Subscriptions`

- File: `vpn-lab/.github/workflows/publish.yml`
- Runner: self-hosted control VM
- Trigger: every 6 hours plus manual dispatch
- Job:
  runs monitor, regenerates `output/` subscription files, refreshes the private repo README, and pushes changes

#### `Publish Public Feed`

- File: `vpn-lab/.github/workflows/publish-public.yml`
- Runner: GitHub-hosted `ubuntu-latest`
- Trigger:
  - push to `output/**`
  - successful `workflow_run` from `Monitor Nodes`
  - successful `workflow_run` from `Publish Subscriptions`
  - manual dispatch
- Job:
  checks out both repositories, syncs only `output/` into `Au1rxx/free-vpn-subscriptions`, and pushes the public update

### Public repository: `Au1rxx/free-vpn-subscriptions`

#### `Deploy Pages`

- File: `.github/workflows/pages.yml`
- Trigger: every push to `main`, plus manual dispatch
- Job:
  copies `site/` and `output/` into a Pages artifact and deploys to GitHub Pages

#### `Release Feed Snapshot`

- File: `.github/workflows/release-feed.yml`
- Trigger:
  - push to `main` when `output/clash.yaml`, `output/singbox.json`, `output/v2ray-base64.txt`, or `output/status.json` changes
  - manual dispatch
- Job:
  creates a tagged GitHub Release containing the four public artifacts, generates release notes from current versus previous status, and dispatches the amplify workflow

#### `Amplify Release Updates`

- File: `.github/workflows/release-amplify.yml`
- Trigger:
  - new published release
  - manual dispatch
- Job:
  regenerates `site/updates.json` and `site/updates.xml`, commits them back to `main`, and publishes an announcement discussion

#### `Validate Public Assets`

- File: `.github/workflows/validate-public-assets.yml`
- Trigger:
  - push affecting `output/**`, `site/**`, README files, or validation workflow/script
  - successful `workflow_run` from `Deploy Pages`, `Release Feed Snapshot`, or `Amplify Release Updates`
  - manual dispatch
  - schedule every 6 hours
- Job:
  validates both repository artifacts and live public endpoints

#### `Seed Discussions`

- File: `.github/workflows/seed-discussions.yml`
- Trigger: manual dispatch
- Job:
  seeds and preserves baseline announcement / Q&A / ideas discussions for the public repo

## Schedules

- Private node health monitoring: hourly
- Private subscription generation: every 6 hours
- Public live validation: every 6 hours
- Public Pages deployment: on every push to `main`
- Public release snapshots: whenever public `output/` changes

## Required secrets and settings

### Private repository secret

The upstream sync requires this secret in `Au1rxx/vpn-lab`:

- `PUBLIC_REPO_SYNC_TOKEN`

Recommended access:

- write access to `Au1rxx/free-vpn-subscriptions`
- `contents: write` is sufficient for the sync step

### Public repository Pages setting

GitHub Pages must use:

- Source: `GitHub Actions`

Once Pages is enabled, `Deploy Pages` handles later deployments automatically.

## What is deployed publicly

- GitHub Pages site from `site/`
- `output/` artifacts copied into the Pages artifact
- GitHub Release assets for each public snapshot
- Update feeds in `site/updates.json` and `site/updates.xml`
- Announcement discussions for new releases

## Manual operations

### Rebuild a release snapshot manually

```bash
gh workflow run release-feed.yml --repo Au1rxx/free-vpn-subscriptions
```

### Re-seed baseline discussions

```bash
gh workflow run seed-discussions.yml --repo Au1rxx/free-vpn-subscriptions
```

### Validate repository artifacts locally

```bash
python3 scripts/validate_public_assets.py --mode local --repo-root .
```

### Validate live public endpoints

```bash
python3 scripts/validate_public_assets.py --mode live
```

## Failure diagnosis

### Public raw files are fresh, but Pages is stale

Check:

- `Deploy Pages`

### Release assets are stale, but raw files changed

Check:

- `Release Feed Snapshot`

### Release exists, but `updates.json` / `updates.xml` is stale

Check:

- `Amplify Release Updates`

### Public repo stopped receiving fresh `output/`

Check upstream:

- `Publish Public Feed`
- `PUBLIC_REPO_SYNC_TOKEN`
- `Monitor Nodes`
- `Publish Subscriptions`

### Pages and releases look healthy, but users still report import failure

Check:

- the client-specific guide in `docs/`
- `output/status.json`
- `site/verification.html`
- whether the user imported a raw remote URL or a release snapshot file

## Maintainer checklist

Before merging workflow or deployment changes:

1. Confirm the public/private boundary is still intact.
2. Confirm direct raw URLs remain stable unless a breaking change is intentional.
3. Confirm release generation and update-feed rendering still work together.
4. Confirm the validator still checks both raw and release outputs.
5. Confirm Pages still includes both `site/` and `output/`.
