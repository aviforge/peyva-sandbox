package content

var Chapter09 = ChapterContent{
	Number:     9,
	Slug:       "chapter-9",
	Title:      "Giving Up Well: Retries, Timeouts and Backoff",
	Subtitle:   "A call that has not answered is not a call that failed. Knowing when to stop waiting is the whole skill.",
	Category:   "Reliability",
	Difficulty: "Intermediate",
	QuickTip:   "A retry without an idempotency key is a duplicate payment with extra steps.",

	HeroImage:   "images/chapter-9.webp",
	HeroCaption: "Every call over a network has three outcomes: it worked, it failed, and nobody knows. The third is the one you design for.",

	Why: []string{
		"A call has three endings, not two: it worked, it failed, or you never heard back. The third may be a slow success.",
		"Every call to another program needs a time limit, set from how long a good answer really takes.",
		"Only retry with the same reference. Without one, a retry is a second payment.",
		"Trying again at once piles onto a service already in trouble. Wait longer each time, plus a random bit so callers do not return together.",
		"Never retry a 4xx. It will be a 4xx again.",
		"Stop after a few tries. The honest answer is then 'we do not know, quote this reference', not an endless loop.",
	},

	Concepts: []ConceptItem{
		{Term: "Timeout", Description: "How long a caller waits before giving up on an answer. Set it from how long a good answer really takes."},
		{Term: "Unknown Outcome", Description: "The call timed out. It may have worked, failed, or be about to work. The caller cannot tell, and must not guess."},
		{Term: "Retry", Description: "Sending the same request again. Safe only when it carries a reference the other side uses to spot repeats."},
		{Term: "Exponential Backoff", Description: "Waiting longer before each try: one second, then two, then four. Gives the other side time to recover."},
		{Term: "Jitter", Description: "A small random amount added to each wait, so that many callers do not all try again at the same instant."},
		{Term: "Retry Storm", Description: "Everyone giving up and trying again at once, piling load on a service exactly when it is struggling."},
	},

	BuildIt: BuildIt{
		Technique: "Plan-and-Solve",
		Why:       "Ask for the plan first, then build from it. A table of what can go wrong is a specification beforehand and a story about the code afterwards.",
		Source:    "The Prompt Report: Zero-Shot, Plan-and-Solve",
		Prompts: []Prompt{
			{Label: "Think", Thinking: true, Text: `A page sends a payment over a network and gets a reference back. The other side already spots repeats by that reference.

Before any code, a table. One row per way the call can end: answer lost after the money moved, connection refused, call hangs, 4xx, 5xx. Columns: can the caller tell whether money moved, is retrying safe, is retrying useful, what to do.

Then: how long to wait before giving up, and how you would find that number rather than guess it. How many retries, how the wait grows, and why everyone should not wait the same.

Done when I have the table and a retry rule that never retries what is unsafe or useless.`},
			{Label: "Build", Text: `The Gateway waits for the Teller forever, and the page resends when it hears nothing.

Carry out your plan. A time limit on every call to another program. On no answer, retry with the same reference, waiting longer each time plus a random bit, a few times only. Never retry a 4xx. After the last try, say 'outcome unknown, quote this reference'.

Prove it: make the Teller sleep past the limit on the first try only, send one payment, and show me one Ledger pair and several attempts in the log.

Done when a slow first try pays once, a 400 is never retried, and the gaps in the log grow.`},
			{Label: "Portal", Portal: true, Text: `Send waits forever, and a customer who taps again cannot tell whether the first tap paid.

Give Send a deadline. Past it, the page says the outcome is unknown, shows the reference, and offers to check rather than send again.

Done when a slow server leaves the customer with a reference and a way to find out, never a second payment.`},
		},
	},
}
