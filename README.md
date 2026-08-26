# Peyva Sandbox

Distributed systems, taught by building PEYVA — a peer-to-peer wallet — one
honest layer at a time.

**[Read it here](https://aviforge.github.io/peyva-sandbox/)**, or open
`index.html` from a clone. Nothing to install or build.

## Contents

<!-- toc:start -->
- **0.** [What Are We Building?](https://aviforge.github.io/peyva-sandbox/docs/chapter-0.html)
- **1.** [Inside One Computer](https://aviforge.github.io/peyva-sandbox/docs/chapter-1.html)
- **2.** [Finding Peyva (Processes & Ports)](https://aviforge.github.io/peyva-sandbox/docs/chapter-2.html)
- **3.** [Across the Wire (Networking)](https://aviforge.github.io/peyva-sandbox/docs/chapter-3.html)
- **4.** [Designing the API](https://aviforge.github.io/peyva-sandbox/docs/chapter-4.html)
- **5.** [Storing Money (Databases)](https://aviforge.github.io/peyva-sandbox/docs/chapter-5.html)
- **6.** [Finding Things Fast (Indexes)](https://aviforge.github.io/peyva-sandbox/docs/chapter-6.html)
- **7.** [Making It Safe (Transactions)](https://aviforge.github.io/peyva-sandbox/docs/chapter-7.html)
- **8.** [Exactly Once (Idempotency)](https://aviforge.github.io/peyva-sandbox/docs/chapter-8.html)
- **9.** [How Big Is PEYVA? (Capacity Estimation)](https://aviforge.github.io/peyva-sandbox/docs/chapter-9.html)
- **10.** [Growing the Team: Scale Out](https://aviforge.github.io/peyva-sandbox/docs/chapter-10.html)
- **11.** [Sharing Work: Caching](https://aviforge.github.io/peyva-sandbox/docs/chapter-11.html)
- **12.** [Decoupling with Messages: Queues](https://aviforge.github.io/peyva-sandbox/docs/chapter-12.html)
- **13.** [Reliability Patterns: Transactional Outbox](https://aviforge.github.io/peyva-sandbox/docs/chapter-13.html)
- **14.** [Big Changes Safely: Sagas](https://aviforge.github.io/peyva-sandbox/docs/chapter-14.html)
- **15.** [Data Copies: Replication](https://aviforge.github.io/peyva-sandbox/docs/chapter-15.html)
- **16.** [When Things Fail: CAP / Consistency](https://aviforge.github.io/peyva-sandbox/docs/chapter-16.html)
- **17.** [See Everything: Observability](https://aviforge.github.io/peyva-sandbox/docs/chapter-17.html)
- **18.** [Lock It Down: Security](https://aviforge.github.io/peyva-sandbox/docs/chapter-18.html)
- **19.** [Operating in Production](https://aviforge.github.io/peyva-sandbox/docs/chapter-19.html)
- **20.** [Putting It All Together](https://aviforge.github.io/peyva-sandbox/docs/chapter-20.html)
<!-- toc:end -->

## Contributing

The published site is committed, so reading it needs nothing installed.
Rebuilding it needs [Go](https://go.dev/dl/) 1.21+.

Chapters live in `site/content` as Go values, not markdown. Edit one, then
regenerate from the repo root:

    go run ./site/cmd/generate

That rewrites `index.html` and `docs/`, and refreshes the Contents list above
from the chapter registry, so titles cannot drift out of sync. Commit the
regenerated output alongside the content change. The generator resolves its
templates and assets by relative path, so it only works from the repo root.

    go test ./site/...

`go test ./...` fails from the root: `go.work`, no root module.

## Layout

    index.html      entry point — the published site's front door
    docs/           the published pages, images and styles
    site/content    chapter text, as Go values
    site/templates  page and sidebar templates
    site/assets     source images, CSS and JS
    site/cmd        the generator
