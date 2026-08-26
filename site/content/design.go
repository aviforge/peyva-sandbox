package content

// DesignBrief is what the reader saves as peyva/portal/design.md.
//
// It is a third setup file rather than text inside a prompt because thirteen
// chapters touch the Portal. Repeating three hundred words of design language
// in each of them would cost more than every other prompt in the book put
// together, and the thirteen copies would drift. The portal preamble points at
// this file instead, which costs one line.
//
// The reference section names real wallets on purpose. "Make it polished" is
// not an instruction an assistant can act on; "the balance is the largest thing
// on the screen and money is set in tabular figures" is. The section after it
// exists to stop the result being a tribute act: the rules produce competence,
// and competence without a point of view is the house style of every framework
// starter template.
const DesignBrief = `# The look

The Portal is a wallet someone opens several times a day to check one number.
Calm, quick, exact.

## What good wallets already agree on

Monzo, Revolut, Wise, Cash App and Apple Wallet converge on the same handful of
moves. Take these. Do not reproduce any of them.

- The balance is the largest thing on the screen. Nothing competes with it.
- Money is set in tabular figures, so digits do not shift as values change,
  aligned on the decimal point, always two places.
- A transaction row reads across: who, what it was for, then the amount hard
  right. The date is the quietest thing in the row.
- Money in and money out differ by sign and weight, not by colour alone.
- History is grouped under date headings, most recent first.
- Anything not yet settled looks pending. It is never silently missing.
- One accent colour, spent on the primary action and almost nothing else.

## Then make it yours

Choose one visual idea and carry it through every screen: a distinctive way of
setting numerals, a single strong accent against near neutral, a card shape, a
motif that repeats. Write the idea in a comment at the top of the stylesheet so
it survives the next chapter.

## Rules

- A type scale of at most five sizes. One family for prose, one for figures.
- One corner radius, one border colour, one shadow. Reuse them everywhere.
- Spacing from one scale, used generously. Crowding reads as unfinished.
- Every action says what happened, and to what.
- Every screen has an empty state that says what to do next, an error state that
  says what went wrong, and a loading state that does not move the layout.
- Legible in light and dark, following the reader's system setting.
- Anything tappable is at least 44 pixels.
- Nothing changes position after the page has loaded.
`
