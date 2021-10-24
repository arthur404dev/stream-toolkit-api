package message

type RawMessage struct {
	Action string `json:"action,omitempty"`
}

type Message struct {
	Action    string      `json:"action,omitempty"`
	Payload   ChatPayload `json:"payload,omitempty"`
	Stats     Stats       `json:"stats,omitempty"`
	Timestamp int         `json:"timestamp,omitempty"`
	Type      string      `json:"type,omitempty"`
}

type ChatPayload struct {
	ConnectionIdentifier string             `json:"connectionIdentifier,omitempty"`
	EventIdentifier      string             `json:"eventIdentifier,omitempty"`
	EventPayload         ChatMessagePayload `json:"eventPayload,omitempty"`
	EventSourceId        int                `json:"eventSourceId,omitempty"`
	EventTypeId          int                `json:"eventTypeId,omitempty"`
	UserId               int                `json:"userId,omitempty"`
}

type ChatMessagePayload struct {
	Author struct {
		Avatar        string `json:"avatar,omitempty"`
		AvatarUrl     string `json:"avatarUrl,omitempty"`
		Picture       string `json:"picture,omitempty"`
		Color         string `json:"color,omitempty"`
		DisplayName   string `json:"displayName,omitempty"`
		Id            string `json:"id,omitempty"`
		Name          string `json:"name,omitempty"`
		SubscribedFor int    `json:"subscribedFor,omitempty"`
		Badges        []struct {
			Title    string `json:"title,omitempty"`
			ImageUrl string `json:"imageUrl,omitempty"`
			ClickUrl string `json:"clickUrl,omitempty"`
		} `json:"badges,omitempty"`
	} `json:"author,omitempty"`
	Bot              bool   `json:"bot,omitempty"`
	Text             string `json:"text,omitempty"`
	ContentModifiers struct {
		Me      bool `json:"me,omitempty"`
		Whisper bool `json:"whisper,omitempty"`
	} `json:"contentModifiers,omitempty"`
	Replaces []struct {
		From    int    `json:"from,omitempty"`
		To      int    `json:"to,omitempty"`
		Type    string `json:"type,omitempty"`
		Payload struct {
			Url string `json:"url,omitempty"`
		} `json:"payload,omitempty"`
	} `json:"replaces,omitempty"`
}

type Stats struct {
	UserId            int    `json:"userId,omitempty"`
	PlatformId        int    `json:"platformId,omitempty"`
	ChannelId         int    `json:"channelId,omitempty"`
	CreatedAt         int    `json:"createdAt,omitempty"`
	UpdatedAt         int    `json:"updatedAt,omitempty"`
	ChannelIdentifier string `json:"channelIdentifier,omitempty"`
	EventIdentifier   string `json:"eventIdentifier,omitempty"`
	ChannelViews      int    `json:"channelViews,omitempty"`
	Followers         int    `json:"followers,omitempty"`
	GameTitle         string `json:"gameTitle,omitempty"`
	Online            bool   `json:"online,omitempty"`
	StreamViews       int    `json:"streamViews,omitempty"`
	Title             string `json:"title,omitempty"`
	Viewers           int    `json:"viewers,omitempty"`
}
