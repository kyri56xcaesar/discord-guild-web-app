package user

type User struct {
	Guild         string   `json:"string"`
	ID            int      `json:"id"`
	User          string   `json:"user"`
	Nick          string   `json:"nick"`
	Avatar        string   `json:"avatar"`
	DisplayAvatar string   `json:"display_avatar"`
	Banner        string   `json:"banner"`
	DisplayBanner string   `json:"display_banner"`
	User_color    string   `json:"user_color"`
	JoinedAt      string   `json:"joined_at"`
	Status        string   `json:"status"`
	Roles         []Role   `json:"roles"`
	Messages      []string `json:"messages"`
	MsgCount      int      `json:"msg_count"`
}

type Role struct {
	Role_name string `json:"role_name"`
	Color     string `json:"role_color"`
}

type Member struct {
	User
}

type Bot struct {
	User
	IsBot bool `json:"bot"`
}
