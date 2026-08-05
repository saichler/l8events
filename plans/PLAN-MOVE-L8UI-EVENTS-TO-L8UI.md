# Plan: Move `l8ui/events/` from l8events into the canonical `../l8ui` library (with mobile parity)

## Context

`l8events` currently ships 9 reusable UI files at `l8events/l8ui/events/` (alarm/event/maintenance
enums, tables, detail views, archive viewer, state actions, CSS). Per the l8events README and
`plans/PLAN-L8EVENTS-SHARED-LIBRARY.md`, these were designed to be **copy-distributed**: each
consuming project runs `cp -r l8events/l8ui/events/ <project>/.../web/l8ui/events/` and wires the
scripts into its own `app.html`.

Since `l8events` is required infrastructure for every Layer 8 project
(`events-service-required.md`), the call is that its UI belongs inside the canonical `l8ui`
library itself — the same way `l8ui/sys/security`, `l8ui/sys/health`, `l8ui/sys/logs`, and
`l8ui/sys/dataimport` are built-in, domain-specific-but-universally-needed modules that ship with
`l8ui`, not something every consumer hand-copies. Once merged into `l8ui`, any project that adds
`l8ui` as a submodule (`l8ui-copy-to-new-project.md`) gets the events components for free — no
separate copy step, no drift.

### Revision history (why this plan looks the way it does)

1. **First draft**: pure file relocation, desktop-only. Deferred mobile as "out of scope, no
   equivalent exists."
2. **Corrected**: `mobile-rules.md` requires desktop/mobile parity — "no mobile file exists yet" is
   a gap to close during this touch, not a reason to defer it. Revised to add a Phase 2 that builds
   7 new `Layer8M*`-suffixed files (`-m.js`), one per existing desktop file, mirroring the
   `l8agent/m/` split pattern.
3. **Corrected again, this revision**: that 7-file mobile fork would have duplicated ~600+ lines of
   near-identical enum/column logic (differing only in which `Renderers` module built the `render`
   map, or a global-namespace swap) — well past the 100-line threshold where
   `plan-requirements.md`'s Duplication Audit makes an extraction phase mandatory. Investigating the
   actual runtime dependencies found something better than an extraction: `Layer8DRenderers` is
   already loaded on mobile pages (`mobile-script-loading-order.md` loads `layer8d-renderers.js`
   "needed for ... renderers"), and `l8security-enums.js` — confirmed shipped identically to both
   `app.html` and `m/app.html` in `l8erp`, no mobile fork exists for it anywhere — already proves
   this pattern works in production. **No mobile files need to be built at all.** The 9 existing
   files are reused directly on both platforms; only small in-place additions (mobile card
   `primary`/`secondary` markers) are needed. Real visual confirmation that this renders
   acceptably is out of scope for this plan (no runnable UI exists in `l8events` to check it in —
   see Phase 5) and happens later, by the user, in a consuming project.

### Verification already done
- **No naming collision**: none of the 9 filenames exist anywhere under `../l8ui` today.
- **Theme compliance**: `l8events.css` already uses `var(--layer8d-*)` exclusively (57 references,
  0 hardcoded colors, no `[data-theme]` blocks) — passes `l8ui-theme-compliance.md` as-is.
- **Mobile runtime dependency confirmed available**: `Layer8DRenderers` (and `Layer8EnumFactory`,
  `Layer8ColumnFactory`, `Layer8FormFactory`) are loaded on mobile pages per
  `mobile-script-loading-order.md` — not mobile-exclusive APIs, genuinely shared.
- **Real precedent, not assumption**: `l8security-enums.js` ships unforked to both `app.html` and
  `m/app.html` in `l8erp` today, built the same way `l8events-enums.js` is (`Layer8EnumFactory` +
  `layer8d-status-*` class names), confirming this is an established, running pattern — not a
  novel risk being introduced here.
- **Consumer drift**: `l8alarms` and `l8vendingmachine` already carry their own ad-hoc copies of
  the desktop files (not via `setup-l8ui-submodule.sh` — `l8alarms`'s `web/l8ui/` is a
  manually-vendored, independently-committed copy, not a real l8ui submodule). Reconciling those
  projects stays **out of scope** for this plan (separate repos, separate task) — see "Follow-ups."

## Design Decision: How Mobile Parity Is Achieved (Without Forking)

