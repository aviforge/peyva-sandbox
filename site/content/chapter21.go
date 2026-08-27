package content

var Chapter21 = ChapterContent{
	Number:     21,
	Slug:       "chapter-21",
	Title:      "Putting It All Together",
	Subtitle:   "Every piece from the last twenty-one chapters works together to make peyva fast, reliable, and safe.",
	Category:   "System Design",
	Difficulty: "Advanced",
	QuickTip:   "No single piece makes a system reliable: only how they cooperate under load and failure.",

	HeroImage:   "images/chapter-21.webp",
	HeroCaption: "System design is not one clever box. It's how all the pieces cooperate under load, failure and change.",

	Why: []string{
		"Each piece guards one rule at one crossing: transactions, keys, the outbox, the lease. The first crossing without a guard is where money goes wrong.",
		"Reconciling assumes everything else has a bug. Balances and the Ledger should always agree; the Reconciler is for the day they do not.",
		"Report a difference, never fix it automatically. Something that quietly corrects balances has become a second thing writing them.",
		"Find the weak points by walking a payment's path and asking, at each step, what happens if this dies now.",
		"Every part here earned its place by a failure some chapter showed you. One that cannot name its failure is the one to question.",
		"Saying the system back in your own words is the cheapest test of whether you understand it.",
	},

	Concepts: []ConceptItem{
		{Term: "Reconciler", Description: "The component that checks balances against the Ledger and reports where they disagree."},
		{Term: "Reconciliation", Description: "Comparing two records that were written separately, to catch the day one of them is wrong."},
		{Term: "Discrepancy", Description: "An account whose Ledger entries do not add up to its balance, and by how much."},
		{Term: "Single Point of Failure", Description: "A part with no stand-in, whose loss stops the system rather than slowing it."},
		{Term: "Conservation Check", Description: "Every balance, plus money in flight, adds up to what you started with. Every Ledger entry adds up to zero. If either fails, money was created or lost."},
	},

	BuildIt: BuildIt{
		Technique: "Rephrase and Respond (RaR)",
		Why:       "The restatement is the cheapest look you will get at what it actually understood.",
		Source:    "The Prompt Report: Zero-Shot, Rephrase and Respond",
		Prompts: []Prompt{
			{Label: "Restate", Thinking: true, Text: `peyva is a payments system built from these parts: a proxy that spreads requests across copies, a Gateway in each copy that takes requests from outside, a Teller that runs one payment end to end, two Vault shards that are the only things that change a balance and hold the Ledger recording every movement, a replica following a Vault's log, a Warden that says which Vault may write, a Courier that carries out work after a payment clears, and a Portal a customer uses.

Before writing anything, restate that back to me in your own words. Describe what each part is for, and how a payment travels from the proxy to the point where the recipient has been told. Don't repeat my names for things back at me: say what each one actually does.

At each hop, say what would happen if that part failed right then, and which mechanism, if any, makes that failure safe.

Done when I have your description of every part and the path a payment takes, and I can tell you where yours and mine disagree.`},
			{Label: "Build", Text: `peyva has two Vault shards holding balances and Ledger entries, cross-shard payments in flight between them, and nothing checks that any of it still agrees.

Build the Reconciler. For every account on every shard, the sum of its Ledger entries must equal the balance the shard reports, and any account where they don't is reported with the size of the gap. Also check conservation across the whole system: the sum of all balances plus money in flight equals what was seeded plus what was opened, and the sum of all Ledger entries is zero. Report any payment in flight longer than a threshold as stuck. Where a shard has a replica, run against it too and report any difference. It reports; it never corrects.

Then answer three things from the code you have, not from general knowledge: which single part failing would hurt customers most, which failure the system currently handles worst, and which piece is over-engineered for its actual load.

Where the code doesn't match how a real payments system would do it, say so plainly.

Done when the Reconciler reports no discrepancies on a healthy two-shard system, reports the exact gap after I kill a shard mid-payment with a deliberately broken transaction, reports a stuck payment when I stop the recipient's shard mid-saga, and I have your three answers.`},
			{Label: "Portal restate", Portal: true, Thinking: true, Text: `A customer's wallet page has grown a screen at a time: a balance, sending money to a handle, a history, a note that a message was delivered, and a sign-in in front of all of it.

Describe that page back to me: every screen, what a customer can do on each, and which part of the system answers it. Use your own words, not mine. Say which parts you are unsure of.

Done when I have your description of every screen, and I can tell you where yours and mine disagree.`},
			{Label: "Portal", Portal: true, Text: `The Portal has a screen for each thing it learned to do, added a chapter at a time and looking like it.

Finish it: one menu from which a customer opens an account, sees what they hold, sends money, reads their history and knows a message was delivered, with the switcher deciding whose wallet it all belongs to.

Done when someone who has never seen peyva can use it without being told how.`},
		},
	},
}
