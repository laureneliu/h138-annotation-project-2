package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Spec struct {
	Rooms map[int]Room
	Riddles []Riddle
}

type Room struct {
	Id int
	Name string
	Desc string
	Links map[string]int
	Objects map[string]string
}

type Riddle struct {
	Q string
	A string
	Room int
}

func ParseSpec(specFile string) (Spec, error) {
	spec := Spec{}
	file, err := ioutil.ReadFile(specFile)
	if (err != nil ) {
		return spec, err
	}
	warn := yaml.UnmarshalStrict(file, &spec)
	if warn != nil {
		return spec, warn
	}

	oppositeDirs := map[string]string {
		"N": "S",
		"S": "N",
		"W": "E",
		"E": "W",
	}

	// Make sure everything's connected properly
	// Also check that ids match
	for id, room := range spec.Rooms {
		if id != room.Id {
			return spec, fmt.Errorf("Id is wrong for room %d", id)
		}
		for dir, link := range room.Links {
			opp, ok := oppositeDirs[dir]
			if !ok {
				return spec, fmt.Errorf("Invalid direction %s in room %d", dir, id)
			}
			linkedRoom := spec.Rooms[link]
			if linkedRoom.Links[opp] != id {
				return spec, fmt.Errorf("Rooms %d and %d incorrectly linked", id, link)
			}
		}
	}

	return spec, nil
}
