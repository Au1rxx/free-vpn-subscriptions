#!/usr/bin/env python3

from __future__ import annotations

import argparse
import json
from collections import Counter
from pathlib import Path

SITE_URL = "https://au1rxx.github.io/free-vpn-subscriptions/"
UPDATES_URL = "https://au1rxx.github.io/free-vpn-subscriptions/updates.html"
STATUS_URL = "https://au1rxx.github.io/free-vpn-subscriptions/status.html"
VERIFICATION_URL = "https://au1rxx.github.io/free-vpn-subscriptions/verification.html"
DISCUSSIONS_URL = "https://github.com/Au1rxx/free-vpn-subscriptions/discussions"


def load_nodes(path: Path) -> list[dict]:
    if not path.exists():
        return []

    with path.open("r", encoding="utf-8") as handle:
        data = json.load(handle)

    return data if isinstance(data, list) else []


def format_timestamp(value: str | None) -> str:
    if not value:
        return "无"

    return value.replace("T", " ").replace("Z", " UTC")


def format_counter(counter: Counter) -> str:
    if not counter:
        return "无"

    parts = [f"{key} x{counter[key]}" for key in sorted(counter)]
    return "、".join(parts)


def format_names(values: list[str]) -> str:
    return "、".join(values) if values else "无"


def summarize(nodes: list[dict]) -> dict:
    active = [node for node in nodes if str(node.get("status", "")).lower() == "active"]
    latest_check = sorted(
        node.get("last_check_at") for node in nodes if node.get("last_check_at")
    )

    return {
        "total": len(nodes),
        "active": len(active),
        "regions": Counter(node.get("region", "unknown") for node in nodes),
        "protocols": Counter(node.get("protocol", "unknown") for node in nodes),
        "latest_check": latest_check[-1] if latest_check else None,
    }


def compare_nodes(current: list[dict], previous: list[dict]) -> dict:
    current_map = {node.get("name"): node for node in current if node.get("name")}
    previous_map = {node.get("name"): node for node in previous if node.get("name")}

    current_names = set(current_map)
    previous_names = set(previous_map)

    added = sorted(current_names - previous_names)
    removed = sorted(previous_names - current_names)

    status_changes = []
    endpoint_changes = []

    for name in sorted(current_names & previous_names):
        current_node = current_map[name]
        previous_node = previous_map[name]

        current_status = str(current_node.get("status", "unknown")).lower()
        previous_status = str(previous_node.get("status", "unknown")).lower()
        if current_status != previous_status:
            status_changes.append(f"{name}: {previous_status} -> {current_status}")

        current_endpoint = (
            current_node.get("protocol"),
            current_node.get("public_ip"),
            current_node.get("port"),
        )
        previous_endpoint = (
            previous_node.get("protocol"),
            previous_node.get("public_ip"),
            previous_node.get("port"),
        )
        if current_endpoint != previous_endpoint:
            endpoint_changes.append(
                f"{name}: "
                f"{previous_node.get('protocol', 'unknown')} "
                f"{previous_node.get('public_ip', 'unknown')}:{previous_node.get('port', 'unknown')} -> "
                f"{current_node.get('protocol', 'unknown')} "
                f"{current_node.get('public_ip', 'unknown')}:{current_node.get('port', 'unknown')}"
            )

    current_active = sum(
        1 for node in current if str(node.get("status", "")).lower() == "active"
    )
    previous_active = sum(
        1 for node in previous if str(node.get("status", "")).lower() == "active"
    )

    return {
        "added": added,
        "removed": removed,
        "status_changes": status_changes,
        "endpoint_changes": endpoint_changes,
        "active_delta": current_active - previous_active,
        "current_active": current_active,
        "previous_active": previous_active,
    }


def build_markdown(current: list[dict], previous: list[dict]) -> str:
    current_summary = summarize(current)
    comparison = compare_nodes(current, previous) if previous else None

    headline = (
        f"当前公开节点 {current_summary['total']} 个，"
        f"active {current_summary['active']} 个，"
        f"覆盖 {len(current_summary['regions'])} 个地区，"
        f"协议分布为 {format_counter(current_summary['protocols'])}。"
    )

    sections = [
        headline,
        "",
        "## 本次公开快照摘要",
        f"- 当前可用节点：{current_summary['active']} / {current_summary['total']}",
        f"- 协议分布：{format_counter(current_summary['protocols'])}",
        f"- 地区覆盖：{format_counter(current_summary['regions'])}",
        f"- 最近状态检查：{format_timestamp(current_summary['latest_check'])}",
    ]

    if comparison:
        delta_prefix = "+" if comparison["active_delta"] > 0 else ""
        sections.extend(
            [
                "",
                "## 与上一快照相比",
                (
                    f"- Active 变化：{comparison['previous_active']} -> "
                    f"{comparison['current_active']} ({delta_prefix}{comparison['active_delta']})"
                ),
                f"- 新增节点：{format_names(comparison['added'])}",
                f"- 移除节点：{format_names(comparison['removed'])}",
                f"- 状态变化：{format_names(comparison['status_changes'])}",
                f"- 端点变化：{format_names(comparison['endpoint_changes'])}",
            ]
        )
    else:
        sections.extend(
            [
                "",
                "## 与上一快照相比",
                "- 这是当前公开仓库检测到的首个快照，暂时没有更早的 Release 用来对比。",
            ]
        )

    sections.extend(
        [
            "",
            "## 为什么主订阅链接保持稳定",
            "- 主订阅直链继续服务客户端自动刷新，不会因为历史快照发布而变更。",
            "- Release 负责历史回看、手动下载、发布提醒和增长回访。",
            f"- 更新记录页：`{UPDATES_URL}`",
            f"- 状态页：`{STATUS_URL}`",
            "",
            "## 为什么值得回来看",
            "- 看今天有没有新快照，而不是等主链接变化。",
            "- 看当前状态是否波动，再决定要不要排障或手动回滚。",
            "- 看有没有新的客户端教程、兼容性说明和 Discussions 讨论。",
            "",
            "## 继续跟踪",
            f"- 站点首页：`{SITE_URL}`",
            f"- 更新记录：`{UPDATES_URL}`",
            f"- 验证说明：`{VERIFICATION_URL}`",
            f"- Discussions：`{DISCUSSIONS_URL}`",
        ]
    )

    return "\n".join(sections).strip() + "\n"


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Build release notes for feed snapshots.")
    parser.add_argument("--current", required=True, type=Path)
    parser.add_argument("--previous", required=False, type=Path)
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    current_nodes = load_nodes(args.current)
    previous_nodes = load_nodes(args.previous) if args.previous else []
    print(build_markdown(current_nodes, previous_nodes), end="")


if __name__ == "__main__":
    main()
