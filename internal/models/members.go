package models

type Member struct {
	Guild         string   `json:"userguild"`
	ID            int      `json:"userid"`
	Username      string   `json:"username"`
	Nick          string   `json:"nickname"`
	Avatar        string   `json:"avatarurl"`
	DisplayAvatar string   `json:"displayavatarurl"`
	Banner        string   `json:"bannerurl"`
	DisplayBanner string   `json:"displaybannerurl"`
	User_color    string   `json:"usercolor"`
	JoinedAt      string   `json:"joinedat"`
	Status        string   `json:"userstatus"`
	Roles         []Role   `json:"userroles"`
	Messages      []string `json:"usermessages"`
	MsgCount      int      `json:"messagecount"`
}

type Role struct {
	Role_name string `json:"role_name"`
	Color     string `json:"role_color"`
}

type Bot struct {
	Member
	IsBot bool `json:"bot"`
}
