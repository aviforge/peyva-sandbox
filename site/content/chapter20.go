package content

var Chapter20 = ChapterContent{
	Number:     20,
	Slug:       "chapter-20",
	Title:      "Putting It All Together",
	Subtitle:   "Every piece from the last twenty chapters works together to make peyva fast, reliable, and safe.",
	Category:   "System Design",
	Difficulty: "Advanced",
	EstTime:    "30 min",
	QuickTip:   "No single piece makes a system reliable: only how they cooperate under load and failure.",

	HeroImage:   "images/chapter-20.webp",
	HeroCaption: "System design is not one clever box. It's how all the pieces cooperate under load, failure and change.",

	Concepts: []ConceptItem{
		{Term: "Reconciler", Description: "The component that checks the Vault against the Ledger and reports where they disagree."},
		{Term: "Reconciliation", Description: "Comparing two records that were written independently, to catch the case where one of them is wrong."},
		{Term: "Discrepancy", Description: "An account whose Ledger entries do not sum to the balance the Vault reports, and the size of the gap."},
		{Term: "Single Point of Failure", Description: "A part with no stand-in, whose loss stops the system rather than slowing it."},
	},

	BuildIt: BuildIt{
		Technique: "Rephrase and Respond (RaR)",
		Why:       "The restatement is the cheapest look you will get at what it actually understood.",
		Source:    "The Prompt Report: Zero-Shot, Rephrase and Respond",
		Prompts: []Prompt{
			{Label: "Restate", Thinking: true, Text: `peyva is a payments system built from these parts: a Gateway that takes requests from outside, a Teller that runs one payment end to end, a Vault that is the only thing that changes a balance, a Ledger that records every movement, a Courier that carries out work after a payment clears, and a Portal a customer uses.

Before writing anything, restate that back to me in your own words. Describe what each part is for, and how a payment travels from the Gateway to the point where the recipient has been told. Don't repeat my names for things back at me: say what each one actually does.

At each hop, say what would happen if that part failed right then.

Done when I have your description of every part and the path a payment takes, and I can tell you where yours and mine disagree.`},
			{Label: "Build", Text: `peyva has a Vault holding balances and a Ledger recording every movement of money, and nothing checks that the two still agree.

Build the Reconciler. For every account, the sum of its Ledger entries must equal the balance the Vault reports, and any account where they don't is reported with the size of the gap.

Then answer three things from the code you have, not from general knowledge: which single part failing would hurt customers most, which failure the system currently handles worst, and which piece is over-engineered for its actual load.

Where the code doesn't match how a real payments system would do it, say so plainly.

Done when the Reconciler reports no discrepancies on a healthy system, reports the exact gap after I kill the process mid-payment, and I have your three answers.`},
			{Label: "Portal restate", Portal: true, Thinking: true, Text: `A customer's wallet page has grown a screen at a time: a balance, sending money to a handle, a history, a note that a message was delivered, and a sign-in in front of all of it.

Describe that page back to me: every screen, what a customer can do on each, and which part of the system answers it. Use your own words, not mine. Say which parts you are unsure of.

Done when I have your description of every screen, and I can tell you where yours and mine disagree.`},
			{Label: "Portal", Portal: true, Text: `The Portal has a screen for each thing it learned to do, added a chapter at a time and looking like it.

Finish it: one menu from which a customer opens an account, sees what they hold, sends money, reads their history and knows a message was delivered, with the switcher deciding whose wallet it all belongs to.

Done when someone who has never seen peyva can use it without being told how.`},
		},
	},
}
