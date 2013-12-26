package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

// Info: Client <1> <User: Robo Handsm> disconnected
// Info: Client <1> <User: Robo Handsm> connected
// Warn: Perf: Spawner::update millis: 176
var consoleMsg = regexp.MustCompile("([^:]+): (.*)\n")

func parseLine(line string) (e *event, err error) {
	fmt.Print(line)
	matches := consoleMsg.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("Line format unknown.")
	}
	// Discards the message priority.
	return processLine(matches[2]), nil
}

func main() {
	flag.Parse()
	// The minecraft server actually writes to stderr, but reading from
	// stdin makes things easier since I can use bash and a pipe.
	stdin := bufio.NewReader(os.Stdin)
	for {
		line, err := stdin.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		if err != nil || len(line) <= 1 {
			continue
		}
		_, err = parseLine(line)
		if err != nil {
			fmt.Println("parseLine error:", err)
			continue
		}
	}
	os.Exit(0)
}