Read all 9 source files in full. None of them call `Layer8DTable`, `Layer8DPopup`, or any other
genuinely desktop-only widget class — they either produce pure data (`Layer8ColumnFactory`/
`Layer8FormFactory` output, platform-agnostic per `shared-schemas.md`) or plain HTML-string
generation via `container.innerHTML` (platform-agnostic DOM APIs). The only per-file dependency
that looked desktop-specific was `l8events-enums.js`/`l8events-category-enums.js` importing
`Layer8DRenderers` — and that dependency is satisfied at runtime on mobile too (see above).

| File | Mobile treatment |
|---|---|
| `l8events-enums.js` | **Reused as-is.** `Layer8DRenderers` is loaded on mobile; matches the `l8security-enums.js` precedent exactly. |
| `l8events-category-enums.js` | **Reused as-is.** Same reasoning. |
| `l8events-alarm-table.js` | **Reused as-is** + add `primary`/`secondary` markers to `getColumns()`'s column defs so `Layer8MTable`/`Layer8MEditTable` render sensible cards (desktop ignores these extra properties). |
| `l8events-event-viewer.js` | Same — reused as-is + `primary`/`secondary` markers. |
| `l8events-archive-viewer.js` | Same — reused as-is + `primary`/`secondary` markers on both `getArchivedAlarmColumns()` and `getArchivedEventColumns()`. |
| `l8events-maintenance.js` | Same — reused as-is + `primary`/`secondary` markers. |
| `l8events-alarm-detail.js` | **Reused as-is**, no changes — pure HTML-string generation, works inside a `Layer8MPopup` body the same as a `Layer8DPopup` body. |
| `l8events-state-actions.js` | **Reused as-is**, no changes — only touches shared `layer8d-btn`/`layer8d-btn-small` theme classes. |
| `l8events.css` | **Reused as-is**, no changes — already 100% `var(--layer8d-*)`-based. |

**Net result: zero new files, zero forked logic.** Four files get a small in-place edit (adding
`primary`/`secondary` properties to existing column arrays — a few lines each, not new files or
duplicated logic). This is what closes the `mobile-rules.md` parity gap without triggering
`plan-requirements.md`'s duplication-extraction requirement, because there's nothing to extract —
the files were never forked.

**Risk being carried forward, NOT resolved by this plan:** reusing `Layer8DRenderers`-built status
badges inside a mobile card context is the same choice `l8security-enums.js` already makes in
production (and `layer8d-btn` on `l8events-state-actions.js` has equally solid precedent — both
`l8ui/m/js/layer8m-table-touch.js` and a live consumer's mobile JS, `l8vendingmachine`'s
`m/js/routes/routes-generate.js`, use `layer8d-btn` directly), but neither guarantees it looks good
in *this* component's specific card layout. `l8events` has no web app of its own to test in — no
`app.html`, no `run-local.sh`, by design (UI is fully separated from the backend here). Real,
visual, in-browser verification is explicitly **out of scope for this plan** and will be done by
the user separately, in a consuming project, once `l8ui/events/` is wired in (see Phase 5). This
plan's own Phase 4 is limited to what can actually be checked without a running app: file content
and structure.

## What Moves

From `l8events/l8ui/events/` → `l8ui/events/` (top-level category, parallel to `l8ui/sys/`,
`l8ui/chart/`, etc.) — used directly by both desktop and mobile, no per-platform subdirectory:

```
l8events-enums.js            (must load first — others depend on it)                — unchanged
l8events-category-enums.js   (depends on enums)                                     — unchanged
l8events-state-actions.js    (loads before alarm-detail)                            — unchanged
l8events-alarm-table.js                                                             — + primary/secondary markers
l8events-alarm-detail.js                                                            — unchanged
l8events-event-viewer.js                                                            — + primary/secondary markers
l8events-archive-viewer.js                                                          — + primary/secondary markers
l8events-maintenance.js                                                             — + primary/secondary markers
l8events.css                                                                        — unchanged
```

Global names (`window.L8EventsEnums`, etc.) and desktop script load order are preserved exactly.
No mobile-specific script load order is needed — the same 9 files serve both `app.html` and
`m/app.html` includes.

## Platform Audit

