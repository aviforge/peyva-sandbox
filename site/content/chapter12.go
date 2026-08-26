package content

var Chapter12 = ChapterContent{
	Number:     12,
	Slug:       "chapter-12",
	Title:      "Decoupling with Messages: Queues",
	Subtitle:   "The order queue helps the shop accept many orders quickly, even if the baker is busy.",
	Category:   "System Design",
	Difficulty: "Intermediate",
	EstTime:    "20 min",
	QuickTip:   "Only defer work the caller doesn't actually need to wait for.",

	HeroImage:   "images/chapter-12.webp",
	HeroCaption: "Queues let different parts of the system work at their own speed and not depend on each other directly. This is called decoupling.",

	Intuition: []string{
		"Not everything a transfer triggers needs to happen before Alice sees 'success'. Bob's notification can happen a moment later.",
		"A queue is an order-ticket rail: the handler drops a message and moves on.",
		"A worker picks it up when ready.",
	},

	Concepts: []ConceptItem{
		{Term: "Message Queue", Description: "A buffer that holds messages until a worker is ready to process them."},
		{Term: "Producer", Description: "The part of the system that adds messages to the queue: peyva, after a transfer completes."},
		{Term: "Consumer / Worker", Description: "The part that reads messages off the queue and does the actual work, at its own pace."},
		{Term: "Decoupling", Description: "Producer and consumer don't call each other directly or wait on each other. The queue sits between them."},
		{Term: "Courier", Description: "The component that carries out work after a payment clears, reading from the queue at its own pace."},
	},

	UnderTheHood: []string{
		"The Teller -> 1. Send message -> Queue -> 2. Consume message -> the Courier -> 3. Perform work (e.g. send notification).",
		"Messages stay in the queue until a worker picks them up. Nothing is lost if the worker is briefly offline.",
	},

	BuildIt: BuildIt{
		Intro:     "Build the Courier. The component that carries out work after a payment clears.",
		Technique: "Allow it to say it can't",
		Why:       "An assistant would rather produce something than report that your constraints don't fit, so it quietly reaches for the dependency you ruled out. Giving explicit permission to stop and say so turns a silent workaround into a question you can answer.",
		Source:    "Anthropic: Reduce hallucinations, Allow Claude to say I don't know",
		Prompt: "The Teller notifies the recipient inline, so the caller waits for that notification before getting their response.\n\n" +
			"Build the Courier. The component that carries out work after a payment has cleared. The Teller hands it the notification and returns immediately; the Courier delivers on its own schedule.\n\n" +
			"Scope: this is a single-process learning project on a laptop, one user, no deployment. Use only the Go standard library: a buffered channel and a goroutine. Do not introduce Kafka, RabbitMQ, NATS, Redis, Docker, or any broker or queue library. No retry policies, no dead-letter handling, no backpressure tuning.\n\n" +
			"If the standard library genuinely can't express this, stop and tell me so. I would rather hear that than receive a dependency I ruled out, or a workaround that hides the problem. Saying you can't do it within these constraints is an acceptable answer here, and a more useful one than a guess.\n\n" +
			"Done when a deliberately slow notification doesn't delay the payment response, and work handed over while the Courier is stopped is delivered once it starts again.",
	},

	BreakIt: BreakIt{
		Intro: "Prove the API no longer waits on the slow part.",
		Exercises: []string{
			"Make the notification step artificially slow (sleep 5 seconds) and confirm /transfer still responds to Alice instantly.",
			"Stop the worker entirely, send several transfers, then start the worker again: every queued notification still gets delivered, just late.",
			"Compare this to doing the notification inline in Chapter 4's handler, where a slow notification would have made Alice wait too.",
		},
	},
}
