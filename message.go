package main

import "regexp"

type Message struct {
	Content string
	Tags    []string
}

func NewMessage(content string) Message {
	return Message{
		Content: content,
		Tags:    tagContent(content),
	}
}

func (m Message) hasTag(tag string) bool {
	for _, t := range m.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

type PlayerStatus struct {
	Name  string
	Level int
	HP    int
	MP    int
	AP    int
}

type TeammateStatus struct {
	PlayerStatus
	Row    int
	Column int
}

func tagContent(c string) []string {
	var tags []string

	r := regexp.MustCompile(`\s[[({]\w+[\])}]:\s`)
	if r.MatchString(c) {
		tags = append(tags, "chat")
	}

	return tags
}
