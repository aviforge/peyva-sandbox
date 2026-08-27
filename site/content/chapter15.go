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
			{Label: "Analogy", Thinking: true, Text: `A system keeps its records in one file on one disk. If that disk dies, every record dies with it. I want a second copy elsewhere that follows the first, change by change.

Before designing anything, give me a real-world comparison with no computers in it. Who writes first, who copies, how the copier knows where it got to, how far behind it runs, and what happens if the original is destroyed while the copier is behind.

Then name the one part of it that actually matters here.

Done when I have your comparison and the single part you say carries over.`},
			{Label: "Build", Text: `The Vault is one process holding one file, and a payment is saved there and nowhere else.

Give it a log: every saved change, numbered, written in the same transaction as the change. Then let the Vault run as a replica instead. Started with PEYVA_PRIMARY set to the primary's port, it applies the primary's log in order from its own last number, and reports its number and the primary's latest.

Add a promotion by hand: something that tells the replica to stop following and start taking writes, reporting the number it was on. A primary told this has happened refuses every write from then on.

Fill in the runner's START_REPLICA line. The replica must survive being stopped and restarted while payments carry on.

Done when a payment reaches both copies with the same number, stopping the replica during ten payments and restarting it loses none, and you can show me how far behind it got.`},
			{Label: "Critique", Thinking: true, Text: `You gave me a real-world comparison for keeping a second copy, then built replication from it.

Where does the comparison break down for real databases, and did it lead you into a mistake in the code?

Then, directly: if the primary dies with the replica three changes behind and the replica is promoted, which rules in goal.md break from the customer's side, and how would anyone find out?

Done when I know which parts of the comparison to stop trusting, whether any reached the code, and exactly what a promotion loses.`},
		},
	},
}
