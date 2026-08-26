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

	Intuition: []string{
		"Twenty chapters ago, peyva was one file holding one account.",
		"Every chapter since added one honest piece: a port, an API, a database, a queue, a replica, security, and more.",
		"A restaurant's Front Door, Hosts, Kitchen, Cold Storage, and Order Tickets are really just all of those pieces working together.",
	},

	Concepts: []ConceptItem{
		{Term: "Load Balancer", Description: "The Front Door: routes users to available peyva instances."},
		{Term: "peyva Instances", Description: "The Hosts: stateless, scaled-out copies of the app from Chapter 10."},
		{Term: "Database & Cache", Description: "The Kitchen and Cold Storage: durable data from Chapter 5, fast lookups from Chapter 11."},
		{Term: "Queue, Outbox & Workers", Description: "Order Tickets. Asynchronous, reliable work from Chapters 12 and 13."},
		{Term: "Reconciler", Description: "The Till Count: proves the Vault and the Ledger still agree, and names any account where they don't."},
	},

	UnderTheHood: []string{
		"Engineer view: Users -> Load Balancer -> peyva Instances (stateless) -> Cache / Database / Queue / Outbox -> Workers, with Security wrapping every edge and Observability watching every part.",
		"The chain: Networking gets the request in -> the API handles it -> Scale Out shares the load -> Cache serves the hot path -> the Database holds the truth -> Queue & Outbox handle the rest asynchronously.",
		"Restaurant operations underneath it all: health checks, deployments, alerts, runbooks, and incident response (Chapter 19's practices, running continuously).",
	},

	BuildIt: BuildIt{
		Intro:     "Build the Reconciler, and have the whole system explained back to you.",
		Technique: "Rephrase and Respond (RaR)",
		Why:       "Make the assistant restate the task in its own words before it answers. The restatement is the cheapest look you'll get at what it actually understood, and comparing it against the story in your head shows you what you didn't actually learn.",
		Source:    "The Prompt Report: Zero-Shot, Rephrase and Respond",
		Prompt: `I want you to build the Reconciler. The component that proves the Vault and the Ledger agree. For every account, the sum of its Ledger entries must equal the balance the Vault reports, and any account where they don't is reported with the size of the gap.

Before you write it, restate the job back to me in your own words. Read the whole system first, then describe what each component is for and how a payment travels from the Gateway through the Teller, the Vault, the Ledger and the Courier. Don't repeat my names for things back at me. Describe what the code actually does. At each hop, say what would happen if that component failed right then.

I'll tell you where your restatement doesn't match mine before you write any code.

Then build the Reconciler, and answer three things from the code rather than from general knowledge: which single component failing would hurt customers most, which failure the system currently handles worst, and which piece is over-engineered for its actual load.

Where the code doesn't match how a real payments system would do it, say so plainly.

Done when the Reconciler reports no discrepancies on a healthy system, reports the exact gap after I kill the process mid-payment, and your restatement and mine agree.`,
		UIIntro: "The whole portal, explained back to you.",
		UIPrompt: `Before touching it, describe the portal back to me: every screen, what a customer can do on each, and which component answers it. Use your own words, not mine.

Where your description and mine differ, one of us has misunderstood the system. Say which parts you are unsure of.

Then finish it: one place a customer can open an account, see a balance, send money, read their history and know a message was delivered.

Done when someone who has never seen peyva can use it without being told how.`,
	},

	BreakIt: BreakIt{
		Intro: "Break one piece deliberately and confirm the rest of the system absorbs it gracefully, the way earlier chapters promised.",
		Exercises: []string{
			"Kill one instance mid-traffic. The load balancer routes around it (Ch. 10). Nothing is lost.",
			"Disconnect the replica. Writes are refused, reads keep working (Ch. 16's CP choice).",
			"Send the exact same transfer twice. Bob is paid once, not twice (Ch. 8's idempotency, still holding under the full system's weight).",
		},
	},
}
