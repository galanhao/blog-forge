# AGENTS.md

Project harness for reliable agent-assisted development of a Go static site generator.

## Startup Workflow

Before writing code:

1. **Confirm working directory** with `pwd` — must be `/root/blog-forge`
2. **Read this file** completely
3. **Read project docs**: `ARCHITECTURE.md`, `README.md`
4. **Run `./init.sh`** to verify environment is healthy (Go 1.24 via `g`, tests pass, build succeeds)
5. **Read `feature_list.json`** to see current feature state
6. **Review recent commits** with `git log --oneline -5`

If baseline verification is failing, repair that first before adding new scope.

## Working Rules

- **One feature at a time**: Pick exactly one `not-started` or `active` feature from `feature_list.json`
- **Verification required**: Don't claim done without running `./init.sh` AND the feature's specific verification command
- **Update artifacts**: Before ending session, update `progress.md` and `feature_list.json`
- **Stay in scope**: Don't modify files unrelated to the current feature
- **Leave clean state**: Next session must be able to run `./init.sh` immediately
- **No refactoring before verified**: Core functionality must be verified working before any style/perf changes

## Tech Stack

- **Language**: Go 1.24 (managed by `g`, source `/root/.g/env` before use)
- **Module**: `github.com/galanhao/blog-forge`
- **ORM**: None (stdlib `database/sql` for future use; current build is pure static)
- **Markdown**: goldmark + extensions (GFM, chroma syntax highlighting)
- **Templates**: `html/template` with composition pattern (`page-start` / `page-end` partials)
- **Proxy**: `GOPROXY=https://goproxy.cn,direct GONOSUMCHECK=*`

## Architecture Boundaries

```
cmd/blog-forge/       → CLI entry point only; no business logic
internal/config/   → Site config loading; depends on nothing internal
internal/content/  → Post model + loader; depends on nothing internal
internal/render/   → Markdown → HTML rendering; depends on nothing internal
internal/permalink/→ URL computation; depends on nothing internal
internal/theme/    → Theme loading + template engine + context structs; depends on config
internal/site/     → Build pipeline orchestrator; depends on ALL internal packages
internal/serve/    → Local dev server; depends on nothing internal (currently dead code)
```

- `internal/site/` is the only package that imports other internal packages
- `internal/content/`, `internal/render/`, `internal/permalink/` must NEVER import from each other or from `theme`/`site`
- Template context structs live in `internal/theme/context.go`
- Build pipeline steps live in `internal/site/builder.go`

## Naming Conventions

- **post** (not article/entry) — `Post`, `PostCtx`, `LoadPosts`, `writePosts`
- **site** (not website/blog) — `SiteCtx`, `SiteConfig`, `buildSiteCtx`
- **page** — ambiguous: "standalone page" (`content/pages/`) vs "paginated page number" (`PaginationCtx`). Context disambiguates.
- **slug** — URL-friendly filename component
- **permalink** — full computed URL path (e.g., `/2026/05/24/hello-world/`)

## Hard Constraints

1. Never commit `checkin.db`, `.env`, or any file containing secrets/tokens
2. Never use `math/rand` for anything security-relevant
3. All exported functions must have doc comments
4. Error wrapping: always use `fmt.Errorf("context: %w", err)`
5. No `panic` in library code; only in `main` for truly unrecoverable states

## Required Artifacts

- `feature_list.json` — Feature state tracker (source of truth)
- `progress.md` — Session continuity log
- `init.sh` — Standard startup and verification path
- `session-handoff.md` — Optional, for larger sessions

## Definition of Done

A feature is done only when ALL of the following are true:

- [ ] Target behavior is implemented
- [ ] `go test ./...` passes
- [ ] `go build ./...` succeeds
- [ ] `./blog-forge build` produces valid output in `dist/`
- [ ] Evidence recorded in `feature_list.json`
- [ ] Repository remains restartable from standard startup path

## End of Session

Before ending a session:

1. Update `progress.md` with current state
2. Update `feature_list.json` with new feature status
3. Record any unresolved risks or blockers
4. Commit with descriptive message once work is in safe state
5. Leave repo clean enough for next session to run `./init.sh` immediately

## Verification Commands

```bash
# Full verification (recommended)
./init.sh
```

Required checks:
- `source /root/.g/env && go test ./...`
- `go build -o blog-forge ./cmd/blog-forge/`
- `./blog-forge build`

## Escalation

If you encounter:
- **Architecture decisions**: Consult `ARCHITECTURE.md`, otherwise ask user
- **Unclear requirements**: Check `README.md` and `feature_list.json`, otherwise ask user
- **Repeated test failures**: Update `progress.md`, flag for human review
- **Scope ambiguity**: Re-read `feature_list.json` for definition of done
