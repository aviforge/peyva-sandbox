package content

var Chapter10 = ChapterContent{
	Number:     10,
	Slug:       "chapter-10",
	Title:      "Growing the Team: Scale Out",
	Subtitle:   "When more customers come, we open more counters and add more staff so everyone gets served fast.",
	Category:   "System Design",
	Difficulty: "Intermediate",
	QuickTip:   "Copies can only be interchangeable if none of them hold state the others do not have.",

	HeroImage:   "images/chapter-10.webp",
	HeroCaption: "Scale out = add more servers/instances to handle more load in parallel.",

	Why: []string{
		"Copies are interchangeable only if nothing that matters lives inside one of them. A counter in a variable, a cache in a map, a database file on one copy's disk: each makes that copy different from the others, and a request routed to the wrong one gets a wrong answer.",
		"So the store cannot be inside the copies. Three copies each holding a database file are three banks with three sets of balances. The Vault becomes its own process, one of it, and the copies become clients of it over the network.",
		"A transaction cannot span processes. The debit, the credit, the Ledger entries and the idempotency record were one atomic unit because they were in one database in one process. That unit moves into the Vault's process whole, and the Teller drives it with one call rather than several.",
		"Scaling out the stateless part does not scale the store. Every payment still ends at one Vault, so the copies buy you concurrency in parsing, validation and waiting, not in writing balances. Knowing which part is the bottleneck is what makes scaling out worth doing.",
		"A load balancer is a process like any other, and it is now the one that holds the port everyone knows. It can only spread requests it can see fail, which means a copy dying mid-request looks to the balancer like a timeout, and the caller needs the retry rules from before.",
		"Capacity is arithmetic, done before building. A hundred thousand users making three payments a day is three hundred thousand a day: about three and a half a second on average, thirty-five at a ten-times peak. At a kilobyte per payment, two years of Ledger is two hundred gigabytes. Those numbers say whether you need a second copy of anything.",
	},

	Concepts: []ConceptItem{
		{Term: "Instance", Description: "One running copy of the stateless part of peyva, identical to every other. Also just called a copy."},
		{Term: "Load Balancer", Description: "Sits in front of all the copies and spreads incoming requests across them. Holds the port the outside world knows."},
		{Term: "Stateless", Description: "Keeps no unique data of its own. Any copy can handle any request, because the state lives in the Vault's own process."},
		{Term: "Horizontal Scaling", Description: "Adding more copies to handle more load, instead of making one machine bigger."},
		{Term: "Capacity Estimate", Description: "Users, actions per user, a peak factor and a record size, multiplied out. Tells you which part will fill up first, before you build any of it."},
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
