package content

var Chapter16 = ChapterContent{
	Number:     16,
	Slug:       "chapter-16",
	Title:      "When Things Fail: CAP / Consistency",
	Subtitle:   "During a network failure, you can keep only two of these three: Consistency, Availability, Partition Tolerance.",
	Category:   "Distributed Systems",
	Difficulty: "Advanced",
	QuickTip:   "For money, choose consistency. Refuse the payment rather than risk moving it wrongly.",

	HeroImage:   "images/chapter-16.webp",
	HeroCaption: "Different situations call for different trade-offs. peyva has to choose.",

	Why: []string{
		"When part of the system cannot reach another part, you either refuse some requests or answer some of them wrongly.",
		"Consistency means every read sees the newest value. Availability means every request gets an answer rather than an error.",
		"Choose per action. A balance can be a little old if the page says so. A payment cannot.",
		"When everything can talk, the trade is speed against freshness. For money, wait for the second copy to confirm.",
		"A lease is permission to be the primary for a fixed time, given by someone else. Lose contact and the clock takes it away.",
		"The Warden is one program, so it can fail. Real systems use a group that votes, like Raft or Paxos. This book stops at one.",
	},

	Concepts: []ConceptItem{
		{Term: "Consistency (C)", Description: "Every read gives the newest saved value, as though there were only one copy."},
		{Term: "Availability (A)", Description: "Every request to a working copy gets an answer rather than an error, even if the answer is old."},
		{Term: "Partition (P)", Description: "Part of the system cannot reach another part. Not a choice: it happens, and the design must say what it does meanwhile."},
		{Term: "CAP Theorem", Description: "While parts cannot reach each other, you keep either consistency or availability, not both. When they can, the trade is speed against freshness instead."},
		{Term: "Lease", Description: "Permission to be the primary, for a fixed time, given by someone else and renewed before it runs out."},
		{Term: "Fencing", Description: "Making sure a copy that lost its lease cannot write, even if it thinks it still holds one. Each lease has a number that only goes up, and an old number is refused."},
		{Term: "Warden", Description: "The component that hands out the lease. It says which Vault may write, for how long, and never names two at once."},
	},

	BuildIt: BuildIt{
		Technique: "Tree-of-Thought",
		What:      "Asking for several possible approaches, then having the assistant compare them and drop the weak ones.",
		Why:       "Three strategies scored in the open leave you a choice you can disagree with.",
		Source:    "The Prompt Report: Decomposition, Tree-of-Thought",
		Prompts: []Prompt{
			{Label: "Think", Thinking: true, Text: `Balances live in a primary that copies to a follower a moment later. Promotion is by hand. I need to decide what happens when the primary, the follower and the request handlers cannot all reach each other.

No code. At least three genuinely different strategies. For each: what a customer sees while it lasts, what happens to a payment taken during it, whether two stores could both take writes, and what a person must do afterwards to put things right.

Then cut. Say what rules each rejected one out for money. Recommend one, and name the condition that would change your mind.

Done when I have three real options, a reason each was rejected, and a recommendation.`},
			{Label: "Build", Text: `The Vault's primary copies to a follower a moment later, promotion is by hand, and you recommended what should happen when they cannot reach each other.

Build the Warden: a process on PEYVA_PORT that grants a lease to one Vault at a time, for a few seconds, renewed at half that. Each lease carries a number that only goes up. Vaults read PEYVA_WARDEN. Fill in START_WARDEN in the runner.

A Vault writes only while its lease is good, and stamps every write with the lease number. Promotion is now the follower getting the lease. The copies ask the Warden which Vault to use.

Then do what you recommended for writes. If that means waiting for the follower before answering, do that, and refuse payments while it is unreachable.

Done when stopping the Warden stops all writes within one lease, stopping the primary makes the follower take over with nothing lost, and the old primary comes back as a follower.`},
			{Label: "Try", Reader: true, Text: `Take the Warden away. With everything running, run this: it kills the Warden on 9302, waits ten seconds for every lease to run out, then tries to pay. Then it stops and starts the runner to put the Warden back.

You should see: the payment refused, with a reason that names the lease. Nobody is allowed to write while nobody can say who may. That is the rule you chose, and this is what it costs.`,
				Commands: CommandsSplit(
					`Get-NetTCPConnection -LocalPort 9302 -State Listen | ForEach-Object { Stop-Process -Id $_.OwningProcess -Force }
Start-Sleep -Seconds 10
curl.exe -s -m 30 -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}' -w ' -> %{http_code}\n'
.\peyva\run.ps1 stop
.\peyva\run.ps1 start 3`,
					`for /f "tokens=5" %p in ('netstat -ano ^| findstr ":9302 " ^| findstr LISTENING') do taskkill /PID %p /F
timeout /t 10 /nobreak >nul
curl.exe -s -m 30 -X POST http://127.0.0.1:9310/pay -H "Content-Type: application/json" -d "{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}" -w " -> %{http_code}\n"
peyva\run.bat stop
peyva\run.bat start 3`,
					`kill $(lsof -ti tcp:9302)
sleep 10
curl -s -m 30 -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w ' -> %{http_code}\n'
./peyva/run.sh stop
./peyva/run.sh start 3`,
					`fuser -k 9302/tcp
sleep 10
curl -s -m 30 -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w ' -> %{http_code}\n'
./peyva/run.sh stop
./peyva/run.sh start 3`,
				)},
			{Label: "Try", Reader: true, Text: `Now take the primary away. This kills the Vault on 9300, waits ten seconds, pays, and reads alice from the follower on 9301. Watch the runner's terminal while it waits. Stop and start the runner afterwards.

You should see: the follower's log say it holds the lease, the payment succeed, and the follower's copy of alice showing it. Promotion was nobody's decision. The lease moved, and the writes moved with it.`,
				Commands: CommandsSplit(
					`Get-NetTCPConnection -LocalPort 9300 -State Listen | ForEach-Object { Stop-Process -Id $_.OwningProcess -Force }
Start-Sleep -Seconds 10
curl.exe -s -m 30 -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}' -w ' -> %{http_code}\n'
curl.exe -s http://127.0.0.1:9301/accounts/alice -w '\n'`,
					`for /f "tokens=5" %p in ('netstat -ano ^| findstr ":9300 " ^| findstr LISTENING') do taskkill /PID %p /F
timeout /t 10 /nobreak >nul
curl.exe -s -m 30 -X POST http://127.0.0.1:9310/pay -H "Content-Type: application/json" -d "{\"from\":\"alice\",\"to\":\"bob\",\"amount\":1}" -w " -> %{http_code}\n"
curl.exe -s http://127.0.0.1:9301/accounts/alice -w "\n"`,
					`kill $(lsof -ti tcp:9300)
sleep 10
curl -s -m 30 -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w ' -> %{http_code}\n'
curl -s http://127.0.0.1:9301/accounts/alice -w '\n'`,
					`fuser -k 9300/tcp
sleep 10
curl -s -m 30 -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":1}' -w ' -> %{http_code}\n'
curl -s http://127.0.0.1:9301/accounts/alice -w '\n'`,
				)},
			{Label: "Portal", Portal: true, Text: `When the copies disagree, the balance the Portal shows may be behind, and the server now says how far.

Three genuinely different ways for the page to handle that. For each, what a customer believes after reading it. Recommend one, say what it costs, and build it.

Done when an old balance is visibly old and nobody is misled about their money.`},
		},
	},
}
