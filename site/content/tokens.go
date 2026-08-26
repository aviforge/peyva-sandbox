package content

import "strings"

// EstimateTokens is a rough token count for a prompt.
//
// Four characters per token is the usual heuristic for English, and it is
// close enough for what this is used for: comparing one prompt against
// another, and showing that a chapter's asking price is a few hundred tokens
// rather than a few thousand. It is not close enough to bill anyone with, and
// the site says so where it shows the number.
//
// Nothing here can know the real cost of a chapter. The prompt is the small
// end of it: the assistant then reads code, writes code, and often reasons at
// length before either. That number lives in the reader's own tool, which is
// what CostNote points them at.
func EstimateTokens(text string) int {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return 0
	}
	return (len(trimmed) + 3) / 4
}

// PromptTokens totals a chapter's turns.
func (b BuildIt) PromptTokens() int {
	total := 0
	for _, p := range b.Prompts {
		total += EstimateTokens(p.Text)
	}
	return total
}

// TotalPromptTokens is every prompt in the book, for the setup note.
func TotalPromptTokens() int {
	total := 0
	for _, c := range All {
		total += c.BuildIt.PromptTokens()
	}
	return total
}

// CostNote explains where the reader's tokens actually go, and how to see it.
// It is shown in setup, beside the files, because the answer changes what a
// reader does from chapter 1 rather than at the end.
const CostNote = `Every prompt in this book adds up to about %d tokens, across
all twenty-one chapters. That is not where the money goes.

The cost is the assistant reading your code and writing more of it, which is
tens of times larger and grows as peyva does. Nothing on this page can measure
it. Your tool can:

  Claude Code   /cost, or /context for what is loaded right now
  Codex         the usage line at the end of a turn
  Copilot       your account's usage page
  An API        the usage object on every response

Two things make it cheaper, and you control both.

Caching. Your assistant sends the same opening every turn: its own instructions,
then goal.md, then the rules file. Charging full price for that each time would
be absurd, so it does not: a prefix it has seen recently costs a fraction of the
usual rate. The catch is that the match has to be exact. Editing either file
halfway through a session throws the cache away, and the next turn pays full
price for all of it again. Save them once, here, then leave them alone. If you
do need to change one, do it between chapters rather than mid-chapter.

Context. No prompt here asks the assistant to go through the project first. Each
one says what it starts from, so it opens the two or three files that turn
touches instead of all of them. The difference grows with every chapter: by the
end, reading everything costs more than every prompt in this book put together.`
