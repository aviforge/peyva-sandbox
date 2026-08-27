package content

var Chapter13 = ChapterContent{
	Number:     13,
	Slug:       "chapter-13",
	Title:      "Reliability Patterns: Transactional Outbox",
	Subtitle:   "We don't send the email right away. First, we save the order and write the email into the outbox, in the same step.",
	Category:   "Reliability",
	Difficulty: "Advanced",
	QuickTip:   "If it isn't written in the same transaction as the data change, it can be lost.",

	HeroImage:   "images/chapter-13.webp",
	HeroCaption: "Transactional Outbox makes sure important work is never lost, even if something fails right after saving.",

	Why: []string{
		"Save the payment, then send the message: a crash in the gap between those two lines loses the message, with no record it was owed.",
		"The outbox writes the message in the same transaction as the payment. Both exist, or neither does.",
		"Crash between sending and marking it done, and it goes twice. You have traded 'maybe never' for 'maybe twice'.",
		"Deliver at least once, then spot the repeats. The same idea as chapter 8's key, now applied to the work after a payment.",
		"Three copies reading one table can grab the same job. Have each claim it in a single update, and let the database decide.",
	},

	Concepts: []ConceptItem{
		{Term: "Outbox Table", Description: "A table in the same database as the payment, holding messages still to be sent."},
		{Term: "Same Transaction", Description: "The payment and the message are written together, in the transaction from chapter 7, never as two steps."},
		{Term: "Outbox Publisher", Description: "A worker that reads unsent rows, sends them, then marks them done."},
		{Term: "At Least Once", Description: "A crash after sending and before marking done sends it twice. The receiver has to cope, which is what chapter 8 built."},
		{Term: "Claiming", Description: "Taking a row in a single update before working on it, so two workers do not both take it. Not the same as marking it done."},
	},

	BuildIt: BuildIt{
		Technique: "Contrastive Chain-of-Thought",
		Why:       "Naming the naive design stops the assistant rediscovering it, and forces it to say what its version does differently.",
		Source:    "The Prompt Report: Few-Shot CoT, Contrastive CoT",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Courier takes work from a queue in memory inside each copy. If a copy dies with work in it, that work is gone and nothing knows.

Reasoning to reject: save the payment, then hand the work over, because the gap is too small to matter. It is not a probability, it is a gap, and a crash inside it loses work with no record anything was owed.

Reasoning to follow: anything that must happen because a payment happened is written in the same transaction as the payment.

Build the second. The Vault writes the Courier's pending work in the transaction that moves the money. The Courier in each copy takes work from the Vault, claims each item in a single update so two copies cannot both take it, delivers it, and marks it done.

Done when killing a copy right after a payment leaves the work saved and untaken, restarting any copy delivers it, and three copies collecting at once deliver each item once.`},
			{Label: "Contrast", Thinking: true, Text: `You had the Vault write the Courier's pending work in the same transaction as the payment, instead of handing it over after the commit.

Name the exact instant at which the rejected design loses work and yours does not.

Then: what happens if the Courier dies after delivering but before marking the item done, and is that acceptable for a notification? What happens if a copy claims an item and dies before delivering it, and what would it take for someone else to deliver it?

Done when I can point at the single instant that separates the two, and I know the two cases yours still handles imperfectly.`},
		},
	},
}
