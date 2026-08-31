package content

var Chapter10 = ChapterContent{
	Number:     10,
	Slug:       "chapter-10",
	Title:      "Growing the Team: Scale Out",
	Subtitle:   "When more customers come, we open more counters and add more staff so everyone gets served fast.",
	Category:   "System Design",
	Difficulty: "Intermediate",
	QuickTip:   "The moment one copy holds something the others lack, the same request starts getting different answers.",

	HeroImage:   "images/chapter-10.webp",
	HeroCaption: "Scale out = add more servers/instances to handle more load in parallel.",

	Why: []string{
		"Copies are only swappable if none of them holds something the others lack.",
		"A database file inside each copy is three banks with three sets of balances. So the Vault becomes a program of its own.",
		"A transaction cannot stretch across two programs. The whole payment now runs inside the Vault, and the Teller makes one call.",
		"More copies help with reading requests and waiting around. They do not make writing faster: every payment still ends at one Vault.",
		"The load balancer now holds the port everyone knows. A copy dying mid-request looks to it like a timeout.",
		"Whether any of this is needed is arithmetic. The sidebar below does it.",
	},

	Aside: &Aside{
		Title:       "How Big Is PEYVA? (Capacity Estimation)",
		HeroImage:   "images/sidebar-10.webp",
		HeroCaption: "Capacity estimation helps us choose the right technology, plan scaling, and control cost: before we build.",
		Why: []string{
			"An estimate is a few guesses multiplied out: how many users, how often, how much busier the peak is, how big a record is.",
			"100,000 users at three payments a day is 300,000 a day: about 3.5 a second, 35 at a ten times peak.",
			"At a kilobyte a payment, two years of Ledger is around 220 GB.",
			"Size for the peak, not the average. The busiest hour is the one that matters.",
			"Roughly right before you build beats exactly right afterwards. Know which guess moves the answer most.",
		},
		BuildIt: BuildIt{
			Technique: "Self-Ask",
			Why:       "Have it ask itself the questions first, and answer them. The guesses behind a capacity estimate come out as a labelled list you can argue with.",
			Source:    "The Prompt Report: Zero-Shot, Self-Ask",
			Prompts: []Prompt{
				{Label: "Think", Thinking: true, Text: `I need to size a payments system. One process today, and no idea what it needs to survive.

No numbers yet. Write out the questions the estimate depends on. Answer each with a stated guess, and label where it came from: industry norm, your guess, or arithmetic.

Done when I have your questions and a labelled answer to each.`},
				{Label: "Think", Thinking: true, Text: `You listed the questions a capacity estimate depends on, and answered each with a labelled guess.

Work out payments per second at peak, Ledger growth over two years, and network traffic at peak. Show each sum with the numbers filled in.

Then: which guess moves the answer most, what happens if it is off by double, and which one you most want me to confirm.

Done when I have a peak figure and a two-year storage figure, and know which guess to revisit first.`},
			},
		},
	},

	Concepts: []ConceptItem{
		{Term: "Instance", Description: "One running copy of the part of peyva that holds nothing of its own. Also just called a copy."},
		{Term: "Load Balancer", Description: "Sits in front of the copies and spreads requests across them. It holds the port the outside world knows."},
		{Term: "Stateless", Description: "Keeps no data of its own, so any copy can answer any request. The data lives in the Vault."},
		{Term: "Horizontal Scaling", Description: "Adding more copies to handle more load, instead of buying one bigger machine."},
		{Term: "Capacity Estimate", Description: "How many users, how often they act, how much busier the peak is, how big a record is. Multiply out, and you know what fills up first."},
	},

	BuildIt: BuildIt{
		Technique: "Step-Back Prompting",
		Why:       "Ask the general question before the specific one. What makes any service replaceable is a better place to start than code that was never going to scale.",
		Source:    "The Prompt Report: Thought Generation, Step-Back Prompting",
		Prompts: []Prompt{
			{Label: "Think", Thinking: true, Text: `I have a service handling payments, and I want several copies of it behind a router.

Without looking at code: what lets any service run as interchangeable copies? What may live inside one process, what may not, and why. What does that mean for a database file inside the process, and for a transaction that spans the handler and the database?

Done when I have the principle in general terms, with nothing about my project in it.`},
			{Label: "Build", Text: `The Gateway, the Teller, the Vault's file and the Ledger all sit in one process. I want several copies of the Gateway and Teller behind a router.

Show me every place that breaks the principle you just gave, the database file included. Then fix it: the Vault becomes its own process, the only one with the database, and the whole payment happens inside it in one transaction from one call by the Teller. The copies hold nothing of their own. Add a small round-robin proxy in front.

Settings come from the environment only: PEYVA_PORT for everything, PEYVA_VAULT for the copies, PEYVA_PEERS for the proxy. A missing one means say so and exit.

Done when the runner starts the Vault, three copies and the proxy, ten payments across them leave correct balances and one Ledger, and killing a copy mid-traffic fails no request.`},
			{Label: "Try", Reader: true, Text: `Kill a copy while money is moving. Start everything with the runner, then run this: five payments through the proxy on 9310, then it kills the first copy on 9311, then five more. Then ask the runner what is alive.

You should see: ten answers, each with a reference and none an error, and the runner's status showing one process gone. The proxy stopped sending to a copy that stopped answering, and the other two took the rest.`,
				Commands: CommandsSplit(
					`1..5 | ForEach-Object { curl.exe -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}' -w ' -> %{http_code}\n' }
Get-NetTCPConnection -LocalPort 9311 -State Listen | ForEach-Object { Stop-Process -Id $_.OwningProcess -Force }
1..5 | ForEach-Object { curl.exe -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}' -w ' -> %{http_code}\n' }
.\peyva\run.ps1 status`,
					`for /l %i in (1,1,5) do @curl.exe -s -X POST http://127.0.0.1:9310/pay -H "Content-Type: application/json" -d "{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}" -w " -> %{http_code}\n"
for /f "tokens=5" %p in ('netstat -ano ^| findstr ":9311 " ^| findstr LISTENING') do taskkill /PID %p /F
for /l %i in (1,1,5) do @curl.exe -s -X POST http://127.0.0.1:9310/pay -H "Content-Type: application/json" -d "{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}" -w " -> %{http_code}\n"
peyva\run.bat status`,
					`for i in 1 2 3 4 5; do curl -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w ' -> %{http_code}\n'; done
kill $(lsof -ti tcp:9311)
for i in 1 2 3 4 5; do curl -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w ' -> %{http_code}\n'; done
./peyva/run.sh status`,
					`for i in 1 2 3 4 5; do curl -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w ' -> %{http_code}\n'; done
fuser -k 9311/tcp
for i in 1 2 3 4 5; do curl -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w ' -> %{http_code}\n'; done
./peyva/run.sh status`,
				)},
		},
	},
}
