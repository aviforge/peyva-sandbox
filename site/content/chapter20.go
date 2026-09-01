package content

var Chapter20 = ChapterContent{
	Number:     20,
	Slug:       "chapter-20",
	Title:      "Splitting the Vault: Sharding",
	Subtitle:   "One store can only grow so far. Split the accounts across several, and everything that relied on one database has to be rebuilt.",
	Category:   "Distributed Systems",
	Difficulty: "Advanced",
	QuickTip:   "A payment inside one shard is a transaction. A payment across two is a saga. Design so most are the first kind.",

	HeroImage:   "images/chapter-20.webp",
	HeroCaption: "Sharding = each store holds some of the accounts. Routing decides which, and cross-shard payments pay for the split.",

	Why: []string{
		"More copies never made writing faster. Splitting the accounts does: writes to different accounts go to different stores.",
		"The shard key is a hash of the handle, so any copy works it out. It also puts alice and bob apart half the time.",
		"A payment across two shards cannot be one transaction. It is the saga from chapter 14, now on the everyday path.",
		"Between taking and giving, the money is in progress. The totals only add up if you count it.",
		"Adding a shard moves accounts between stores. A careful choice of hash keeps that number small; it does not make the move free.",
		"Each shard needs its own second copy and its own lease. Two shards are two of every failure. That is why this comes last.",
	},

	Concepts: []ConceptItem{
		{Term: "Shard", Description: "One Vault holding some of the accounts, with its own database and its own Ledger for them."},
		{Term: "Shard Key", Description: "What decides which shard an account lives on, here a hash of the handle. Any copy can work it out, so no lookup table is needed."},
		{Term: "Routing", Description: "Sending each request to the shard that owns the account. A payment names two, and they may live apart."},
		{Term: "Cross-Shard Payment", Description: "Taking from one shard and giving on another. It cannot be one transaction, so it runs as a saga."},
		{Term: "In Progress", Description: "Money taken from one shard and not yet given on the other. Count it, or the totals look wrong during every such payment."},
		{Term: "Rebalancing", Description: "Moving accounts when a shard is added or removed, while payments carry on. A careful choice of hash keeps the number that move small."},
	},

	BuildIt: BuildIt{
		Technique: "Decomposed Prompting",
		What:      "Splitting one request into the separate problems inside it and solving each one on its own.",
		Why:       "Routing, the cross-shard payment and the accounting are three jobs wearing one name.",
		Source:    "The Prompt Report: Decomposition, DECOMP",
		SourceURL: PromptReportURL,
		Prompts: []Prompt{
			{Label: "Think", Thinking: true, Text: `A payments system keeps every account in one store. I want to split them across two, so writes to different accounts go to different stores.

Break that into the separate problems inside it. For each: what it must guarantee, and which earlier mechanism it reuses. Say which to build first and why, and which not to build here at all.

Done when I have the problems, an order, and a reason for the order.`},
			{Label: "Build", Text: `The Vault is one process holding every account, and the copies find it by PEYVA_VAULT.

Run two Vaults, each owning the accounts whose handle hashes to it, each with its own file and Ledger. The copies work out which from the handle. The second shard listens on 9303. Extend the runner to start both. GET /accounts/carol now includes "shard": the port of the Vault that owns her.

A payment within one shard runs as before. A payment across shards is refused for now, with an error that says so.

Done when accounts spread across both shards, a same-shard payment works, and a cross-shard payment is refused clearly.`},
			{Label: "Build", Text: `Two Vault shards each own some accounts, and a payment between them is refused.

Run it in three steps: take from the payer's shard, recording the payment as in progress; give on the recipient's shard with the same reference; mark it complete. If the second step fails for good, put the money back as new Ledger entries. If it times out, retry with the same reference. If the copy dies part-way, carry on from the record after restart. A shard that has already handled a reference treats it as done.

Done when a cross-shard payment leaves balanced entries on both shards, a closed recipient puts the money back, a kill part-way through finishes on restart, and no amount of retrying gives twice.`},
			{Label: "Try", Reader: true, Text: `Send money across the split. With both shards running, run this: it opens carol, dave and erin, reads all four accounts to see which shard owns each, pays one from alice to each of the three, and prints alice's history and then each of theirs.

You should see: the four accounts spread across 9300 and 9303. If all three landed with alice, open a few more until one does not. For a recipient on the other shard, alice's history holds the debit and theirs holds the credit, with one reference across two files. For one on her own shard, both entries sit in one file as before.`,
				Commands: Commands(
					`curl.exe -s -X POST http://127.0.0.1:9310/accounts -H 'Content-Type: application/json' -d '{\"handle\":\"carol\"}' -w '\n'
curl.exe -s -X POST http://127.0.0.1:9310/accounts -H 'Content-Type: application/json' -d '{\"handle\":\"dave\"}' -w '\n'
curl.exe -s -X POST http://127.0.0.1:9310/accounts -H 'Content-Type: application/json' -d '{\"handle\":\"erin\"}' -w '\n'
'alice', 'carol', 'dave', 'erin' | ForEach-Object { curl.exe -s http://127.0.0.1:9310/accounts/$_ -w '\n' }
'carol', 'dave', 'erin' | ForEach-Object { curl.exe -s -m 30 -X POST http://127.0.0.1:9310/pay -H "Authorization: Bearer $env:PEYVA_TOKEN_ALICE" -H 'Content-Type: application/json' -d ('{\"from\":\"alice\",\"to\":\"' + $_ + '\",\"amount\":1}') -w ' -> %{http_code}\n' }
'alice', 'carol', 'dave', 'erin' | ForEach-Object { curl.exe -s http://127.0.0.1:9310/accounts/$_/history -w '\n' }`,
					`curl.exe -s -X POST http://127.0.0.1:9310/accounts -H "Content-Type: application/json" -d "{\"handle\":\"carol\"}" -w "\n"
curl.exe -s -X POST http://127.0.0.1:9310/accounts -H "Content-Type: application/json" -d "{\"handle\":\"dave\"}" -w "\n"
curl.exe -s -X POST http://127.0.0.1:9310/accounts -H "Content-Type: application/json" -d "{\"handle\":\"erin\"}" -w "\n"
for %h in (alice carol dave erin) do @curl.exe -s http://127.0.0.1:9310/accounts/%h -w "\n"
for %h in (carol dave erin) do @curl.exe -s -m 30 -X POST http://127.0.0.1:9310/pay -H "Authorization: Bearer %PEYVA_TOKEN_ALICE%" -H "Content-Type: application/json" -d "{\"from\":\"alice\",\"to\":\"%h\",\"amount\":1}" -w " -> %{http_code}\n"
for %h in (alice carol dave erin) do @curl.exe -s http://127.0.0.1:9310/accounts/%h/history -w "\n"`,
					`curl -s -X POST http://127.0.0.1:9310/accounts -H 'Content-Type: application/json' -d '{"handle":"carol"}' -w '\n'
curl -s -X POST http://127.0.0.1:9310/accounts -H 'Content-Type: application/json' -d '{"handle":"dave"}' -w '\n'
curl -s -X POST http://127.0.0.1:9310/accounts -H 'Content-Type: application/json' -d '{"handle":"erin"}' -w '\n'
for h in alice carol dave erin; do curl -s http://127.0.0.1:9310/accounts/$h -w '\n'; done
for h in carol dave erin; do curl -s -m 30 -X POST http://127.0.0.1:9310/pay -H "Authorization: Bearer $PEYVA_TOKEN_ALICE" -H 'Content-Type: application/json' -d "{\"from\":\"alice\",\"to\":\"$h\",\"amount\":1}" -w ' -> %{http_code}\n'; done
for h in alice carol dave erin; do curl -s http://127.0.0.1:9310/accounts/$h/history -w '\n'; done`,
				)},
			{Label: "Build", Text: `Two Vault shards own the accounts, and payments between them run in recorded steps.

From each shard, expose the total of its balances and of its Ledger entries. From the records, expose every payment still in progress with its amount, reference and step, and flag any that has been in progress too long.

Then run ten cross-shard payments with one deliberately stalled, and show that both shards plus the money in progress add up to what was seeded plus what was opened.

Done when a healthy system adds up, a stalled payment is flagged, and the total still adds up while it is stalled.`},
			{Label: "Think", Thinking: true, Text: `Two Vault shards own the accounts, each one process with one file. The single Vault before them had a follower and a lease.

From the code: what would it take to give each shard its own follower and lease, and what would the copies have to do when one shard is promoted? Which would you do first? Do not build it.

Done when I have your answer and the order you would do it in.`},
		},
	},
}
