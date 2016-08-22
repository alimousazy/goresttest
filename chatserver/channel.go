package chatserver

type Channel struct {
	Name         string
	UserSessions map[string]*ChatSession
}

func NewChannel(name string) *Channel {
	return &Channel{
		Name:         name,
		UserSessions: make(map[string]*ChatSession),
	}
}
