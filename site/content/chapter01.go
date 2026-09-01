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
		What:      "Telling the assistant who to be, such as a systems engineer, before you ask your question.",
		Why:       "A systems engineer at your shoulder explains memory differently from a textbook, and the role decides what it assumes you already know.",
		Source:    "Anthropic: Prompting best practices, Give Claude a role",
		SourceURL: AnthropicBestPracticeURL + "#give-claude-a-role",
		Prompts: []Prompt{
			{Label: "Build", Text: `You are a systems engineer beside me, pointing at my screen. I have a program running: a Vault holding one balance in memory.

Show me how the operating system sees it: its process id, its CPU and memory, what each column means. Then say where in that memory the balance lives, and why it is gone when I stop the process.

I can program. Skip what a program is.

Done when I can state the process id and its memory use, and say in one sentence why nothing survives a restart.`},
			{Label: "Try", Reader: true, Text: `The assistant just read the process table for you. Read it yourself. With the Vault running in one terminal, run this in a second one: it lists your newest processes with their memory. Find the Vault: the one running your language, started a moment ago. Note its id and its memory. Then press Ctrl+C in the Vault's terminal, start it again, and run the list again.

You should see: a different process id for the same program, and alice back at 100. The number lived in the process you stopped, and the new one never heard of it.`,
				Commands: Commands(
					`Get-CimInstance Win32_Process | Sort-Object CreationDate -Descending | Select-Object -First 12 ProcessId, Name, @{n='MB'; e={[int]($_.WorkingSetSize / 1MB)}}`,
					`powershell -Command "Get-CimInstance Win32_Process | Sort-Object CreationDate -Descending | Select-Object -First 12 ProcessId, Name, @{n='MB'; e={[int]($_.WorkingSetSize / 1MB)}}"`,
					`ps -u "$USER" -o pid,rss,etime,comm | tail -15`,
				)},
		},
	},
}
