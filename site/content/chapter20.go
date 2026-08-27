package content

var Chapter20 = ChapterContent{
	Number:     20,
	Slug:       "chapter-20",
	Title:      "Putting It All Together",
	Subtitle:   "Every piece from the last twenty chapters works together to make peyva fast, reliable, and safe.",
	Category:   "System Design",
	Difficulty: "Advanced",
	QuickTip:   "No single piece makes a system reliable: only how they cooperate under load and failure.",

	HeroImage:   "images/chapter-20.webp",
	HeroCaption: "System design is not one clever box. It's how all the pieces cooperate under load, failure and change.",

	Why: []string{
		"Every mechanism in this book protects an invariant at one boundary. Transactions at the database, idempotency at the API, the outbox at the queue, the lease at the primary. A system is correct when every boundary a payment crosses has its mechanism, and wrong at the first one that does not.",
		"Reconciliation is the check that assumes everything else has a bug. The Vault's balances and the Ledger's entries are written together, so they should always agree, and the Reconciler exists for the day they do not. It is how real payments systems find out, and it runs constantly, not after an incident.",
		"A discrepancy is a fact to report, not a thing to fix automatically. The Reconciler says which account, how much, and since when; a person decides which side is right. A reconciler that silently corrects balances is a second writer of balances, which the first invariant forbids.",
		"Single points of failure are found by walking the payment's path and asking, at each hop, what happens if this dies now. The Warden is one. The proxy is one. Knowing them is more useful than removing all of them, because removing each costs something and some are worth keeping.",
		"Over-engineering is real cost with no failure it prevents. Every component here earned its place by a failure a previous chapter demonstrated; a component that cannot name its failure is the one to question.",
		"Restating a system in your own words is the cheapest test of whether you understand it. Where your description and the code disagree, one of them is wrong, and it is worth finding out which before the next change.",
	},

	Concepts: []ConceptItem{
		{Term: "Reconciler", Description: "The component that checks the Vault against the Ledger and reports where they disagree."},
		{Term: "Reconciliation", Description: "Comparing two records that were written independently, to catch the case where one of them is wrong."},
		{Term: "Discrepancy", Description: "An account whose Ledger entries do not sum to the balance the Vault reports, and the size of the gap."},
		{Term: "Single Point of Failure", Description: "A part with no stand-in, whose loss stops the system rather than slowing it."},
		{Term: "Conservation Check", Description: "The sum of every balance equals what was seeded plus what was opened, and the sum of every Ledger entry is zero. If either fails, money was created or lost."},
	},

	BuildIt: BuildIt{
		Technique: "Rephrase and Respond (RaR)",
		Why:       "The restatement is the cheapest look you will get at what it actually understood.",
		Source:    "The Prompt Report: Zero-Shot, Rephrase and Respond",
		Prompts: []Prompt{
			{Label: "Restate", Thinking: true, Text: `peyva is a payments system built from these parts: a proxy that spreads requests across copies, a Gateway in each copy that takes requests from outside, a Teller that runs one payment end to end, a Vault that is the only thing that changes a balance and holds the Ledger recording every movement, a replica following the Vault's log, a Warden that says which Vault may write, a Courier that carries out work after a payment clears, and a Portal a customer uses.

Before writing anything, restate that back to me in your own words. Describe what each part is for, and how a payment travels from the proxy to the point where the recipient has been told. Don't repeat my names for things back at me: say what each one actually does.

At each hop, say what would happen if that part failed right then, and which mechanism, if any, makes that failure safe.

Done when I have your description of every part and the path a payment takes, and I can tell you where yours and mine disagree.`},
			{Label: "Build", Text: `peyva has a Vault holding balances and a Ledger recording every movement of money, and nothing checks that the two still agree.

Build the Reconciler. For every account, the sum of its Ledger entries must equal the balance the Vault reports, and any account where they don't is reported with the size of the gap. Also check conservation: the sum of all balances equals what was seeded plus what was opened, and the sum of all Ledger entries is zero. Run it against the replica as well as the primary and report any difference between the two. It reports; it never corrects.

Then answer three things from the code you have, not from general knowledge: which single part failing would hurt customers most, which failure the system currently handles worst, and which piece is over-engineered for its actual load.

Where the code doesn't match how a real payments system would do it, say so plainly.

Done when the Reconciler reports no discrepancies on a healthy system, reports the exact gap after I kill the Vault mid-payment with a deliberately broken transaction, and I have your three answers.`},
			{Label: "Portal restate", Portal: true, Thinking: true, Text: `A customer's wallet page has grown a screen at a time: a balance, sending money to a handle, a history, a note that a message was delivered, and a sign-in in front of all of it.

Describe that page back to me: every screen, what a customer can do on each, and which part of the system answers it. Use your own words, not mine. Say which parts you are unsure of.

Done when I have your description of every screen, and I can tell you where yours and mine disagree.`},
			{Label: "Portal", Portal: true, Text: `The Portal has a screen for each thing it learned to do, added a chapter at a time and looking like it.

Finish it: one menu from which a customer opens an account, sees what they hold, sends money, reads their history and knows a message was delivered, with the switcher deciding whose wallet it all belongs to.

Done when someone who has never seen peyva can use it without being told how.`},
		},
	},
}
