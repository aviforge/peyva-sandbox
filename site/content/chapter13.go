package content

var Chapter13 = ChapterContent{
	Number:     13,
	Slug:       "chapter-13",
	Title:      "Reliability Patterns: Transactional Outbox",
	Subtitle:   "We don't send the email right away. First, we save the order and write the email into the outbox, in the same step.",
	Category:   "Reliability",
	Difficulty: "Advanced",
	EstTime:    "20 min",
	QuickTip:   "If it isn't written in the same transaction as the data change, it can be lost.",

	HeroImage:   "images/chapter-13.webp",
	HeroCaption: "Transactional Outbox makes sure important work is never lost, even if something fails right after saving.",

	Concepts: []ConceptItem{
		{Term: "Outbox Table", Description: "A table in the same database as the transfer, holding messages that still need to be published."},
		{Term: "Same Transaction", Description: "The transfer and the outbox row are written together, inside the transaction from Chapter 7: never as two separate steps."},
		{Term: "Outbox Publisher", Description: "A background worker that reads unsent outbox rows and publishes them to the queue, then marks them sent."},
		{Term: "At Least Once", Description: "A crash between publishing and marking the row done sends the message twice. The receiver has to cope, which is what Chapter 8 built."},
	},

	BuildIt: BuildIt{
		Technique: "Contrastive Chain-of-Thought",
		Why:       "Naming the naive design stops the assistant rediscovering it, and forces it to say what its version does differently.",
		Source:    "The Prompt Report: Few-Shot CoT, Contrastive CoT",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Courier picks up work from memory. If the process dies between the payment committing and the work reaching the Courier, that work is gone and nothing in the system knows it's missing.

Reasoning I want you to reject: commit the payment, then hand the work to the Courier, because the gap between them is too small to matter. The gap is not a probability, it is a window, and a crash inside it loses work with no record that anything is owed.

Reasoning I want you to follow: anything that must happen because a payment happened is recorded in the same atomic unit as the payment. One commit, or nothing.

Build the second one. The Teller records the Courier's pending work in the same atomic unit that moves the money. The Courier collects from that record, and marks each item done once it's delivered.

Done when killing the process right after a payment leaves the work durable and uncollected, and restarting the Courier still delivers it.`},
			{Label: "Contrast", Thinking: true, Text: `You built the Teller recording the Courier's pending work in the same atomic unit as the payment, instead of handing it over after the commit.

Contrast the two designs directly: name the exact instant at which the rejected one loses work and yours doesn't.

Then tell me what happens if the Courier dies after delivering but before marking the item done, and whether that is acceptable for a notification.

Done when I can point at the single instant that separates the two designs, and I know what your version does in the one case it still handles imperfectly.`},
		},
	},
}
