package content

var Chapter07 = ChapterContent{
	Number:     7,
	Slug:       "chapter-7",
	Title:      "Making It Safe (Transactions)",
	Subtitle:   "A transfer is three writes, not one, and all three need to succeed together.",
	Category:   "Databases",
	Difficulty: "Intermediate",
	EstTime:    "20 min",
	QuickTip:   "Nothing is final until COMMIT. That's what makes ROLLBACK possible.",

	HeroImage:   "images/chapter-7.webp",
	HeroCaption: "All or nothing. If something fails, nothing changes.",

	Concepts: []ConceptItem{
		{Term: "Transaction", Description: "A group of database changes that succeed or fail together, never partially."},
		{Term: "BEGIN / COMMIT", Description: "Marks the start and the permanent end of a transaction. Nothing is final until COMMIT."},
		{Term: "Rollback", Description: "Undoing every change in a transaction because one step failed."},
		{Term: "ACID", Description: "Atomicity, Consistency, Isolation, Durability: the four guarantees a real transaction gives you."},
		{Term: "Ledger", Description: "The append-only record of every movement of money, the proof behind each balance the Vault reports."},
	},

	BuildIt: BuildIt{
		Intro:     "Build the Ledger: the proof behind every balance the Vault reports.",
		Technique: "Chain-of-thought prompting",
		Why:       "The failure modes are the hard part here, and thinking first surfaces the hole while it is still cheap.",
		Source:    "The Prompt Report: Thought Generation, Chain-of-Thought",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Teller moves money by updating two balances in the Vault. Nothing records that the movement happened, so if a balance looks wrong there's no way to prove how it got that way.

Build the Ledger: an append-only, double-entry record. Every payment writes two balanced entries (a debit and a credit sharing one reference and timestamp), and the Vault's balances update in the same atomic unit. All of it commits, or none of it happened.

Before writing code: walk me through what the data looks like if the process dies between the debit, the credit, and the Ledger write. Three scenarios, and for each tell me exactly who is owed what.

Then make those states unreachable, and tell me which of your three scenarios is still possible afterwards, if any.

Done when a completed payment leaves two balanced Ledger entries and updated balances, a forced mid-payment failure leaves neither, and summing alice's Ledger entries equals her Vault balance.`},
			{Label: "Decide", Portal: true, Thinking: true, Intro: "Decide what a failed payment looks like on the page.", Text: `A customer's wallet page shows a balance and can send money. I want to add a history: every movement in and out of their own account.

Before building it, walk me through what that view should show for a payment that failed halfway. Should it appear at all? If it should, how does a customer tell it apart from one that worked?

Done when I know what a half-failed payment looks like on the page, and why.`},
			{Label: "Portal", Portal: true, Intro: "The history view.", Text: `The Portal shows a balance and can send money, and says nothing about how the balance came to be what it is.

Add History to the menu: every movement in and out of the customer's own account, newest first, each with its reference, amount and the other party. Show a failed payment the way you just described.

Done when History explains alice's balance without me reading the database.`},
		},
	},
}
