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
		"A payment is at least three writes: debit one account, credit another, record that it happened. A crash between any two of them leaves money created, destroyed or unexplained. A transaction is the mechanism that makes three writes one event.",
		"Atomicity and durability are what the database gives you. Consistency in ACID is your invariant, not the database's: it will happily commit a negative balance if your code asks it to. The database keeps the writes together; you decide what they may say.",
		"Isolation is the property most people skip and the one that loses money. Two payments from alice, each reading her balance as 100, each checking 20 fits, each debiting 20 from what it read, leave her at 80 with 40 gone. That is a lost update, and it happens inside two perfectly atomic transactions.",
		"The fix is an isolation level, or a lock, that makes the read and the write one unit. Serialisable isolation makes concurrent transactions behave as if they ran one after another; weaker levels permit named anomalies. SQLite is serialisable because it holds one write lock, which is why this is easy here and hard in a server database.",
		"Double-entry is the oldest correctness check in finance. Every movement is a debit and a matching credit, so the sum of all entries is always zero and the sum of an account's entries is always its balance. It turns 'is the money right' into arithmetic anyone can run.",
		"The append-only record is the proof; the balance is a cached answer derived from it. When they disagree, the record wins, which is why nothing ever edits or deletes an entry. A mistake is corrected by another entry.",
	},

	Concepts: []ConceptItem{
		{Term: "Transaction", Description: "A group of database changes that succeed or fail together, never partially."},
		{Term: "BEGIN / COMMIT", Description: "Marks the start and the permanent end of a transaction. Nothing is final until COMMIT."},
		{Term: "Rollback", Description: "Undoing every change in a transaction because one step failed."},
		{Term: "ACID", Description: "Atomicity, Isolation and Durability are what the database provides. Consistency is your invariants, which the database cannot know."},
		{Term: "Isolation Level", Description: "How much one running transaction can see of another. Serialisable means the outcome equals some one-at-a-time ordering; weaker levels each permit named anomalies."},
		{Term: "Lost Update", Description: "Two transactions read the same balance, each subtracts from what it read, and one's subtraction vanishes. Atomicity alone does not prevent it."},
		{Term: "Ledger", Description: "The append-only record of every movement of money, the proof behind each balance the Vault reports."},
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
