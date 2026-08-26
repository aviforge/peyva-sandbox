package content

var Chapter10 = ChapterContent{
	Number:     10,
	Slug:       "chapter-10",
	Title:      "Growing the Team: Scale Out",
	Subtitle:   "When more customers come, we open more counters and add more staff so everyone gets served fast.",
	Category:   "System Design",
	Difficulty: "Intermediate",
	EstTime:    "20 min",
	QuickTip:   "Copies can only be interchangeable if none of them hold state the others do not have.",

	HeroImage:   "images/chapter-10.webp",
	HeroCaption: "Scale out = add more servers/instances to handle more load in parallel.",

	Intuition: []string{
		"Chapter 9's math said peyva needs to sustain real load. One process can't do that alone.",
		"A busy café doesn't make one barista work faster; it opens more counters.",
		"Scaling out runs several peyva copies side by side, with something in front routing each request.",
	},

	Concepts: []ConceptItem{
		{Term: "Instance", Description: "One running copy of the peyva process, identical to every other. This book usually just says copy."},
		{Term: "Load Balancer", Description: "Sits in front of all the copies and spreads incoming requests across them."},
		{Term: "Stateless", Description: "Keeps no unique data of its own. Any copy can handle any request, because the state lives in the Vault."},
		{Term: "Horizontal Scaling", Description: "Adding more copies to handle more load, instead of making one machine bigger."},
	},

	UnderTheHood: []string{
		"Users -> Load Balancer -> Instance 1 / Instance 2 / Instance N, all backed by the same Database.",
		"Before scaling: one copy handles everyone, slow when busy. After: a load balancer spreads requests across many copies, fast even when busy.",
	},

	BuildIt: BuildIt{
		Intro:     "Run several Gateways and Tellers side by side, sharing one Vault.",
		Technique: "Step-Back Prompting",
		Why:       "Make the assistant abstract to the general principle before it touches the specifics. Reasoning down from what makes any service replaceable beats reasoning up from this code, which tends to optimise a design that was never going to scale.",
		Source:    "The Prompt Report: Thought Generation, Step-Back Prompting",
		Prompt: `I want to run several copies of the Gateway and Teller behind a router, so load spreads across them.

Step back before you look at my code. In a few sentences, state the general property that lets any service run as several interchangeable copies. What may live inside one process, what may not, and why. Don't mention this codebase yet.

Now apply that principle here. Audit the code against it and show me every place it currently fails: anything cached, counted, or held in a package-level variable.

Then fix what you found, make the port configurable, and add a small round-robin reverse proxy in front using only the standard library.

I start the copies with a script I already have, so both the copies and the proxy take their settings from the environment and nowhere else:

  PEYVA_PORT   the port to listen on. Both read it.
  PEYVA_PEERS  the proxy only: a comma separated list of the copies' ports.

No flags, no config file, no defaults that hide a missing value. If either is absent, say so and exit rather than guessing a port.

No service discovery, no external load balancer, no configuration framework.

Done when three copies started this way spread ten payments between them with correct final balances, and killing one mid-traffic fails no request.`,
	},

	BreakIt: BreakIt{
		Intro: "Prove the copies really are interchangeable.",
		Exercises: []string{
			"Start three copies with the runner, send ten payments through the proxy, and confirm from the tagged output that they landed on different copies while all updating the same balances.",
			"Kill one copy while traffic is flowing. The proxy routes around it and requests keep succeeding. Start the runner with one copy instead of three and the same traffic now queues behind a single process.",
			"This only works because Chapter 5 already moved state into a shared store. A copy holding state of its own could not be killed this safely.",
			"Push writes through several copies at once and watch them queue behind each other. Scaling out multiplied the copies, not the single Vault they all write to, so the bottleneck moved rather than disappeared.",
			"Close the terminal without stopping anything, then run stop from a new one. Every copy goes, the ports are free, and starting again works first time.",
		},
	},
}
