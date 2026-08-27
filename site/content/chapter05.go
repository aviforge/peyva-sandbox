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
