const RELEASES_ENDPOINT =
  "https://api.github.com/repos/Au1rxx/free-vpn-subscriptions/releases?per_page=6";

const dashboardState = {
  activeCountText: "",
  latestCheck: "",
  latestReleaseStamp: "",
};

function formatUtcStamp(value) {
  return value ? value.replace("T", " ").replace("Z", " UTC") : "";
}

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

function extractReleaseSummary(body) {
  if (!body) {
    return "";
  }

  const blocks = body
    .split(/\n\s*\n/)
    .map((block) => block.replace(/\s+/g, " ").trim())
    .filter(Boolean);

  return blocks[0] || "";
}

function renderFreshnessBanner() {
  const freshnessTitle = document.getElementById("freshness-title");
  const freshnessDetail = document.getElementById("freshness-detail");

  if (!freshnessTitle || !freshnessDetail) {
    return;
  }

  if (!dashboardState.latestCheck) {
    freshnessTitle.textContent = "公开状态时间未知";
    freshnessDetail.textContent =
      "先打开状态页确认最近检查时间，再决定是否导入、刷新，或去 Releases 看最近快照。";
    return;
  }

  const diffHours = (Date.now() - Date.parse(dashboardState.latestCheck)) / (1000 * 60 * 60);
  const statusStamp = formatUtcStamp(dashboardState.latestCheck);
  const releaseStamp = dashboardState.latestReleaseStamp
    ? formatUtcStamp(dashboardState.latestReleaseStamp)
    : "";
  const today = new Date().toISOString().slice(0, 10);
  const statusDay = dashboardState.latestCheck ? dashboardState.latestCheck.slice(0, 10) : "";
  const releaseDay = dashboardState.latestReleaseStamp
    ? dashboardState.latestReleaseStamp.slice(0, 10)
    : "";

  if (releaseDay === today) {
    freshnessTitle.textContent = "今天已更新";
  } else if (statusDay === today && diffHours <= 3) {
    freshnessTitle.textContent = "今天状态已刷新";
  } else if (Number.isNaN(diffHours)) {
    freshnessTitle.textContent = "公开状态已发布";
  } else if (diffHours <= 3) {
    freshnessTitle.textContent = "公开状态在线";
  } else if (diffHours <= 12) {
    freshnessTitle.textContent = "公开状态略有延迟";
  } else {
    freshnessTitle.textContent = "公开状态可能滞后";
  }

  const parts = [];
  if (dashboardState.activeCountText) {
    parts.push(dashboardState.activeCountText);
  }
  parts.push(`最近状态检查 ${statusStamp}`);
  if (releaseStamp) {
    parts.push(`最近快照发布 ${releaseStamp}`);
  }
  parts.push("想持续跟踪更新就去看 Releases / Watch，或订阅 Feed");
  freshnessDetail.textContent = `${parts.join("。")}。`;
}

async function loadStatus() {
  const activeCount = document.getElementById("active-count");
  const statusStamp = document.getElementById("status-stamp");
  const nodeList = document.getElementById("node-list");

  if (!activeCount || !statusStamp || !nodeList) {
    return;
  }

  try {
    const response = await fetch("./output/status.json", { cache: "no-store" });
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const nodes = await response.json();
    const active = nodes.filter((node) => String(node.status).toLowerCase() === "active");
    const latestCheck = nodes
      .map((node) => node.last_check_at)
      .filter(Boolean)
      .sort()
      .pop();

    const activeText = `${active.length} / ${nodes.length} active`;
    dashboardState.activeCountText = activeText;
    dashboardState.latestCheck = latestCheck || "";

    activeCount.textContent = activeText;
    statusStamp.textContent = latestCheck
      ? `Last check ${formatUtcStamp(latestCheck)}`
      : "No recent status timestamp";

    nodeList.innerHTML = nodes
      .map((node) => {
        const activeNode = String(node.status).toLowerCase() === "active";
        const latency = node.latency_ms ? `${node.latency_ms} ms` : "Protocol-level check";
        return `
          <article class="node-item">
            <div class="node-meta">
              <strong>${node.name}</strong>
              <span>${node.region}</span>
            </div>
            <div class="node-meta">
              <strong>${node.protocol}</strong>
              <span>${node.public_ip}:${node.port}</span>
            </div>
            <div class="node-meta">
              <strong>${latency}</strong>
              <span>${node.last_check_at ? formatUtcStamp(node.last_check_at) : "No timestamp"}</span>
            </div>
            <div>
              <span class="status-pill ${activeNode ? "active" : "offline"}">
                ${activeNode ? "Active" : "Offline"}
              </span>
            </div>
          </article>
        `;
      })
      .join("");
    renderFreshnessBanner();
  } catch (error) {
    dashboardState.activeCountText = "";
    dashboardState.latestCheck = "";
    activeCount.textContent = "Unavailable";
    statusStamp.textContent = "Could not load current status";
    nodeList.innerHTML = `<p class="muted">The public status feed is not available yet: ${error.message}</p>`;
    renderFreshnessBanner();
  }
}

