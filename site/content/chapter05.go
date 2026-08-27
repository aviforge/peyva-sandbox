package content

var Chapter05 = ChapterContent{
	Number:     5,
	Slug:       "chapter-5",
	Title:      "Storing Money (Databases)",
	Subtitle:   "A bank keeps records in a safe place so nothing is lost.",
	Category:   "Databases",
	Difficulty: "Beginner",
	QuickTip:   "A file on disk beats a variable in memory, and SQLite needs no server to give you one.",

	HeroImage:   "images/chapter-5.webp",
	HeroCaption: "Durable. Persistent. Survives restarts.",

	Why: []string{
		"Durable means the disk has confirmed the write. Until then it sits in memory, and a power cut takes it.",
		"A database is what makes 'saved' mean something: it orders writes, forces them to disk, and repairs half-finished ones after a crash.",
		"Once the data lives outside the program, every read is a trip and every write is a wait.",
		"Keep one official copy. Two copies of a balance disagree the first time only one of them is updated.",
		"SQLite is a real database in a single file. No server to install, and enough for everything here.",
	},

	Concepts: []ConceptItem{
		{Term: "Durable", Description: "Written where it survives the program being killed. Only true once the disk has confirmed it."},
		{Term: "Persistent", Description: "Kept outside the program's memory, so it is still there next run."},
		{Term: "Source of Truth", Description: "The one place a value is official. Anything else holding it is a copy that can be wrong."},
		{Term: "Migration", Description: "Moving data you already have to a new home without losing any, here from memory to disk."},
	},

	BuildIt: BuildIt{
		Technique: "Explicit success criteria",
		Why:       "Specify the finish line rather than the steps, and the assistant works out the delta, including what you would have forgotten.",
		Source:    "Anthropic: Prompting best practices, Provide clear success criteria",
		Prompts: []Prompt{
			{Label: "Build", Text: `Now: the Vault holds balances in memory, so every restart resets alice.

Target: the Vault keeps accounts in a file on disk. It creates the file on first run and seeds alice only if she is not already there. Every read and write of a balance goes to that file, and no balance is kept in memory.

Use SQLite. If your language's standard library has no SQLite, say so and use the most widely used driver, preferring one that needs no C toolchain. Name it: that is the one exception to standard library only. Do not change how the Gateway or the Teller look from outside.

Tell me what happens between my write returning and the bytes being safe on disk, and what a power cut in that gap costs.

Done when a restart still reports alice's balance, and deleting the file is the only thing that loses it.`},
			{Label: "Portal", Portal: true, Text: `The Portal reads balances from the Vault's file, not from whatever was in memory when the page was written.

Done when I make a payment, restart the process, reload the page, and the new balance is still there.`},
		},
	},
}
