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

func (m *Member) VerifyMember() error {
	if !isValidUTF8String(m.Guild) {
		return &FieldError{Field: "Guild", Message: "must contain letters, numbers and symbols"}
	}
	if !isValidUTF8String(m.Username) {
		return &FieldError{Field: "Username", Message: "must contain only letters and numbers"}
	}
	if !isValidUTF8String(m.Nick) {
		return &FieldError{Field: "Nick", Message: "must contain only letters and numbers"}
	}

	allowedStatuses := map[string]bool{"online": true, "offline": true, "idle": true, "dnd": true}
	if !allowedStatuses[m.Status] {
		return &FieldError{Field: "Status", Message: "invalid status value"}
	}

	allowedColors := map[string]bool{"red": true, "blue": true, "yellow": true, "green": true, "black": true, "white": true, "pink": true, "purple": true}
	if !isValidColor(m.User_color, allowedColors) {
		return &FieldError{Field: "User_color", Message: "must be a valid hex or rgba color"}
	}

	if !isValidURLOrBase64(m.Avatar) {
		return &FieldError{Field: "Avatar", Message: "must be a valid URL or base64 string"}
	}

	if !isValidURLOrBase64(m.DisplayAvatar) {
		return &FieldError{Field: "DisplayAvatar", Message: "must be a valid URL or base64 string"}
	}
	if !isValidURLOrBase64(m.Banner) {
		return &FieldError{Field: "Banner", Message: "must be a valid URL or base64 string"}
	}
	if !isValidURLOrBase64(m.DisplayBanner) {
		return &FieldError{Field: "DisplayBanner", Message: "must be a valid URL or base64 string"}
	}

	for _, message := range m.Messages {
		if !isValidUTF8String(message) {
			return &FieldError{Field: "Message", Message: "must contain letters, numbers or symbols"}

		}
	}

	return nil
}
