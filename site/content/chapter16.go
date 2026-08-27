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
		Why:       "Several strategies scored against stated criteria puts the choice in the open, where you can disagree with it.",
		Source:    "The Prompt Report: Decomposition, Tree-of-Thought",
		Prompts: []Prompt{
			{Label: "Decide", Thinking: true, Text: `A payments system keeps balances in a primary store that replicates asynchronously, by a sequenced log, to a second copy. Promotion is manual. I need to decide what the system does when the primary, the second copy and the request handlers can't all reach each other.

Don't write code yet. Propose at least three genuinely distinct strategies for handling the partition: not three variations on one idea. For each, work out what a customer experiences during the partition, what happens to a payment accepted mid-partition, whether two stores could both accept writes, and what manual work recovery needs.

Then prune. Eliminate the ones that are unacceptable for money movement specifically, and say what disqualified each. Recommend one of the survivors, say whether it needs the primary to wait for the replica before answering, and name the condition that would change your recommendation.

Done when I have three real options, a reason each rejected one was rejected, and a recommendation with the condition that would overturn it.`},
			{Label: "Build", Text: `The Vault's primary replicates asynchronously to a replica by a sequenced log, promotion is manual, and you recommended a strategy for what happens when the two cannot reach each other.

Build the Warden: a small process on PEYVA_PORT that grants a lease to one Vault at a time. A lease has a holder, a number that only increases, and an expiry a few seconds out. A Vault asks for the lease on start and renews it at half the lease length; the Warden grants to the current holder if it renews in time, otherwise to whichever Vault asks next once the old lease has expired. Each Vault reads PEYVA_WARDEN for the Warden's port. Fill in the runner's START_WARDEN line so it starts with everything else.

A Vault accepts writes only while it holds an unexpired lease, stopping a margin before expiry, and stamps every write with its lease number. The replica keeps following the primary's log, and promotion is now the replica obtaining the lease. The copies ask the Warden which Vault holds the lease, cache the answer, and ask again when a request to it fails.

Implement your recommended strategy for writes. If it was to wait for the replica, have the primary answer a payment only after the replica acknowledges the sequence, and refuse payments when the replica is unreachable. Serve balance enquiries from whichever copy is reachable, saying which and how far behind it is.

Done when stopping the Warden makes both Vaults refuse writes within one lease length, stopping the primary makes the replica hold the lease and accept payments with no restart of anything else and no committed payment lost, restarting the old primary makes it a follower, and enquiries keep answering throughout.`},
			{Label: "Portal", Portal: true, Text: `When the Vault's copies disagree, the balance the Portal shows may be behind, and the server now says how far.

Propose three genuinely different ways for the page to handle that, not three wordings of one. For each, say what a customer believes after reading it, and what they do next. Then recommend one and say what it costs.

Build the one you recommend.

Done when a stale balance is visibly stale and a customer is not misled about their money.`},
		},
	},
}
