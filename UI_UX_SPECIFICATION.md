# UI DEVELOPMENT PLAN

**Stack:** Go (Gin) + HTML templates + HTMX (no heavy JS) + Tailwind CSS (utility-first, small footprint).
**Goal:** Deliver an operator-friendly web app that covers the whole workflow with real-time feedback and minimum clicks.

─────────────────────────────────────────────

## FOUNDATION (½ week)
─────────────────────────────────────────────
1.1 Create folders
    templates/          – HTML files
    templates/partials/ – nav, modals, cards…
    assets/css/         – Tailwind build
    web/static/         – logo, sounds, ico

1.2 Tooling
    • npm init -y → tailwindcss + autoprefixer + postcss-cli
    • make ui-dev → tailwind --watch, Gin hot-reload

1.3 Base layout `base.tmpl`
    ├─ <head> Tailwind build, htmx.min.js, _hyperscript.js
    └─ Header, side-nav, main <div id="page"> {{ template "body" . }}
    Dark-mode toggle stored in localStorage.

─────────────────────────────────────────────
## 2. AUTH & ROLE ROUTING (½ week)
─────────────────────────────────────────────
• `/login` returns `hx-redirect` on success.
• Middleware injects `.CurrentUser` and `.Role`; side-nav partial renders items permitted by role.

─────────────────────────────────────────────
## 3. DASHBOARD (real-time) (1 week)
─────────────────────────────────────────────
3.1 Cards (HTMX polling every 5 s, `hx-get="/api/stats"`, `hx-swap="outerHTML"`)
    • Online modems / Total
    • SIMs low balance
    • Live calls (P-bar with ASR/ACD)
    • Spam calls blocked today

3.2 Live CDR ticker
    ```html
    <tbody hx-get="/ui/cdr/stream" hx-trigger="load, every 3s" hx-swap="beforeend scroll:bottom">
    </tbody>
    ```

─────────────────────────────────────────────
## 4. CRUD MODULES (2 weeks)
─────────────────────────────────────────────
4.1 Modems
    `list_modems.tmpl` (server-side pagination)
    Row actions: “Details”, “Reboot USB”, “Disable” via `hx-post`.

4.2 SIMs
    `list_sims.tmpl` + filter pills (status, operator).
    Modal becomes form target `hx-post="/ui/sim/{{.ID}}/edit"`.

4.3 Customers
    `list_customers.tmpl`, rate-sheet upload (`hx-put` multipart).

4.4 CDR explorer
    Date-range picker → `hx-get` with query params returns table partial.

─────────────────────────────────────────────
## 5. OPERATIONS TOOLS (1 week)
─────────────────────────────────────────────
• Bulk recharge wizard
    Stepper component → HTMX `hx-boost` forms.
• Blacklist editor
    Simple textarea with `hx-put`.
• Alert centre
    Streams Alertmanager JSON → SSE endpoint → htmx-sse extension.

─────────────────────────────────────────────
## 6. SETTINGS & PROFILE (½ week)
─────────────────────────────────────────────
• User profile, 2-FA QR, token regeneration.
• System settings (env overrides) behind Super-Admin role.

─────────────────────────────────────────────
## 7. POLISH & ACCESSIBILITY (½ week)
─────────────────────────────────────────────
• Form error snippets (`hx-target="this"`).
• Tab order, aria-labels, keyboard shortcuts (Hotkeys lib or _hyperscript).
• Responsive tweaks for ≥ 1366px + tablet.

─────────────────────────────────────────────
## 8. DELIVERY & CI/CD (½ week)
─────────────────────────────────────────────
• `npm run build` → minified Tailwind to `assets/css/bundle.css`.
• `go:embed` templates + assets in binary for single-file deploy.
• GitHub Actions: on push run `go vet/test`, `tailwind build`, `govulncheck`, then `docker build & push`.

**Total:** ~6 weeks, 8 sprints (4-day each).
After each sprint the UI is functional and demo-ready; HTMX allows incremental server-side additions without JS refactors.

### Key HTMX patterns to follow
• `hx-boost="true"` on `<body>` → links become AJAX by default.
• `hx-target="#page" hx-swap="innerHTML transition:true"` for page-level nav.
• Re-use partials: every server endpoint returns full or fragment depending on `Accept` header (`text/html` vs `hx-request`).

### Design language
• Tailwind + Flowbite components for tables, forms, modals.
• Colour code: green = healthy, amber = warning, red = blocked.

*With this plan backend and UI evolve together, zero React, clear increments, and your staff can operate the system by the end of Sprint 3 even while later modules are still polished.*
