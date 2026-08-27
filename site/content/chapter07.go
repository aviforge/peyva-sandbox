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
		{Term: "ACID", Description: "The database keeps changes together, apart, and safe on disk. The C, consistency, is your own rules, which it cannot know."},
		{Term: "Isolation Level", Description: "How much one running transaction sees of another. The strictest gives the same answer as running them one after another. Weaker ones are faster and let known bugs through."},
		{Term: "Lost Update", Description: "Two transactions read the same balance, each subtracts from what it read, and one subtraction disappears."},
		{Term: "Ledger", Description: "The record of every movement of money, only ever added to. The proof behind each balance."},
	},

	BuildIt: BuildIt{
		Technique: "Chain-of-thought prompting",
		Why:       "The failure modes are the hard part here, and thinking first surfaces the hole while it is still cheap.",
		Source:    "The Prompt Report: Thought Generation, Chain-of-Thought",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Teller moves money by updating two balances in the Vault. Nothing records that the movement happened, so if a balance looks wrong there's no way to prove how it got that way.

Build the Ledger: an append-only, double-entry record. Every payment writes two balanced entries (a debit and a credit sharing one reference and timestamp), and the Vault's balances update in the same atomic unit. All of it commits, or none of it happened.

Before writing code: walk me through what the data looks like if the process dies between the debit, the credit, and the Ledger write. Three scenarios, and for each tell me exactly who is owed what.

Then make those states unreachable, and tell me which of your three scenarios is still possible afterwards, if any.

Done when a completed payment leaves two balanced Ledger entries and updated balances, a forced mid-payment failure leaves neither, and summing alice's Ledger entries equals her Vault balance.`},
			{Label: "Isolation", Text: `The Teller checks that the payer has enough and then debits them, inside one transaction.

Before changing anything, reason through this: two payments from alice of 60 each arrive at the same instant, and she holds 100. Walk through the interleaving where both read 100, both decide 60 fits, and both debit. Say what her balance is afterwards, whether the Ledger and the Vault still agree, and which invariant in goal.md was broken.

Then tell me which isolation level, or which lock, the database you are using applies to that transaction, and whether it prevents that interleaving. Do not assume; find out from the database's own documentation or by testing it.

Then fire those two payments concurrently, for real, and show me the result. If the second one is not refused, fix it so the read and the debit are one unit, and run it again.

Done when two concurrent 60s from a balance of 100 leave exactly one applied, one refused, alice at 40, and the Ledger summing to her balance.`},
			{Label: "Decide", Portal: true, Thinking: true, Text: `A customer's wallet page shows a balance and can send money. I want to add a history: every movement in and out of their own account.

Before building it, walk me through what that view should show for a payment that failed halfway. Should it appear at all? If it should, how does a customer tell it apart from one that worked?

Done when I know what a half-failed payment looks like on the page, and why.`},
			{Label: "Portal", Portal: true, Text: `The Portal shows a balance and can send money, and says nothing about how the balance came to be what it is.

Add History to the menu: every movement in and out of the customer's own account, newest first, each with its reference, amount and the other party. Show a failed payment the way you just described.

Done when History explains alice's balance without me reading the database.`},
		},
	},
}
