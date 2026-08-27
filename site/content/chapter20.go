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
		"Between taking and giving, the money is in flight. The totals only add up if you count it.",
		"Adding a shard moves accounts. Consistent hashing keeps the number small; it does not make the move free.",
		"Each shard needs its own second copy and its own lease. Two shards are two of every failure. That is why this comes last.",
	},

	Concepts: []ConceptItem{
		{Term: "Shard", Description: "One Vault holding some of the accounts, with its own database and its own Ledger for them."},
		{Term: "Shard Key", Description: "What decides which shard an account lives on, here a hash of the handle. Any copy can work it out, so no lookup table is needed."},
		{Term: "Routing", Description: "Sending each request to the shard that owns the account. A payment names two, and they may live apart."},
		{Term: "Cross-Shard Payment", Description: "Taking from one shard and giving on another. It cannot be one transaction, so it runs as a saga."},
		{Term: "In Flight", Description: "Money taken from one shard and not yet given on the other. Count it, or the totals look wrong during every such payment."},
		{Term: "Rebalancing", Description: "Moving accounts when a shard is added or removed, while payments carry on. Consistent hashing keeps the number small."},
	},

	BuildIt: BuildIt{
		Technique: "Decomposed Prompting",
		Why:       "Routing, the cross-shard saga and in-flight accounting are three problems wearing one name. Solved separately, each can be checked on its own.",
		Source:    "The Prompt Report: Decomposition, DECOMP",
		Prompts: []Prompt{
			{Label: "Decompose", Thinking: true, Text: `A payments system keeps every account in one store. I want to split them across two, so writes to different accounts go to different stores.

Before designing anything, break that into the separate problems inside it. For each: what it takes in, what it must guarantee, and which earlier mechanism it reuses. Then say which you would build first and why, and which one you would not build at all here.

Done when I have the problems with their guarantees, an order, and a reason for the order.`},
			{Label: "Route", Text: `The Vault is one process holding every account, and the copies find it by PEYVA_VAULT.

Solve the first problem, routing. Run two Vaults, each owning the accounts whose handle hashes to it, each with its own file and Ledger. The copies work out the shard from the handle and send each request to the Vault that owns it. Extend the runner to start both, and say what you changed.

Opening an account and asking a balance go to one shard. A payment whose accounts share a shard runs as before, in one transaction. A payment across shards is refused for now, with an error that says so.

Done when accounts spread across both shards, a same-shard payment works, and a cross-shard payment is refused clearly.`},
			{Label: "Cross-shard", Text: `Two Vault shards each own some accounts, and a payment between them is refused.

Solve the second problem. Run a cross-shard payment in three steps: take from the payer's shard with the payment reference, recording it as in flight; give on the recipient's shard with the same reference; mark it complete. A permanent failure of the second step puts the money back as a new pair of Ledger entries. A timeout retries with the same reference. A copy dying between the first two steps carries on from the record after a restart. Both shards ignore a reference they have already handled.

Done when a cross-shard payment leaves balanced entries on both shards, a closed recipient puts the money back, a kill part-way through finishes on restart, and no amount of retrying gives twice.`},
			{Label: "In flight", Text: `Two Vault shards own the accounts, and payments between them run in recorded steps.

Solve the third problem. From each shard, expose the total of its balances and the total of its Ledger entries. From the records, expose every payment in flight with its amount, reference and the step it reached, and flag any in flight longer than a threshold as stuck.

Then show me, with ten cross-shard payments running and one deliberately stalled, that both shards plus the money in flight add up to what was seeded plus what was opened.

Done when a healthy two-shard system adds up, a stalled payment is flagged with its reference and step, and the total still adds up while it is stuck.`},
			{Label: "Replicate", Thinking: true, Text: `Two Vault shards own the accounts, each one process with its own file, while the single Vault before them had a replica following its log and a lease saying which copy may write.

Answer from the code: what it would take to give each shard its own replica and lease, what the copies would have to do about routing when one shard is promoted, and which of those two you would do first. Do not build it.

Done when I have your answer on replicating a shard and the order you would do it in.`},
		},
	},
}
