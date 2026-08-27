package content

var Chapter16 = ChapterContent{
	Number:     16,
	Slug:       "chapter-16",
	Title:      "When Things Fail: CAP / Consistency",
	Subtitle:   "During a network failure, you can keep only two of these three: Consistency, Availability, Partition Tolerance.",
	Category:   "Distributed Systems",
	Difficulty: "Advanced",
	QuickTip:   "For money, prefer CP. Reject a request rather than show the wrong balance.",

	HeroImage:   "images/chapter-16.webp",
	HeroCaption: "Different situations call for different trade-offs. peyva has to choose.",

	Why: []string{
		"CAP is narrower than its slogan. It says that during a network partition a system cannot both answer every request and guarantee every answer reflects the latest write. Partitions are not optional, so the real choice is what to do while one lasts: refuse some requests, or answer some of them wrongly.",
		"Consistency here means linearisability: every read returns the most recent committed write, as if there were one copy. Availability means every request to a live node gets a non-error response. A system that refuses payments during a partition has chosen consistency; a system that accepts them on both sides has chosen availability and will have to reconcile two histories afterwards.",
		"The choice is per operation, not per system. A balance enquiry can be served from a replica and labelled stale; a payment cannot, because a debit against a stale balance can send money that has already been spent. Money writes are CP. Money reads can be AP if they say so.",
		"When there is no partition, the trade-off is latency against consistency instead. Waiting for the replica to acknowledge every write costs a round trip per payment and means promotion loses nothing; not waiting is faster and promotion can lose the tail. That is the PACELC extension, and for money the answer on both sides is: wait.",
		"Who is primary has to be decided by something other than the primary. A lease is a promise, granted by a third party for a fixed time, that one node may accept writes; it renews before expiry or stops. A node cut off from the grantor loses its lease by the clock, so an isolated old primary stops writing.",
		"A lease depends on clocks agreeing roughly, which is why it carries a margin: the holder stops writing before the grantor considers it expired. A holder that keeps writing past its lease because its clock is slow is split-brain with a timestamp.",
		"The grantor in this chapter is one process, and one process is a single point of failure. Real systems make the grantor a majority of several nodes agreeing, which is what Raft and Paxos are for. This book stops at one grantor, and says so, because a correct majority protocol is a book of its own.",
	},

	Concepts: []ConceptItem{
		{Term: "Consistency (C)", Description: "Every read returns the latest committed write, as if there were a single copy. Also called linearisability."},
		{Term: "Availability (A)", Description: "Every request to a working node gets a non-error answer, even if the answer is out of date."},
		{Term: "Partition (P)", Description: "Part of the system cannot reach another part. Not a choice: it happens, and the design has to say what it does while it lasts."},
		{Term: "CAP Theorem", Description: "During a partition, a system can keep consistency or availability, not both. Outside a partition, the same design trades latency against consistency instead."},
		{Term: "Lease", Description: "Time-limited permission to be primary, granted by a third party and renewed before it expires. A holder that cannot renew stops writing when the clock runs out."},
		{Term: "Fencing", Description: "Making sure a node that has lost its lease cannot write, even if it believes it still holds one. Each lease carries a number that only goes up, and a write with an old number is refused."},
		{Term: "Warden", Description: "The component that grants the lease. It says which Vault is primary, for how long, and refuses to name two at once."},
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
