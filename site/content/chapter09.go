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
		Why:       "The plan is a table of failure and response. Written first, it is the code's specification; written after, it is a story about the code.",
		Source:    "The Prompt Report: Zero-Shot, Plan-and-Solve",
		Prompts: []Prompt{
			{Label: "Plan", Thinking: true, Text: `A page sends a payment over a network and gets a reference back. The request already carries a reference the other side uses to spot repeats.

Before any code, give me a table. One row for each way that call can end: the answer lost after the money moved, the connection refused, the call hanging, a 4xx, a 5xx. Columns: can the caller tell whether money moved, is retrying safe, is retrying useful, what should the caller do.

Then say how long the time limit should be, and how you would find that number from a running system rather than guessing. Say how many retries, how the wait grows, and why every caller should not wait the same.

Done when I have the table, a time limit with a method behind it, and a retry rule that never retries what is unsafe or useless.`},
			{Label: "Build", Text: `The Gateway waits for the Teller with no deadline, and the page resends when it hears nothing.

Carry out your plan. Put a time limit on every call to another program. When the outcome is unknown, retry with the same reference, waiting longer each time plus a random amount, a few times only. Never retry a 4xx. After the last try, answer 'outcome unknown, quote this reference' rather than reporting failure.

Then prove it. Make the Teller sleep past the limit on the first attempt only, send one payment, and show me one pair of Ledger entries and more than one attempt in the log.

Done when a slow first attempt pays once, a 400 is never retried, and the log shows the gaps growing.`},
			{Label: "Portal", Portal: true, Text: `Send waits with no limit, and a customer who gives up and taps again cannot tell whether the first tap paid.

Give Send a deadline. Past it, the page says the outcome is unknown, shows the reference, and offers to check rather than to send again. Checking asks the server what happened to that reference.

Done when a slow server leaves the customer with a reference and a way to find out, never a second payment.`},
		},
	},
}
