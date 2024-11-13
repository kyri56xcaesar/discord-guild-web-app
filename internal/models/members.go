package models

import (
	"kyri56xcaesar/discord_bots_app/internal/utils"
)

type Member struct {
	Guild         string    `json:"userguild"`
	Username      string    `json:"username"`
	Nick          string    `json:"nickname"`
	Avatar        string    `json:"avatarurl"`
	DisplayAvatar string    `json:"displayavatarurl"`
	Banner        string    `json:"bannerurl"`
	DisplayBanner string    `json:"displaybannerurl"`
	User_color    string    `json:"usercolor"`
	JoinedAt      string    `json:"joinedat"`
	Status        string    `json:"userstatus"`
	Roles         []Role    `json:"userroles"`
	Messages      []Message `json:"usermessages"`
	MsgCount      int       `json:"messagecount"`
	ID            int       `json:"userid"`
}

type Role struct {
	Role_name string `json:"rolename"`
	Color     string `json:"rolecolor"`
	ID        int    `json:"id"`
	UID       int    `json:"userid"`
}

type Message struct {
	Content   string `json:"content"`
	Channel   string `json:"channel"`
	CreatedAt string `json:"createdat"`
	ID        int    `json:"messageid"`
	UID       int    `json:"userid"`
}

func (mb *Member) PtrFieldsDB() []any {
	return []any{&mb.ID, &mb.Guild, &mb.Username, &mb.Nick, &mb.Avatar, &mb.DisplayAvatar, &mb.Banner, &mb.DisplayBanner, &mb.User_color, &mb.JoinedAt, &mb.Status, &mb.MsgCount}
}

func (r *Role) PtrFieldsDB() []any {
	return []any{&r.Role_name, &r.Color, &r.ID, &r.UID}
}

func (msg *Message) PtrFieldsDB() []any {
	return []any{&msg.Content, &msg.Channel, &msg.CreatedAt, &msg.ID, &msg.UID}
}

func (msg *Message) VerifyMessage() error {
	if !utils.IsValidUTF8String(msg.Content) {
		return &utils.FieldError{Field: "Content", Message: "must contain letters, numbers or symbols"}
	}

	if !utils.IsValidUTF8String(msg.Channel) {
		return &utils.FieldError{Field: "Channel", Message: "must contain letters, numbers or symbols"}
	}

	return nil
}

func (m *Member) VerifyMember() error {
	if !utils.IsValidUTF8String(m.Guild) {
		return &utils.FieldError{Field: "Guild", Message: "must contain letters, numbers and symbols"}
	}
	if !utils.IsValidUTF8String(m.Username) {
		return &utils.FieldError{Field: "Username", Message: "must contain only letters and numbers"}
	}
	if !utils.IsValidUTF8String(m.Nick) {
		return &utils.FieldError{Field: "Nick", Message: "must contain only letters and numbers"}
	}

	allowedStatuses := map[string]bool{"online": true, "offline": true, "idle": true, "dnd": true}
	if !allowedStatuses[m.Status] {
		return &utils.FieldError{Field: "Status", Message: "invalid status value"}
	}

	allowedColors := map[string]bool{"red": true, "blue": true, "yellow": true, "green": true, "black": true, "white": true, "pink": true, "purple": true}
	if !utils.IsValidColor(m.User_color, allowedColors) {
		return &utils.FieldError{Field: "User_color", Message: "must be a valid hex or rgba color"}
	}

	if !utils.IsValidURLOrBase64(m.Avatar) {
		return &utils.FieldError{Field: "Avatar", Message: "must be a valid URL or base64 string"}
	}

	if !utils.IsValidURLOrBase64(m.DisplayAvatar) {
		return &utils.FieldError{Field: "DisplayAvatar", Message: "must be a valid URL or base64 string"}
	}
	if !utils.IsValidURLOrBase64(m.Banner) {
		return &utils.FieldError{Field: "Banner", Message: "must be a valid URL or base64 string"}
	}
	if !utils.IsValidURLOrBase64(m.DisplayBanner) {
		return &utils.FieldError{Field: "DisplayBanner", Message: "must be a valid URL or base64 string"}
	}

	for _, message := range m.Messages {
		if !utils.IsValidUTF8String(message.Content) {
			return &utils.FieldError{Field: "Message", Message: "must contain letters, numbers or symbols"}
		}
	}

	return nil
}
