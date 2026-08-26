package content

var Chapter01 = ChapterContent{
	Number:     1,
	Slug:       "chapter-1",
	Title:      "Inside One Computer",
	Subtitle:   "The Chapter 0 process has to actually run somewhere — this is where.",
	Category:   "Foundations",
	Difficulty: "Beginner",
	EstTime:    "10 min",
	QuickTip:   "Right now, everything peyva needs is inside this one computer.",

	HeroImage:   "images/chapter-1.webp",
	HeroCaption: "A computer is like a house. Many things work together inside it to get one job done.",

	Intuition: []string{
		"The peyva program from Chapter 0 isn't magic — it's a room in a house.",
		"The OS is the caretaker; the hardware is the house.",
		"CPU does the work, memory (RAM) is today's desk, disk is the long-term filing cabinet.",
	},

	Concepts: []ConceptItem{
		{Term: "Process", Description: "A running program — the room peyva's app is currently working in."},
		{Term: "Operating System (OS)", Description: "The caretaker that manages resources, files, and hardware for every process."},
		{Term: "CPU", Description: "The worker who executes instructions and does calculations."},
		{Term: "Memory (RAM)", Description: "The desk holding data the app is actively using — fast, but empty when power goes off."},
		{Term: "Disk (Storage)", Description: "The filing cabinet where data is kept long term, even after a restart."},
	},

	UnderTheHood: []string{
		"One computer = peyva App Process -> Operating System -> Hardware.",
		"Running `go run .` asks the OS to create a process, load it into memory, and give it CPU time.",
	},

	BuildIt: BuildIt{
		Intro:     "No component this chapter — see the Vault's process the way the OS sees it.",
		Technique: "Role prompting",
		Why:       "A role focuses the assistant's tone and its judgement about what to include. Unassigned, it explains what a program is; cast as an engineer sitting beside you, it hands you the command and tells you what the columns mean.",
		Source:    "Anthropic — Prompting best practices, Give Claude a role",
		Prompt: "You are a systems engineer sitting next to me, teaching by pointing at what is actually on my screen. You explain a column when I meet it, not before.\n\n" +
			"I have a Go program running — a Vault holding one account balance in memory.\n\n" +
			"Walk me through inspecting it as the operating system sees it: how to find its process id, and how to read its CPU and memory use. Give me the command for my OS and explain what each column means.\n\n" +
			"Then tell me which part of that memory holds the Vault's balances, and why that number goes to zero the moment I stop the process.\n\n" +
			"I know Go but I've never looked at a process from the outside. Skip the explanation of what a program is.\n\n" +
			"Done when I can state the process id and its memory use, and explain in one sentence why nothing survives a restart.",
	},

	BreakIt: BreakIt{
		Intro: "Processes are isolated from each other — prove it.",
		Exercises: []string{
			"Start two copies of peyva at once. Each gets its own PID and its own memory — they don't share the wallet variable.",
			"Kill one with Ctrl+C. The other keeps running untouched.",
			"That isolation is why Chapter 10 needs a load balancer once there's more than one instance to route between.",
		},
	},
}
