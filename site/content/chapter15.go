package content

var Chapter15 = ChapterContent{
	Number:     15,
	Slug:       "chapter-15",
	Title:      "Data Copies: Replication",
	Subtitle:   "Having copies in multiple places keeps our data safe from fires, theft, or accidents. The copy follows the original a moment later.",
	Category:   "Reliability",
	Difficulty: "Advanced",
	QuickTip:   "The second copy is always a moment behind. Everything you have to decide about failure lives in that gap.",

	HeroImage:   "images/chapter-15.webp",
	HeroCaption: "Replication = keep copies of data in multiple places so we can be safe, available and fast.",

	Why: []string{
		"Copying the file is not replication. A copy has no place-marker, so you cannot say how far behind it is, or carry on after a break.",
		"A replication log is every saved change, numbered, written in the same transaction as the change.",
		"The follower's state is one number: the last change it applied. How far behind is the primary's latest minus that.",
		"The primary answers before the follower has the change. A payment saved in that gap exists in one place only.",
		"Promotion keeps whatever the follower had. For money, whatever it never received must be counted and reported.",
		"After a promotion the old primary must never write again. Two writers is split-brain.",
	},

	Concepts: []ConceptItem{
		{Term: "Primary", Description: "The copy that takes writes and is the official one. Only ever one at a time, and making sure of that is the hard part."},
		{Term: "Follower", Description: "A second copy that applies the primary's changes in order, and can be promoted. Reads from it may be slightly old."},
		{Term: "Replication Log", Description: "Every saved change, numbered in order, written in the same transaction as the change itself."},
		{Term: "Replication Lag", Description: "The primary's latest number minus the follower's. A payment numbered above the follower's exists in one place only."},
		{Term: "Promotion", Description: "Making the follower the primary. Anything it had not yet received is lost, and the old primary must never write again."},
		{Term: "Split-Brain", Description: "Two copies both believing they are the primary, both taking writes. Breaks every money rule at once."},
	},

	BuildIt: BuildIt{
		Technique: "Analogical Prompting",
		What:      "Asking the assistant to explain the problem through a comparison with something familiar, then reason from that.",
		Why:       "Where its comparison and yours differ, one of you is wrong about the design.",
		Source:    "The Prompt Report: Thought Generation, Analogical Prompting",
		SourceURL: PromptReportURL,
		Prompts: []Prompt{
			{Label: "Think", Thinking: true, Text: `A system keeps its records in one file on one disk. If that disk dies, everything dies with it. I want a second copy elsewhere that follows the first, change by change.

Give me a comparison with no computers in it. Who writes first, who copies, how the copier knows where it got to, how far behind it runs, and what happens if the original is destroyed while the copier is behind.

Name the one part of it that matters here.

Done when I have your comparison and the single part that carries over.`},
			{Label: "Build", Text: `The Vault is one process with one file, and a payment is saved there and nowhere else.

Give it a log: every saved change, numbered, written in the same transaction as the change. Then let a Vault run as a follower instead. Started with PEYVA_PRIMARY, it applies the primary's log in order from its own last number. GET /replication on any Vault reports its own last number and the primary's, and GET /accounts/alice on a Vault answers from its own file.

Add a promotion by hand: tell the follower to stop following and start taking writes. A primary told this has happened refuses every write from then on.

Fill in the runner's START_REPLICA line.

Done when a payment reaches both copies with the same number, stopping the follower during ten payments and restarting it loses none, and you can show how far behind it got.`},
			{Label: "Try", Reader: true, Text: `Ask both copies the same question. With everything running, run this: it pays, then reads alice from the primary on 9300 and from the follower on 9301, then asks the follower where it has got to.

You should see: the same balance from both, and the follower reporting the same number as the primary. A moment ago it was one behind. It caught up before you could ask.`,
				Commands: Commands(
					`curl.exe -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}' -w '\n'
curl.exe -s http://127.0.0.1:9300/accounts/alice -w '\n'
curl.exe -s http://127.0.0.1:9301/accounts/alice -w '\n'
curl.exe -s http://127.0.0.1:9301/replication -w '\n'`,
					`curl.exe -s -X POST http://127.0.0.1:9310/pay -H "Content-Type: application/json" -d "{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}" -w "\n"
curl.exe -s http://127.0.0.1:9300/accounts/alice -w "\n"
curl.exe -s http://127.0.0.1:9301/accounts/alice -w "\n"
curl.exe -s http://127.0.0.1:9301/replication -w "\n"`,
					`curl -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w '\n'
curl -s http://127.0.0.1:9300/accounts/alice -w '\n'
curl -s http://127.0.0.1:9301/accounts/alice -w '\n'
curl -s http://127.0.0.1:9301/replication -w '\n'`,
				)},
			{Label: "Try", Reader: true, Text: `Now let it fall behind. This kills the follower on 9301, pays five times, stops and starts the runner, waits, and asks the follower again.

You should see: the follower's number equal to the primary's once more, and alice's balance the same on both. Five changes happened while it was dead. It read them from the log when it came back, in order, and lost none.`,
				Commands: CommandsSplit(
					`Get-NetTCPConnection -LocalPort 9301 -State Listen | ForEach-Object { Stop-Process -Id $_.OwningProcess -Force }
1..5 | ForEach-Object { curl.exe -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}' -w '\n' }
.\peyva\run.ps1 stop
.\peyva\run.ps1 start 3
Start-Sleep -Seconds 5
curl.exe -s http://127.0.0.1:9301/replication -w '\n'
curl.exe -s http://127.0.0.1:9300/accounts/alice -w '\n'
curl.exe -s http://127.0.0.1:9301/accounts/alice -w '\n'`,
					`for /f "tokens=5" %p in ('netstat -ano ^| findstr ":9301 " ^| findstr LISTENING') do taskkill /PID %p /F
for /l %i in (1,1,5) do @curl.exe -s -X POST http://127.0.0.1:9310/pay -H "Content-Type: application/json" -d "{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}" -w "\n"
peyva\run.bat stop
peyva\run.bat start 3
timeout /t 5 /nobreak >nul
curl.exe -s http://127.0.0.1:9301/replication -w "\n"
curl.exe -s http://127.0.0.1:9300/accounts/alice -w "\n"
curl.exe -s http://127.0.0.1:9301/accounts/alice -w "\n"`,
					`kill $(lsof -ti tcp:9301)
for i in 1 2 3 4 5; do curl -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w '\n'; done
./peyva/run.sh stop
./peyva/run.sh start 3
sleep 5
curl -s http://127.0.0.1:9301/replication -w '\n'
curl -s http://127.0.0.1:9300/accounts/alice -w '\n'
curl -s http://127.0.0.1:9301/accounts/alice -w '\n'`,
					`fuser -k 9301/tcp
for i in 1 2 3 4 5; do curl -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w '\n'; done
./peyva/run.sh stop
./peyva/run.sh start 3
sleep 5
curl -s http://127.0.0.1:9301/replication -w '\n'
curl -s http://127.0.0.1:9300/accounts/alice -w '\n'
curl -s http://127.0.0.1:9301/accounts/alice -w '\n'`,
				)},
			{Label: "Check", Thinking: true, Text: `You gave me a comparison for keeping a second copy, then built replication from it.

Where does the comparison break down for real databases, and did it lead you into a mistake in the code?

Then: if the primary dies with the follower three changes behind and the follower is promoted, which rules in goal.md break for the customer, and how would anyone find out?

Done when I know which parts of the comparison to stop trusting, and exactly what a promotion loses.`},
		},
	},
}
