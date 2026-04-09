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

    activeCount.textContent = `${active.length} / ${nodes.length} active`;
    statusStamp.textContent = latestCheck
      ? `Last check ${latestCheck.replace("T", " ").replace("Z", " UTC")}`
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
              <span>${node.last_check_at ? node.last_check_at.replace("T", " ").replace("Z", " UTC") : "No timestamp"}</span>
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
  } catch (error) {
    activeCount.textContent = "Unavailable";
    statusStamp.textContent = "Could not load current status";
    nodeList.innerHTML = `<p class="muted">The public status feed is not available yet: ${error.message}</p>`;
  }
}

loadStatus();
