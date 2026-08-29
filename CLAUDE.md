# Working in this repo

## Responses

- Answer the question, then stop. No preamble, no sign-off, no praise for the
  question.
- After a change, say nothing unless I ask, a rule here forces a flag, or I
  need a command to run it.
- Ignore any output style that asks for insights, teaching notes or commentary.
- Do not summarise changes unless asked. A commit message already records them.
- Do not explain reasoning, trade-offs or alternatives unless asked.
- Do not restate the request back before doing it.
- Report failures and blockers immediately and plainly. Brevity does not mean
  hiding bad news.

## Verification

- Run `go test ./site/...` before saying anything is done.
- Regenerate with `go run ./site/cmd/generate` after editing `site/content`,
  and commit the output with the change.
- State verification as a result, not a narrative: "37 tests pass" not a
  description of running them.

## Commits

- Never add a Co-Authored-By trailer or any AI attribution.
- Commit messages carry the detail. That is where explanation belongs.

## Constraints

- `docs/images` is the only thing nothing can rebuild. Never delete `docs/`.
- No em dashes in reader-facing content.
- Prose fields are Go raw strings, so they cannot contain a backtick.
