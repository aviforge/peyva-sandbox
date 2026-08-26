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

	Concepts: []ConceptItem{
		{Term: "Instance", Description: "One running copy of the peyva process, identical to every other. Also just called a copy."},
		{Term: "Load Balancer", Description: "Sits in front of all the copies and spreads incoming requests across them."},
		{Term: "Stateless", Description: "Keeps no unique data of its own. Any copy can handle any request, because the state lives in the Vault."},
		{Term: "Horizontal Scaling", Description: "Adding more copies to handle more load, instead of making one machine bigger."},
	},

	BuildIt: BuildIt{
		Technique: "Step-Back Prompting",
		Why:       "Reasoning down from what makes any service replaceable beats reasoning up from code that was never going to scale.",
		Source:    "The Prompt Report: Thought Generation, Step-Back Prompting",
		Prompts: []Prompt{
			{Label: "Principle", Thinking: true, Text: `I have a service that handles payment requests, and I want to run several copies of it behind a router so load spreads across them.

Don't look at any code yet, and don't ask for it. In a few sentences, state the general property that lets any service run as several interchangeable copies. What may live inside one process, what may not, and why.

Done when I have the principle in general terms, with nothing about my project in it.`},
			{Label: "Build", Text: `The Gateway takes payment requests and the Teller acts on them, both in one process, with the Vault behind them holding balances. I want several copies of that process behind a router.

Audit it against the property you just described. Show me every place it fails: anything cached, counted, or held in a variable that outlives one request.

Then fix what you found, make the port configurable, and add a small round-robin reverse proxy in front.

I start the copies with a script I already have, so both the copies and the proxy take their settings from the environment and nowhere else:

  PEYVA_PORT   the port to listen on. Both read it.
  PEYVA_PEERS  the proxy only: a comma separated list of the copies' ports.

No flags, no config file, no service discovery, no configuration framework. If either variable is absent, say so and exit rather than guessing.

Done when three copies started this way spread ten payments between them with correct final balances, and killing one mid-traffic fails no request.`},
		},
	},
}
