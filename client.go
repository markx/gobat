package main

type Client struct {
	conn *Conn
	ui   *UI
}

func NewClient(addr string) (*Client, error) {
	conn, err := Dial(addr)
	if err != nil {
		return nil, err
	}

	ui := NewUI()

	c := &Client{
		conn,
		ui,
	}

	return c, nil
}

func (c *Client) handleInput(content string) {
	c.conn.Write([]byte(content + "\n"))
}

func (c *Client) Run() error {
	c.ui.SetInputHandler(c.handleInput)

	errChan := make(chan error, 1)

	go func() {
		errChan <- c.ui.Run()
	}()

	for {
		select {
		case err := <-errChan:
			return err
		default:
			line, err := c.conn.ReadLine()
			if err != nil {
				c.ui.Stop()
				return err
			}
			c.handleMessage(NewMessage(line))
		}
	}
}

func (c *Client) handleMessage(m Message) {
	if m.hasTag("chat") {
		c.ui.SendToWindow("chat", m.Content)
		return
	}
	c.ui.SendToWindow("general", m.Content)
}
