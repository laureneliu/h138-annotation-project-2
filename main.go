package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"

	"github.com/logrusorgru/aurora"
)

var trim = " \n\t"

var dirValid = map[string]bool {
	"N": true,
	"S": true,
	"E": true,
	"W": true,
}

var reader = bufio.NewReader(os.Stdin)

func main() {
	spec, err := ParseSpec("spec.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	rooms := spec.Rooms
	riddles := spec.Riddles

	room := rooms[8] // Start in entrance hall
	riddlesI := 0 // Start with first riddle

	PrintPlain("--------------------------------------------------------------------------------\n")
	PrintPlain("It's 10am. You just came back from a quick run, but as you walk up to the front door you see there are some post it notes and an envelope stuck to the front door. You wonder what your dumbass housemates got up to while you were gone...\n")
	PrintPlain("'To thank you for being such a wonderful housemate, we've made a fun scavenger hunt, just for you!'\n")
	PrintPlain("This is one of the most suspicious things you've seen in your life. Whatever, you'll play along. You open the envelope.\n\n")

	PrintPlain(fmt.Sprintf("Here's your first riddle:\n%s\n\n", riddles[riddlesI].Q))

	PrintHelp()
	PrintPlain("Type 'help' at any time to see the list of valid moves.\n\n")

	PrintRoom(room, rooms)

	for riddlesI < len(riddles) {
		inp := GetInput("")

		if strings.HasPrefix(inp, "move ") {
			// Move around the board
			i := strings.Index(inp, " ")
			dir := strings.Trim(inp[i+1:], trim)
			if !dirValid[dir] {
				PrintRed(fmt.Sprintf("Invalid move direction %s.\n", dir))
			} else {
				nextId, ok := room.Links[dir]
				if !ok {
					PrintRed("Can't move in that direction.\n")
				} else {
					room = rooms[nextId]
					PrintRoom(room, rooms)
				}
			}

		} else if strings.HasPrefix(inp, "examine ") {
			// Examine a specific object in the room
			i := strings.Index(inp, " ")
			obj := strings.Trim(inp[i+1:], trim)
			desc, ok := room.Objects[obj]
			if !ok {
				PrintRed("Unknown object. Check capitalization.\n")
			} else {
				PrintPlain(fmt.Sprintf("%s\n", desc))
			}

		} else if inp == "answer" {
			// Input answer to current riddle
			inp := GetInput("Type the exact name of the object you think answers the riddle. You must be in the same room as the object!")
			if room.Id == riddles[riddlesI].Room && inp == riddles[riddlesI].A {
				riddlesI += 1
				if riddlesI == len(riddles) {
					break
				}
				fmt.Println(aurora.BrightGreen("Correct! Here's your next riddle:"))
				PrintPlain(fmt.Sprintf("%s\n", riddles[riddlesI].Q))
			} else {
				PrintRed("Incorrect. Check for typos or keep thinking.\n")
			}

		} else if inp == "where" {
			PrintRoom(room, rooms)
		} else if inp == "riddle" {
			PrintPlain(fmt.Sprintf("Current riddle: %s\n", riddles[riddlesI].Q))
		} else if inp == "quit" {
			break
		} else if inp == "help" {
			PrintHelp()
		} else {
			PrintRed("Unknown input. Type 'help' to see valid inputs.\n")
		}
	}

	if riddlesI == len(riddles) {
		PrintPlain("There's a note on the back of the kitchen rules:\n")
		PrintPlain("'Congratulations, you've completed the scavenger hunt! We lost one of your earbuds, by the way. Sorry. We're out getting replacemenets right now. If you're reading this, you must've completed everything faster than we expected.'\n\n")
		PrintPlain("You contemplate ")
		PrintRed("murder.\n\n")
		PrintPlain("THE END. Thanks for playing :)\n")
	}
}

func PrintRoom(room Room, rooms map[int]Room) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Print(aurora.Bold(room.Name))
	if len(room.Desc) != 0 {
		fmt.Print(aurora.Bold(": "), room.Desc)
	}
	fmt.Printf("\n\n")

	if len(room.Links) != 0 {
		fmt.Println(aurora.Bold("Connecting rooms:"))
		for dir, id := range room.Links {
			PrintPlain(fmt.Sprintf("- %s: %s\n", dir, rooms[id].Name))
		}
		fmt.Printf("\n")
	}

	if len(room.Objects) != 0 {
		fmt.Println(aurora.Bold("Objects:"))
		for obj, _ := range room.Objects {
			PrintPlain(fmt.Sprintf("- %s\n", obj))
		}
		fmt.Printf("\n")
	}
}

func PrintHelp() {
	fmt.Println(`Inputs are case-sensitive because I am lazy.
- To get room information: Type 'where'.
- To move: Type 'move [direction]'. For example, to move north, type 'move N'.
- To examine objects: Type 'examine [object]'. For example, to examine the sofa, type 'examine sofa'.
- To see riddle: Type 'riddle'.
- To answer the riddle: Type 'answer'. You will be prompted for your answer.
- To quit: Type 'quit'.`)
}

func PrintRed(s string) {
	fmt.Print(aurora.BrightRed(s))
}

func PrintPlain(s string) {
	fmt.Print(s)
}

func GetInput(prompt string) string {
	if len(prompt) != 0 {
		prompt += "\n"
	}
	PrintPlain(fmt.Sprintf("%s> ", prompt))
	inp, err := reader.ReadString('\n')
	if err != nil {
		PrintRed(fmt.Sprintf("Technical error occured while reading input string, go yell at L: %s\n", err))
		os.Exit(1)
	}
	return inp[:len(inp)-1] // Remove newline
}
