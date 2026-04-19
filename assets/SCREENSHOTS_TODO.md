# Screenshot Checklist

These screenshots dramatically improve star-conversion. Take them from a real
client with a real subscription imported from this repo, save as PNG into this
directory, then delete this file (or keep for future replacements).

All screenshots: **1600×1000** ideal, can be 1280×800 minimum. PNG, < 400 KB
each (optimize with `pngquant` or `tinypng`).

## 1. `assets/screenshot-clash-verge-connected.png`

**Priority: highest — goes at top of README.**

- Client: Clash Verge (dark or light, pick the cleaner one)
- Subscription: imported from `output/clash.yaml`
- State: connected to a **low-latency node** (< 50ms, ideally HK/JP/US)
- Composition: show the proxy group with the active node highlighted + RTT
  badge visible, with at least 5-6 other nodes listed below showing varied
  countries + latencies.
- Crop to remove OS chrome if on macOS/Windows (keep client window only).

## 2. `assets/screenshot-by-country.png`

- Client: any (Clash Verge or v2rayN fine)
- State: showing the **Proxy Groups** view where users can pick a country
  group, OR showing multiple country-specific subscriptions imported side by
  side in the profile list.
- Goal: visually communicate "we give you country-targeted subscriptions".

## 3. `assets/screenshot-v2rayng-mobile.png` (optional but nice)

- Client: v2rayNG on Android
- State: subscription imported, nodes listed with ping results.
- Portrait phone-shaped screenshot (e.g. 450×950), used in "Mobile" section.

## 4. `assets/demo.gif` (optional, but extremely high conversion)

- Length: 6-10 seconds, looped
- Content:
  1. Show README "One-Click Subscribe" URL
  2. Copy the URL (clipboard-ish animation)
  3. Paste into Clash Verge import field
  4. Import succeeds → nodes appear
  5. Click a node → "Connected" appears
- Tool: `peek` (Linux), `Kap` (macOS), `ScreenToGif` (Windows)
- < 2 MB (use gifsicle or convert to APNG if too large)

---

Until real screenshots exist, the README falls back to the SVG workflow
diagram at `assets/workflow.svg` (generated).
