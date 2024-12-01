package models

import (
	"log"
	"reflect"

	"kyri56xcaesar/discord-guild-web-app/internal/utils"
)

type Bot struct {
	Guild     string `json:"guild"`
	Username  string `json:"username"`
	Avatarurl string `json:"avatarurl"`
	Bannerurl string `json:"bannerurl"`
	Createdat string `json:"createdat"`
	Author    string `json:"author"`
	Status    string `json:"status"`
	Lines     []Line `json:"linewords"`
	Id        int    `json:"id"`
	Issinger  bool   `json:"issinger"`
}

type Line struct {
	Phrase    string `json:"phrase"`
	Author    string `json:"author"`
	Toid      string `json:"toid"`
	Ltype     string `json:"ltype"`
	Createdat string `json:"createdat"`
	Id        int    `json:"id"`
	Bid       int    `json:"bid"`
}

func (l *Line) PtrFieldsDB() []any {
	return []any{&l.Id, &l.Bid, &l.Phrase, &l.Author, &l.Toid, &l.Ltype, &l.Createdat}
}

func (l *Line) PtrFieldsSpecific(fields []string) []any {
	if len(fields) == 0 {
		return l.PtrFieldsDB()
	}

	ptrFields := []any{}

	value := reflect.ValueOf(l).Elem()
	typ := reflect.TypeOf(*l)

	for _, field := range fields {
		structField, found := typ.FieldByName(utils.ToUpperFirstLetter(field))
		if !found {
			log.Printf("Field %s does not exist in Line struct", field)
			continue
		}

		fieldValue := value.FieldByIndex(structField.Index)

		if fieldValue.CanAddr() {
			ptrFields = append(ptrFields, fieldValue.Addr().Interface())
		} else {
			log.Printf("Field %s is not addressable", field)
		}
	}

	return ptrFields
}

func (b *Bot) PtrFieldsDB() []any {
	return []any{&b.Id, &b.Guild, &b.Username, &b.Avatarurl, &b.Bannerurl, &b.Createdat, &b.Author, &b.Status, &b.Issinger}
}

func (b *Bot) PtrFieldsSpecific(fields []string) []any {
	if len(fields) == 0 {
		return b.PtrFieldsDB()
	}

	ptrFields := []any{}

	// Use reflection to access struct fields dynamically
	value := reflect.ValueOf(b).Elem()
	typ := reflect.TypeOf(*b)

	for _, field := range fields {
		// Find the field by name
		structField, found := typ.FieldByName(utils.ToUpperFirstLetter(field))
		if !found {
			log.Printf("Field %s does not exist in Bot struct", field)
			continue // Skip non-existing fields
		}

		// Get the pointer to the field's value
		fieldValue := value.FieldByIndex(structField.Index)

		// Append the address of the field value to the slice
		if fieldValue.CanAddr() {
			ptrFields = append(ptrFields, fieldValue.Addr().Interface())
		} else {
			log.Printf("Field %s is not addressable", field)
		}
	}

	return ptrFields
}

func (l *Line) VerifyLine() error {
	if !utils.IsValidUTF8String(l.Phrase) {
		return &utils.FieldError{Field: "Phrase", Message: "must contain letters, numbers or symbols"}
	}
	if !utils.IsValidUTF8String(l.Author) {
		return &utils.FieldError{Field: "Author", Message: "must contain letters, numbers or symbols"}
	}
	if !utils.IsValidUTF8String(l.Toid) {
		return &utils.FieldError{Field: "To", Message: "must contain letters, numbers or symbols"}
	}
	if !utils.IsValidUTF8String(l.Ltype) {
		return &utils.FieldError{Field: "LineType", Message: "must contain letters, numbers or symbols"}
	}

	return nil
}

func (b *Bot) VerifyBot() error {
	if !utils.IsValidUTF8String(b.Username) {
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

	if !utils.IsValidURLOrBase64(b.Avatarurl) {
		return &utils.FieldError{Field: "Avatar", Message: "must be a valid URL or base64 string"}
	}

	if !utils.IsValidURLOrBase64(b.Bannerurl) {
		return &utils.FieldError{Field: "Banner", Message: "must be a valid URL or base64 string"}
	}

	for _, l := range b.Lines {
		if err := l.VerifyLine(); err != nil {
			return err
		}
	}

	return nil
}
