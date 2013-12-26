package main

// Process serverbound output lines and to cool things with them.

import (
	"fmt"
	"regexp"
	"time"
)

type event struct {
	// t is the time when the parser saw the event line. Serverbound itself does
	// not show timestamps, so this is not precise.
	time time.Time
	kind eventKind
	// Details are parsed details after regexp extraction with a matching
	// eventHandler. The slice members are different depending on the
	// eventKind.
	details []string
}

// processLine looks at one line of output from the serverbound console. If the line matches with
// one defined in eventHandlers, an optional function for that event is executed and the event is
// returned. If the line doesn't match any known events, returns nil.
func processLine(line string) *event {
	for _, h := range eventHandlers {
		p := h.re.FindStringSubmatch(line)
		if p == nil {
			continue
		}
		ev := &event{
			kind: h.kind,
			// Best-effort event timestamp.
			time: time.Now(),
			// Save the matched substrings into the event.
			details: p,
		}
		// Execute this event kind's function.
		h.f(ev)
		return ev
	}
	return nil
}

type eventKind int

const (
	userConnected    eventKind = 1
	userDisconnected eventKind = iota
)

var (
	// Client <1> <User: Robo Handsm> connected
	userConnectedRegexp = regexp.MustCompile(`Client <\d+> <User: ([^>]+)> connected`)
	// Client <1> <User: Robo Handsm> disconnected
	userDisconnectedRegexp = regexp.MustCompile(`Client <\d+> <User: ([^>]+)> disconnected`)
)

type eventHandler struct {
	re   *regexp.Regexp
	f    func(*event)
	kind eventKind
}

var eventHandlers = []eventHandler{
	announceConnected,
	twitJournal,
}

// announceConnected will post to twitter when a user logs in to the server.
var announceConnected = eventHandler{
	userConnectedRegexp,
	func(ev *event) {
		msg := fmt.Sprintf("User %v logged in!", ev.details[1])
		fmt.Println("twit:", msg)
		/*
			err := twitter.Update(msg)
			if err != nil {
				fmt.Println("Twitter error:", err)
			}
		*/
	},
	userConnected,
}

// twitJournal will post a message from a user to twitter.
var twitJournal = eventHandler{
	// 2011-04-22 12:12:49 [INFO] nictuku tried command: twit foo
	regexp.MustCompile(`([^ ]+) tried command: twit (.+)$`),
	func(ev *event) {
		msg := fmt.Sprintf("<%v> %v", ev.details[1], ev.details[2])
		fmt.Println("twit:", msg)
		/*
			err := twitter.Update(msg)
			if err != nil {
				fmt.Println("Twitter error:", err)
			}
		*/
	},
	userDisconnected,
}
