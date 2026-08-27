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
		Why:       "Enumerating how duplicates arise gives the design somewhere to put each case, before production does it for you.",
		Source:    "The Prompt Report: Generated Knowledge",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Teller will move money twice if it receives the same payment request twice, a retry after a timeout is indistinguishable from a genuine second payment.

First, list the distinct ways a duplicate request reaches a payment system in the real world. For each, say whether the caller knows the first attempt succeeded.

Then make the Teller recognise a repeat: the caller supplies a reference with the request, and a reference the Teller has already handled returns the original result without moving money again. Store the reference and its response in the same atomic unit that moves the money. Not before, not after. Let the database refuse a second insert of the same reference; do not check-then-insert in code.

Then go back over the list you wrote and tell me, case by case, which ones your design now handles and which it doesn't. Include what happens if two requests carrying the same brand-new reference arrive at the same instant, and prove it by sending them.

Done when the same reference twice pays once, two different references pay twice, two simultaneous sends of one new reference pay once, and every duplicate response is byte-identical to the first.`},
			{Label: "Portal", Portal: true, Text: `Send posts the form and takes whatever comes back, so a customer who taps twice pays twice. Attach the same reference to a resubmission, and show the original result.

Done when double-submitting the form leaves one payment in History, and the page looks the same both times.`},
		},
	},
}
