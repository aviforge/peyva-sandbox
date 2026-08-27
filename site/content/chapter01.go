package content

var Chapter01 = ChapterContent{
	Number:     1,
	Slug:       "chapter-1",
	Title:      "Inside One Computer",
	Subtitle:   "A running program is one room in a house. The OS is the caretaker, the hardware is the house.",
	Category:   "Foundations",
	Difficulty: "Beginner",
	QuickTip:   "Right now, everything peyva needs is inside this one computer.",

	HeroImage:   "images/chapter-1.webp",
	HeroCaption: "A computer is like a house. Many things work together inside it to get one job done.",

	Why: []string{
		"A process is the thing that fails. When it dies, everything in its memory dies with it.",
		"The OS can pause your program between any two lines, at any time.",
		"Memory is fast and forgets on restart. Disk is far slower and remembers.",
		"One core does one thing at a time. Two jobs that look simultaneous are taking turns, and you do not choose the order.",
		"A distributed system is several of these. Adding a network makes none of them more reliable.",
	},

	Concepts: []ConceptItem{
		{Term: "Process", Description: "A program that is running, with its own memory that no other program can reach."},
		{Term: "Operating System (OS)", Description: "Decides which program gets the CPU, how much memory it may hold, and what it may touch."},
		{Term: "CPU", Description: "Runs instructions. One core does one thing at a time, however many programs are open."},
		{Term: "Memory (RAM)", Description: "Where a program keeps what it is working on. Fast, and empty again the moment the program stops."},
		{Term: "Disk (Storage)", Description: "Where data outlives the program that wrote it. Far slower than memory."},
		{Term: "Scheduler", Description: "The part of the OS that picks which program runs next, and for how long. It can pause yours between any two instructions."},
	},

	BuildIt: BuildIt{
		Technique: "Role prompting",
		Why:       "Who you tell it to be decides what it assumes you already know.",
		Source:    "Anthropic: Prompting best practices, Give Claude a role",
		Prompts: []Prompt{
			{Label: "Build", Text: `You are a systems engineer sitting beside me, pointing at what is on my screen. Explain a column when I meet it, not before.

I have a program running: a Vault holding one balance in memory.

Show me how the operating system sees it. How to find its process id, how to read its CPU and memory, what each column means. Then tell me which part of that memory holds the balances, and why the number goes to zero when I stop the process.

I can program, but I have never looked at a process from outside. Skip what a program is.

Done when I can state the process id and its memory use, and say in one sentence why nothing survives a restart.`},
		},
	},
}
