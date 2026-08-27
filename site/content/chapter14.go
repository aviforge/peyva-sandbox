package content

var Chapter14 = ChapterContent{
	Number:     14,
	Slug:       "chapter-14",
	Title:      "Big Changes Safely: Sagas",
	Subtitle:   "We can't do everything in one shot. So we do it in steps. If one step fails, we undo the earlier steps so nothing is left half-done.",
	Category:   "Reliability",
	Difficulty: "Advanced",
	QuickTip:   "Undo in reverse: the last finished step is the first one undone.",

	HeroImage:   "images/chapter-14.webp",
	HeroCaption: "Sagas help us complete a big task in steps and roll back safely if something goes wrong.",

	Why: []string{
		"A transaction covers one database. Two systems with their own databases share no COMMIT.",
		"A saga is small steps, each saved on its own, plus a record of how far you got. Between steps, the world sees a half-done payment.",
		"Undo is a new entry forward, never a deletion. The history shows the money leaving and coming back.",
		"Only undo after a failure that can never succeed. A timeout gets tried again, because the step may have worked.",
		"Each step is safe to repeat and written down, so a crash picks up at the last finished step.",
		"An undo can fail too. A payment stuck like that goes to a person, not into an endless retry.",
	},

	Concepts: []ConceptItem{
		{Term: "Saga", Description: "A long job done as a sequence of small steps, each saved on its own, that can be undone step by step."},
		{Term: "Local Transaction", Description: "One step's own transaction, inside its own database, not shared with the other steps."},
		{Term: "Compensating Action", Description: "The undo for a step that already worked, run when a later step fails for good. A new entry, never a deletion."},
		{Term: "Saga Record", Description: "Which steps of a payment are done, saved to disk, so a crash can pick up where it stopped."},
		{Term: "Permanent vs Retryable", Description: "A closed account will always be closed, so undo. A timeout might have worked, so try again instead."},
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
