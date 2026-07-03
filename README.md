# go-learning — DX Go Learning Path

The standard onboarding curriculum for Go developers joining the **Data Exchange (DX) platform** at CDPG. It takes an engineer from Go beginner to a productive platform contributor in **10–12 weeks part-time** (~8–10 hours/week), and doubles as a reference for engineers from other teams contributing to the Go service fleet.

Built with [Docusaurus](https://docusaurus.io/), styled to match [cdpg-docs](https://github.com/datakaveri/cdpg-docs).

## What's inside

| Module | Weeks | Covers |
|---|---|---|
| 0 — Setup & Orientation | 1 | Toolchain, running the stack locally, platform mental map |
| 1 — Go Fundamentals | 2–4 | Language: syntax → interfaces → errors → generics |
| 2 — Intermediate Go | 5–7 | Concurrency, context, DI, config, logging, testing, profiling |
| 3 — Advanced Go | 8–10 | HTTP/chi, REST, auth, pgx, transactions, RabbitMQ, workers, K8s, CI/CD |
| 4 — The DX Platform | 11–12 | Architecture, dx-common-go, service anatomy, standards, security, deployment |
| 5 — Capstone | +1–2 | Build a full DX-style service, land your first real PR |

Exercises are **standalone Go programs in Modules 0–2** and **hands-on tasks against the real local stack (`make dev-up`) in Modules 3–5**.

## Run locally

```bash
npm ci
npm start        # dev server at http://localhost:3000/go-learning/
npm run build    # production build (fails on broken links)
npm run typecheck
```

## Contributing a page

Every topic page follows the same template (see any existing page):

1. **Learning objectives** — bulleted, testable
2. **Prerequisites** — links to prior pages
3. **Time estimate** — hours, feeds the weekly schedule in the roadmap
4. **Concepts** — explanations with runnable Go snippets
5. **Platform connection** — admonition linking the topic to real DX code
6. **Exercises** — standalone early, platform-tied later
7. **Mini-project** — where meaningful
8. **Check yourself** — 3–5 self-assessment questions
9. **References** — official docs, style guides, platform docs

House rules for content:

- Platform claims must be **verified against the actual code/docs** (dx-common-go, GO-SERVICE-STANDARDS.md) — don't teach from memory.
- Teach what the platform actually practices (chi, pgx v5, zap, viper, amqp091-go; table-driven tests; no ORM, no testify).
- Go snippets should compile — check with `go vet` before committing.
