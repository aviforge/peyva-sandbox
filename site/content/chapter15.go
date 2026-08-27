package content

var Chapter15 = ChapterContent{
	Number:     15,
	Slug:       "chapter-15",
	Title:      "Data Copies: Replication",
	Subtitle:   "Having copies in multiple places keeps our data safe from fires, theft, or accidents. The copy follows the original a moment later.",
	Category:   "Reliability",
	Difficulty: "Advanced",
	QuickTip:   "Async replication means the replica can be briefly behind. Everything you have to decide about a partition lives in that gap.",

	HeroImage:   "images/chapter-15.webp",
	HeroCaption: "Replication = keep copies of data in multiple places so we can be safe, available and fast.",

	Why: []string{
		"Copying the file is not replication. A file copy is a snapshot with no position, so you cannot say how far behind it is, cannot resume it after an interruption, and cannot tell whether a given payment is in it. Replication ships changes, in order, with a number on each.",
		"A replication log is the Ledger's idea applied to the whole store: every committed change appended with a sequence number, in the same transaction as the change. The replica's state is defined by one number, the last sequence it applied, and that number is what every question about lag or loss comes down to.",
		"Asynchronous means the primary answers the caller before the replica has the change. In the window between, a committed payment exists in exactly one place, and if the primary dies in that window the payment is gone from the replica's point of view. The window is measurable: primary's latest sequence minus the replica's.",
		"Promotion is the moment async replication charges you. The replica becomes the primary with whatever it has, and every committed change past its position is lost. For money, that loss must be counted and reported, not discovered later by a customer.",
		"After promotion the old primary must never accept a write again. If it comes back believing it is still primary, two stores accept payments independently, and that is split-brain, the one failure that breaks every invariant at once.",
		"Replication gives you a second copy of the data. It does not decide who is primary, when to promote, or how the copies find the new primary. Those are separate problems, and this chapter leaves promotion as a manual act so the next one can automate it honestly.",
	},

	Concepts: []ConceptItem{
		{Term: "Primary", Description: "The store that accepts writes and is the source of truth. There is one at a time, and enforcing that is harder than it sounds."},
		{Term: "Replica", Description: "A second copy of the store that applies the primary's changes in order and can be promoted. Reads may be served from it if the reader accepts they may be behind."},
		{Term: "Replication Log", Description: "Every committed change, appended with a sequence number in the same transaction. The replica applies it in order; the last applied number is its position."},
		{Term: "Replication Lag", Description: "Primary's latest sequence minus the replica's position. A committed payment with a sequence above the replica's position exists in one place only."},
		{Term: "Promotion", Description: "Turning the replica into the primary. With async replication, everything past its position at that moment is lost, and the old primary must be fenced off from writing again."},
		{Term: "Split-Brain", Description: "Two stores each believing they are the primary and both accepting writes. Breaks every money invariant at once."},
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
