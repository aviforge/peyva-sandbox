package content

var Chapter07 = ChapterContent{
	Number:     7,
	Slug:       "chapter-7",
	Title:      "Making It Safe (Transactions)",
	Subtitle:   "A transfer is three writes, not one — and all three need to succeed together.",
	Category:   "Databases",
	Difficulty: "Intermediate",
	EstTime:    "20 min",
	QuickTip:   "Nothing is final until COMMIT — that's what makes ROLLBACK possible.",

	HeroImage:   "images/chapter-7.webp",
	HeroCaption: "All or nothing. If something fails, nothing changes.",

	Intuition: []string{
		"A transfer isn't one change, it's three: debit Alice, credit Bob, record it happened.",
		"If peyva crashes between steps, money vanishes.",
		"A bank teller doesn't hand over a receipt until every step succeeds — if one fails, they undo everything.",
	},

	Concepts: []ConceptItem{
		{Term: "Transaction", Description: "A group of database changes that succeed or fail together, never partially."},
		{Term: "BEGIN / COMMIT", Description: "Marks the start and the permanent end of a transaction — nothing is final until COMMIT."},
		{Term: "Rollback", Description: "Undoing every change in a transaction because one step failed."},
		{Term: "ACID", Description: "Atomicity, Consistency, Isolation, Durability — the four guarantees a real transaction gives you."},
		{Term: "Ledger", Description: "The append-only record of every movement of money — the proof behind each balance the Vault reports."},
	},

	UnderTheHood: []string{
		"Transaction flow: BEGIN -> Debit Alice -> Credit Bob -> Insert Transfer Record -> COMMIT.",
		"If any step fails, ROLLBACK undoes every change made since BEGIN — Alice never loses money to a step that never finished.",
	},

	BuildIt: BuildIt{
		Intro:     "Build the Ledger — the proof behind every balance the Vault reports.",
		Technique: "Chain-of-thought prompting",
		Why:       "Ask for the reasoning before the answer and the reasoning becomes inspectable. When the failure modes are the hard part, code-first gets you a plausible design with a hole in it; thinking first surfaces the hole while it's still cheap.",
		Source:    "The Prompt Report — Thought Generation, Chain-of-Thought",
		Prompt: "The Teller moves money by updating two balances in the Vault. Nothing records that the movement happened, so if a balance looks wrong there's no way to prove how it got that way.\n\n" +
			"Build the Ledger: an append-only, double-entry record. Every payment writes two balanced entries — a debit and a credit sharing one reference and timestamp — and the Vault's balances update in the same atomic unit. All of it commits, or none of it happened.\n\n" +
			"Before writing code: walk me through what the data looks like if the process dies between the debit, the credit, and the Ledger write. Three scenarios, and for each tell me exactly who is owed what.\n\n" +
			"Then make those states unreachable, and tell me which of your three scenarios is still possible afterwards, if any.\n\n" +
			"Done when a completed payment leaves two balanced Ledger entries and updated balances, a forced mid-payment failure leaves neither, and summing alice's Ledger entries equals her Vault balance.",
	},

	BreakIt: BreakIt{
		Intro: "Force a failure partway through and confirm nothing is left half-done.",
		Exercises: []string{
			"Deliberately make the credit step fail (e.g. an invalid 'to' account) and confirm Alice's debit was rolled back too — her balance is unchanged.",
			"Compare this to Chapter 4's version without a transaction, where a mid-transfer crash really would leave Alice debited with no credit to Bob.",
			"Kill the process between BEGIN and COMMIT (if you can time it) — on restart, the transaction never happened at all.",
		},
	},
}
