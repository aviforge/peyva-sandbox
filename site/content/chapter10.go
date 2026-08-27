package content

var Chapter10 = ChapterContent{
	Number:     10,
	Slug:       "chapter-10",
	Title:      "Growing the Team: Scale Out",
	Subtitle:   "When more customers come, we open more counters and add more staff so everyone gets served fast.",
	Category:   "System Design",
	Difficulty: "Intermediate",
	QuickTip:   "The moment one copy holds something the others lack, the same request starts getting different answers.",

	HeroImage:   "images/chapter-10.webp",
	HeroCaption: "Scale out = add more servers/instances to handle more load in parallel.",

	Why: []string{
		"Copies are only swappable if none of them holds something the others lack.",
		"A database file inside each copy is three banks with three sets of balances. So the Vault becomes a program of its own.",
		"A transaction cannot stretch across two programs. The whole payment now runs inside the Vault, and the Teller makes one call.",
		"More copies help with reading requests and waiting around. They do not make writing faster: every payment still ends at one Vault.",
		"The load balancer now holds the port everyone knows. A copy dying mid-request looks to it like a timeout.",
		"Whether any of this is needed is arithmetic. The sidebar below does it.",
	},

	Aside: &Aside{
		Title:       "How Big Is PEYVA? (Capacity Estimation)",
		HeroImage:   "images/sidebar-10.webp",
		HeroCaption: "Capacity estimation helps us choose the right technology, plan scaling, and control cost: before we build.",
		Why: []string{
			"An estimate is a few guesses multiplied out: how many users, how often, how much busier the peak is, how big a record is.",
			"100,000 users at three payments a day is 300,000 a day: about 3.5 a second, 35 at a ten times peak.",
			"At a kilobyte a payment, two years of Ledger is around 220 GB.",
			"Size for the peak, not the average. The busiest hour is the one that matters.",
			"Roughly right before you build beats exactly right afterwards. Know which guess moves the answer most.",
		},
		BuildIt: BuildIt{
			Technique: "Self-Ask",
			Why:       "Hidden guesses become a written list of assumptions you can argue with.",
			Source:    "The Prompt Report: Zero-Shot, Self-Ask",
			Prompts: []Prompt{
				{Label: "Assumptions", Thinking: true, Text: `I need to size the infrastructure for a payments system that moves money between accounts. One process today, and I have no idea what it needs to survive.

Don't give me a number yet. Work out which follow-up questions the estimate genuinely depends on, and write them out. Answer each one yourself with an explicit assumption, and label where that assumption came from: industry norm, your own guess, or arithmetic from an earlier answer.

Done when I have your questions and an answer to each, every one labelled with where it came from.`},
				{Label: "Estimate", Thinking: true, Text: `You worked out the questions a capacity estimate for this payments system depends on, and answered each with a labelled assumption.

Use your own answers to work out peak payments per second, Ledger growth over two years, and peak network throughput. Show each formula with the numbers substituted in, so I can check the arithmetic rather than trust it.

Then tell me which single assumption the estimate is most sensitive to, what the number becomes if you're wrong about it by a factor of two, and which of your assumptions you most want me to confirm.

Done when I have a peak-throughput figure and a two-year storage figure I can defend, and I know which assumption to revisit first.`},
			},
		},
	},

	Concepts: []ConceptItem{
		{Term: "Instance", Description: "One running copy of the part of peyva that holds nothing of its own. Also just called a copy."},
		{Term: "Load Balancer", Description: "Sits in front of the copies and spreads requests across them. It holds the port the outside world knows."},
		{Term: "Stateless", Description: "Keeps no data of its own, so any copy can answer any request. The data lives in the Vault."},
		{Term: "Horizontal Scaling", Description: "Adding more copies to handle more load, instead of buying one bigger machine."},
		{Term: "Capacity Estimate", Description: "How many users, how often they act, how much busier the peak is, how big a record is. Multiply out, and you know what fills up first."},
	},

	BuildIt: BuildIt{
		Technique: "Step-Back Prompting",
		Why:       "Reasoning down from what makes any service replaceable beats reasoning up from code that was never going to scale.",
		Source:    "The Prompt Report: Thought Generation, Step-Back Prompting",
		Prompts: []Prompt{
			{Label: "Principle", Thinking: true, Text: `I have a service that handles payment requests, and I want to run several copies of it behind a router so load spreads across them.

Don't look at any code yet, and don't ask for it. In a few sentences, state the general property that lets any service run as several interchangeable copies. What may live inside one process, what may not, and why. Then say what that implies for a database file that currently sits inside the process, and for a transaction that currently spans the request handler and the database.

Done when I have the principle in general terms, with nothing about my project in it.`},
			{Label: "Build", Text: `The Gateway takes payment requests and the Teller acts on them, and the Vault with its SQLite file and the Ledger sit in the same process. I want several copies of the Gateway and Teller behind a router.

Audit it against the property you just described. Show me every place it fails: anything cached, counted, held in a variable that outlives one request, and the database file itself.

Then fix it. The Vault becomes its own process, the only one holding the database, listening on PEYVA_PORT and speaking HTTP to the copies. The whole atomic unit of a payment, debit, credit, Ledger entries and the idempotency record, runs inside the Vault in one transaction, and the Teller drives it with one call carrying the reference. The copies read PEYVA_VAULT for the Vault's port on localhost, and hold no balance, no file and no idempotency record of their own. Make the copies' port configurable, and add a small round-robin reverse proxy in front of them.

I start everything with a script I already have, so every process takes its settings from the environment and nowhere else:

  PEYVA_PORT   the port to listen on. Everything reads it.
  PEYVA_VAULT  the copies only: the Vault's port.
  PEYVA_PEERS  the proxy only: a comma separated list of the copies' ports.

No flags, no config file, no service discovery, no configuration framework. If a required variable is absent, say so and exit rather than guessing. A copy must not need the Vault to be up in order to start; it connects per request.

Done when the runner starts the Vault, three copies and the proxy, ten payments spread across the copies with correct final balances and one Ledger, and killing one copy mid-traffic fails no request.`},
		},
	},
}
