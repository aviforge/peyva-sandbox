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

	Intuition: []string{
		"Chapter 0 proved a hard truth: stop the process and Alice's balance is gone.",
		"A database is a bank record room: a safe place outside process memory.",
		"A restart or crash never erases what Alice has.",
	},

	Concepts: []ConceptItem{
		{Term: "Write", Description: "Saving new or updated data to the database, inserting an account or updating a balance."},
		{Term: "Read", Description: "Retrieving data back out, checking a balance before showing it to Alice."},
		{Term: "Durable", Description: "Once written, the data survives even if the process (or the machine) restarts."},
		{Term: "Persistent", Description: "Data stored on disk, outside process memory, so it outlives any single run."},
	},

	UnderTheHood: []string{
		"peyva App writes to and reads from a Database instead of an in-memory struct.",
		"SQLite fits a local setup best: a real database that's just a file on disk, no server required.",
	},

	BuildIt: BuildIt{
		Intro:     "The Vault learns to remember. Balances that survive a restart.",
		Technique: "Explicit success criteria",
		Why:       "For a migration the useful thing to specify isn't the steps, it's the finish line: what must be true when it's done. State the target precisely and the assistant works out the delta, including the parts you'd have forgotten to list.",
		Source:    "Anthropic: Prompting best practices, Provide clear success criteria",
		Prompt: "Now: the Vault holds balances in memory. They vanish when the process stops, so every restart resets alice to her seeded amount.\n\n" +
			"Target: the Vault keeps accounts in a file on disk. It creates its storage on first run and seeds alice only if she isn't already there. Every read and every write of a balance goes to that file. No balance is cached in a Go variable anywhere.\n\n" +
			"Use SQLite with a driver that needs no cgo. Don't change what the Gateway or the Teller look like from the outside.\n\n" +
			"Done when restarting the process still reports alice's balance, and deleting the Vault's file is the only thing that loses it.",
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
