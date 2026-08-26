# Peyva Sandbox

Distributed systems, taught by building PEYVA (a peer-to-peer wallet) one
honest layer at a time.

**[Read it here](https://aviforge.github.io/peyva-sandbox/)**, or open
`docs/index.html` from a clone. Nothing to install or build.

Each chapter ends with a prompt you hand to your own AI assistant rather than
code to copy. Chapter 0 gives you a spec to save as `peyva/goal.md`, holding
the goal, the rules money must never break, and the constraints. Every prompt
after that points at it. Chapter 0 asks which language you want to build in, and every
prompt after it asks for that one, so what you build in chapter 12 still fits
what you built in chapter 4.

Chapters that grow the portal carry two prompts: one for the component, one
for the page, each copied on its own. The portal is plain HTML and CSS with no
dependencies, so it is the same portal whichever backend language you pick.

The choice is offered once, in chapter 0. Later chapters show which language is
in force and link back there to change it. Twelve are available: Go, Python,
JavaScript, TypeScript, Java, C#, C++, Rust, Ruby, PHP, Kotlin and Swift.

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

### Adding a language

Add an entry to `site/content/languages.go` and regenerate. That is the whole
change, because no prompt names a language: they describe what the component
has to do and leave the idiom to the assistant. A test fails if a prompt starts
naming one, which is what keeps a single entry enough.

The language is written into every page at build time, not added by script. A
prompt that reaches an assistant without naming a language is worse than one
naming the wrong language: the assistant picks its own, can pick differently on
different chapters, and the reader finds out several chapters later when
nothing fits together.

## Layout

    docs/                      the published site, served by Pages as the root
    docs/index.html            entry point, so a bare / lands on chapter 0
    docs/images                the chapter illustrations
    site/content               chapter text, as Go values
    site/content/languages.go  the languages the picker offers
    site/content/goal.go       the spec chapter 0 hands the reader
    site/templates             page and sidebar templates
    site/assets                CSS and JS
    site/cmd                   the generator

Everything under `docs/` is rewritten by the generator except `docs/images`.
Those are the only files in the repository that nothing can rebuild, and they
live there because GitHub Pages publishes `docs/` and nothing outside it.

The generator reads that folder and never writes to it. If it is missing the
build stops rather than publishing a site of empty panels, so deleting `docs/`
to force a clean build fails loudly instead of quietly.
