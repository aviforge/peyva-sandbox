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
			{Label: "Decompose", Thinking: true, Text: `A payments system keeps every account in one store, and I want to split the accounts across two stores so that writes to different accounts can go to different stores.

Before designing anything, break that into the separate sub-problems it contains. For each, say what it takes as input, what it must guarantee, and which of the earlier mechanisms in this system it reuses: transactions, idempotency, the saga record. Then say which sub-problem you would build first and why, and which one you would not build in this chapter at all.

Done when I have a list of sub-problems with their guarantees, an order, and a reason for the order.`},
			{Label: "Route", Text: `The Vault runs as one process holding every account, and the copies find it by PEYVA_VAULT.

Solve the first sub-problem, routing. Run two Vaults, each owning the accounts whose handle hashes to it, each with its own file and Ledger. The copies compute the shard from the handle and send each request to the owning Vault. Extend the runner to start both, and say what you changed in it. Opening an account and enquiring a balance go to one shard. A payment whose two accounts share a shard runs as it always did, one transaction. A payment whose accounts differ is refused for now, with an error that says so.

Done when accounts spread across both shards, a same-shard payment works, and a cross-shard payment is refused with a clear message.`},
			{Label: "Cross-shard", Text: `Two Vault shards each own some accounts, and a payment between shards is refused.

Solve the second sub-problem. Run a cross-shard payment as a saga: debit on the payer's shard with the payment reference, recording it in flight; credit on the recipient's shard with the same reference; mark complete. A permanent failure of the credit reverses the debit as a new Ledger entry pair. A timeout on the credit retries with the same reference. A copy dying between debit and credit resumes from the record after restart. Everything is idempotent by reference at both shards.

Done when a cross-shard payment leaves balanced Ledger entries on both shards, a closed recipient account puts the money back, a kill between debit and credit completes on restart, and no sequence of retries credits twice.`},
			{Label: "In flight", Text: `Two Vault shards own the accounts, and cross-shard payments run as sagas with an in-flight record.

Solve the third sub-problem. Expose, from each shard, the sum of its balances and the sum of its Ledger entries, and from the saga record every payment currently in flight with its amount, reference and the step it reached. Report any payment in flight longer than a threshold as stuck.

Then show me, with ten cross-shard payments running and one deliberately stalled, that the balances on both shards plus the money in flight add up to what was seeded plus what was opened.

Done when a healthy two-shard system adds up, a stalled payment is reported as stuck with its reference and step, and the total still adds up while it is stuck.`},
			{Label: "Replicate", Thinking: true, Text: `Two Vault shards own the accounts, each a single process with its own file, while the earlier single Vault had a replica following its log and a lease saying which copy may write.

Answer from the code: what it would take to give each shard its own replica and lease, what the copies' routing would have to do when one shard is promoted, and which of those two changes you would make first. Do not build it.

Done when I have your answer on replicating a shard and the order you would do it in.`},
		},
	},
}
