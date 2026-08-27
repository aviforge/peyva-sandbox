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
			{Label: "Build", Text: `A payment is one step inside the Vault. I want it to run in stages that can be undone if a later stage fails for good.

Build it in five pieces. Stop after each and show me it works.

1. A saved record, per payment, of which stages are done.
2. The existing money movement as stage one.
3. A stand-in for an outside ledger: a small process with its own file that credits an account, refuses for good on a closed handle, and can be switched off. Stage two credits the recipient there.
4. An undo for stage one when stage two fails for good: new Ledger entries pointing at the original, never a deletion.
5. Kill the copy between the stages, restart, and have the payment carry on.

A timeout is not a permanent failure. It retries with the same reference.

Done when a closed recipient puts the money back, and an interrupted payment finishes after a restart.`},
			{Label: "Portal", Portal: true, Text: `A reversed payment looks like two unrelated rows. Show it as the original and the reversal that answers it, tied together, with the reason.

Build it in three steps and show me each.

Done when a customer can see that money left and came back, and why.`},
		},
	},
}
