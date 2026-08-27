package content

var Chapter15 = ChapterContent{
	Number:     15,
	Slug:       "chapter-15",
	Title:      "Data Copies: Replication",
	Subtitle:   "Having copies in multiple places keeps our data safe from fires, theft, or accidents. The copy follows the original a moment later.",
	Category:   "Reliability",
	Difficulty: "Advanced",
	QuickTip:   "The second copy is always a moment behind. Everything you have to decide about failure lives in that gap.",

	HeroImage:   "images/chapter-15.webp",
	HeroCaption: "Replication = keep copies of data in multiple places so we can be safe, available and fast.",

	Why: []string{
		"Copying the file is not replication. A copy has no place-marker, so you cannot say how far behind it is, or carry on after a break.",
		"A replication log is every saved change, numbered, written in the same transaction as the change.",
		"The replica's state is one number: the last change it applied. How far behind is the primary's latest minus that.",
		"The primary answers before the replica has the change. A payment saved in that gap exists in one place only.",
		"Promotion keeps whatever the replica had. For money, whatever it never received must be counted and reported.",
		"After a promotion the old primary must never write again. Two writers is split-brain.",
	},

	Concepts: []ConceptItem{
		{Term: "Primary", Description: "The copy that takes writes and is the official one. Only ever one at a time, and making sure of that is the hard part."},
		{Term: "Replica", Description: "A second copy that applies the primary's changes in order, and can be promoted. Reads from it may be slightly old."},
		{Term: "Replication Log", Description: "Every saved change, numbered in order, written in the same transaction as the change itself."},
		{Term: "Replication Lag", Description: "The primary's latest number minus the replica's. A payment numbered above the replica's exists in one place only."},
		{Term: "Promotion", Description: "Making the replica the primary. Anything it had not yet received is lost, and the old primary must never write again."},
		{Term: "Split-Brain", Description: "Two copies both believing they are the primary, both taking writes. Breaks every money rule at once."},
	},

	BuildIt: BuildIt{
		Technique: "Analogical Prompting",
		Why:       "Compare its analogy against yours. Where the two differ, one of you is wrong about the design.",
		Source:    "The Prompt Report: Thought Generation, Analogical Prompting",
		Prompts: []Prompt{
			{Label: "Analogy", Thinking: true, Text: `A system keeps its records in one file on one disk. If that disk dies, every record dies with it. I want a second copy somewhere else that follows the first, change by change.

Before designing anything, give me a real-world analogy for keeping a second copy of records elsewhere, something with no computers in it. Say who writes first, who copies, how the copier knows where it got up to, how far behind the copy runs, and what happens when the original is destroyed while the copier is behind.

Then name the one part of your analogy that actually matters for this design.

Done when I have your analogy and the single part of it you say carries over.`},
			{Label: "Build", Text: `The Vault runs as one process holding one SQLite file, and a payment is committed there and nowhere else.

Give the Vault a replication log: every committed change appended with a sequence number in the same transaction as the change. Then make the Vault able to run as a replica. Started with PEYVA_PRIMARY set to the primary's port, it applies the primary's log in order, fetching everything after its own last applied sequence, and reports its position and the primary's latest sequence on an endpoint. Keep the copying asynchronous: the primary answers the caller without waiting for the replica.

Add a manual promotion: an endpoint or signal that tells the replica to stop following and start accepting writes, reporting the sequence it was at when promoted. An old primary that is told a promotion has happened refuses every write from then on.

The runner starts the replica with everything else when its START_REPLICA line is filled in. Fill it in. The replica must be able to be stopped and restarted while the primary keeps taking payments, and catch up from its position.

Done when a payment appears in both copies with the same sequence, stopping the replica during ten payments and restarting it catches it up with none missing, and you can show me a number for how far behind it was.`},
			{Label: "Critique", Thinking: true, Text: `You gave me a real-world analogy for keeping a second copy of records, then built log-based asynchronous replication for the Vault from it.

Tell me where that analogy breaks down for real databases, and whether it led you into any mistake in the code you wrote.

Then answer directly: if the primary dies with the replica three sequences behind and the replica is promoted, which invariants in goal.md are broken from the customer's point of view, and how would anyone find out?

Done when I know which parts of your analogy to stop trusting, whether any of them reached the code, and exactly what a promotion loses.`},
		},
	},
}
