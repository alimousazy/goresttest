package chatserver

import (
	"fmt"
)

type Hub struct {
	Sessions      map[string]*ChatSession
	Messages      chan Message
	SendChannel   chan Message
	AddSession    chan *ChatSession
	ListUsers     chan *ChatSession
	Channel       chan *Channel
	RemoveSession chan *ChatSession
	ChannelMap    map[string]*Channel
	UnSubChannel  chan *ChannelSub
	SubChannel    chan *ChannelSub
	ListChannel   chan *ChatSession
	logger        *Logger
}

func (h *Hub) start() {
	for {
		select {
		//Create user connection
		case ch := <-h.AddSession:
			ch.in <- "Welcom " + ch.User + "\n"
			h.Sessions[ch.User] = ch
			//Send message to a user
		case ch := <-h.Messages:
			if _, ok := h.Sessions[ch.from]; !ok {
				continue
			}
			if to, ok := h.Sessions[ch.to]; ok {
				h.Sessions[ch.from].in <- "Message deliverd\n"
				msg := fmt.Sprintf("@%s: %s\n", ch.from, ch.body)
				to.in <- msg
				h.logger.In <- msg
			} else {
				h.Sessions[ch.from].in <- "User doesn't exist.\n"
			}
			//List down connected users
		case ch := <-h.ListUsers:
			for _, sh := range h.Sessions {
				ch.in <- fmt.Sprintf("- %s\n", sh.User)
			}
			//Create a new channel
		case ch := <-h.Channel:
			h.ChannelMap[ch.Name] = ch
			for _, user := range ch.UserSessions {
				user.in <- fmt.Sprintf("Channel created %s\n", ch.Name)
			}
			//Close user session.
		case ch := <-h.RemoveSession:
			for _, channel := range h.ChannelMap {
				if _, ok := channel.UserSessions[ch.User]; ok {
					delete(channel.UserSessions, ch.User)
				}
			}
			delete(h.Sessions, ch.User)
			//Send message to a channel
		case ch := <-h.SendChannel:
			if _, ok := h.Sessions[ch.from]; !ok {
				continue
			}
			if to, ok := h.ChannelMap[ch.to]; ok {
				msg := fmt.Sprintf("#%s-@%s: %s\n", ch.to, ch.from, ch.body)
				for _, x := range to.UserSessions {
					x.in <- msg
				}
				h.logger.In <- msg
				h.Sessions[ch.from].in <- "Message deliverd\n"
			} else {
				h.Sessions[ch.from].in <- "Channel doesn't exist.\n"
			}
			//list created channels
		case ch := <-h.ListChannel:
			ch.in <- fmt.Sprintf("Channel list: \n")
			for name, _ := range h.ChannelMap {
				ch.in <- fmt.Sprintf("- %s\n", name)
			}
			//subscribe to a channel
		case ch := <-h.SubChannel:
			var userSes *ChatSession
			var isFound bool
			if userSes, isFound = h.Sessions[ch.UserName]; !isFound {
				continue
			}
			if channel, ok := h.ChannelMap[ch.ChannelName]; ok {
				channel.UserSessions[userSes.User] = userSes
				userSes.in <- fmt.Sprintf("You subscribed to a %s channel.\n", channel.Name)
			} else {
				userSes.in <- fmt.Sprintf("Channel (%s) doesn't exist.\n", ch.ChannelName)
			}
			//unsbuscribe from a channel
		case ch := <-h.UnSubChannel:
			var userSes *ChatSession
			var isFound bool
			if userSes, isFound = h.Sessions[ch.UserName]; !isFound {
				continue
			}
			if channel, ok := h.ChannelMap[ch.ChannelName]; ok {
				delete(channel.UserSessions, userSes.User)
				userSes.in <- fmt.Sprintf("You un-subscribed from a %s channel.\n", channel.Name)
			} else {
				userSes.in <- fmt.Sprintf("Channel (%s) doesn't exist.\n", ch.ChannelName)
			}
		}
	}
}

func NewHub(log *Logger) *Hub {
	return &Hub{
		Sessions:      make(map[string]*ChatSession),
		Messages:      make(chan Message),
		SendChannel:   make(chan Message),
		ListUsers:     make(chan *ChatSession),
		AddSession:    make(chan *ChatSession),
		RemoveSession: make(chan *ChatSession),
		ChannelMap:    make(map[string]*Channel),
		Channel:       make(chan *Channel),
		SubChannel:    make(chan *ChannelSub),
		UnSubChannel:  make(chan *ChannelSub),
		ListChannel:   make(chan *ChatSession),
		logger:        log,
	}
}
