# Peyva Sandbox

## Contents

<!-- toc:start -->
- **0.** What Are We Building?
- **1.** Inside One Computer
- **2.** Finding Peyva (Processes & Ports)
- **3.** Across the Wire (Networking)
- **4.** Designing the API
- **5.** Storing Money (Databases)
- **6.** Finding Things Fast (Indexes)
- **7.** Making It Safe (Transactions)
- **8.** Exactly Once (Idempotency)
- **9.** How Big Is PEYVA? (Capacity Estimation)
- **10.** Growing the Team: Scale Out
- **11.** Sharing Work: Caching
- **12.** Decoupling with Messages: Queues
- **13.** Reliability Patterns: Transactional Outbox
- **14.** Big Changes Safely: Sagas
- **15.** Data Copies: Replication
- **16.** When Things Fail: CAP / Consistency
- **17.** See Everything: Observability
- **18.** Lock It Down: Security
- **19.** Operating in Production
- **20.** Putting It All Together
<!-- toc:end -->

## Prerequisites

- [Go](https://go.dev/dl/) 1.21+

## Build and view the site

From the repo root, run:

    go run ./site/cmd/generate

The generator resolves its templates and assets by relative path, so it only
works from the repo root. Then open `site/dist/index.html` in your browser.

It also refreshes the Contents list above from `site/content`, so chapter
titles cannot drift out of sync.

## Run tests

    go test ./site/...

`go test ./...` fails from the root: `go.work`, no root module.
