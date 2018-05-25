package main

type Message struct {
	Content string
	Tags    []string
}

func NewMessage(content string) Message {
	return Message{
		Content: content,
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
