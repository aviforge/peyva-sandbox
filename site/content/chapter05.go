package content

var Chapter05 = ChapterContent{
	Number:     5,
	Slug:       "chapter-5",
	Title:      "Storing Money (Databases)",
	Subtitle:   "A bank keeps records in a safe place so nothing is lost.",
	Category:   "Databases",
	Difficulty: "Beginner",
	EstTime:    "15 min",
	QuickTip:   "A file on disk beats a variable in memory. SQLite needs no server process to give you that.",

	HeroImage:   "images/chapter-5.webp",
	HeroCaption: "Durable. Persistent. Survives restarts.",

	Concepts: []ConceptItem{
		{Term: "Durable", Description: "Once written, the data survives the process dying, not just ending."},
		{Term: "Persistent", Description: "Held outside process memory, so it is still there for the next run."},
		{Term: "Migration", Description: "Moving existing data to a new home without losing any of it, here from memory to disk."},
	},

	BuildIt: BuildIt{
		Intro:     "The Vault learns to remember. Balances that survive a restart.",
		Technique: "Explicit success criteria",
		Why:       "Specify the finish line rather than the steps, and the assistant works out the delta, including what you would have forgotten.",
		Source:    "Anthropic: Prompting best practices, Provide clear success criteria",
		Prompts: []Prompt{
			{Label: "Build", Text: `Now: the Vault holds balances in memory. They vanish when the process stops, so every restart resets alice to her seeded amount.

Target: the Vault keeps accounts in a file on disk. It creates its storage on first run and seeds alice only if she isn't already there. Every read and every write of a balance goes to that file. No balance is cached in memory anywhere.

Use SQLite with a driver that needs no cgo. Don't change what the Gateway or the Teller look like from the outside.

Done when restarting the process still reports alice's balance, and deleting the Vault's file is the only thing that loses it.`},
			{Label: "Portal", Portal: true, Intro: "The Portal stops forgetting.", Text: `The Portal reads balances from the Vault's file rather than from whatever was in memory when the page was written.

Done when I make a payment, stop the process, start it again, reload the page, and the new balance is still there.`},
		},
	},

	BreakIt: BreakIt{
		Intro: "Prove durability actually holds, the same way Chapter 0 proved its absence.",
		Exercises: []string{
			"Start peyva, confirm Alice's balance reads 100.",
			"Stop the process entirely, then start it again. The balance still reads 100, read fresh from the Vault's file.",
			"Delete the Vault's file while the process is stopped and restart. Now the account is genuinely gone, because the file itself went, not just the process.",
		},
	},
}
