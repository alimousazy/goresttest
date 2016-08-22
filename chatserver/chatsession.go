package chatserver

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
)

type commandHandler func([]string) bool

type ChatSession struct {
	conn   net.Conn
	User   string
	Reader *bufio.Reader
	in     chan string
	hub    *Hub
}

func (c *ChatSession) Start() {
	defer c.conn.Close()
	defer close(c.in)
	c.Reader = bufio.NewReader(c.conn)
	go c.StartInput()
	for {
		line, err := c.Reader.ReadString('\n')
		if err != nil {
			fmt.Println("%v", err)
		}
		if c.commandParser(line) != true {
			break
		}
	}
}
func (c *ChatSession) StartInput() {
	writer := bufio.NewWriter(c.conn)
	for resp := range c.in {
		_, err := writer.WriteString(resp)
		if err != nil {
			fmt.Println("%v", err)
		}
		writer.Flush()
	}
}

func (c *ChatSession) commandParser(line string) bool {
	if strings.TrimSpace(line) == "" {
		return true
	}
	pattern := map[*regexp.Regexp]commandHandler{
		regexp.MustCompile("/connect ([a-zA-Z_]+)"): func(list []string) bool {
			if c.User != "" {
				c.in <- fmt.Sprintf("You already connected as user(%s), please close session to start again.", c.User)
			}
			c.User = list[0]
			c.hub.AddSession <- c
			return true
		},
		regexp.MustCompile("/message #([a-zA-Z]+) (.*)"): func(list []string) bool {
			if c.User == "" {
				c.in <- "Please connect as user to send message.\n"
			}
			c.hub.SendChannel <- Message{
				from: c.User,
				to:   list[0],
				body: list[1],
			}
			return true
		},

		regexp.MustCompile("/message ([a-zA-Z]+) (.*)"): func(list []string) bool {
			if c.User == "" {
				c.in <- "Please connect as user to send message.\n"
				return true
			}
			c.hub.Messages <- Message{
				from: c.User,
				to:   list[0],
				body: list[1],
			}
			return true
		},
		regexp.MustCompile("/close"): func(list []string) bool {
			if c.User == "" {
				c.in <- "Please connect as user to send message.\n"
				return true
			}
			c.hub.RemoveSession <- c
			c.in <- "Closing connection"
			return false
		},
		regexp.MustCompile("/list"): func(list []string) bool {
			if c.User == "" {
				c.in <- "Please connect as user to send message.\n"
				return true
			}
			c.hub.ListUsers <- c
			return true
		},
		regexp.MustCompile("/channel create ([a-xA-Z_]+)"): func(list []string) bool {
			if c.User == "" {
				c.in <- "Please connect as user to send message.\n"
				return true
			}
			ch := NewChannel(list[0])
			ch.UserSessions[c.User] = c
			c.hub.Channel <- ch
			return true
		},
		regexp.MustCompile("/channel list"): func(list []string) bool {
			if c.User == "" {
				c.in <- "Please connect as user to send message.\n"
				return true
			}
			c.hub.ListChannel <- c
			return true
		},
		regexp.MustCompile("/channel unsub ([a-zA-Z_]+)"): func(list []string) bool {
			if c.User == "" {
				c.in <- "Please connect as user to send message.\n"
				return true
			}
			sub := NewChannelSub(c.User, list[0])
			c.hub.UnSubChannel <- sub
			return true
		},
		regexp.MustCompile("/channel sub ([a-zA-Z_]+)"): func(list []string) bool {
			if c.User == "" {
				c.in <- "Please connect as user to send message.\n"
				return true
			}
			sub := NewChannelSub(c.User, list[0])
			c.hub.SubChannel <- sub
			return true
		},
		regexp.MustCompile("/help"): func(list []string) bool {
			cmd := `
Avalible commands:
  /connect [username] #Must be called at start of a session, Ex: /connect foo
  /list #List down connected user, Ex: /list
  /channel create [chanel-name] create a channel, Ex: /channel create foo
  /channel list #List down created channels, Ex: /channel list
  /channel sub #subscribe to a channel
  /channel unsub #unsubscribe from a channel
  /message #[channel-name] [message-text] #Send message to a channel 
  /message [user-name] [message-text] #Send message to a user 
      
`
			c.in <- cmd
			return true
		},
	}
	for key, callback := range pattern {
		if key.MatchString(line) {
			match := key.FindAllStringSubmatch(line, -1)
			return callback(match[0][1:])
		}
	}
	c.in <- "Command not found, try help for list of command.\n"
	return true
}

func NewChatSession(con net.Conn, hub *Hub) *ChatSession {
	return &ChatSession{
		conn: con,
		hub:  hub,
		in:   make(chan string),
	}
}
