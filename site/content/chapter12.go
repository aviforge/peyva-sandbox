package content

var Chapter12 = ChapterContent{
	Number:     12,
	Slug:       "chapter-12",
	Title:      "Decoupling with Messages: Queues",
	Subtitle:   "The order queue helps the shop accept many orders quickly, even if the baker is busy.",
	Category:   "System Design",
	Difficulty: "Intermediate",
	QuickTip:   "Only defer work the caller doesn't actually need to wait for.",

	HeroImage:   "images/chapter-12.webp",
	HeroCaption: "Queues let different parts of the system work at their own speed and not depend on each other directly. This is called decoupling.",

	Why: []string{
		"A reply should wait for what the caller needs and nothing more. They need the money moved, not the message sent.",
		"A queue separates deciding that work must happen from actually doing it.",
		"A queue held in memory dies with the program, and nothing records what was owed. This chapter builds that flaw on purpose.",
		"The number of jobs waiting is the most honest health signal you have.",
		"The page shows the payment as done when the money moved, and the message as sent separately, later.",
	},

	Concepts: []ConceptItem{
		{Term: "Message Queue", Description: "A holding line for jobs, until a worker is ready to take one."},
		{Term: "Producer", Description: "The part that puts jobs on the queue: peyva, once a payment is done."},
		{Term: "Consumer / Worker", Description: "The part that takes jobs off the queue and does them, at its own pace."},
		{Term: "Decoupling", Description: "Producer and worker never call or wait on each other. The queue sits between them."},
		{Term: "Backlog", Description: "How many jobs are waiting. Growing means the worker is falling behind, and it is the first number to watch."},
		{Term: "Courier", Description: "The component that does the work following a payment, taking it from the queue at its own pace."},
	},

	BuildIt: BuildIt{
		Technique: "Allow it to say it can't",
		Why:       "An assistant would rather produce something than tell you your constraints do not fit. Permission to stop changes that.",
		Source:    "Anthropic: Reduce hallucinations, Allow Claude to say I don't know",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Teller sends the recipient's notification itself, so the caller waits for it.

Build the Courier: it does the work after a payment clears. The Teller hands it the notification and returns at once.

A queue in memory and a worker, inside each copy. No Kafka, RabbitMQ, Redis, or any queue library. No retries, no tuning. Show how many jobs are waiting.

Say what a queue in memory costs when its copy dies with work in it. Do not fix that here. If the standard library cannot do this, say so.

Done when a slow notification does not delay the payment answer, and work handed over while the Courier is stopped is delivered once it starts.`},
			{Label: "Portal", Portal: true, Text: `The page waits for the notification before it answers. Show the payment as done when the money moves, and the message as sent later, once the Courier has sent it.

If plain HTML and CSS cannot show something arriving after the page loads, say so and give me the smallest honest option.

Done when a slow notification does not delay what the customer sees.`},
		},
	},
}
