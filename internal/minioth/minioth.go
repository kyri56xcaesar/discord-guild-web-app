package minioth

import (
	"crypto/sha256"
)

// I want to implement a simplistic but handy user/group state system.
// Let's copy the UNIX model
//
// /etc/passwd, /etc/shadow, /etc/group

const (
	MINIOTH_PASSWD string = "minioth/mpass"
	MINIOTH_GROUP  string = "minioth/mgroup"
	MINIOTH_SHADOW string = "minioth/mshadow"

	MINIOTH_DB string = "minioth.db"
)

type userspace interface {
	useradd() error
	userdel() error
	usermod() error

	groupadd() error
	groupdel() error
	groupmod() error

	passwd() error
}

type User struct {
	Name     string
	Password Password
	Groups   []Group
	Uid      int64
}
type Password struct {
	Hashpass       string
	ExpirationDate string
	Length         int
}
type Group struct {
	Name     string
	Password Password
	Users    []User
	Gid      int64
}

func Hash(password string) []byte {
	hasher := sha256.New()

	return hasher.Sum([]byte(password))
}

type Minioth struct {
	dbPath string

	root       User
	usercount  int
	groupcount int

	useDB bool
}

func NewMinioth(rootname string, useDb bool, dbPath string) Minioth {
	return Minioth{
		dbPath: dbPath,
		root: User{
			Name: rootname,
			Password: Password{
				Hashpass:       "",
				ExpirationDate: "",
				Length:         0,
			},
			Groups: nil,
			Uid:    1,
		},
		usercount:  0,
		groupcount: 0,
		useDB:      useDb,
	}
}

func (m *Minioth) sync() error {
	return nil
}

func (m *Minioth) useradd(username, password string, groups []Group) error {
	return nil
}

func (m *Minioth) userdel(username string) error {
	return nil
}

func (m *Minioth) usermod(username, password string, groups []Group) error {
	return nil
}

func (m *Minioth) groupadd(groupname, password string, users []User) error {
	return nil
}

func (m *Minioth) groupdel(groupname string) error {
	return nil
}

func (m *Minioth) groupmod(groupname, password string, users []User) error {
	return nil
}

func (m *Minioth) passwd(username, password string) error {
	return nil
}
