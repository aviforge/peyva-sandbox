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

	Intuition: []string{
		"Chapter 15 gave peyva two locations.",
		"What happens when the network between them goes down and Alice still wants to pay?",
		"Two branches cut off from each other can't stay both perfectly accurate and perfectly available. One has to give.",
	},

	Concepts: []ConceptItem{
		{Term: "Consistency (C)", Description: "All users see the same data at the same time."},
		{Term: "Availability (A)", Description: "Every request gets a response, no errors or timeouts."},
		{Term: "Partition Tolerance (P)", Description: "The system keeps working even if parts of it can't talk to each other."},
		{Term: "CAP Theorem", Description: "The choice only binds while the network is broken. On a healthy day peyva has all three, which is why the decision is easy to forget until it matters."},
	},

	UnderTheHood: []string{
		"CP: choose consistency over availability. peyva rejects requests during a partition to keep data correct.",
		"AP: choose availability over consistency. peyva keeps serving, even with stale data.",
		"CA: only possible with no partition, a single region on a normal day.",
		"peyva's rule of thumb: CP for money and account operations, AP for browsing and notifications. The choice depends on what each request actually needs.",
	},

	BuildIt: BuildIt{
		Intro:     "The Teller learns what to do when the Vault's two copies can't reach each other.",
		Technique: "Tree-of-Thought",
		Why:       "Branch, evaluate each branch, prune, then commit. Ask for an implementation and you get one design with its reasoning hidden; ask for several scored against stated criteria and the choice comes into the open where you can disagree with it.",
		Source:    "The Prompt Report: Decomposition, Tree-of-Thought",
		Prompt: "The Vault's primary copy replicates asynchronously to a second copy. I need to decide what the Teller does when the two can't reach each other.\n\n" +
			"Don't write code yet. Propose at least three genuinely distinct strategies for handling the partition: not three variations on one idea. For each, work out what a customer experiences during the partition, what happens to a payment accepted mid-partition, and what manual work recovery needs.\n\n" +
			"Then prune. Eliminate the ones that are unacceptable for money movement specifically, and say what disqualified each. Recommend one of the survivors, and name the condition that would change your recommendation.\n\n" +
			"Then implement your recommendation. Treat an unreachable second copy as a partition rather than ignoring it, and keep balance enquiries served from whichever copy is reachable even while payments are refused.\n\n" +
			"Done when disconnecting the second copy makes payments fail loudly while enquiries still succeed, and reconnecting restores payments with no restart.",
	},

	BreakIt: BreakIt{
		Intro: "Force a partition and confirm peyva picks the trade-off you intended.",
		Exercises: []string{
			"Disconnect the replica from Chapter 15 and make a transfer. Confirm it succeeds anyway, because async replication never waits for the replica.",
			"Read the same balance from the disconnected replica. Confirm it is stale, and that gap is the consistency peyva gave up to stay available.",
			"Work out what it would take to reject that transfer instead. The primary would have to wait for the replica, and that wait is the price of choosing CP.",
		},
	},
}
