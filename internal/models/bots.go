package models

import (
	"kyri56xcaesar/discord_bots_app/internal/utils"
)

type Bot struct {
	Guild     string `json:"botguild"`
	Name      string `json:"botname"`
	Avatar    string `json:"avatarurl"`
	Banner    string `json:"bannerurl"`
	CreatedAt string `json:"createdat"`
	Author    string `json:"author"`
	Status    string `json:"botstatus"`
	Lines     []Line `json:"linewords"`
	ID        int    `json:"botid"`
	IsSinger  bool   `json:"isSinger"`
}

type Line struct {
	Phrase    string `json:"phrase"`
	Author    string `json:"author"`
	To        string `json:"toid"`
	LineType  string `json:"ltype"`
	CreatedAt string `json:"createdat"`
	ID        int    `json:"lineid"`
	BID       int    `json:"bid"`
}

func (l *Line) PtrsFieldsDB() []any {
	return []any{&l.ID, &l.BID, &l.Phrase, &l.Author, &l.To, &l.LineType, &l.CreatedAt}
}

func (b *Bot) PtrsFieldsDB() []any {
	return []any{&b.ID, &b.Guild, &b.Name, &b.Avatar, &b.Banner, &b.CreatedAt, &b.Author, &b.Status, &b.IsSinger}
}

func (l *Line) VerifyLine() error {
	if !utils.IsValidUTF8String(l.Phrase) {
		return &utils.FieldError{Field: "Phrase", Message: "must contain letters, numbers or symbols"}
	}
	if !utils.IsValidUTF8String(l.Author) {
		return &utils.FieldError{Field: "Author", Message: "must contain letters, numbers or symbols"}
	}
	if !utils.IsValidUTF8String(l.To) {
		return &utils.FieldError{Field: "To", Message: "must contain letters, numbers or symbols"}
	}
	if !utils.IsValidUTF8String(l.LineType) {
		return &utils.FieldError{Field: "LineType", Message: "must contain letters, numbers or symbols"}
	}

	return nil
}

func (b *Bot) VerifyBot() error {
	if !utils.IsValidUTF8String(b.Name) {
		return &utils.FieldError{Field: "Name", Message: "must contain letters, numbers or symbols"}
	}

	if !utils.IsValidUTF8String(b.Guild) {
		return &utils.FieldError{Field: "Guild", Message: "must contain letters, numbers or symbols"}
	}

	if !utils.IsValidUTF8String(b.Author) {
		return &utils.FieldError{Field: "Author", Message: "must contain letters, numbers or symbols"}
	}

	if !utils.IsValidUTF8String(b.Status) {
		return &utils.FieldError{Field: "Status", Message: "must contain letters, numbers or symbols"}
	}

	if !utils.IsValidURLOrBase64(b.Avatar) {
		return &utils.FieldError{Field: "Avatar", Message: "must be a valid URL or base64 string"}
	}

	if !utils.IsValidURLOrBase64(b.Banner) {
		return &utils.FieldError{Field: "Banner", Message: "must be a valid URL or base64 string"}
	}

	for _, l := range b.Lines {
		if err := l.VerifyLine(); err != nil {
			return err
		}
	}

	return nil
}
