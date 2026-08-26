package content

var Chapter10 = ChapterContent{
	Number:     10,
	Slug:       "chapter-10",
	Title:      "Growing the Team: Scale Out",
	Subtitle:   "When more customers come, we open more counters and add more staff so everyone gets served fast.",
	Category:   "System Design",
	Difficulty: "Intermediate",
	EstTime:    "20 min",
	QuickTip:   "Instances can only be interchangeable if none of them hold state the others don't have.",

	HeroImage:   "images/chapter-10.webp",
	HeroCaption: "Scale out = add more servers/instances to handle more load in parallel.",

	Intuition: []string{
		"Chapter 9's math said peyva needs to sustain real load — one process can't do that alone.",
		"A busy café doesn't make one barista work faster; it opens more counters.",
		"Scaling out runs several peyva copies side by side, with something in front routing each request.",
	},

	Concepts: []ConceptItem{
		{Term: "Instance", Description: "One running copy of the peyva process, identical to every other instance."},
		{Term: "Load Balancer", Description: "Sits in front of all instances and distributes incoming requests across them."},
		{Term: "Stateless", Description: "Keeps no unique data of its own — any instance can handle any request, since state lives in the database."},
		{Term: "Horizontal Scaling", Description: "Adding more machines/instances to handle more load, instead of making one machine bigger."},
	},

	UnderTheHood: []string{
		"Users -> Load Balancer -> Instance 1 / Instance 2 / Instance N, all backed by the same Database.",
		"Before scaling: one instance handles everyone — slow when busy. After: a load balancer spreads requests across many instances — fast even when busy.",
	},

	BuildIt: BuildIt{
		Intro:     "Run several Gateways and Tellers side by side, sharing one Vault.",
		Technique: "Step-Back Prompting",
		Why:       "Make the assistant abstract to the general principle before it touches the specifics. Reasoning down from what makes any service replaceable beats reasoning up from this code, which tends to optimise a design that was never going to scale.",
		Source:    "The Prompt Report — Thought Generation, Step-Back Prompting",
		Prompt: "I want to run several copies of the Gateway and Teller behind a router, so load spreads across them.\n\n" +
			"Step back before you look at my code. In a few sentences, state the general property that lets any service run as several interchangeable copies — what may live inside one process, what may not, and why. Don't mention this codebase yet.\n\n" +
			"Now apply that principle here. Audit the code against it and show me every place it currently fails — anything cached, counted, or held in a package-level variable.\n\n" +
			"Then fix what you found, make the port configurable so I can run three copies against one Vault, and add a small round-robin reverse proxy in front using only the standard library.\n\n" +
			"No service discovery, no external load balancer, no configuration framework.\n\n" +
			"Done when ten payments spread across all three copies with correct final balances, and killing one copy mid-traffic fails no request.",
	},

	BreakIt: BreakIt{
		Intro: "Prove the instances really are interchangeable.",
		Exercises: []string{
			"Send ten transfer requests through the load balancer and confirm they land on different instances, but all update the same shared balances correctly.",
			"Kill one instance while traffic is flowing — the load balancer should route around it, and requests keep succeeding.",
			"This only works because Chapter 5 already moved state into a shared database — a stateful instance couldn't be killed this safely.",
		},
	},
}
