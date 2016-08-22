package chatserver

type ChannelSub struct {
	UserName    string
	ChannelName string
}

func NewChannelSub(name string, ch string) *ChannelSub {
	return &ChannelSub{
		UserName:    name,
		ChannelName: ch,
	}
}
