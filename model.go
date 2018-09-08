package main

import (
	"regexp"
	"strconv"
)

type Game struct {
	//Messages []Message
	//channels map[string][]Message
	Player Player
}

type Player struct {
	Name    string
	Level   int
	HP      int
	HPTotal int
	SP      int
	SPTotal int
	EP      int
	EPTotal int
	EXP     int
}

func PlayerReducer(state Player, msg string) Player {
	// Hp:388/388 Sp:557/557 Ep:239/239 Exp:1133595
	r := regexp.MustCompile(`^Hp:(\d+)/(\d+) Sp:(\d+)/(\d+) Ep:(\d+)/(\d+) Exp:(\d+) >`)
	if !r.MatchString(msg) {
		return state
	}

	match := r.FindStringSubmatch(msg)[1:]
	state.HP, _ = strconv.Atoi(match[0])
	state.HPTotal, _ = strconv.Atoi(match[1])
	state.SP, _ = strconv.Atoi(match[2])
	state.SPTotal, _ = strconv.Atoi(match[3])
	state.EP, _ = strconv.Atoi(match[4])
	state.EPTotal, _ = strconv.Atoi(match[5])

	if len(match) > 7 {
		state.EXP, _ = strconv.Atoi(match[6])
	}

	return state
}

type Teammate struct {
	Player
	Row    int
	Column int
}
