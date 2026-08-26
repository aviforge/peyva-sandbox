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
// speaks, so the citation is a required field, not a nicety.
type BuildIt struct {
	Intro     string
	Technique string
	Why       string
	Source    string

	// Prompts are the chapter's turns, in the order they are asked. A reader
	// copies one, reads the answer, and copies the next.
	//
	// Most chapters have more than one because most techniques here have a
	// shape: reason before building, or build then take the work apart. As a
	// single block an assistant answers all of it in one pass, and the
	// reasoning it was meant to do first comes back as a paragraph justifying
	// code it had already written.
	//
	// Technique shapes every turn in a chapter, the portal's included. The
	// chapters that grow the portal are the ones whose technique suits a page
	// as well as a component: chapter 11's self-critique loop is what turns a
	// working portal into a presentable one, and chapter 12's permission to
	// refuse is what stops a plain page quietly acquiring a framework.
	Prompts []Prompt
}

// Prompt is one turn of a chapter's Build It.
//
// Every one states the state it starts from, because a prompt is copied out of
// the page on its own. Nothing tells the assistant to go and read what earlier
// chapters built: on a large codebase that is thousands of tokens per chapter
// to rediscover what one sentence can say, and on a small context window it
// does not fit at all.
type Prompt struct {
	// Label names the turn: Build, Design, Review. Shown above the block.
	Label string
	// Intro is the one line of why, above the block.
	Intro string
	Text  string
	// Portal marks a turn that builds the page rather than a component, so it
	// carries the portal's rules instead of the component ones.
	Portal bool

	// Thinking marks a turn that produces an answer rather than a change: an
	// analogy, a list of questions, three strategies scored against each
	// other. It carries only the pointer to the spec, because the language to
	// build in and how to round money are rules about work it is not doing.
	Thinking bool
}

// First returns the chapter's opening prompt, or the zero Prompt if it somehow
// has none. Callers that only want to know what a chapter asks for use this
// rather than indexing, which panics on the empty case.
func (b BuildIt) First() Prompt {
	if len(b.Prompts) == 0 {
		return Prompt{}
	}
	return b.Prompts[0]
}

// PortalPrompts and ComponentPrompts split the turns by what they build.
func (b BuildIt) PortalPrompts() []Prompt    { return b.filter(true) }
func (b BuildIt) ComponentPrompts() []Prompt { return b.filter(false) }

func (b BuildIt) filter(portal bool) []Prompt {
	var out []Prompt
	for _, p := range b.Prompts {
		if p.Portal == portal {
			out = append(out, p)
		}
	}
	return out
}

// RecognisedSources are the corpora a BuildIt.Source may cite. Both are
// public and checkable:
//
//   - "The Prompt Report": Schulhoff et al., arXiv:2406.06608, a systematic
//     survey taxonomising 58 text-based prompting techniques into six
//     families (In-Context Learning, Zero-Shot, Thought Generation,
//     Decomposition, Ensembling, Self-Criticism).
//   - "Anthropic": the vendor's own prompt-engineering documentation at
//     platform.claude.com/docs/en/build-with-claude/prompt-engineering.
var RecognisedSources = []string{
	"The Prompt Report",
	"Anthropic",
}

type ChapterContent struct {
	Number      int
	Slug        string
	Title       string
	Subtitle    string
	Category    string
	Difficulty  string
	EstTime     string
	QuickTip    string
	HeroImage   string
	HeroCaption string
	Concepts    []ConceptItem
	BuildIt     BuildIt
}
