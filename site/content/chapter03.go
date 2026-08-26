package content

var Chapter03 = ChapterContent{
	Number:     3,
	Slug:       "chapter-3",
	Title:      "Across the Wire (Networking)",
	Subtitle:   "A message reaches someone when you know where to go, which door to use, and how the package travels.",
	Category:   "Foundations",
	Difficulty: "Beginner",
	EstTime:    "15 min",
	QuickTip:   "TCP hands your program a complete, in-order stream or fails the connection. It never delivers half a request.",

	HeroImage:   "images/chapter-3.webp",
	HeroCaption: "peyva sends data between machines over the network.",

	Concepts: []ConceptItem{
		{Term: "IP Address", Description: "Where a machine lives on the network, the equivalent of a street address."},
		{Term: "Protocol", Description: "The agreed-upon way of talking. peyva speaks TCP/IP."},
		{Term: "Packets", Description: "Chunks a message is split into for the trip, reassembled on arrival."},
	},

	BuildIt: BuildIt{
		Technique: "Add context and motivation",
		Why:       "Told where you are and why, it picks commands that run on your machine and warns you about the firewall prompt.",
		Source:    "Anthropic: Prompting best practices, Add context to improve performance",
		Prompts: []Prompt{
			{Label: "Build", Text: `My environment: {os}. The Gateway is listening on TCP port 9310 on this machine, and I have a phone on the same Wi-Fi.

I want to reach the Gateway from the phone instead of from localhost. I'm doing this to understand how a service becomes reachable beyond the machine it runs on, so tell me what is actually happening rather than only what to type.

Tell me the command for my OS to find this machine's network address, and how to know which of the addresses it prints is the right one. Then tell me whether the Gateway needs a code change to accept connections from another machine, or whether it already does.

If a firewall prompt is likely on my OS, warn me before I hit it.

Done when the phone's browser reaches the Gateway and the Gateway logs that connection.`},
			{Label: "Portal", Portal: true, Text: `The Gateway accepts connections but serves nothing. Have it return peyva/portal/index.html to anything that connects, so the page I could only open locally is now reachable from my phone on the same Wi-Fi.

Tell me what changes about how the page loads its stylesheet once it arrives over a connection instead of from disk, before I hit it.

Done when the phone's browser shows alice's balance.`},
		},
	},
}
