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
// demonstrates, What says in one plain sentence what that approach is, and
// Why explains why it suits this chapter. Teaching a range of prompting
// techniques is a goal of the book in its own right, so no two chapters use
// the same one.
//
// What exists because the name is jargon. Generated Knowledge Prompting means
// nothing to a reader meeting it for the first time, and the sentence that
// followed used to assume they had already looked it up. What is the
// definition a reader can act on without leaving the page.
//
// Source cites where the technique is documented, and must name one of the
// corpora in RecognisedSources. A technique the reader can't look up is a
// technique we invented, and inventing one teaches a vocabulary nobody else
// speaks, so the citation is a required field, not a nicety.
//
// SourceURL is where the citation points. A citation the reader has to
// search for is a citation most readers will not follow, so the Source text
// is rendered as a link to it. It must sit on the cited corpus's own site.
type BuildIt struct {
	Technique string
	What      string
	Why       string
	Source    string
	SourceURL string

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
	// Label names what kind of turn it is, from a fixed set of five so the
	// same word means the same thing on every page: Think produces an answer
	// and no code, Build makes a component, Check takes apart what was just
	// built, Portal builds the Portal, Try is the reader's own turn at the
	// keyboard. Shown above the block.
	Label string
	Text  string
	// Portal marks a turn that builds the page rather than a component, so it
	// carries the portal's rules instead of the component ones.
	Portal bool

	// Thinking marks a turn that produces an answer rather than a change: an
	// analogy, a list of questions, three strategies scored against each
	// other. It carries only the pointer to the spec, because the language to
	// build in and how to round money are rules about work it is not doing.
	Thinking bool

	// Reader marks a turn the reader performs rather than the assistant: a
	// command to run, a process to stop, a page to reload, and one thing to
	// watch for. It is never sent anywhere, so it carries no rules, no token
	// count and no {os}. Labelled Try.
	//
	// It exists because the assistant was running every proof and reporting
	// it, and a reader who never sees the two reads or the refused second
	// copy has watched a chapter rather than done one.
	Reader bool

	// Commands are what the reader runs, one per operating system, verbatim.
	// The reader authors nothing: anything typed is here with a copy button,
	// and a step that needs a pause in the middle pauses itself. Empty when
	// the turn is done in the browser or with Ctrl+C alone.
	Commands []SystemCommand
}

// SystemCommand is one Try command for one operating system. SystemID
// matches System.ID.
type SystemCommand struct {
	SystemID string
	Command  string
}

// Commands builds a Try turn's list from a PowerShell command, a batch
// command and a shell command, with macOS and Linux sharing the shell one.
// Most turns are curl and a pause, which the two Unixes spell the same way.
func Commands(ps, bat, sh string) []SystemCommand {
	return CommandsSplit(ps, bat, sh, sh)
}

// CommandsSplit is Commands for the turns where macOS and Linux differ,
// which is whenever a process has to be found by its port: macOS ships lsof
// and Linux ships fuser, and neither reliably has the other.
func CommandsSplit(ps, bat, mac, linux string) []SystemCommand {
	return []SystemCommand{
		{SystemID: "windows", Command: ps},
		{SystemID: "windows-bat", Command: bat},
		{SystemID: "macos", Command: mac},
		{SystemID: "linux", Command: linux},
	}
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

// Where the two corpora live, for SourceURL. The Prompt Report is one paper,
// so every citation of it links to the same page. Anthropic's prompting
// advice was once a page per topic and is now one page with a section per
// topic, so a citation of it appends the section's anchor.
const (
	PromptReportURL          = "https://arxiv.org/abs/2406.06608"
	AnthropicDocsURL         = "https://platform.claude.com/docs/"
	AnthropicBestPracticeURL = AnthropicDocsURL + "en/build-with-claude/prompt-engineering/claude-prompting-best-practices"
)

type ChapterContent struct {
	Number      int
	Slug        string
	Title       string
	Subtitle    string
	Category    string
	Difficulty  string
	QuickTip    string
	HeroImage   string
	HeroCaption string

	// Why is the chapter's teaching, stated before the vocabulary and the
	// prompts. Each entry is one claim about how systems behave, short enough
	// to be checked against a real system and wrong in an obvious way if it
	// is wrong. The assistant the reader prompts will say a great deal; this
	// is the part the book stands behind.
	//
	// Bullets rather than paragraphs. A reader who has just been told the
	// fact does not need it three more times in different clothes.
	Why []string

	Concepts []ConceptItem
	BuildIt  BuildIt

	// Aside is a sidebar: a short illustrated topic that belongs beside this
	// chapter rather than in a chapter of its own. It has its own image, its
	// own claims and its own prompts, and its own technique, because it is a
	// lesson in its own right that happened not to earn a number.
	Aside *Aside
}

// Aside is the sidebar a chapter may carry. See ChapterContent.Aside.
type Aside struct {
	Title       string
	HeroImage   string
	HeroCaption string
	Why         []string
	BuildIt     BuildIt
}
