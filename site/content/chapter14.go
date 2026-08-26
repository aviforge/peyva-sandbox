package content

var Chapter14 = ChapterContent{
	Number:     14,
	Slug:       "chapter-14",
	Title:      "Big Changes Safely: Sagas",
	Subtitle:   "We can't do everything in one shot. So we do it in steps. If one step fails, we undo the earlier steps so nothing is left half-done.",
	Category:   "Reliability",
	Difficulty: "Advanced",
	EstTime:    "25 min",
	QuickTip:   "Compensations run in strict reverse order: the last completed step is undone first.",

	HeroImage:   "images/chapter-14.webp",
	HeroCaption: "Sagas help us complete a big task in steps and roll back safely if something goes wrong.",

	Intuition: []string{
		"A single transaction works when everything lives in one database.",
		"Once separate services are involved (order, payment, kitchen, delivery), no single BEGIN/COMMIT can span them.",
		"A saga runs the steps one at a time, and undoes completed steps if a later one fails.",
	},

	Concepts: []ConceptItem{
		{Term: "Saga", Description: "A sequence of local transactions across services, coordinated so the whole workflow either completes or is undone."},
		{Term: "Local Transaction", Description: "Each step's own transaction, scoped to its own service and database: not shared with the other steps."},
		{Term: "Compensating Action", Description: "The 'undo' for a step that already succeeded, run in reverse order when a later step fails."},
		{Term: "Saga Coordinator", Description: "Tracks which steps have completed and triggers compensations if the saga needs to unwind."},
	},

	UnderTheHood: []string{
		"Order Saga: Order Service (Create Order) -> Payment Service (Charge Payment) -> Kitchen Service (Prepare Food) -> Delivery Service (Deliver Order), tracked by a Saga Coordinator.",
		"If step 3 (Kitchen) fails: run compensating actions in reverse order, refund the payment (undo step 2), cancel the order (undo step 1). Step 4 never executes.",
		"Sagas leave data eventually consistent across services, without a shared transaction holding them together.",
	},

	BuildIt: BuildIt{
		Intro:     "The Teller learns to run a payment in stages, and to unwind one that fails.",
		Technique: "Least-to-Most Prompting",
		Why:       "One prompt for a whole multi-stage workflow gets you a sketch of all of it and a working version of none.",
		Source:    "The Prompt Report: Decomposition, Least-to-Most",
		Prompts: []Prompt{
			{Label: "Build", Text: `A payment is currently one atomic step. I want the Teller to run payments that span several stages and can unwind if a later stage fails permanently.

Build it in four stages, each on what the last one left. Stop after each and tell me it works before starting the next.

1. A record, per payment reference, of which stages have completed.
2. Wire the existing money movement in as stage one, recording its completion.
3. Add stage two: crediting the recipient at a second ledger peyva does not own, which can refuse for good (a closed account) or simply be unreachable.
4. A reversal for stage one (put the money back, as a new pair of Ledger entries rather than by deleting the old ones), triggered when stage two fails in a way that can never succeed.

Distinguish permanent failures from retryable ones. Only permanent failures reverse.

Done when a permanently failing stage two puts the money back, the Ledger shows both the original payment and its reversal, and the payment's record shows every stage it passed through.`},
			{Label: "Portal", Portal: true, Intro: "The Portal has to explain a payment that was undone.", Text: `A reversed payment currently looks like two unrelated rows. Show it as what it is: the original, and the reversal that answers it, tied together.

Build it in stages. First mark a reversed payment as reversed. Then link the two rows. Then say why it was reversed. Stop after each and show me before starting the next.

Done when a customer can see that money left and came back, and why.`},
		},
	},

	BreakIt: BreakIt{
		Intro: "Force a late step to fail and confirm the earlier ones actually get undone.",
		Exercises: []string{
			"Make the recipient ledger refuse for good (a closed account) and confirm the saga reverses stage one, leaving the payer whole.",
			"Now make it merely unreachable and confirm the saga waits and retries instead of reversing. A temporary failure that triggers a reversal throws away a payment that would have gone through.",
			"Compare this to Chapter 7's transaction: that couldn't span two separate services, but a saga can.",
		},
	},
}
