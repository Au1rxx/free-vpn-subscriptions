#!/usr/bin/env python3

from __future__ import annotations

import argparse
import base64
import hashlib
import json
import sys
import urllib.request
from pathlib import Path

import yaml


REQUIRED_STATUS_KEYS = {
    "last_check_at",
    "name",
    "port",
    "protocol",
    "public_ip",
    "region",
    "status",
}


def ensure(condition: bool, message: str) -> None:
    if not condition:
        raise RuntimeError(message)


def fetch_bytes(url: str) -> bytes:
    with urllib.request.urlopen(url, timeout=30) as response:
        ensure(response.status == 200, f"{url} returned HTTP {response.status}")
        return response.read()


def sha256_bytes(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def validate_clash_bytes(data: bytes) -> dict[str, object]:
    parsed = yaml.safe_load(data.decode("utf-8"))
    ensure(isinstance(parsed, dict), "clash.yaml did not parse into a mapping")
    proxies = parsed.get("proxies")
    ensure(isinstance(proxies, list) and proxies, "clash.yaml has no proxies")
    return {
        "proxies": len(proxies),
        "keys": sorted(parsed.keys()),
    }


def validate_singbox_bytes(data: bytes) -> dict[str, object]:
    parsed = json.loads(data.decode("utf-8"))
    ensure(isinstance(parsed, dict), "singbox.json did not parse into an object")
    outbounds = parsed.get("outbounds")
    ensure(isinstance(outbounds, list) and outbounds, "singbox.json has no outbounds")
    return {
        "outbounds": len(outbounds),
        "keys": sorted(parsed.keys()),
    }


def validate_v2ray_bytes(data: bytes) -> dict[str, object]:
    decoded = base64.b64decode(data.decode("utf-8").strip() + "===")
    lines = [line for line in decoded.decode("utf-8", errors="replace").splitlines() if line.strip()]
    ensure(lines, "v2ray-base64.txt decoded into zero lines")
    ensure(
        all("://" in line for line in lines),
        "v2ray-base64.txt decoded lines do not look like subscription URIs",
    )
    return {
        "entries": len(lines),
        "sample_prefix": lines[0][:24],
    }


def validate_status_bytes(data: bytes) -> dict[str, object]:
    parsed = json.loads(data.decode("utf-8"))
    ensure(isinstance(parsed, list) and parsed, "status.json did not parse into a non-empty list")
    for index, node in enumerate(parsed):
        ensure(isinstance(node, dict), f"status.json node {index} is not an object")
        missing = REQUIRED_STATUS_KEYS - set(node.keys())
        ensure(not missing, f"status.json node {index} is missing keys: {sorted(missing)}")
    return {
        "nodes": len(parsed),
        "sample_keys": sorted(parsed[0].keys()),
    }


def validate_updates_json_bytes(data: bytes) -> dict[str, object]:
    parsed = json.loads(data.decode("utf-8"))
    ensure(isinstance(parsed, dict), "updates.json did not parse into an object")
    ensure(isinstance(parsed.get("items"), list), "updates.json has no items list")
    return {
        "items": len(parsed["items"]),
        "keys": sorted(parsed.keys()),
    }


def local_asset_paths(repo_root: Path) -> dict[str, Path]:
    output = repo_root / "output"
    return {
        "clash": output / "clash.yaml",
        "singbox": output / "singbox.json",
        "v2ray": output / "v2ray-base64.txt",
        "status": output / "status.json",
    }


def run_local(repo_root: Path) -> None:
    paths = local_asset_paths(repo_root)
    validators = {
        "clash": validate_clash_bytes,
        "singbox": validate_singbox_bytes,
        "v2ray": validate_v2ray_bytes,
        "status": validate_status_bytes,
    }

    print("LOCAL VALIDATION")
    for name, path in paths.items():
        ensure(path.exists(), f"missing expected file: {path}")
        data = path.read_bytes()
        details = validators[name](data)
        print(f"- {name}: ok sha256={sha256_bytes(data)} details={json.dumps(details, ensure_ascii=True)}")


def run_live(owner: str, repo: str, site_base: str) -> None:
    repo_base = f"https://raw.githubusercontent.com/{owner}/{repo}/main/output"
    release_base = f"https://github.com/{owner}/{repo}/releases/latest/download"
    pages = {
        "site_home": f"{site_base}/",
        "site_updates": f"{site_base}/updates.html",
        "site_status": f"{site_base}/status.html",
        "site_verification": f"{site_base}/verification.html",
    }
    assets = {
        "clash": (f"{repo_base}/clash.yaml", f"{release_base}/clash.yaml", validate_clash_bytes),
        "singbox": (f"{repo_base}/singbox.json", f"{release_base}/singbox.json", validate_singbox_bytes),
        "v2ray": (f"{repo_base}/v2ray-base64.txt", f"{release_base}/v2ray-base64.txt", validate_v2ray_bytes),
        "status": (f"{repo_base}/status.json", f"{release_base}/status.json", validate_status_bytes),
    }

    print("LIVE VALIDATION")
    for name, url in pages.items():
        data = fetch_bytes(url)
        ensure(len(data) > 0, f"{url} returned an empty body")
        print(f"- {name}: ok bytes={len(data)}")

    updates_data = fetch_bytes(f"{site_base}/updates.json")
    updates_details = validate_updates_json_bytes(updates_data)
    print(f"- updates_json: ok details={json.dumps(updates_details, ensure_ascii=True)}")

    for name, (raw_url, rel_url, validator) in assets.items():
        raw_data = fetch_bytes(raw_url)
        rel_data = fetch_bytes(rel_url)
        raw_details = validator(raw_data)
        rel_details = validator(rel_data)
        raw_hash = sha256_bytes(raw_data)
        rel_hash = sha256_bytes(rel_data)
        ensure(raw_hash == rel_hash, f"{name} raw/release hash mismatch: {raw_hash} != {rel_hash}")
        print(
            f"- {name}: ok raw_sha256={raw_hash} "
            f"details={json.dumps(raw_details, ensure_ascii=True)} "
            f"release_details={json.dumps(rel_details, ensure_ascii=True)}"
        )


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Validate public subscription artifacts and live endpoints.")
    parser.add_argument("--mode", choices=["local", "live"], required=True)
    parser.add_argument("--repo-root", default=".")
    parser.add_argument("--owner", default="Au1rxx")
    parser.add_argument("--repo", default="free-vpn-subscriptions")
    parser.add_argument("--site-base", default="https://au1rxx.github.io/free-vpn-subscriptions")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    try:
        if args.mode == "local":
            run_local(Path(args.repo_root).resolve())
        else:
            run_live(args.owner, args.repo, args.site_base.rstrip("/"))
    except Exception as exc:  # noqa: BLE001
        print(f"VALIDATION FAILED: {exc}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
