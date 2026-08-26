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
		Purpose: "What the system is for, and what must never happen to money.",
		Content: GoalSpec,
	},
	{
		Path:    "peyva/portal/design.md",
		Purpose: "What the Portal should look like, and what it should never look like.",
		Content: DesignBrief,
	},
	{
		Path:    "Your assistant's instruction file",
		Purpose: "How your assistant should work: what it may decide alone, when to stop and ask, and how much to say back.",
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
// reader uses.
//
// It deliberately holds nothing about what to build. The goal, the invariants
// and the constraints on the code all live in goal.md, and stating them twice
// means two places to change and one of them going stale. This file is the
// other half: how the assistant should behave, which goal.md says nothing
// about.
const AgentRules = `# Working with me on peyva

Read peyva/goal.md before your first change. It holds the goal, the rules money
must never break, and the constraints on what you build. Treat it as settled: if
a change would break something in it, the change is wrong.

Before touching the Portal, read peyva/portal/design.md. It holds the look, and
the one visual idea chapter 0 committed to. Keep to it rather than restyling.

## Scope

- Do what the prompt asks, and nothing it did not ask for.
- Build what the current chapter needs, not what a later one might.
- Make the smallest change that satisfies the prompt.
- Leave code you were not asked to touch alone.

## Honesty

- If a constraint makes the task impossible, say so rather than working around
  it. That is a useful answer, not a failure.
- If the task seems to need something goal.md rules out, stop and tell me.
  Do not add it and mention it afterwards.
- Do not tell me something works when you have not run it. Say what you checked.
- When you are unsure, ask. A guess that reads as certainty costs me more than
  the question would have.

## Answers

- Say what changed and why, in a line or two.
- No summary of work I just watched you do, and no restating my request back.
- Give reasoning, trade-offs and alternatives only when I ask for them.
`
