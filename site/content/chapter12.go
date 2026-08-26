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

	Concepts: []ConceptItem{
		{Term: "Message Queue", Description: "A buffer that holds messages until a worker is ready to process them."},
		{Term: "Producer", Description: "The part of the system that adds messages to the queue: peyva, after a transfer completes."},
		{Term: "Consumer / Worker", Description: "The part that reads messages off the queue and does the actual work, at its own pace."},
		{Term: "Decoupling", Description: "Producer and consumer don't call each other directly or wait on each other. The queue sits between them."},
		{Term: "Courier", Description: "The component that carries out work after a payment clears, reading from the queue at its own pace."},
	},

	BuildIt: BuildIt{
		Technique: "Allow it to say it can't",
		Why:       "An assistant would rather produce something than tell you your constraints do not fit. Permission to stop changes that.",
		Source:    "Anthropic: Reduce hallucinations, Allow Claude to say I don't know",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Teller notifies the recipient inline, so the caller waits for that notification before getting their response.

Build the Courier. The component that carries out work after a payment has cleared. The Teller hands it the notification and returns immediately; the Courier delivers on its own schedule.

Scope: an in-process queue and a worker, inside each copy. Do not introduce Kafka, RabbitMQ, NATS, Redis, Docker, or any broker or queue library. No retry policies, no dead-letter handling, no backpressure tuning.

Several copies run at once now. Say what an in-process queue costs when the copy holding it dies with work still in it, and do not fix it here.

If the standard library genuinely cannot express this, say so rather than reaching for a dependency I ruled out.

Done when a deliberately slow notification doesn't delay the payment response, and work handed over while the Courier is stopped is delivered once it starts again.`},
			{Label: "Portal", Portal: true, Text: `The page currently waits for the notification before it responds. Have it show the payment as done the moment the money has moved, and the message as delivered separately once the Courier has sent it.

If plain HTML and CSS cannot show something arriving after the page has loaded without me adding a dependency, say so and tell me what the smallest honest option is.

Done when a slow notification does not delay what the customer sees.`},
		},
	},
}
