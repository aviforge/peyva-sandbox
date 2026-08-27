package content

var Chapter02 = ChapterContent{
	Number:     2,
	Slug:       "chapter-2",
	Title:      "Finding Peyva (Processes & Ports)",
	Subtitle:   "To reach someone, you need the right address and the receptionist lets you in.",
	Category:   "Foundations",
	Difficulty: "Beginner",
	QuickTip:   "One program holds an address and port at a time. That is what makes it a door you can knock on.",

	HeroImage:   "images/chapter-2.webp",
	HeroCaption: "A process is the running app. A port is the door callers knock on.",

	Why: []string{
		"A port is a door number. One program holds it at a time, and a second copy asking for it is refused.",
		"That refusal is the first sign that running several copies will need something in front of them.",
		"127.0.0.1 can only be reached from this machine. 0.0.0.0 can be reached by anything on the network.",
		"Every open connection costs memory until it is closed. A program that never closes them runs out.",
		"One front door is one place to later count, refuse or redirect every request.",
	},

	Concepts: []ConceptItem{
		{Term: "Process ID (PID)", Description: "A number the OS gives each running program, to tell one from another."},
		{Term: "Port", Description: "A number, like 9310, that a program claims so the OS knows which one to hand a call to."},
		{Term: "Binding", Description: "Claiming an address and port. The OS gives that pair to one program at a time, so a second copy is refused."},
		{Term: "Loopback (127.0.0.1)", Description: "The address meaning 'this same machine'. Used to reach a program running locally."},
		{Term: "Gateway", Description: "The component that holds the port and takes requests from outside. The system's one front door."},
	},

	BuildIt: BuildIt{
		Technique: "Constrain the scope",
		Why:       "Ask for a listener without stating a ceiling and you get an HTTP server with routes and JSON.",
		Source:    "Anthropic: Prompting best practices, Overeagerness",
		Prompts: []Prompt{
			{Label: "Build", Text: `peyva is one process holding the Vault, and the only way to reach it is to run it yourself. Build the Gateway, the way payment requests reach the system from outside.

For now it does one thing: claim TCP port 9310, accept connections, log one line per connection with the caller's address, then close it.

Do not add HTTP handling, routes, JSON parsing, request or response types, graceful shutdown, or any third-party package. The Gateway only owns a door right now. Everything else arrives in later chapters and would get in the way.

The right amount of complexity is the minimum that does what I just described.

Done when a request to port 9310 makes the Gateway log a connection, and starting a second copy fails with 'address already in use'.`},
		},
	},
}
