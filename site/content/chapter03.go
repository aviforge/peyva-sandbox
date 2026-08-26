package content

var Chapter03 = ChapterContent{
	Number:     3,
	Slug:       "chapter-3",
	Title:      "Across the Wire (Networking)",
	Subtitle:   "A message reaches someone when you know where to go, which door to use, and how the package travels.",
	Category:   "Foundations",
	Difficulty: "Beginner",
	EstTime:    "15 min",
	QuickTip:   "TCP guarantees packets arrive in order with none missing: exactly what moving money needs.",

	HeroImage:   "images/chapter-3.webp",
	HeroCaption: "peyva sends data between machines over the network.",

	Intuition: []string{
		"Chapter 2 got peyva listening on your own machine, but Alice and Bob aren't on the same one.",
		"Reaching another machine needs an address (IP), a door (port), a shared language (protocol), and packets to carry the message.",
	},

	Concepts: []ConceptItem{
		{Term: "IP Address", Description: "Where a machine lives on the network, the equivalent of a street address."},
		{Term: "Port", Description: "Which door on that machine to use: same idea as Chapter 2, reached remotely."},
		{Term: "Protocol", Description: "The agreed-upon way of talking. peyva speaks TCP/IP."},
		{Term: "Packets", Description: "Chunks a message is split into for the trip, reassembled on arrival."},
	},

	UnderTheHood: []string{
		"Your App -> TCP/IP -> Internet -> the other machine's App, following IP addresses on both ends.",
		"TCP re-sends anything lost and reassembles what arrives out of order, so peyva never sees a half-delivered request.",
		"net.Listen from Chapter 2 already speaks TCP/IP. This chapter is about reaching a different machine, not just your own.",
	},

	BuildIt: BuildIt{
		Intro:     "No component this chapter: reach the Gateway from a different machine.",
		Technique: "Add context and motivation",
		Why:       "Say what your situation is and why you're asking, and the assistant generalises from it. Picking commands that run on your OS and warning you about the firewall prompt before you hit it. Two lines of context saves three rounds of that didn't work.",
		Source:    "Anthropic: Prompting best practices, Add context to improve performance",
		Prompt: "My environment: <your OS and shell>. The Gateway is listening on TCP port 9310 on this machine, and I have a phone on the same Wi-Fi.\n\n" +
			"I want to reach the Gateway from the phone instead of from localhost. I'm doing this to understand how a service becomes reachable beyond the machine it runs on, so tell me what is actually happening rather than only what to type.\n\n" +
			"Tell me the command for my OS to find this machine's network address, and how to know which of the addresses it prints is the right one. Then tell me whether the Gateway needs a code change to accept connections from another machine, or whether it already does.\n\n" +
			"If a firewall prompt is likely on my OS, warn me before I hit it.\n\n" +
			"Done when the phone's browser reaches the Gateway and the Gateway logs that connection.",
	},

	BreakIt: BreakIt{
		Intro: "Networks aren't perfect. See what happens when the path breaks.",
		Exercises: []string{
			"Turn off Wi-Fi on the device trying to reach peyva. The connection fails immediately with no route to the machine.",
			"Try reaching peyva at the wrong port (e.g. 9311 instead of 9310): connection refused, even though the IP address is right.",
			"The IP gets you to the right machine, but the port still decides whether anyone answers.",
		},
	},
}
