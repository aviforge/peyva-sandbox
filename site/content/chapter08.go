package content

var Chapter08 = ChapterContent{
	Number:     8,
	Slug:       "chapter-8",
	Title:      "Exactly Once (Idempotency)",
	Subtitle:   "Same request, same ID, same result.",
	Category:   "Reliability",
	Difficulty: "Intermediate",
	QuickTip:   "The client genuinely cannot tell a duplicate response from the real one.",

	HeroImage:   "images/chapter-8.webp",
	HeroCaption: "Idempotency = same request (same key) -> same effect (once).",

	Why: []string{
		"No network can promise a message arrives exactly once. A caller whose request timed out cannot know whether it went through.",
		"So make the effect happen once even when the message arrives twice: the caller attaches a key, and you spot the repeat.",
		"Save the key and the reply in the same transaction as the money. Save it before or after, and a crash breaks it.",
		"Send back the saved reply word for word. Working the answer out again can give a different one.",
		"Two copies of a brand new key arriving at the same instant: let the database refuse the second, not your own code.",
		"A key belongs to one caller. Two customers who pick the same key must not collide.",
	},

	Concepts: []ConceptItem{
		{Term: "Idempotency Key", Description: "A unique ID the caller puts on a request, so a repeat can be spotted as a repeat."},
		{Term: "Duplicate Request", Description: "The same key arriving twice, usually because the first reply was slow or lost and the caller tried again."},
		{Term: "Idempotent", Description: "An action with the same result whether it runs once or ten times with the same key."},
		{Term: "At Least Once", Description: "The best a network offers: a message may arrive more than once, and 'never arrived' looks the same as 'the reply was lost'."},
		{Term: "Stored Result", Description: "The reply saved next to the key, so a repeat gets exactly the first answer back rather than a fresh one."},
	},

	BuildIt: BuildIt{
		Technique: "Generated Knowledge Prompting",
		Why:       "Have it write down what it knows before it designs. The list of ways a payment arrives twice is what the design then has to answer.",
		Source:    "The Prompt Report: Generated Knowledge",
		Prompts: []Prompt{
			{Label: "Build", Text: `Send the same payment twice and the Teller pays twice. A retry after a timeout looks just like a second payment.

First, list the ways a duplicate reaches a payments system in real life.

Then have the caller send a reference with each payment. A reference already handled returns the first answer without moving money. Store the reference and its answer in the same transaction as the money. Let the database refuse the second insert; do not check first in code.

Done when the same reference twice pays once, two references pay twice, two sends of one new reference at the same instant pay once, and every repeat answer is identical.`},
			{Label: "Portal", Portal: true, Text: `A customer who taps Send twice pays twice. Attach the same reference to a resend and show the original result.

Done when double-submitting leaves one payment in History, and the page looks the same both times.`},
		},
	},
}