| Component | Desktop File | Desktop Status | Mobile Status |
|---|---|---|---|
| Event/alarm/maintenance enums | `l8events-enums.js`, `l8events-category-enums.js` | Exists — relocating unchanged | Reused directly — `Layer8DRenderers` confirmed available on mobile pages |
| Alarm table | `l8events-alarm-table.js` | Exists — relocating + primary/secondary markers | Reused directly via the same file |
| Alarm detail | `l8events-alarm-detail.js` | Exists — relocating unchanged | Reused directly via the same file |
| Event viewer | `l8events-event-viewer.js` | Exists — relocating + primary/secondary markers | Reused directly via the same file |
| Archive viewer | `l8events-archive-viewer.js` | Exists — relocating + primary/secondary markers | Reused directly via the same file |
| Maintenance windows | `l8events-maintenance.js` | Exists — relocating + primary/secondary markers | Reused directly via the same file |
| State transition actions | `l8events-state-actions.js` | Exists — relocating unchanged | Reused directly via the same file |
| Styling | `l8events.css` | Exists — relocating unchanged | Reused directly via the same file; visual mobile-viewport check deferred to the user (Phase 5) |

## Traceability Matrix

| # | Section | Gap / Action Item | Platform | Phase |
|---|---------|-------------------|----------|-------|
| 1 | What Moves | Copy 9 files from `l8events/l8ui/events/` to `../l8ui/events/` | Both | Phase 1 |
| 2 | What Moves | Diff-verify the 5 unchanged files are byte-for-byte identical to source | Both | Phase 1 |
| 3 | Design Decision | Add `primary`/`secondary` markers to `l8events-alarm-table.js` | Mobile | Phase 1 |
| 4 | Design Decision | Add `primary`/`secondary` markers to `l8events-event-viewer.js` | Mobile | Phase 1 |
| 5 | Design Decision | Add `primary`/`secondary` markers to `l8events-archive-viewer.js` (both column sets) | Mobile | Phase 1 |
| 6 | Design Decision | Add `primary`/`secondary` markers to `l8events-maintenance.js` | Mobile | Phase 1 |
| 7 | Verification already done | Document the component (desktop + mobile usage) in `l8ui`'s own README (`Shared & Utilities` list) | Both | Phase 2 |
| 8 | Verification already done | Add `../l8ui/rules/l8events-ui.md` (prerequisites, load order, API surface) | Both | Phase 2 |
| 9 | Verification already done | Add files to `desktop-script-loading-order.md` and `mobile-script-loading-order.md`, conditional on those rule files existing in the l8ui repo | Both | Phase 2 |
| 10 | Context | Delete `l8events/l8ui/` (now redundant at the source) | Both | Phase 3 |
| 11 | Context | Rewrite `l8events/README.md` §"l8ui Components" to point at the `l8ui` submodule instead of copy-distribution, covering both surfaces | Both | Phase 3 |
| 12 | Context | Rewrite the integration steps section (`cp -r ...` instruction) | Both | Phase 3 |
| 13 | Context | Update `plans/PLAN-L8EVENTS-SHARED-LIBRARY.md` with a superseded/pointer note | Both | Phase 3 |
| 14 | Design Decision (risk) | Confirm in an actual mobile viewport that `Layer8DRenderers`-built status badges render acceptably inside `Layer8MTable` cards and the `Layer8MPopup` detail view | Mobile | Phase 5 — deferred to the user, in a separate project; any fix comes back here first |
| 15 | Follow-ups | Reconcile `l8alarms` / `l8vendingmachine`'s drifted local copies onto the new canonical source | Both | Deferred — separate repos, separate task, not requested |

Every row above resolves to a phase in this plan or an explicitly-deferred follow-up; nothing is
unaccounted for.

## Steps

### Phase 1 — Copy into `../l8ui`, add mobile card markers
1. Copy all 9 files from `l8events/l8ui/events/` to `../l8ui/events/` (new directory).
2. Diff-verify the 5 unchanged files (`l8events-enums.js`, `l8events-category-enums.js`,
   `l8events-state-actions.js`, `l8events-alarm-detail.js`, `l8events.css`) are byte-for-byte
   identical to the pre-move originals.
3. In `l8events-alarm-table.js`'s `getColumns()`: add `primary: true` to `name`, `secondary: true`
   to `state` (confirm during implementation which pairing reads best as a card title/subtitle —
   these are the working defaults, not locked in).
4. In `l8events-event-viewer.js`'s `getColumns()`: add `primary: true` to `eventType`,
   `secondary: true` to `severity`.
5. In `l8events-archive-viewer.js`: add `primary`/`secondary` markers to both
   `getArchivedAlarmColumns()` (e.g. `name` primary, `state` secondary) and
   `getArchivedEventColumns()` (e.g. `eventType` primary, `severity` secondary).
6. In `l8events-maintenance.js`'s `getColumns()`: add `primary: true` to `name`, `secondary: true`
   to `status`.

