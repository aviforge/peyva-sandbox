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
		Why:       "Describe the finished state, not the steps to it. Given the finish line, it works out what has to change, including the parts you would have forgotten.",
		Source:    "Anthropic: Prompting best practices, Provide clear success criteria",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Vault holds balances in memory, so every restart resets alice to her starting amount.

Keep accounts in a file on disk instead, with SQLite. The file is peyva/vault/peyva.db. Create it on first run, seed alice only if she is missing, and keep no balance in memory. If your language has no SQLite built in, say so and use the most common driver, preferring one with no C toolchain. Name it. The Gateway and Teller do not change from outside.

Say what happens between my write returning and the bytes being safe on disk.

Done when a restart still shows alice's balance, and deleting the file is the only thing that loses it.`},
			{Label: "Try", Reader: true, Text: `Pay, then read. With the program running, run this: it pays bob 20 and reads alice.

You should see: alice at 80.00, or 20 less than she had.`,
				Commands: Commands(
					`curl.exe -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{\"from\":\"alice\",\"to\":\"bob\",\"amount\":20}' -w '\n'
curl.exe -s http://127.0.0.1:9310/accounts/alice -w '\n'`,
					`curl.exe -s -X POST http://127.0.0.1:9310/pay -H "Content-Type: application/json" -d "{\"from\":\"alice\",\"to\":\"bob\",\"amount\":20}" -w "\n"
curl.exe -s http://127.0.0.1:9310/accounts/alice -w "\n"`,
					`curl -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":20}' -w '\n'
curl -s http://127.0.0.1:9310/accounts/alice -w '\n'`,
				)},
			{Label: "Try", Reader: true, Text: `Now stop the program with Ctrl+C, start it again, and read alice once more.

You should see: the same balance as before the restart. Nothing lived in the process, so stopping it lost nothing.`,
				Commands: Commands(
					`curl.exe -s http://127.0.0.1:9310/accounts/alice -w '\n'`,
					`curl.exe -s http://127.0.0.1:9310/accounts/alice -w "\n"`,
					`curl -s http://127.0.0.1:9310/accounts/alice -w '\n'`,
				)},
			{Label: "Try", Reader: true, Text: `Stop it again. Delete the file with this, start the program, and read alice with the command above.

You should see: alice back at 100.00. Everything lived in the file, so deleting it is the one thing that loses money here.`,
				Commands: Commands(
					`Remove-Item peyva\vault\peyva.db`,
					`del peyva\vault\peyva.db`,
					`rm peyva/vault/peyva.db`,
				)},
			{Label: "Portal", Portal: true, Text: `The Portal's balances are written into the page when the program starts, so they are a second copy of what the Vault holds and they are wrong from the first payment onwards. A restart hides that by rewriting the file.

Have the page ask for a balance, and the Vault answer from its file on disk. No balance is written into the page ahead of time.

Done when a payment shows on reload with nothing restarted, and opening the file with nothing running shows the empty state rather than an old number.`},
		},
	},
}
