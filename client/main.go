package main

import (
	"fmt"
	"github.com/musica/go-osc"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var port int

func bundle(time time.Time) *osc.Bundle {
	return osc.NewBundle().SetTimeTag(time)
}

func message(address string) *osc.Message {
	msg := osc.NewMessage()
	msg.SetAddress(address)
	return msg
}

func intArg(n int) *osc.Argument {
	return osc.NewArgument().SetInt32(int32(n))
}

func stringArg(s string) *osc.Argument {
	return osc.NewArgument().SetString(s)
}

func systemPlayMsg() *osc.Message {
	return message("/system/play")
}

func systemStopMsg() *osc.Message {
	return message("/system/stop")
}

func midiPatchMsg(track int, offset int, patch int) *osc.Message {
	return message(
		fmt.Sprintf("/track/%d/midi/patch", track),
	).AddArguments(
		intArg(offset), intArg(patch),
	)
}

func midiPercussionMsg(track int) *osc.Message {
	return message(fmt.Sprintf("/track/%d/midi/percussion", track)).AddArguments(
		intArg(0),
	)
}

func midiNoteMsg(
	track int, offset int, note int, duration int, audibleDuration int,
	velocity int) *osc.Message {
	return message(
		fmt.Sprintf("/track/%d/midi/note", track),
	).AddArguments(
		intArg(offset),
		intArg(note),
		intArg(duration),
		intArg(audibleDuration),
		intArg(velocity))
}

func patternMsg(track int, offset int, pattern string, times int) *osc.Message {
	return message(fmt.Sprintf("/track/%d/pattern", track)).AddArguments(
		intArg(offset),
		stringArg(pattern),
		intArg(times),
	)
}

func patternMidiNoteMsg(
	pattern string, offset int, note int, duration int, audibleDuration int,
	velocity int) *osc.Message {
	return message(fmt.Sprintf("/pattern/%s/midi/note", pattern)).AddArguments(
		intArg(offset),
		intArg(note),
		intArg(duration),
		intArg(audibleDuration),
		intArg(velocity),
	)
}

func patternClearMsg(pattern string) *osc.Message {
	return message(fmt.Sprintf("/pattern/%s/clear", pattern))
}

func oneNote() *osc.Bundle {
	return bundle(time.Now()).AddElements(
		midiPatchMsg(1, 0, 30),
		midiNoteMsg(1, 0, 45, 1000, 1000, 127),
		systemPlayMsg(),
	)
}

func sixteenFastNotes() *osc.Bundle {
	bundle := bundle(time.Now()).AddElements(midiPatchMsg(1, 0, 70))

	interval := 100
	audibleDuration := 80

	noteNumber := 30 + rand.Intn(60)

	for offset := 0; offset <= interval*16; offset += interval {
		bundle.AddElements(
			midiNoteMsg(1, offset, noteNumber, interval, audibleDuration, 127))
	}

	bundle.AddElements(systemPlayMsg())

	return bundle
}

func playPattern(times int) *osc.Bundle {
	pattern := "simple"
	return bundle(time.Now()).AddElements(
		patternClearMsg(pattern),
		patternMidiNoteMsg(pattern, 0, 57, 500, 500, 127),
		patternMidiNoteMsg(pattern, 500, 60, 500, 500, 127),
		patternMidiNoteMsg(pattern, 1000, 62, 500, 500, 127),
		patternMidiNoteMsg(pattern, 1500, 64, 500, 500, 127),
		patternMidiNoteMsg(pattern, 2000, 67, 500, 500, 127),
		midiPatchMsg(1, 0, 60),
		patternMsg(1, 0, pattern, times),
		systemPlayMsg(),
	)
}

func playPatternOnce() *osc.Bundle {
	return playPattern(1)
}

func playPatternTwice() *osc.Bundle {
	return playPattern(2)
}

func changePattern() *osc.Bundle {
	pattern := "simple"
	bundle := bundle(time.Now()).AddElements(patternClearMsg(pattern))

	interval := 500
	audibleDuration := 250

	for offset := 0; offset <= interval*3; offset += interval {
		noteNumber := 30 + rand.Intn(60)

		bundle.AddElements(
			patternMidiNoteMsg(
				pattern, offset, noteNumber, interval, audibleDuration, 127))
	}

	return bundle
}

func printUsage() {
	fmt.Printf("Usage: %s PORT EXAMPLE\n", os.Args[0])
}

func main() {
	rand.Seed(time.Now().Unix())

	numArgs := len(os.Args[1:])

	if numArgs < 1 || numArgs > 2 {
		printUsage()
		os.Exit(1)
	}

	port, err := strconv.ParseInt(os.Args[1], 10, 32)
	if err != nil {
		fmt.Println(err)
		printUsage()
		os.Exit(1)
	}

	var example string
	if numArgs < 2 {
		example = "1"
	} else {
		example = os.Args[2]
	}

	client := osc.NewClient()
	client.Connect("udp", fmt.Sprintf("localhost:%d", port))

	switch example {
	case "play":
		client.Send(systemPlayMsg())
	case "stop":
		client.Send(systemStopMsg())
	case "perc":
		client.Send(midiPercussionMsg(1))
	case "1":
		client.Send(oneNote())
	case "16fast":
		client.Send(sixteenFastNotes())
	case "pat1":
		client.Send(playPatternOnce())
	case "pat2":
		client.Send(playPatternTwice())
	case "patchange":
		client.Send(changePattern())
	case "patx":
		client.Send(patternClearMsg("simple"))
	default:
		fmt.Printf("No such example: %s\n", example)
		os.Exit(1)
	}
}