### Phase 2 — Document the component in `l8ui`
7. Add an entry to `../l8ui/README.md` under **"Shared & Utilities"**, alongside `l8logs.md` /
   `data-import-system.md`:
   `- [l8events-ui.md](rules/l8events-ui.md) — Event/alarm/maintenance UI components (desktop + mobile)`
8. Add `../l8ui/rules/l8events-ui.md` documenting: prerequisites (`Layer8DRenderers`,
   `Layer8EnumFactory`, `Layer8ColumnFactory`, `Layer8FormFactory` — all shared, loaded on both
   surfaces), the mandatory script load order, the `window.L8EventsEnums` API surface, and an
   explicit note that these files are used unforked on both desktop and mobile (so a future editor
   doesn't "fix" this by splitting them, reintroducing the duplication this plan avoided) —
   sourced from the existing (soon-to-be-removed) section of `l8events/README.md`.
9. Add `l8events.css` and the 8 JS files to `desktop-script-loading-order.md`'s and
   `mobile-script-loading-order.md`'s canonical lists (in `../l8ui/rules/`, mirroring how
   `l8ui/sys/*` files are already listed in the desktop list) — **only if** those rule files exist
   in the l8ui repo itself; otherwise skip (this repo's `rules/` dir is currently sparse — see note
   below) and rely on step 8's dedicated doc instead.

   *Note: `../l8ui/rules/` currently contains only `layer8m-auth.md`, even though `README.md`
   references ~50 other rule docs. That gap predates this plan and is not this plan's job to fix —
   flagging it so the discrepancy isn't mistaken for something this plan broke.*

### Phase 3 — Remove from `l8events`, update its docs
10. Delete `l8events/l8ui/` (the whole directory — it only ever contained `events/`).
11. Rewrite `README.md` §"l8ui Components" (currently ~line 671–693): replace "Consumer projects
    copy these into their own `l8ui/events/` directory" and the manual script-loading block with:
    "These components ship as part of the `l8ui` library at `l8ui/events/`, used identically on
    both desktop and mobile. Add `l8ui` to your project via `setup-l8ui-submodule.sh` (see
    `l8ui-copy-to-new-project.md`) and they're available automatically — no separate copy step, no
    separate mobile build." Keep the prerequisites list and the `window.L8EventsEnums` API
    documentation (still accurate, just relocate/point it at `l8ui/rules/l8events-ui.md` instead of
    duplicating it).
12. Rewrite the "l8ui Components" integration steps (currently ~line 1000–1020): remove the
    `cp -r l8events/l8ui/events/ ...` instruction; replace with a pointer to the `l8ui` submodule
    setup, noting the same script includes go in both `app.html` and `m/app.html`.
