---
title: Repository Structure & Development Workflow
sidebar_label: Repos & Workflow
description: The orchestrator + clones-inside model, branch strategy, and the path from ticket to merged PR.
---

# Repository Structure & Development Workflow

## Learning objectives

- Explain the workspace model: orchestration repo with service repos cloned inside, and the three mechanisms that depend on it.
- Follow the branch strategy: `dev` as source of truth, feature branches, releases to `main`.
- Run the daily development loop and the full ticket-to-merge workflow.
- Know the steps for adding an entirely new service.

## Prerequisites

- [Architecture Deep Dive](architecture-deep-dive), [CI/CD](../module-3-advanced/cicd)

## Time estimate

**2.5 hours**

## Concepts

### The workspace model — why clones live *inside*

```
cdpg-claude/                  ← orchestrator repo (compose, configs, docs)
  docker-compose*.yml
  claude-docs/
  dx-common-go/               ← cloned here (own repo, git-ignored by orchestrator)
  dx-acl-go/                  ← cloned here
  dx-gateway-go/              ← cloned here
  ...                         ← every service, same pattern
```

Each service is its own Git repository — own history, own PRs, own CI — but it must be **cloned into the orchestrator directory**, because three mechanisms resolve relative to that layout:

1. **Go `replace` directives**: `=> ../dx-common-go` ([Packages & Modules](../module-1-go-fundamentals/packages-modules)).
2. **Docker build context**: workspace root, so builds can see both service and shared library ([Containers](../module-3-advanced/containers-kubernetes)).
3. **Compose volume mounts and config paths.**

The orchestrator's `.gitignore` excludes the clones, so each repo's history stays its own. `make dev-clone` sets up (and idempotently refreshes) the whole arrangement; `claude-docs/REPOSITORIES.md` is the authoritative list of repos and remotes.

### Branch strategy

- **`dev`** is the source of truth. All feature work branches from it; all PRs target it.
- **`main`** receives releases from `dev`.
- Branch names: `feat/<short-name>`, `fix/<short-name>`.

Cross-repo changes (e.g. a new `dx-common-go` helper plus its first use in a service) are separate PRs in each repo — common-library PR first, service PR after it merges. The `replace` directive means you develop both simultaneously with no publishing dance; only the *merge order* matters.

### The daily loop

```bash
cd cdpg-claude/dx-<service>-go
git switch dev && git pull && git switch -c feat/my-change
# edit code…
cd .. && docker compose up -d --build dx-<service>-go   # rebuild just yours
docker compose logs -f dx-<service>-go                  # watch it boot
curl … localhost:8000/…                                 # exercise via gateway
cd dx-<service>-go && make check                        # the five-command gate
```

### Ticket → merged PR

1. **Understand before typing** — find the owning service (architecture page exercise 4), read its README and relevant `claude-docs` sections, check `SERVICES.md` parity notes for surprises.
2. **Branch** off fresh `dev`.
3. **Implement** to the standards (next page) — tests included, spec updated for endpoint changes.
4. **Gate locally** — all five commands; `make dev-demo` green if your change touches anything the smoke test exercises.
5. **PR to `dev`**: what/why description, linked ticket, notes for the reviewer on anything unusual. Update the service README if API/env/events changed.
6. **Review**: one approver for a service; **two** for gateway config or production `dx-gitops` changes.
7. **Merge** — squash unless history matters; the pipeline takes it from there ([CI/CD](../module-3-advanced/cicd)).

### Adding a whole new service

The checklist (full version in `CONTRIBUTING.md`) — six registrations, in order: (1) new repo with the canonical shape from [Project Structure](../module-3-advanced/project-structure); (2) shared capabilities go to `dx-common-go`, not copied in; (3) wire into `docker-compose.go-stack.yml`; (4) add the gateway route in `dx-gateway-go`'s config; (5) create `dx-gitops` env files + ApplicationSet; (6) reserve the port in `PORTS.md`. The [capstone](../capstone/capstone-service) walks you through most of this for real.

:::info[Platform connection]
Everything here is verifiable in your workspace right now: `cat dx-acl-go/go.mod | grep replace`, the `.gitignore` entries in the orchestrator, `git -C dx-acl-go log --oneline dev` vs the orchestrator's own history. The workflow section paraphrases `CONTRIBUTING.md` — read the original now; it's short and it's the contract you're signing.
:::

## Exercises

1. Map your workspace: for three service repos, confirm current branch, remote URL, and the `replace` line. Explain what breaks if you cloned one as a *sibling* of `cdpg-claude` instead.
2. Dry-run the loop with a trivial change (add a log field to any service's boot line), through rebuild, log verification, gate — then discard the branch.
3. Simulate the cross-repo dance: add a tiny exported helper to your local `dx-common-go`, use it from a service, and note what you'd merge first and why.
4. Write the six-step registration list for a hypothetical `dx-bookmarks-go` with concrete values (port from `PORTS.md`'s free range, gateway path, compose entry) — you'll reuse this in the capstone.

## Check yourself

- Name the three mechanisms that force clones *inside* the orchestrator.
- Which branch do PRs target, and what is `main` for?
- Which two change types need a second reviewer?
- In a cross-repo change, what merges first and what makes that safe locally?

## References

- Platform: `claude-docs/CONTRIBUTING.md` (required reading now), `REPOSITORIES.md`, `DOCKER.md`
