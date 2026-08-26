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
// length before either. Only the reader's own tool can report that, so the page
// says the number is the prompt alone rather than implying it is the total.
func EstimateTokens(text string) int {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return 0
	}
	return (len(trimmed) + 3) / 4
}
