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

	Intuition: []string{
		"Chapter 12's queue has a gap: if peyva commits the transfer but crashes before the message reaches the queue, Bob never gets told.",
		"The Outbox pattern writes the transfer and the message in the same transaction.",
		"Either both happen or neither does.",
	},

	Concepts: []ConceptItem{
		{Term: "Outbox Table", Description: "A table in the same database as the transfer, holding messages that still need to be published."},
		{Term: "Same Transaction", Description: "The transfer and the outbox row are written together, inside the transaction from Chapter 7: never as two separate steps."},
		{Term: "Outbox Publisher", Description: "A background worker that reads unsent outbox rows and publishes them to the queue, then marks them sent."},
		{Term: "At Least Once", Description: "The publisher can send a message it has already sent, if it crashes after publishing but before marking the row done. The receiver has to cope with that, which is what Chapter 8 built."},
	},

	UnderTheHood: []string{
		"1. Service saves the order and writes the message to the outbox, in one transaction. 2. A publisher reads unsent outbox rows. 3. Publisher pushes them to the queue.",
		"If the service crashes right after saving, the message is still sitting in the outbox. The publisher sends it later.",
		"No lost messages, and automatic recovery: the background worker retries until every outbox row is sent.",
	},

	BuildIt: BuildIt{
		Intro:     "The Courier learns to never lose work handed to it.",
		Technique: "Contrastive Chain-of-Thought",
		Why:       "Show the wrong reasoning alongside the right reasoning, not just the right one. Naming the naive design and why it fails stops the assistant rediscovering it, and forces it to say what its version does differently.",
		Source:    "The Prompt Report: Few-Shot CoT, Contrastive CoT",
		Prompts: []Prompt{
			{Label: "Build", Intro: "Two designs, one of them named as wrong.", Text: `The Courier picks up work from memory. If the process dies between the payment committing and the work reaching the Courier, that work is gone and nothing in the system knows it's missing.

Reasoning I want you to reject: commit the payment, then hand the work to the Courier. Both nearly always succeed, so the gap between them is too small to matter. That is wrong because the gap isn't a probability, it's a window, and a crash inside it loses work silently with no record that anything is owed.

Reasoning I want you to follow: anything that must happen because a payment happened is recorded in the same atomic unit as the payment. One commit, or nothing.

Build the second one. The Teller records the Courier's pending work in the same atomic unit that moves the money. The Courier collects from that record, and marks each item done once it's delivered.

Done when killing the process right after a payment leaves the work durable and uncollected, and restarting the Courier still delivers it.`},
			{Label: "Contrast", Thinking: true, Intro: "Where the two designs part company.", Text: `You built the Teller recording the Courier's pending work in the same atomic unit as the payment, instead of handing it over after the commit.

Contrast the two designs directly: name the exact instant at which the rejected one loses work and yours doesn't.

Then tell me what happens if the Courier dies after delivering but before marking the item done, and whether that is acceptable for a notification.

Done when I can point at the single instant that separates the two designs, and I know what your version does in the one case it still handles imperfectly.`},
		},
	},

	BreakIt: BreakIt{
		Intro: "Prove the message survives a crash that Chapter 12's plain queue wouldn't have.",
		Exercises: []string{
			"Simulate a crash right after the transfer transaction commits, before the publisher runs: restart peyva and confirm the outbox row is still there, unsent.",
			"Start the publisher and confirm it picks up that leftover row and delivers it. Nothing was lost.",
			"Compare this to Chapter 12's design: if peyva crashed between committing the transfer and enqueueing the message there, the notification would be gone forever.",
		},
	},
}
