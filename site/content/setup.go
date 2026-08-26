package content

// SetupFile is one file the reader saves before starting chapter 0.
//
// Both files are rendered into the page with a copy button rather than kept in
// this repository. A reader following the published site should never have to
// leave it to fetch something, and a file on GitHub is a second place the site
// cannot reach.
type SetupFile struct {
	Path    string
	Purpose string
	Content string

	// Names is set when the file has no single agreed filename, and lists the
	// name each assistant reads. The content saved under each is the same.
	Names []AgentFile
}

// AgentFile is the name one assistant reads its standing instructions from.
type AgentFile struct {
	Tool string
	Path string
}

// SetupFiles are the two files chapter 0 hands the reader.
var SetupFiles = []SetupFile{
	{
		Path:    "peyva/goal.md",
		Purpose: "What the system is for, and what must never happen to money. Every prompt from here on points at it instead of repeating it.",
		Content: GoalSpec,
	},
	{
		Path:    "Your assistant's instruction file",
		Purpose: "The rules that do not change between chapters: how to build, and how to answer. Saving it once is what keeps every prompt after this one short.",
		Content: AgentRules,
		Names:   AgentFiles,
	},
}

// AgentFiles is where each assistant looks for standing instructions. No single
// name is read by every tool, so the content is the same in all of them and
// only the filename changes. AGENTS.md is the closest thing to a convention,
// which is why it is also the answer for anything not named here.
var AgentFiles = []AgentFile{
	{Tool: "Claude Code", Path: "peyva/CLAUDE.md"},
	{Tool: "Codex", Path: "peyva/AGENTS.md"},
	{Tool: "GitHub Copilot", Path: "peyva/.github/copilot-instructions.md"},
	{Tool: "Anything else", Path: "peyva/AGENTS.md"},
}

// AgentRules is the standing instruction file for whichever assistant the
// reader uses. It holds what does not change between chapters, so a prompt can
// be a chapter's worth of work rather than a restatement of the rules.
//
// The prompts still carry the constraints in their preamble. That is
// deliberate: a prompt is copied on its own, and may reach an assistant that
// never read this file. A prompt that assumes it did is a prompt that quietly
// builds the wrong thing.
const AgentRules = `# Working on peyva

Read peyva/goal.md first. It holds the goal and the rules money must never
break. If a change would break one of those, the change is wrong.

## Building

- Standard library only. No frameworks, no brokers, no third-party libraries.
- One process on a laptop. No deployment.
- Code lives in peyva/<component>/, one folder per component.
- Money is exact: a decimal type, or integer minor units where the language has
  none. Never floating point. Two decimal places.
- The portal is plain HTML and CSS. No build step, no dependencies.
- The portal is one customer's own wallet, not an operator's view of everyone.
- Structure is earned. Build what this chapter needs, not what a later one
  might.

## Answering

- Do what the prompt asks, and nothing it did not ask for.
- If a constraint makes the task impossible, say so rather than working around
  it. That is a useful answer, not a failure.
- Make the smallest change that satisfies the prompt.
- Say what changed and why, in a line or two. No summary of work I just
  watched you do.
- When you are unsure, ask rather than guess.
`
