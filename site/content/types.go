package content

type ConceptItem struct {
	Term        string
	Description string
}

// BuildIt is a chapter's hands-on exercise. The reader hands Prompt to
// their own AI coding assistant rather than typing Go by hand, so Prompt
// must stand alone: it cannot refer to anything visible only on this page.
//
// Technique names the prompt-engineering approach this chapter
// demonstrates, and Why explains when to reach for it. Teaching a range of
// prompting techniques is a goal of the book in its own right, so no two
// chapters use the same one.
//
// Source cites where the technique is documented, and must name one of the
// corpora in RecognisedSources. A technique the reader can't look up is a
// technique we invented, and inventing one teaches a vocabulary nobody else
// speaks — so the citation is a required field, not a nicety.
type BuildIt struct {
	Intro     string
	Technique string
	Why       string
	Source    string
	Prompt    string
}

// RecognisedSources are the corpora a BuildIt.Source may cite. Both are
// public and checkable:
//
//   - "The Prompt Report" — Schulhoff et al., arXiv:2406.06608, a systematic
//     survey taxonomising 58 text-based prompting techniques into six
//     families (In-Context Learning, Zero-Shot, Thought Generation,
//     Decomposition, Ensembling, Self-Criticism).
//   - "Anthropic" — the vendor's own prompt-engineering documentation at
//     platform.claude.com/docs/en/build-with-claude/prompt-engineering.
var RecognisedSources = []string{
	"The Prompt Report",
	"Anthropic",
}

type BreakIt struct {
	Intro     string
	Exercises []string
}

type ChapterContent struct {
	Number       int
	Slug         string
	Title        string
	Subtitle     string
	Category     string
	Difficulty   string
	EstTime      string
	QuickTip     string
	HeroImage    string
	HeroCaption  string
	Intuition    []string
	Concepts     []ConceptItem
	UnderTheHood []string
	BuildIt      BuildIt
	BreakIt      BreakIt
}
