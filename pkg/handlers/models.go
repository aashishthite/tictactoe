package handlers

type StartGameRequest struct {
	Player1ID string `json:"player_1"`
	Player2ID string `json:"player_2"`
}

type MoveRequest struct {
	PlayerID string `json:"player"`
	Position string `json:"position"`
}

type SlashCommand struct {
	Token       string
	TeamId      string
	TeamDomain  string
	ChannelId   string
	ChannelName string
	UserId      string
	UserName    string
	Command     string
	Text        string
	ResponseURL string
}

type Response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

type SlashCommandMessage struct {
	ResponseType string       `json:"response_type,omitempty"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments,omitempty"`
}

// Attachment contains all the information for an attachment
type Attachment struct {
	Color    string `json:"color,omitempty"`
	Fallback string `json:"fallback"`

	AuthorName    string `json:"author_name,omitempty"`
	AuthorSubname string `json:"author_subname,omitempty"`
	AuthorLink    string `json:"author_link,omitempty"`
	AuthorIcon    string `json:"author_icon,omitempty"`

	Title     string `json:"title,omitempty"`
	TitleLink string `json:"title_link,omitempty"`
	Pretext   string `json:"pretext,omitempty"`
	Text      string `json:"text"`

	ImageURL string `json:"image_url,omitempty"`
	ThumbURL string `json:"thumb_url,omitempty"`

	Fields     []AttachmentField `json:"fields,omitempty"`
	MarkdownIn []string          `json:"mrkdwn_in,omitempty"`
}

// AttachmentField contains information for an attachment field
// An Attachment can contain multiple of these
type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type userResponseFull struct {
	Members []User `json:"members,omitempty"`
}

// User contains all the information of a user
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