async function loadReleases() {
  const latestRelease = document.getElementById("latest-release");
  const latestReleaseSummary = document.getElementById("latest-release-summary");
  const releaseList = document.getElementById("release-list");

  if (!latestRelease && !latestReleaseSummary && !releaseList) {
    return;
  }

  try {
    const response = await fetch(RELEASES_ENDPOINT, {
      cache: "no-store",
      headers: {
        Accept: "application/vnd.github+json",
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const releases = (await response.json()).filter(
      (release) => !release.draft && !release.prerelease,
    );

    if (!releases.length) {
      throw new Error("No public releases");
    }

    const [latest] = releases;
    const latestTitle = latest.name || latest.tag_name;
    const latestStamp = formatUtcStamp(latest.published_at);
    const latestSummary = extractReleaseSummary(latest.body);
    dashboardState.latestReleaseStamp = latest.published_at || "";

    if (latestRelease) {
      latestRelease.textContent = latestStamp
        ? `${latestTitle} · ${latestStamp}`
        : latestTitle;
    }

    if (latestReleaseSummary) {
      latestReleaseSummary.textContent = latestSummary || "最近一次快照的变化摘要暂不可用。";
    }

    if (releaseList) {
      releaseList.innerHTML = releases
        .map((release) => {
          const title = escapeHtml(release.name || release.tag_name);
          const stamp = formatUtcStamp(release.published_at);
          const assetCount = Array.isArray(release.assets) ? release.assets.length : 0;
          const summary = extractReleaseSummary(release.body);
          const assetsText = assetCount
            ? `包含 ${assetCount} 个下载文件`
            : "发布页内含当前快照下载入口";

          return `
            <article class="release-item">
              <div class="release-item-copy">
                <p class="release-meta">${escapeHtml(release.tag_name)}</p>
                <h3>${title}</h3>
                ${summary ? `<p class="release-summary">${escapeHtml(summary)}</p>` : ""}
                <p>${stamp ? `发布时间 ${stamp}，` : ""}${assetsText}。适合手动下载、回看历史版本，或通过 Watch releases 跟踪更新。</p>
              </div>
              <div class="stack-actions">
                <a class="button button-secondary" href="${release.html_url}">查看发布</a>
              </div>
            </article>
          `;
        })
        .join("");
    }
    renderFreshnessBanner();
  } catch (error) {
    dashboardState.latestReleaseStamp = "";
    if (latestRelease) {
      latestRelease.textContent = "暂时无法读取最近发布记录";
    }

    if (latestReleaseSummary) {
      latestReleaseSummary.textContent = "稍后可直接打开 GitHub Releases 查看本次变化摘要。";
    }

    if (releaseList) {
      releaseList.innerHTML =
        `<p class="muted">暂时无法加载更新列表：${escapeHtml(error.message)}。你可以直接打开 GitHub Releases 查看历史快照。</p>`;
    }
    renderFreshnessBanner();
  }
}

loadStatus();
loadReleases();
