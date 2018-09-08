package main

import "regexp"

type Message struct {
	Content string
	Tags    []string
}

func NewMessage(content string) Message {
	m := Message{
		Content: content,
	}

	return handleChat(m)
}

func handleChat(m Message) Message {
	r := regexp.MustCompile(`\s[[({]\w+[\])}]:\s`)
	if r.MatchString(m.Content) {
		m.Tags = append(m.Tags, "chat")
	}

	return m
}

func (m Message) hasTag(tag string) bool {
	for _, t := range m.Tags {
		if t == tag {
			return true
		}
	}
	return false
}
