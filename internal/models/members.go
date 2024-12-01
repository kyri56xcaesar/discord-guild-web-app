package models

import (
	"log"
	"reflect"

	"kyri56xcaesar/discord-guild-web-app/internal/utils"
)

type Member struct {
	Guild            string    `json:"guild"`
	Username         string    `json:"username"`
	Nickname         string    `json:"nickname"`
	Avatarurl        string    `json:"avatarurl"`
	Displayavatarurl string    `json:"displayavatarurl"`
	Bannerurl        string    `json:"bannerurl"`
	Displaybannerurl string    `json:"displaybannerurl"`
	Usercolor        string    `json:"usercolor"`
	Joinedat         string    `json:"joinedat"`
	Status           string    `json:"status"`
	Userroles        []Role    `json:"userroles"`
	Usermessages     []Message `json:"usermessages"`
	Msgcount         int       `json:"msgcount"`
	Id               int       `json:"id"`
}

type Role struct {
	Rolename  string `json:"rolename"`
	Rolecolor string `json:"rolecolor"`
	Id        int    `json:"id"`
	Userid    int    `json:"userid"`
}

type Message struct {
	Content   string `json:"content"`
	Channel   string `json:"channel"`
	Createdat string `json:"createdat"`
	Messageid int    `json:"messageid"`
	Userid    int    `json:"userid"`
}

func FilterFields(data []any, cols []string) []map[string]interface{} {
	var filtered []map[string]interface{}

	for _, member := range data {
		memberFiltered := make(map[string]interface{}, len(cols))
		value := reflect.ValueOf(member)
		typ := reflect.TypeOf(member)

		for _, col := range cols {
			field, found := typ.FieldByName(col)
			if !found {
				log.Printf("Field %s does not exist in Member struct", col)
				continue
			}

			fieldValue := value.FieldByIndex(field.Index)
			memberFiltered[col] = fieldValue.Interface()
		}

		filtered = append(filtered, memberFiltered)
	}

	return filtered
}

func (mb *Member) PtrFieldsDB() []any {
	return []any{&mb.Id, &mb.Guild, &mb.Username, &mb.Nickname, &mb.Avatarurl, &mb.Displayavatarurl, &mb.Bannerurl, &mb.Displaybannerurl, &mb.Usercolor, &mb.Joinedat, &mb.Status, &mb.Msgcount}
}

func (mb *Member) PtrFieldsSpecific(fields []string) []any {
	if len(fields) == 0 {
		return mb.PtrFieldsDB()
	}

	ptrFields := []any{}

	// Use reflection to access struct fields dynamically
	value := reflect.ValueOf(mb).Elem()
	typ := reflect.TypeOf(*mb)

	for _, field := range fields {
		// Find the field by name
		structField, found := typ.FieldByName(utils.ToUpperFirstLetter(field))
		if !found {
			log.Printf("Field %s does not exist in Member struct", field)
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

func (r *Role) PtrFieldsDB() []any {
	return []any{&r.Rolename, &r.Rolecolor, &r.Id, &r.Userid}
}

func (msg *Message) PtrFieldsDB() []any {
	return []any{&msg.Content, &msg.Channel, &msg.Createdat, &msg.Messageid, &msg.Userid}
}

func (role *Role) VerifyRole() error {
	if !utils.IsValidUTF8String(role.Rolename) {
		return &utils.FieldError{Field: "Role_name", Message: "must contain letters, numbers or symbols"}
	}

	// 	if !isValidColor(role.Color) {
	// 		return &FieldError{Field: "Color", Message: "must contain color representation"}
	// 	}
	return nil
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
	if !utils.IsValidUTF8String(m.Nickname) {
		return &utils.FieldError{Field: "Nick", Message: "must contain only letters and numbers"}
	}

	allowedStatuses := map[string]bool{"online": true, "offline": true, "idle": true, "dnd": true}
	if !allowedStatuses[m.Status] {
		return &utils.FieldError{Field: "Status", Message: "invalid status value"}
	}

	allowedColors := map[string]bool{"red": true, "blue": true, "yellow": true, "green": true, "black": true, "white": true, "pink": true, "purple": true}
	if !utils.IsValidColor(m.Usercolor, allowedColors) {
		return &utils.FieldError{Field: "User_color", Message: "must be a valid hex or rgba color"}
	}

	if !utils.IsValidURLOrBase64(m.Avatarurl) {
		return &utils.FieldError{Field: "Avatar", Message: "must be a valid URL or base64 string"}
	}

	if !utils.IsValidURLOrBase64(m.Displayavatarurl) {
		return &utils.FieldError{Field: "DisplayAvatar", Message: "must be a valid URL or base64 string"}
	}
	if !utils.IsValidURLOrBase64(m.Bannerurl) {
		return &utils.FieldError{Field: "Banner", Message: "must be a valid URL or base64 string"}
	}
	if !utils.IsValidURLOrBase64(m.Displaybannerurl) {
		return &utils.FieldError{Field: "DisplayBanner", Message: "must be a valid URL or base64 string"}
	}

	for _, message := range m.Usermessages {
		if !utils.IsValidUTF8String(message.Content) {
			return &utils.FieldError{Field: "Message", Message: "must contain letters, numbers or symbols"}
		}
	}

	return nil
}
