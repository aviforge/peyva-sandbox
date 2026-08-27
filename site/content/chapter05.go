package content

var Chapter05 = ChapterContent{
	Number:     5,
	Slug:       "chapter-5",
	Title:      "Storing Money (Databases)",
	Subtitle:   "A bank keeps records in a safe place so nothing is lost.",
	Category:   "Databases",
	Difficulty: "Beginner",
	QuickTip:   "A file on disk beats a variable in memory. SQLite needs no server process to give you that.",

	HeroImage:   "images/chapter-5.webp",
	HeroCaption: "Durable. Persistent. Survives restarts.",

	Why: []string{
		"Durable means flushed to disk and confirmed. A write sitting in an OS buffer is lost if the power goes.",
		"A database is the code that makes 'written' mean something: ordering, flushing, crash recovery.",
		"Storage outside the process turns every read into a trip and every write into a wait.",
		"One place is the source of truth. A second copy of a balance will disagree the first time only one is updated.",
		"SQLite is a real database in one file, enough to teach every storage idea here.",
	},

	Concepts: []ConceptItem{
		{Term: "Durable", Description: "Once written, the data survives the process dying, not just ending. Only true after the database has flushed it to disk."},
		{Term: "Persistent", Description: "Held outside process memory, so it is still there for the next run."},
		{Term: "Source of Truth", Description: "The one place a value is authoritative. Anything else holding that value is a copy that can be wrong."},
		{Term: "Migration", Description: "Moving existing data to a new home without losing any of it, here from memory to disk."},
	},

	BuildIt: BuildIt{
		Technique: "Explicit success criteria",
		Why:       "Specify the finish line rather than the steps, and the assistant works out the delta, including what you would have forgotten.",
		Source:    "Anthropic: Prompting best practices, Provide clear success criteria",
		Prompts: []Prompt{
			{Label: "Build", Text: `Now: the Vault holds balances in memory. They vanish when the process stops, so every restart resets alice to her seeded amount.

Target: the Vault keeps accounts in a file on disk. It creates its storage on first run and seeds alice only if she isn't already there. Every read and every write of a balance goes to that file. No balance is cached in memory anywhere.

Use SQLite. If the language's standard library has no SQLite binding, say so and use the single most widely used driver for it, preferring one that needs no C toolchain; that is the one exception to standard library only, and name it so I can see it. Don't change what the Gateway or the Teller look like from the outside.

Tell me what the database does between my write returning and the bytes being safe on disk, and what happens to a write if the power goes in that gap.

Done when restarting the process still reports alice's balance, and deleting the Vault's file is the only thing that loses it.`},
			{Label: "Portal", Portal: true, Text: `The Portal reads balances from the Vault's file rather than from whatever was in memory when the page was written.

Done when I make a payment, stop the process, start it again, reload the page, and the new balance is still there.`},
		},
	},
}