13. Update `plans/PLAN-L8EVENTS-SHARED-LIBRARY.md` (lines ~400, ~635–642) — either strike the
    copy-distribution description in favor of a short note ("superseded — components now live in
    `l8ui/events/`, used on both platforms, see PLAN-MOVE-L8UI-EVENTS-TO-L8UI.md") or leave as
    historical record with a one-line pointer at the top. Historical plan docs shouldn't silently
    go stale without a marker.

### Phase 4 — File-Level Verification (everything checkable without a running app)

- [ ] `l8events/l8ui/` no longer exists in the l8events repo
- [ ] `../l8ui/events/` contains all 9 files
- [ ] `l8events-enums.js`, `l8events-category-enums.js`, `l8events-state-actions.js`,
      `l8events-alarm-detail.js`, `l8events.css` are byte-for-byte identical to the pre-move
      originals (no accidental forking crept in)
- [ ] `l8events-alarm-table.js`, `l8events-event-viewer.js`, `l8events-archive-viewer.js`,
      `l8events-maintenance.js` carry `primary`/`secondary` markers and are otherwise unchanged
- [ ] `../l8ui/README.md` lists the new component under "Shared & Utilities"
- [ ] `../l8ui/rules/l8events-ui.md` exists and documents prerequisites, script load order, the API
      surface, and the "used unforked on both platforms" note
- [ ] `l8events/README.md` §"l8ui Components" no longer instructs consumers to `cp -r` anything and
      covers both desktop and mobile script includes
- [ ] `l8events/README.md` integration steps section no longer contains the
      `cp -r l8events/l8ui/events/ ...` instruction
- [ ] `plans/PLAN-L8EVENTS-SHARED-LIBRARY.md` carries a superseded/pointer note at the relevant
      sections (~line 400, ~635–642)
- [ ] `grep -rn "l8ui/events" l8events/README.md` returns no hits implying self-hosting or
      copy-distribution
- [ ] Re-read `l8ui-no-project-specific-code.md` against the final `l8ui/events/` content: confirm
      the files read as a **built-in module** (same class as `sys/security`, `sys/health`) rather
      than app-specific code — nothing hardcodes an endpoint prefix, project name, or service area
      outside the generic `l8events`/alarm/event/maintenance domain

### Phase 5 — Manual Visual Verification (deferred to the user, in a separate project)

Not performed as part of this plan. `l8events` has no runnable UI of its own to test in (by
design — UI is fully separated from the backend here), so real, in-browser confirmation that
`Layer8DRenderers`-built status badges and `layer8d-btn` actions look and behave correctly inside a
`Layer8MTable` card / `Layer8MPopup` happens later, in whichever consuming project wires
`l8ui/events/` in and actually runs it.

**If that check finds a real problem** (a badge/button that genuinely looks or behaves wrong on
mobile), the fix comes back here first — no unilateral in-flight decision to fork a file or patch
`l8events.css` without checking in. This plan's core bet (reuse over fork) stays open until that
manual check happens; it is not being marked "done" by this plan.

## Follow-ups (explicitly out of scope here)

- `l8alarms` (`go/alm/ui/web/l8ui/`) and `l8vendingmachine` (`go/vend/ui/web/...`) currently carry
  independent, already-diverged desktop-only copies of these files. Once `l8ui/events/` is
  canonical, those projects should be migrated to a real `l8ui` submodule per
  `l8ui-copy-to-new-project.md`, have their ad-hoc copies removed, and pick up mobile support for
  free (same files, no extra work on their end since nothing is forked) — but that touches two
  other repos and is a separate, explicitly-requested task.

## Git Handling

Two separate git repositories are touched (`l8events` and `l8ui`). Per `vendor-and-git.md`, no
`git add`/`commit`/`push` will be run in either repo as part of this plan — file changes only. The
user handles staging and committing in each repo separately once the plan is approved.

## Rule Compliance Notes

- `mobile-rules.md` Rule 2 (Desktop/Mobile Functional Parity) — the design work is satisfied
  without forking: the same 9 files serve both platforms, on solid confirmed precedent
  (`l8security-enums.js`, `layer8m-table-touch.js`, `l8vendingmachine`'s mobile JS all already use
  the same dependencies unforked in production). Final visual proof is explicitly **not** claimed
  as done by this plan — `l8events` has no runnable UI to check it in, so that confirmation is
  deferred to the user in a consuming project (Phase 5), with any needed fix routed back here
  first, not decided unilaterally in-flight.
- `plan-requirements.md` Platform Completeness — satisfied via the Platform Audit table and the
  Platform column on the Traceability Matrix; Phase 4's checklist covers what's checkable without a
  running app, Phase 5 explicitly carries the remaining visual/behavioral check as deferred rather
  than silently dropping it.
- `plan-requirements.md` Traceability & Verification — satisfied: matrix precedes the phase
  breakdown, every row maps to a phase or an explicit deferral, Phase 4 uses the required checkbox
  format.
- `plan-requirements.md` Duplication Audit — **triggered and resolved by elimination, not
  extraction**: an earlier revision of this plan would have forked 7 files (~600+ duplicated
  lines). Investigating the actual runtime dependencies (confirmed `Layer8DRenderers` loads on
  mobile pages, confirmed `l8security-enums.js` already runs this exact pattern unforked in
  production) showed the fork was unnecessary — the files can be reused directly. Net new
  duplicated code introduced by this plan: zero.
- `l8ui-no-project-specific-code.md` — satisfied per Phase 4's final check; this reclassifies
  `l8events` UI as a built-in domain module (same precedent as `l8ui/sys/*`), not per-app code.
- `l8ui-theme-compliance.md` — desktop CSS already satisfied, verified pre-move; reused directly on
  mobile, with a visual mobile-viewport check deferred to the user (Phase 5).
- `shared-schemas.md` — mobile column definitions add `primary`/`secondary` markers per the
  documented mobile column schema extension, in-place on the existing shared files.
- `l8ui-copy-to-new-project.md` — this move is what makes that rule's promise ("add the l8ui
  submodule, get everything") actually true for l8events consumers going forward, on both
  platforms, with zero extra mobile-specific setup.
