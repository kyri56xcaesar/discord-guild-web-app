package models

type Bot struct {
	Guild     string `json:"botguild"`
	ID        int    `json:"botid"`
	Name      string `json:"botname"`
	Avatar    string `json:"avatarurl"`
	Banner    string `json:"bannerurl"`
	CreatedAt string `json:"createdat"`
	Author    string `json:"author"`
	Status    string `json:"botstatus"`
	IsSinger  bool   `json:"isSinger"`
	Lines     []Line `json:"linewords"`
}

type Line struct {
	ID        int    `json:"lineid"`
	BID       int    `json:"bid"`
	Phrase    string `json:"phrase"`
	Author    string `json:"author"`
	To        string `json:"toid"`
	LineType  string `json:"ltype"`
	CreatedAt string `json:"createdat"`
}

func (l *Line) VerifyLine() error {
	if !isValidUTF8String(l.Phrase) {
		return &FieldError{Field: "Phrase", Message: "must contain letters, numbers or symbols"}
	}
	if !isValidUTF8String(l.Author) {
		return &FieldError{Field: "Author", Message: "must contain letters, numbers or symbols"}
	}
	if !isValidUTF8String(l.To) {
		return &FieldError{Field: "To", Message: "must contain letters, numbers or symbols"}
	}
	if !isValidUTF8String(l.LineType) {
		return &FieldError{Field: "LineType", Message: "must contain letters, numbers or symbols"}
	}

	return nil
}

func (b *Bot) VerifyBot() error {
	if !isValidUTF8String(b.Name) {
		return &FieldError{Field: "Name", Message: "must contain letters, numbers or symbols"}
	}

	if !isValidUTF8String(b.Guild) {
		return &FieldError{Field: "Guild", Message: "must contain letters, numbers or symbols"}
	}

	if !isValidUTF8String(b.Author) {
		return &FieldError{Field: "Author", Message: "must contain letters, numbers or symbols"}
	}

	if !isValidUTF8String(b.Status) {
		return &FieldError{Field: "Status", Message: "must contain letters, numbers or symbols"}
	}

	if !isValidURLOrBase64(b.Avatar) {
		return &FieldError{Field: "Avatar", Message: "must be a valid URL or base64 string"}
	}

	if !isValidURLOrBase64(b.Banner) {
		return &FieldError{Field: "Banner", Message: "must be a valid URL or base64 string"}
	}

	for _, l := range b.Lines {
		if err := l.VerifyLine(); err != nil {
			return err

		}
	}

	return nil
}
