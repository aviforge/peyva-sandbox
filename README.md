# Peyva Sandbox

Distributed systems, taught by building PEYVA — a peer-to-peer wallet — one
honest layer at a time.

**[Read it here](https://aviforge.github.io/peyva-sandbox/)**, or open
`docs/index.html` from a clone. Nothing to install or build.

## Contents

<!-- toc:start -->
- **0.** [What Are We Building?](https://aviforge.github.io/peyva-sandbox/chapter-0.html)
- **1.** [Inside One Computer](https://aviforge.github.io/peyva-sandbox/chapter-1.html)
- **2.** [Finding Peyva (Processes & Ports)](https://aviforge.github.io/peyva-sandbox/chapter-2.html)
- **3.** [Across the Wire (Networking)](https://aviforge.github.io/peyva-sandbox/chapter-3.html)
- **4.** [Designing the API](https://aviforge.github.io/peyva-sandbox/chapter-4.html)
- **5.** [Storing Money (Databases)](https://aviforge.github.io/peyva-sandbox/chapter-5.html)
- **6.** [Finding Things Fast (Indexes)](https://aviforge.github.io/peyva-sandbox/chapter-6.html)
- **7.** [Making It Safe (Transactions)](https://aviforge.github.io/peyva-sandbox/chapter-7.html)
- **8.** [Exactly Once (Idempotency)](https://aviforge.github.io/peyva-sandbox/chapter-8.html)
- **9.** [How Big Is PEYVA? (Capacity Estimation)](https://aviforge.github.io/peyva-sandbox/chapter-9.html)
- **10.** [Growing the Team: Scale Out](https://aviforge.github.io/peyva-sandbox/chapter-10.html)
- **11.** [Sharing Work: Caching](https://aviforge.github.io/peyva-sandbox/chapter-11.html)
- **12.** [Decoupling with Messages: Queues](https://aviforge.github.io/peyva-sandbox/chapter-12.html)
- **13.** [Reliability Patterns: Transactional Outbox](https://aviforge.github.io/peyva-sandbox/chapter-13.html)
- **14.** [Big Changes Safely: Sagas](https://aviforge.github.io/peyva-sandbox/chapter-14.html)
- **15.** [Data Copies: Replication](https://aviforge.github.io/peyva-sandbox/chapter-15.html)
- **16.** [When Things Fail: CAP / Consistency](https://aviforge.github.io/peyva-sandbox/chapter-16.html)
- **17.** [See Everything: Observability](https://aviforge.github.io/peyva-sandbox/chapter-17.html)
- **18.** [Lock It Down: Security](https://aviforge.github.io/peyva-sandbox/chapter-18.html)
- **19.** [Operating in Production](https://aviforge.github.io/peyva-sandbox/chapter-19.html)
- **20.** [Putting It All Together](https://aviforge.github.io/peyva-sandbox/chapter-20.html)
<!-- toc:end -->

## Contributing

The published site is committed, so reading it needs nothing installed.
Rebuilding it needs [Go](https://go.dev/dl/) 1.21+.

Chapters live in `site/content` as Go values, not markdown. Edit one, then
regenerate from the repo root:

    go run ./site/cmd/generate

That rewrites `docs/` and refreshes the Contents list above from the chapter
registry, so titles cannot drift out of sync. Commit the regenerated output
alongside the content change. The generator resolves its templates and assets
by relative path, so it only works from the repo root.

    go test ./site/...

`go test ./...` fails from the root: `go.work`, no root module.

## Layout

    docs/            the published site. GitHub Pages serves this as the root
    docs/index.html  entry point, so a bare / lands on chapter 0
    docs/images      the chapter illustrations
    site/content     chapter text, as Go values
    site/templates   page and sidebar templates
    site/assets      CSS and JS
    site/cmd         the generator

Everything under `docs/` is rewritten by the generator except `docs/images`.
Those are the only files in the repository that nothing can rebuild, and they
live there because GitHub Pages publishes `docs/` and nothing outside it. Do
not delete the folder to force a clean build.
