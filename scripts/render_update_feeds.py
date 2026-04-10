#!/usr/bin/env python3

from __future__ import annotations

import argparse
import html
import json
from pathlib import Path


SITE_URL = "https://au1rxx.github.io/free-vpn-subscriptions/"
UPDATES_PAGE_URL = f"{SITE_URL}updates.html"
ATOM_URL = f"{SITE_URL}updates.xml"
JSON_URL = f"{SITE_URL}updates.json"


def load_releases(path: Path) -> list[dict]:
    with path.open("r", encoding="utf-8") as handle:
        data = json.load(handle)

    return [
        release
        for release in data
        if not release.get("draft") and not release.get("prerelease")
    ]


def first_paragraph(text: str) -> str:
    parts = [part.strip() for part in text.split("\n\n") if part.strip()]
    return parts[0] if parts else ""


def markdown_to_simple_html(text: str) -> str:
    blocks = [block.strip() for block in text.strip().split("\n\n") if block.strip()]
    rendered: list[str] = []

    for block in blocks:
      lines = [line.strip() for line in block.splitlines() if line.strip()]
      if lines and all(line.startswith("- ") for line in lines):
          rendered.append(
              "<ul>"
              + "".join(f"<li>{html.escape(line[2:])}</li>" for line in lines)
              + "</ul>"
          )
          continue

      title = None
      if lines and lines[0].startswith("## "):
          title = html.escape(lines[0][3:])
          lines = lines[1:]

      if title:
          rendered.append(f"<h2>{title}</h2>")

      if lines:
          rendered.append(f"<p>{html.escape(' '.join(lines))}</p>")

    return "\n".join(rendered)


def write_json_feed(releases: list[dict], output: Path) -> None:
    items = []
    for release in releases[:20]:
        body = release.get("body", "").strip()
        items.append(
            {
                "id": release["tag_name"],
                "url": release["html_url"],
                "title": release.get("name") or release["tag_name"],
                "content_text": body,
                "summary": first_paragraph(body),
                "date_published": release.get("published_at"),
            }
        )

    payload = {
        "version": "https://jsonfeed.org/version/1.1",
        "title": "免费 VPN 订阅更新",
        "home_page_url": UPDATES_PAGE_URL,
        "feed_url": JSON_URL,
        "description": "免费 VPN 订阅的最新快照、节点摘要与历史更新记录。",
        "icon": f"{SITE_URL}social-preview.png",
        "authors": [{"name": "Au1rxx"}],
        "language": "zh-CN",
        "items": items,
    }

    output.write_text(
        json.dumps(payload, ensure_ascii=False, indent=2) + "\n",
        encoding="utf-8",
    )


def write_atom_feed(releases: list[dict], output: Path) -> None:
    updated = releases[0].get("published_at") if releases else "2026-04-10T00:00:00Z"

    entries = []
    for release in releases[:20]:
        body = release.get("body", "").strip()
        title = html.escape(release.get("name") or release["tag_name"])
        summary = html.escape(first_paragraph(body))
        content_html = markdown_to_simple_html(body)
        release_url = html.escape(release["html_url"])
        tag = html.escape(release["tag_name"])
        published = release.get("published_at", "")

        entries.append(
            f"""  <entry>
    <title>{title}</title>
    <id>tag:github.com,2026:Au1rxx/free-vpn-subscriptions/{tag}</id>
    <link href="{release_url}" />
    <updated>{published}</updated>
    <published>{published}</published>
    <summary>{summary}</summary>
    <content type="html"><![CDATA[{content_html}]]></content>
  </entry>"""
        )

    feed = f"""<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>免费 VPN 订阅更新</title>
  <subtitle>订阅快照、节点摘要、状态变化和历史更新记录。</subtitle>
  <id>{ATOM_URL}</id>
  <link href="{SITE_URL}" rel="alternate" />
  <link href="{ATOM_URL}" rel="self" />
  <updated>{updated}</updated>
  <author>
    <name>Au1rxx</name>
  </author>
{"\n".join(entries)}
</feed>
"""
    output.write_text(feed, encoding="utf-8")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Render Atom and JSON update feeds.")
    parser.add_argument("--input", required=True, type=Path)
    parser.add_argument("--output-dir", required=True, type=Path)
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    releases = load_releases(args.input)
    args.output_dir.mkdir(parents=True, exist_ok=True)
    write_json_feed(releases, args.output_dir / "updates.json")
    write_atom_feed(releases, args.output_dir / "updates.xml")


if __name__ == "__main__":
    main()
