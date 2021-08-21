package message

type RawMessage struct {
	Action string `json:"action"`
}

type Message struct {
	Action    string      `json:"action"`
	Payload   ChatPayload `json:"payload"`
	Stats     Stats       `json:"status"`
	Timestamp int         `json:"timestamp"`
	Type      string      `json:"type"`
}

type ChatPayload struct {
	ConnectionIdentifier string        `json:"connectionIdentifier"`
	EventIdentifier      string        `json:"eventIdentifier"`
	EventPayload         TwitchPayload `json:"eventPayload"`
	EventSourceId        int           `json:"eventSourceId"`
	EventTypeId          int           `json:"eventTypeId"`
	UserId               int           `json:"userId"`
}

type TwitchPayload struct {
	Author struct {
		Avatar        string `json:"avatar"`
		Color         string `json:"color"`
		DisplayName   string `json:"displayName"`
		Id            string `json:"id"`
		Name          string `json:"name"`
		SubscribedFor string `json:"number"`
		Badges        []struct {
			Title    string `json:"title"`
			ImageUrl string `json:"imageUrl"`
			ClickUrl string `json:"clickUrl"`
		} `json:"badges"`
	} `json:"author"`
	Bot              bool   `json:"bot"`
	Text             string `json:"text"`
	ContentModifiers struct {
		Me      bool `json:"me"`
		Whisper bool `json:"whisper"`
	} `json:"contentModifiers"`
	Replaces []struct {
		From    int    `json:"from"`
		To      int    `json:"to"`
		Type    string `json:"type"`
		Payload struct {
			Url string `json:"url"`
		} `json:"payload"`
	} `json:"replaces"`
}

type Stats struct {
	UserId            int    `json:"userId"`
	PlatformId        int    `json:"platformId"`
	ChannelId         int    `json:"channelId"`
	CreatedAt         int    `json:"createdAt"`
	UpdatedAt         int    `json:"updatedAt"`
	ChannelIdentifier string `json:"channelIdentifier"`
	EventIdentifier   string `json:"eventIdentifier"`
	ChannelViews      int    `json:"channelViews"`
	Followers         int    `json:"followers"`
	GameTitle         string `json:"gameTitle"`
	Online            bool   `json:"online"`
	StreamViews       int    `json:"streamViews"`
	Title             string `json:"title"`
	Viewers           int    `json:"viewers"`
}
