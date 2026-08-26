package content

var Chapter08 = ChapterContent{
	Number:     8,
	Slug:       "chapter-8",
	Title:      "Exactly Once (Idempotency)",
	Subtitle:   "Same request, same ID, same result.",
	Category:   "Reliability",
	Difficulty: "Intermediate",
	EstTime:    "20 min",
	QuickTip:   "The client genuinely can't tell a duplicate response from the real one. That's the point.",

	HeroImage:   "images/chapter-8.webp",
	HeroCaption: "Idempotency = same request (same key) -> same effect (once).",

	Intuition: []string{
		"Alice's phone sends 'transfer $20', the network hiccups, and she taps send again.",
		"Without protection, peyva would debit her twice.",
		"An idempotency key is a reference number. If peyva has seen it before, it returns the same result instead of repeating the work.",
	},

	Concepts: []ConceptItem{
		{Term: "Idempotency Key", Description: "A unique ID the client attaches to a request so retries can be recognized as duplicates."},
		{Term: "Duplicate Request", Description: "The same idempotency key arriving more than once, usually from a retry after a slow or dropped response."},
		{Term: "Idempotent", Description: "An operation that has the same effect whether it runs once or many times with the same key."},
		{Term: "Stored Result", Description: "The response saved next to the key, so a repeat returns exactly what the first attempt returned rather than a freshly computed answer."},
	},

	UnderTheHood: []string{
		"peyva checks the key: new -> process and store the result. Duplicate -> return the stored result, unchanged.",
		"Either way, the client sees the same response, but the money only ever moves once.",
	},

	BuildIt: BuildIt{
		Intro:     "The Teller learns to recognise a payment it has already handled.",
		Technique: "Generated Knowledge Prompting",
		Why:       "Have the assistant produce the relevant facts before it uses them. Making it enumerate how duplicates actually arise gives the design somewhere to put each case, instead of you meeting them in production.",
		Source:    "The Prompt Report: Generated Knowledge",
		Prompt: `The Teller will move money twice if it receives the same payment request twice, a retry after a timeout is indistinguishable from a genuine second payment.

First, list the distinct ways a duplicate request reaches a payment system in the real world. For each, say whether the caller knows the first attempt succeeded.

Then make the Teller recognise a repeat: the caller supplies a reference with the request, and a reference the Teller has already handled returns the original result without moving money again. Store the reference and its response in the same atomic unit that moves the money. Not before, not after.

Then go back over the list you wrote and tell me, case by case, which ones your design now handles and which it doesn't. Include what happens if two requests carrying the same brand-new reference arrive at the same instant.

Done when the same reference twice pays once, two different references pay twice, and both duplicate responses are byte-identical.`,
		UIIntro: "The portal stops punishing an impatient customer.",
		UIPrompt: `A customer who taps send twice must not pay twice. The portal attaches the same reference to a resubmission of the same form, and shows the original result rather than a second payment.

Done when double-submitting the form leaves one payment in the history, and the page looks the same both times.`,
	},

	BreakIt: BreakIt{
		Intro: "Simulate exactly the flaky-network scenario idempotency exists for.",
		Exercises: []string{
			"Send the same transfer request twice with the same idempotency key. Confirm Bob only receives $20 once, not $40.",
			"Send it twice with two different keys: this time it really is two separate $20 transfers, as intended.",
			"Compare the two duplicate responses byte for byte. They're identical.",
		},
	},
}
