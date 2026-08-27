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
			{Label: "Plan", Thinking: true, Text: `A page sends a payment to a service over a network, and the service records the payment and answers with a reference. The request already carries a reference the service uses to recognise a repeat.

Before writing any code, make a plan as a table. Rows: every distinct way that call can end, including the response being lost after the money moved, the connection being refused, the call hanging, a 4xx, and a 5xx. Columns: whether the caller can know if money moved, whether retrying is safe, whether retrying is useful, and what the caller should do.

Then say how long the timeout should be and how you would find that number from a running system rather than guessing it. Say how many retries, how the wait between them grows, and why the wait should not be exactly the same for every caller.

Done when I have the table, a timeout with a method behind it, and a retry policy that never retries a case where retrying is unsafe or useless.`},
			{Label: "Build", Text: `The Gateway forwards a payment to the Teller and waits for an answer with no deadline, and the Portal resends when it hears nothing.

Carry out the plan you wrote. Give every call across a process or network boundary a timeout. On an unknown outcome, retry with the same reference, with exponential backoff and jitter, a small fixed number of times. Never retry a 4xx. After the last attempt, answer the caller with an honest 'outcome unknown, quote this reference' rather than a failure.

Then prove it: make the Teller sleep past the timeout on the first attempt only, send one payment, and show me from the Ledger that it was applied once and from the log that it was attempted more than once.

Done when a slow first attempt results in exactly one Ledger entry pair, a 400 is never retried, and the log shows the growing gaps between attempts.`},
			{Label: "Portal", Portal: true, Text: `Send waits for an answer with no limit, and a customer who gives up and taps again has no idea whether the first tap paid.

Give Send a deadline. Past it, the page says the outcome is unknown and shows the reference, and the button offers to check rather than to send again. Checking asks the server what happened to that reference.

Done when a deliberately slow server leaves the customer with a reference and a way to find out, never with a second payment.`},
		},
	},
}
