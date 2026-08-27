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
			{Label: "Build", Text: `The Vault holds balances in memory, so every restart resets alice to her starting amount.

Keep accounts in a file on disk instead, with SQLite. Create it on first run, seed alice only if she is missing, and keep no balance in memory. If your language has no SQLite built in, say so and use the most common driver, preferring one with no C toolchain. Name it. The Gateway and Teller do not change from outside.

Say what happens between my write returning and the bytes being safe on disk.

Done when a restart still shows alice's balance, and deleting the file is the only thing that loses it.`},
			{Label: "Page", Portal: true, Text: `The Portal still shows balances from memory. Read them from the Vault's file on disk instead.

Done when I make a payment, restart, reload, and the new balance is still there.`},
		},
	},
}
