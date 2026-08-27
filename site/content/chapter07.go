package content

var Chapter07 = ChapterContent{
	Number:     7,
	Slug:       "chapter-7",
	Title:      "Making It Safe (Transactions)",
	Subtitle:   "A transfer is three writes, not one, and all three need to succeed together.",
	Category:   "Databases",
	Difficulty: "Intermediate",
	QuickTip:   "Nothing is final until COMMIT. That's what makes ROLLBACK possible.",

	HeroImage:   "images/chapter-7.webp",
	HeroCaption: "All or nothing. If something fails, nothing changes.",

	Why: []string{
		"A payment is three writes: take from one, give to the other, write it down. A crash between any two creates, destroys or hides money.",
		"The database keeps writes together and safe. It does not know your rules, and will happily save a negative balance.",
		"Two payments both read 100, both agree 60 fits, both write 40. Alice ends at 40 having sent 120.",
		"The fix is to make the check and the debit one step nothing can slip between: a lock, or the strictest isolation level.",
		"Double-entry: every move takes from one account and gives to another. All entries add up to zero, and each account's add up to its balance.",
		"The record is the proof; the balance is worked out from it. Never edit an entry. Fix a mistake with a new one.",
	},

	Concepts: []ConceptItem{
		{Term: "Transaction", Description: "A group of changes that all succeed or all fail, never half."},
		{Term: "BEGIN / COMMIT", Description: "The start and the permanent end of a transaction. Nothing is final until COMMIT."},
		{Term: "Rollback", Description: "Undoing every change in a transaction because one step failed."},
		{Term: "ACID", Description: "Four promises. Three are the database's job: changes land together, do not tread on each other, and survive a crash. The fourth, consistency, is your own rules, which it cannot know."},
		{Term: "Isolation Level", Description: "How much one running transaction sees of another. The strictest gives the same answer as running them one after another. Weaker ones are faster and let known bugs through."},
		{Term: "Lost Update", Description: "Two transactions read the same balance, each subtracts from what it read, and one subtraction disappears."},
		{Term: "Ledger", Description: "The record of every movement of money, only ever added to. The proof behind each balance."},
	},

	BuildIt: BuildIt{
		Technique: "Chain-of-thought prompting",
		Why:       "The failure modes are the hard part here, and thinking first surfaces the hole while it is still cheap.",
		Source:    "The Prompt Report: Thought Generation, Chain-of-Thought",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Teller moves money by updating two balances. Nothing records that it happened.

Build the Ledger: a record only ever added to. Every payment writes two entries that balance, a debit and a credit with one shared reference, and the balances update in the same transaction. All of it, or none of it.

First: what does the data look like if the process dies between the debit, the credit and the Ledger write? Three cases, who is owed what. Then make them impossible.

Done when a payment leaves two balanced entries, a failure mid-payment leaves none, and alice's entries add up to her balance.`},
			{Label: "Check", Text: `The Teller checks the payer has enough, then debits them, in one transaction.

Alice holds 100. Two payments of 60 arrive at the same instant. Walk me through both reading 100, both deciding 60 fits, both debiting. What is her balance, and which rule in goal.md broke?

Find out from the database's own documentation whether it already prevents this. Then send those two payments at once for real. If the second is not refused, make the check and the debit one step nothing can slip between.

Done when two 60s against 100 leave one applied, one refused, and alice at 40.`},
			{Label: "Think", Portal: true, Thinking: true, Text: `A wallet page shows a balance and can send money. I want to add a history of every movement in and out.

What should it show for a payment that failed halfway? Should it appear at all, and how does a customer tell it from one that worked?

Done when I know what a half-failed payment looks like on the page, and why.`},
			{Label: "Page", Portal: true, Text: `The Portal shows a balance and can send money, but not how the balance got that way.

Add History to the menu: every movement in and out, newest first, with reference, amount and the other party. Show a failed payment the way you just described.

Done when History explains alice's balance without me reading the database.`},
		},
	},
}
