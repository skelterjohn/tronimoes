package main

import (
	"strings"
	"unicode"
)

// Name parts (no spaces); combine from different lists to make full names.

// CreateName returns a full name by using each character in initials to choose
// a part from the corresponding list (R→NamesR, G→NamesG, C→NamesC, P→NamesP).
// Unknown letters are skipped. Same initials produce the same name.
func CreateName(initials string) string {
	var parts []string
	var h uint64
	for _, c := range initials {
		h = h*31 + uint64(unicode.ToUpper(c))
	}
	for i, c := range initials {
		c = unicode.ToUpper(c)
		var list []string
		switch c {
		case 'R':
			list = NamesR
		case 'G':
			list = NamesG
		case 'C':
			list = NamesC
		case 'P':
			list = NamesP
		default:
			continue
		}
		if len(list) == 0 {
			continue
		}
		idx := (h + uint64(i)*17) % uint64(len(list))
		parts = append(parts, list[idx])
	}
	return strings.Join(parts, " ")
}

var NamesR = []string{
	"Rumble", "Rinky", "Razzle", "Runcible", "Rutabaga", "Rumple",
	"Rocket", "Ragbag", "Rinkydink", "Rumbustious", "Razzmatazz",
	"Rigmarole", "Ripley", "Rumpelstiltskin", "Rolypoly", "Ruckus",
	"Ragamuffin", "Rickshaw", "Riproaring", "Rumbustion",
}

var NamesG = []string{
	"Gafferty", "Gorbison", "Guzzle", "Gobbledygook", "Gumption",
	"Gadzooks", "Gimcrack", "Gobsmacked", "Gubbins", "Gonk",
	"Gobbledy", "Gimlet", "Gadabout", "Goblin", "Gargle",
	"Gimcrackery", "Gobstopper", "Gadfly", "Gobbledegook", "Gazpacho",
}

var NamesC = []string{
	"Chicken", "Crumble", "Custard", "Cucumber", "Cantilever",
	"Curmudgeon", "Catawampus", "Collywobbles", "Cockamamie",
	"Curlicue", "Cahoots", "Cockatoo", "Cummerbund", "Catawumpus",
	"Cockalorum", "Curmudgeonly", "Collywobble", "Cahoot", "Catawampous",
	"Claptrap", "Corkscrew",
}

var NamesP = []string{
	"Pools", "Pulicker", "Pumpernickel", "Piffle", "Poppycock",
	"Pantaloon", "Pipsqueak", "Pifflepaff", "Prestidigitation",
	"Poppet", "Presto", "Pandemonium", "Pifflepaffle", "Prickle",
	"Persnickety", "Pomegranate", "Palaver", "Pettifogger", "Poodle",
	"Puddle", "Prickly", "Paraphernalia",
}
