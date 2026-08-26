package content

var Chapter15 = ChapterContent{
	Number:     15,
	Slug:       "chapter-15",
	Title:      "Data Copies: Replication",
	Subtitle:   "Having copies in multiple places keeps our data safe from fires, theft, or accidents. The copy follows the original a moment later.",
	Category:   "Reliability",
	Difficulty: "Advanced",
	EstTime:    "20 min",
	QuickTip:   "Async replication means the replica can be briefly behind. Everything you have to decide about a partition lives in that gap.",

	HeroImage:   "images/chapter-15.webp",
	HeroCaption: "Replication = keep copies of data in multiple places so we can be safe, available and fast.",

	Concepts: []ConceptItem{
		{Term: "Primary Database", Description: "The database that accepts writes and is the source of truth."},
		{Term: "Replica", Description: "A copy of the primary's data, kept in sync, usually in a different location."},
		{Term: "Async Replication", Description: "The replica catches up shortly after the primary writes: fast, but the replica can be briefly behind."},
		{Term: "Promotion", Description: "Turning a replica into the new primary when the original primary fails."},
	},

	BuildIt: BuildIt{
		Intro:     "The Vault learns to keep a spare copy in another place.",
		Technique: "Analogical Prompting",
		Why:       "Compare its analogy against yours. Where the two differ, one of you is wrong about the design.",
		Source:    "The Prompt Report: Thought Generation, Analogical Prompting",
		Prompts: []Prompt{
			{Label: "Analogy", Thinking: true, Intro: "An analogy for a second copy, before any design.", Text: `A system keeps its records in one file on one disk. If that disk dies, every record dies with it. I want a second copy somewhere else.

Before designing anything, give me a real-world analogy for keeping a second copy of records elsewhere, something with no computers in it. Say who writes first, who copies, how far behind the copy runs, and what happens when the original is destroyed.

Then name the one part of your analogy that actually matters for this design.

Done when I have your analogy and the single part of it you say carries over.`},
			{Label: "Build", Intro: "Build what the analogy described.", Text: `The Vault keeps every account in one file on one disk, and a payment is committed there and nowhere else.

Give the Vault a second copy, written to after each committed payment, and a way to promote it when the primary is unreachable. Keep the copying asynchronous. The caller must never wait for it.

The runner brings the second copy up with everything else, and can cut it off or restore it while the rest keeps running. Testing a replica by renaming its file behind the system's back proves less than it looks like it does.

Done when a payment appears in both copies, balance enquiries survive the primary being unavailable, and you can show me the window where a committed payment hasn't reached the second copy yet.`},
			{Label: "Critique", Thinking: true, Intro: "Where the analogy breaks down.", Text: `You gave me a real-world analogy for keeping a second copy of records, then built asynchronous replication for the Vault from it.

Tell me where that analogy breaks down for real databases, and whether it led you into any mistake in the code you wrote.

Done when I know which parts of your analogy to stop trusting, and whether any of them reached the code.`},
		},
	},
}
