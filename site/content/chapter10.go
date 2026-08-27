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
				{Label: "Assumptions", Thinking: true, Text: `I need to size a payments system that moves money between accounts. One process today, and no idea what it needs to survive.

No numbers yet. Work out which questions the estimate depends on, and write them out. Answer each yourself with a stated assumption, and label where it came from: industry norm, your own guess, or arithmetic from an earlier answer.

Done when I have your questions and an answer to each, every one labelled.`},
				{Label: "Estimate", Thinking: true, Text: `You worked out the questions a capacity estimate depends on, and answered each with a labelled assumption.

Use your own answers to work out payments per second at peak, Ledger growth over two years, and network traffic at peak. Show each sum with the numbers filled in, so I can check the arithmetic rather than trust it.

Then say which assumption the answer is most sensitive to, what the number becomes if you are wrong about it by double, and which one you most want me to confirm.

Done when I have a peak figure and a two-year storage figure I can defend, and I know which assumption to revisit first.`},
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
			{Label: "Principle", Thinking: true, Text: `I have a service handling payment requests, and I want several copies of it behind a router.

Do not look at any code. In a few sentences, state the general property that lets any service run as interchangeable copies: what may live inside one process, what may not, and why. Then say what that means for a database file sitting inside the process, and for a transaction that spans the handler and the database.

Done when I have the principle in general terms, with nothing about my project in it.`},
			{Label: "Build", Text: `The Gateway, the Teller, the Vault's file and the Ledger all sit in one process. I want several copies of the Gateway and Teller behind a router.

Audit it against the property you just described, and show me every place it fails, the database file included.

Then fix it. The Vault becomes its own process, the only one holding the database, on PEYVA_PORT, speaking HTTP. The whole payment happens inside it in one transaction, driven by one call from the Teller. The copies read PEYVA_VAULT for its port and hold nothing of their own. Add a small round-robin proxy in front.

Settings come from the environment and nowhere else:

  PEYVA_PORT   the port to listen on. Everything reads it.
  PEYVA_VAULT  the copies only: the Vault's port.
  PEYVA_PEERS  the proxy only: the copies' ports, comma separated.

No flags, no config file, no service discovery. A missing variable means say so and exit.

Done when the runner starts the Vault, three copies and the proxy, ten payments spread across the copies with correct balances and one Ledger, and killing a copy mid-traffic fails no request.`},
		},
	},
}
