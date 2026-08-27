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
		"Exactly-once delivery over a network is impossible. A caller whose request timed out cannot know whether it was processed; the response may have been lost after the money moved. Its only options are to give up, or to send again and risk a double payment.",
		"So the goal is not exactly-once delivery but exactly-once effect: deliver at least once, and make the receiver recognise a repeat. That is what an idempotency key is, a name the caller gives the operation so the second arrival can be matched to the first.",
		"The key and the outcome must be stored in the same transaction as the money. Store the key first and crash, and a legitimate retry is refused forever. Store it after and crash, and the retry pays twice. Same atomic unit, or it does not work.",
		"The stored response is returned verbatim on a repeat. Recomputing it, even from the same data, can give a different answer if anything changed in between, and then the caller has two answers to one question.",
		"Two requests with the same brand-new key at the same instant are the hard case. Both find no record, both try to insert. A unique constraint on the key makes the database refuse the second one, and it is the database, not your code, that has to be the referee.",
		"The key has a scope: it is per caller, not global. Two customers who happen to pick the same key must not collide, so the key is stored alongside who sent it.",
	},

	Concepts: []ConceptItem{
		{Term: "Idempotency Key", Description: "A unique ID the client attaches to a request so retries can be recognised as duplicates."},
		{Term: "Duplicate Request", Description: "The same idempotency key arriving more than once, usually from a retry after a slow or dropped response."},
		{Term: "Idempotent", Description: "An operation that has the same effect whether it runs once or many times with the same key."},
		{Term: "At Least Once", Description: "The best a network can offer: a message may arrive more than once, and never arriving is indistinguishable from a lost reply. Exactly-once delivery does not exist."},
		{Term: "Stored Result", Description: "The response saved next to the key, so a repeat returns exactly what the first attempt returned rather than a freshly computed answer."},
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
