package content

var Chapter16 = ChapterContent{
	Number:     16,
	Slug:       "chapter-16",
	Title:      "When Things Fail: CAP / Consistency",
	Subtitle:   "During a network failure, you can keep only two of these three: Consistency, Availability, Partition Tolerance.",
	Category:   "Distributed Systems",
	Difficulty: "Advanced",
	EstTime:    "25 min",
	QuickTip:   "For money, prefer CP. Reject a request rather than show the wrong balance.",

	HeroImage:   "images/chapter-16.webp",
	HeroCaption: "Different situations call for different trade-offs. peyva has to choose.",

	Concepts: []ConceptItem{
		{Term: "Consistency (C)", Description: "All users see the same data at the same time."},
		{Term: "Availability (A)", Description: "Every request gets a response, no errors or timeouts."},
		{Term: "Partition Tolerance (P)", Description: "The system keeps working even if parts of it can't talk to each other."},
		{Term: "CAP Theorem", Description: "The choice only binds while the network is broken. On a healthy day peyva has all three."},
	},

	BuildIt: BuildIt{
		Intro:     "The Teller learns what to do when the Vault's two copies can't reach each other.",
		Technique: "Tree-of-Thought",
		Why:       "Several strategies scored against stated criteria puts the choice in the open, where you can disagree with it.",
		Source:    "The Prompt Report: Decomposition, Tree-of-Thought",
		Prompts: []Prompt{
			{Label: "Decide", Thinking: true, Intro: "Three strategies, pruned to one.", Text: `A payments system keeps balances in a primary copy that replicates asynchronously to a second copy. I need to decide what it does when the two can't reach each other.

Don't write code yet. Propose at least three genuinely distinct strategies for handling the partition: not three variations on one idea. For each, work out what a customer experiences during the partition, what happens to a payment accepted mid-partition, and what manual work recovery needs.

Then prune. Eliminate the ones that are unacceptable for money movement specifically, and say what disqualified each. Recommend one of the survivors, and name the condition that would change your recommendation.

Done when I have three real options, a reason each rejected one was rejected, and a recommendation with the condition that would overturn it.`},
			{Label: "Build", Intro: "Build the one you recommended.", Text: `The Vault's primary copy replicates asynchronously to a second copy, and you recommended a strategy for what the Teller does when the two cannot reach each other.

Implement that recommendation. Treat an unreachable second copy as a partition rather than ignoring it, and keep balance enquiries served from whichever copy is reachable even while payments are refused.

Done when disconnecting the second copy makes payments fail loudly while enquiries still succeed, and reconnecting restores payments with no restart.`},
			{Label: "Portal", Portal: true, Intro: "The Portal has to be honest about a balance it is not sure of.", Text: `When the Vault's copies disagree, the balance the Portal shows may be behind.

Propose three genuinely different ways for the page to handle that, not three wordings of one. For each, say what a customer believes after reading it, and what they do next. Then recommend one and say what it costs.

Build the one you recommend.

Done when a stale balance is visibly stale and a customer is not misled about their money.`},
		},
	},

	BreakIt: BreakIt{
		Intro: "Force a partition and confirm peyva picks the trade-off you intended.",
		Exercises: []string{
			"Cut the replica off with the runner and make a transfer. Confirm it succeeds anyway, because async replication never waits for the replica.",
			"Read the same balance from the cut-off replica. Confirm it is stale, and that gap is the consistency peyva gave up to stay available. Restore it and watch the gap close.",
			"Work out what it would take to reject that transfer instead. The primary would have to wait for the replica, and that wait is the price of choosing CP.",
		},
	},
}
