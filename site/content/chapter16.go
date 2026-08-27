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
			{Label: "Decide", Thinking: true, Text: `Balances live in a primary that copies to a replica a moment later. Promotion is by hand. I need to decide what happens when the primary, the replica and the request handlers cannot all reach each other.

No code yet. Give me at least three genuinely different strategies, not three versions of one. For each: what a customer sees while it lasts, what happens to a payment taken during it, whether two stores could both take writes, and what a person has to do afterwards to put things right.

Then cut. Say what ruled each rejected one out for money specifically. Recommend one, say whether it needs the primary to wait for the replica, and name the one condition that would change your mind.

Done when I have three real options, a reason each was rejected, and a recommendation with the condition that would overturn it.`},
			{Label: "Build", Text: `The Vault's primary copies to a replica a moment later, promotion is by hand, and you recommended what should happen when they cannot reach each other.

Build the Warden: a process on PEYVA_PORT granting a lease to one Vault at a time. A lease has a holder, a number that only goes up, and an expiry a few seconds out. A Vault asks on start and renews at half that. The Warden renews for the holder, or hands the lease on once the old one has expired. Vaults read PEYVA_WARDEN. Fill in the runner's START_WARDEN line.

A Vault writes only while its lease is good, stopping a margin before expiry, and stamps every write with the lease number. Promotion is now the replica getting the lease. The copies ask the Warden which Vault to use, and ask again when a request fails.

Then do what you recommended for writes. If that means waiting for the replica before answering, do that, and refuse payments while it is unreachable. Answer enquiries from whichever copy is reachable, saying which and how far behind.

Done when stopping the Warden makes both Vaults refuse writes within one lease, stopping the primary makes the replica take over with nothing lost, restarting the old primary makes it a follower, and enquiries answer throughout.`},
			{Label: "Portal", Portal: true, Text: `When the copies disagree, the balance the Portal shows may be behind, and the server now says how far.

Give me three genuinely different ways for the page to handle that, not three wordings of one. For each, say what a customer believes after reading it and what they do next. Recommend one and say what it costs.

Build the one you recommend.

Done when an old balance is visibly old and nobody is misled about their money.`},
		},
	},
}
