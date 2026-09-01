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
		What:      "Telling the assistant that 'this cannot be done' is an acceptable answer.",
		Why:       "Without that permission, it would rather build anything than say your constraints do not fit.",
		Source:    "Anthropic: Reduce hallucinations, Allow Claude to say I don't know",
		SourceURL: AnthropicDocsURL + "en/test-and-evaluate/strengthen-guardrails/reduce-hallucinations#basic-hallucination-minimization-strategies",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Teller sends the recipient's notification itself, so the caller waits for it.

Build the Courier: it does the work after a payment clears. The Teller hands it the notification and returns at once.

A queue in memory and a worker, inside each copy. No Kafka, RabbitMQ, Redis, or any queue library. No retries, no tuning. GET /jobs on a copy answers how many are waiting. Started with PEYVA_NOTIFY_MS, a copy takes that many milliseconds over each notification, so the queue can be watched from outside.

Say what a queue in memory costs when its copy dies with work in it. Do not fix that here. If the standard library cannot do this, say so.

Done when a slow notification does not delay the payment answer, and work handed over while the Courier is stopped is delivered once it starts.`},
			{Label: "Try", Reader: true, Text: `Make notifications slow. Stop the runner. In its terminal, set this, then start it again. Every process the runner starts inherits it.

You should see: everything start as usual.`,
				Commands: Commands(
					`$env:PEYVA_NOTIFY_MS = '3000'`,
					`set PEYVA_NOTIFY_MS=3000`,
					`export PEYVA_NOTIFY_MS=3000`,
				)},
			{Label: "Try", Reader: true, Text: `Now pay through the first copy directly, on 9311, and watch its queue. This sends the payment and prints how long the answer took, asks the copy how many jobs are waiting, waits four seconds, and asks again.

You should see: the answer in well under a second while the notification takes three, one job waiting, then none. In the runner's terminal the copy's notification line lands about three seconds after the payment did.`,
				Commands: Commands(
					`curl.exe -s -X POST http://127.0.0.1:9311/pay -H 'Content-Type: application/json' -d '{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}' -w ' in %{time_total}s\n'
curl.exe -s http://127.0.0.1:9311/jobs -w '\n'
Start-Sleep -Seconds 4
curl.exe -s http://127.0.0.1:9311/jobs -w '\n'`,
					`curl.exe -s -X POST http://127.0.0.1:9311/pay -H "Content-Type: application/json" -d "{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}" -w " in %{time_total}s\n"
curl.exe -s http://127.0.0.1:9311/jobs -w "\n"
timeout /t 4 /nobreak >nul
curl.exe -s http://127.0.0.1:9311/jobs -w "\n"`,
					`curl -s -X POST http://127.0.0.1:9311/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w ' in %{time_total}s\n'
curl -s http://127.0.0.1:9311/jobs -w '\n'
sleep 4
curl -s http://127.0.0.1:9311/jobs -w '\n'`,
				)},
			{Label: "Portal", Portal: true, Text: `The page waits for the notification before it answers. Show the payment as done when the money moves, and the message as sent later, once the Courier has sent it.

If plain HTML and CSS cannot show something arriving after the page loads, say so and give me the smallest honest option.

Done when a slow notification does not delay what the customer sees.`},
		},
	},
}
