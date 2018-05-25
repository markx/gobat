package main

import (
	"regexp"
	"strconv"
)

type Triggers struct {
	triggers []Trigger
}

func (ts *Triggers) Register(t Trigger) {
	ts.triggers = append(ts.triggers, t)
}

func (t *Triggers) Match(m *Message, c *Client) {
	for _, trigger := range t.triggers {
		trigger.Match(m, c)
	}
}

type Trigger interface {
	Match(*Message, *Client)
}

type TriggerFunc func(*Message, *Client)

func (f TriggerFunc) Match(m *Message, c *Client) {
	f(m, c)
}

func NewTriggers() *Triggers {
	triggers := &Triggers{}

	triggers.Register(TriggerFunc(handleChat))
	triggers.Register(&Character{})

	return triggers
}

func handleChat(msg *Message, client *Client) {
	r := regexp.MustCompile(`\s[[({]\w+[\])}]:\s`)
	if r.MatchString(msg.Content) {
		msg.Tags = append(msg.Tags, "chat")
	}
}

type Character struct {
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

func (c *Character) Match(msg *Message, client *Client) {
	c.handlePrompt(msg, client)
	c.healSelf(msg, client)
}

func (c *Character) handlePrompt(msg *Message, _ *Client) {
	// Hp:388/388 Sp:557/557 Ep:239/239 Exp:1133595
	r := regexp.MustCompile(`^Hp:(\d+)/(\d+) Sp:(\d+)/(\d+) Ep:(\d+)/(\d+) Exp:(\d+) >`)
	if !r.MatchString(msg.Content) {
		return
	}

	msg.Tags = append(msg.Tags, "prompt")

	match := r.FindStringSubmatch(msg.Content)[1:]
	c.HP, _ = strconv.Atoi(match[0])
	c.HPTotal, _ = strconv.Atoi(match[1])
	c.SP, _ = strconv.Atoi(match[2])
	c.SPTotal, _ = strconv.Atoi(match[3])
	c.EP, _ = strconv.Atoi(match[4])
	c.EPTotal, _ = strconv.Atoi(match[5])

	if len(match) > 7 {
		c.EXP, _ = strconv.Atoi(match[6])
	}
}

func (c *Character) healSelf(_ *Message, client *Client) {
	if c.HPTotal-c.HP > 70 {
		client.Send("clw")
	}
}

type Teammate struct {
	Character
	Row    int
	Column int
}
