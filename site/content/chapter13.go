package content

var Chapter13 = ChapterContent{
	Number:     13,
	Slug:       "chapter-13",
	Title:      "Reliability Patterns: Transactional Outbox",
	Subtitle:   "We don't send the email right away. First, we save the order and write the email into the outbox, in the same step.",
	Category:   "Reliability",
	Difficulty: "Advanced",
	QuickTip:   "If it isn't written in the same transaction as the data change, it can be lost.",

	HeroImage:   "images/chapter-13.webp",
	HeroCaption: "Transactional Outbox makes sure important work is never lost, even if something fails right after saving.",

	Why: []string{
		"Save the payment, then send the message: a crash in the gap between those two lines loses the message, with no record it was owed.",
		"The outbox writes the message in the same transaction as the payment. Both exist, or neither does.",
		"Crash between sending and marking it done, and it goes twice. You have traded 'maybe never' for 'maybe twice'.",
		"Deliver at least once, then spot the repeats. The same idea as chapter 8's key, now applied to the work after a payment.",
		"Three copies reading one table can grab the same job. Have each claim it in a single update, and let the database decide.",
	},

	Concepts: []ConceptItem{
		{Term: "Outbox Table", Description: "A table in the same database as the payment, holding messages still to be sent."},
		{Term: "Same Transaction", Description: "The payment and the message are written together, in the transaction from chapter 7, never as two steps."},
		{Term: "Outbox Publisher", Description: "A worker that reads unsent rows, sends them, then marks them done."},
		{Term: "At Least Once", Description: "A crash after sending and before marking done sends it twice. The receiver has to cope, which is what chapter 8 built."},
		{Term: "Claiming", Description: "Taking a row in a single update before working on it, so two workers do not both take it. Not the same as marking it done."},
	},

	BuildIt: BuildIt{
		Technique: "Contrastive Chain-of-Thought",
		What:      "Showing an example of wrong reasoning next to right reasoning, so the assistant learns the difference.",
		Why:       "Naming the design you do not want stops it being rediscovered, and forces the difference into the open.",
		Source:    "The Prompt Report: Few-Shot CoT, Contrastive CoT",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Courier takes work from a queue in memory. If its copy dies with work in it, the work is gone and nothing knows.

Reasoning to reject: save the payment, then hand the work over, because the gap is too small to matter. A crash in that gap loses work with no record it was owed.

Reasoning to follow: anything that must happen because of a payment is written in the same transaction as the payment.

Build the second. The Vault writes the Courier's pending work with the money. The Courier in each copy takes work from the Vault, claims each item in one update so two copies cannot both take it, delivers, and marks it done.

Done when killing a copy right after a payment leaves the work saved, any copy delivers it after a restart, and three copies collecting at once deliver each item once.`},
			{Label: "Try", Reader: true, Text: `Kill the copy that owes a notification. With PEYVA_NOTIFY_MS still set from the last chapter and everything running, run this: it pays through the first copy on 9311 and kills that copy in the same second, while the notification is still in flight. Then it stops and starts the runner.

You should see: within a few seconds of the start, one copy's line in the runner's terminal delivering that payment's notification, once. The work was saved with the money, so it survived the copy. If nothing delivers it, you have found the case the next turn asks about: a claim that died with its copy.`,
				Commands: CommandsSplit(
					`curl.exe -s -X POST http://127.0.0.1:9311/pay -H 'Content-Type: application/json' -d '{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}' -w '\n'
Get-NetTCPConnection -LocalPort 9311 -State Listen | ForEach-Object { Stop-Process -Id $_.OwningProcess -Force }
.\peyva\run.ps1 stop
.\peyva\run.ps1 start 3`,
					`curl.exe -s -X POST http://127.0.0.1:9311/pay -H "Content-Type: application/json" -d "{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}" -w "\n"
for /f "tokens=5" %p in ('netstat -ano ^| findstr ":9311 " ^| findstr LISTENING') do taskkill /PID %p /F
peyva\run.bat stop
peyva\run.bat start 3`,
					`curl -s -X POST http://127.0.0.1:9311/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w '\n'
kill $(lsof -ti tcp:9311)
./peyva/run.sh stop
./peyva/run.sh start 3`,
					`curl -s -X POST http://127.0.0.1:9311/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w '\n'
fuser -k 9311/tcp
./peyva/run.sh stop
./peyva/run.sh start 3`,
				)},
			{Label: "Check", Thinking: true, Text: `You had the Vault write the Courier's pending work in the same transaction as the payment, instead of handing it over afterwards.

Name the exact instant the rejected design loses work and yours does not.

Then: what if the Courier dies after delivering but before marking it done? What if a copy claims an item and dies before delivering it?

Done when I can point at the instant that separates the two, and I know the two cases yours still handles imperfectly.`},
		},
	},
}
