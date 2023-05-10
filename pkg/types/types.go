package types

type QueryParams struct {
	ChatId   string `form:"chat_id"`
	ChatSign string `form:"chat_sign"`
	UserId   string `form:"user_id"`
	UserSign string `form:"user_sign"`
}

type BxUserInfo struct {
	Id     string `json:"ID"`
	Active bool   `json:"ACTIVE"`
}

type BxResult struct {
	Result []BxUserInfo `json:"result"`
}

type KickMessage struct {
	UserID int    `json:"user_id"`
	Type   string `json:"type"`
}

type ChatRecordShort struct {
	ID          int64  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
}

type ChatsResult struct {
	All         []ChatRecordShort
	City        []ChatRecordShort
	Dev         []ChatRecordShort
	Recommended []ChatRecordShort
}

type ReplyAction struct {
	Type   string `json:"type"` // skip or stay
	LeadId int    `json:"lead_id"`
	UserId int    `json:"user_id"`
}

type EmployersInfo struct {
	Done  chan struct{}
	Users []UserInfo
}

type UserInfo struct {
	UserId   int
	Name     string
	LastName string
	Chats 	 []string
}
