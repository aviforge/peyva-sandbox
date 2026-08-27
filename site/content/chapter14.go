package content

var Chapter14 = ChapterContent{
	Number:     14,
	Slug:       "chapter-14",
	Title:      "Big Changes Safely: Sagas",
	Subtitle:   "We can't do everything in one shot. So we do it in steps. If one step fails, we undo the earlier steps so nothing is left half-done.",
	Category:   "Reliability",
	Difficulty: "Advanced",
	QuickTip:   "Compensations run in strict reverse order: the last completed step is undone first.",

	HeroImage:   "images/chapter-14.webp",
	HeroCaption: "Sagas help us complete a big task in steps and roll back safely if something goes wrong.",

	Why: []string{
		"A transaction reaches as far as one database. The moment a payment involves a second system with its own store, there is no COMMIT that covers both, and the atomic unit that protected every earlier chapter is gone.",
		"A saga replaces one atomic step with a sequence of local ones, each committed in its own system, and a record of how far the sequence got. It is weaker than a transaction on purpose: between steps, the world can see a payment that is half done.",
		"Undo is a new forward step, not a rollback. The Ledger is append-only, so reversing a debit is a credit entry that references the original, never a deletion. The history shows that the money left and came back, which is the truth.",
		"Compensation only makes sense for permanent failure. A closed account will still be closed; a timeout might succeed on retry. Reversing after a timeout can put money back that the other side has already credited, so the saga needs the three-outcome thinking from the retry chapter at every step.",
		"The saga's record is what makes it restartable. Each step is idempotent and recorded, so a coordinator that dies mid-saga can be restarted and continue from the last completed step rather than starting again or giving up.",
		"Compensations can fail too. A reversal that cannot be applied leaves a payment stuck, and the honest design surfaces stuck sagas for a human rather than retrying forever or pretending they completed.",
	},

	Concepts: []ConceptItem{
		{Term: "Saga", Description: "A sequence of local transactions across systems, coordinated so the whole workflow either completes or is undone step by step."},
		{Term: "Local Transaction", Description: "Each step's own transaction, scoped to its own store: not shared with the other steps."},
		{Term: "Compensating Action", Description: "The 'undo' for a step that already succeeded, run in reverse order when a later step fails permanently. A new entry, never a deletion."},
		{Term: "Saga Record", Description: "Which steps of a payment have completed, stored durably, so the saga can be resumed after a crash rather than guessed at."},
		{Term: "Permanent vs Retryable", Description: "A closed account is permanent and triggers compensation. A timeout is retryable and must not, because the step may have succeeded."},
	},

	BuildIt: BuildIt{
		Technique: "Least-to-Most Prompting",
		Why:       "One prompt for a whole multi-stage workflow gets you a sketch of all of it and a working version of none.",
		Source:    "The Prompt Report: Decomposition, Least-to-Most",
		Prompts: []Prompt{
			{Label: "Build", Text: `A payment is currently one atomic step inside the Vault. I want the Teller to run payments that span several stages and can unwind if a later stage fails permanently.

Build it in five stages, each on what the last one left. Stop after each and tell me it works before starting the next.

1. A durable record, per payment reference, of which stages have completed, stored in the Vault's database.
2. Wire the existing money movement in as stage one, recording its completion.
3. A stand-in for a second ledger peyva does not own: a small separate process with its own file, that credits an account, refuses permanently for a handle marked closed, and can be made unreachable. Add stage two: crediting the recipient there.
4. A reversal for stage one (put the money back, as a new pair of Ledger entries referencing the original rather than by deleting the old ones), triggered when stage two fails in a way that can never succeed.
5. Resumption: kill the copy between stage one and stage two, restart, and have the saga continue from its record rather than start again or stall.

Distinguish permanent failures from retryable ones. Only permanent failures reverse; a timeout retries with the same reference.

Done when a permanently failing stage two puts the money back, the Ledger shows both the original payment and its reversal, a saga interrupted between stages completes after a restart, and the payment's record shows every stage it passed through.`},
			{Label: "Portal", Portal: true, Text: `A reversed payment currently looks like two unrelated rows. Show it as what it is: the original, and the reversal that answers it, tied together.

Build it in stages. First mark a reversed payment as reversed. Then link the two rows. Then say why it was reversed. Stop after each and show me before starting the next.

Done when a customer can see that money left and came back, and why.`},
		},
	},
}
