package main

type Triggers struct {
	triggers []Trigger
}

func (ts *Triggers) Register(t Trigger) {
	ts.triggers = append(ts.triggers, t)
}

func (ts *Triggers) RegisterFunc(tf func(Game, *Message, *Client)) {
	ts.triggers = append(ts.triggers, TriggerFunc(tf))
}

func (t *Triggers) Apply(state Game, m *Message, c *Client) {
	for _, trigger := range t.triggers {
		trigger.Apply(state, m, c)
	}
}

type Trigger interface {
	Apply(Game, *Message, *Client)
}

type TriggerFunc func(Game, *Message, *Client)

func (f TriggerFunc) Apply(s Game, m *Message, c *Client) {
	f(s, m, c)
}

func NewTriggers() *Triggers {
	triggers := &Triggers{}

	triggers.Register(TriggerFunc(healSelf))

	return triggers
}

func healSelf(s Game, _ *Message, client *Client) {
	if s.Player.HPTotal-s.Player.HP > 70 {
		client.Send("clw")
	}
}
