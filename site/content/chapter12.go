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
		"Messages wait in the queue until a worker picks them up. A worker that pauses only delays the work, but a process that dies still loses it, which is Chapter 13's problem.",
	},

	BuildIt: BuildIt{
		Intro:     "Build the Courier. The component that carries out work after a payment clears.",
		Technique: "Allow it to say it can't",
		Why:       "An assistant would rather produce something than report that your constraints don't fit, so it quietly reaches for the dependency you ruled out. Giving explicit permission to stop and say so turns a silent workaround into a question you can answer.",
		Source:    "Anthropic: Reduce hallucinations, Allow Claude to say I don't know",
		Prompt: `The Teller notifies the recipient inline, so the caller waits for that notification before getting their response.

Build the Courier. The component that carries out work after a payment has cleared. The Teller hands it the notification and returns immediately; the Courier delivers on its own schedule.

Scope: an in-process queue and a worker, inside each copy. Do not introduce Kafka, RabbitMQ, NATS, Redis, Docker, or any broker or queue library. No retry policies, no dead-letter handling, no backpressure tuning.

Chapter 10 means several copies now run at once, so say what an in-process queue costs when the copy holding it dies with work still in it. Do not fix that here. Name it.

If the standard library genuinely can't express this, say so rather than reaching for a dependency I ruled out. Admitting it is a more useful answer than a guess.

Done when a deliberately slow notification doesn't delay the payment response, and work handed over while the Courier is stopped is delivered once it starts again.`,
		UIIntro: "The portal stops waiting for work the customer does not care about.",
		UIPrompt: `The page currently waits for the notification before it responds. Have it show the payment as done the moment the money has moved, and the message as delivered separately once the Courier has sent it.

If plain HTML and CSS cannot show something arriving after the page has loaded without me adding a dependency, say so and tell me what the smallest honest option is. I would rather hear that than find a library in my project.

Done when a slow notification does not delay what the customer sees.`,
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
