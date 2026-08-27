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
		"Durability is a promise about the moment of failure, not about the file. A write that sits in an OS buffer is on its way to disk; it is not there until the database has called fsync and the drive has confirmed. Kill the power in between and the write never happened.",
		"A database is the code that makes 'written' mean something. It orders writes, flushes them, and on restart replays or discards partial ones so the file is never half-updated. Writing your own file format skips all of that.",
		"Once storage is outside the process, every read is a trip and every write is a wait. Latency that was nanoseconds is now microseconds to milliseconds, which changes what you can afford to do per request.",
		"Where the truth lives matters more than how fast you can read it. Two copies of a balance, one in memory and one on disk, will disagree the first time one is updated and the other is not. One place is the source of truth; everything else is a view of it.",
		"SQLite is a real database in one file: transactions, a write-ahead log, crash recovery. It is not a toy, and it is enough to teach every storage idea in this book without a server to install.",
		"Migration is the ordinary case. Data always outlives the code that first wrote it, so moving records from one home to another without losing any is a skill you will use more often than designing a fresh schema.",
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
