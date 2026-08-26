package content

var Chapter09 = ChapterContent{
	Number:     9,
	Slug:       "chapter-9",
	Title:      "How Big Is PEYVA? (Capacity Estimation)",
	Subtitle:   "Before starting the trip, we estimate the numbers so we have enough car space, fuel, and time. Same for systems.",
	Category:   "System Design",
	Difficulty: "Intermediate",
	EstTime:    "20 min",
	QuickTip:   "Being roughly right before you build beats being exactly right after.",

	HeroImage:   "images/chapter-9.webp",
	HeroCaption: "Capacity estimation helps us choose the right technology, plan scaling, and control cost: before we build.",

	Concepts: []ConceptItem{
		{Term: "Assumptions", Description: "The starting numbers you estimate: users, transfers per user per day, peak factor, time window."},
		{Term: "Peak Factor", Description: "How much busier the system gets at its busiest moment compared to its average. Traffic isn't flat."},
		{Term: "TPS", Description: "Transactions per second: the rate the system must actually sustain at peak, not on average."},
		{Term: "Storage Growth", Description: "How much disk space the system needs over time, based on transaction volume and record size."},
	},

	BuildIt: BuildIt{
		Intro:     "No component this chapter. Size what you have before you grow it.",
		Technique: "Self-Ask",
		Why:       "Hidden guesses become a written list of assumptions you can argue with.",
		Source:    "The Prompt Report: Zero-Shot, Self-Ask",
		Prompts: []Prompt{
			{Label: "Assumptions", Thinking: true, Intro: "The questions the estimate depends on, and an answer to each.", Text: `I need to size the infrastructure for a payments system that moves money between accounts. One process today, and I have no idea what it needs to survive.

Don't give me a number yet. Work out which follow-up questions the estimate genuinely depends on, and write them out. Answer each one yourself with an explicit assumption, and label where that assumption came from: industry norm, your own guess, or arithmetic from an earlier answer.

Done when I have your questions and an answer to each, every one labelled with where it came from.`},
			{Label: "Estimate", Thinking: true, Intro: "The arithmetic, with the numbers shown.", Text: `You worked out the questions a capacity estimate for this payments system depends on, and answered each with a labelled assumption.

Use your own answers to work out peak payments per second, Ledger growth over two years, and peak network throughput. Show each formula with the numbers substituted in, so I can check the arithmetic rather than trust it.

Then tell me which single assumption the estimate is most sensitive to, what the number becomes if you're wrong about it by a factor of two, and which of your assumptions you most want me to confirm.

Done when I have a peak-throughput figure and a two-year storage figure I can defend, and I know which assumption to revisit first.`},
		},
	},
}
