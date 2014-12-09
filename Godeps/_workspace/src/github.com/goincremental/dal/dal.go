package dal

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("ErrNotFound")

type BSON map[string]interface{}

type DAL interface {
	Connect(string) (Session, error)
	IsObjectIdHex(string) bool
}

type Session interface {
	Clone() Session
	Close()
	DB(s string) Database
}

type Database interface {
	C(string) Collection
}

type Collection interface {
	Find(BSON) Query
	EnsureIndex(Index) error
	FindId(interface{}) Query
	RemoveId(interface{}) error
	UpsertId(interface{}, interface{}) (*ChangeInfo, error)
}

type Query interface {
	One(interface{}) error
	Sort(...string) Query
	Iter() Iter
}

type Iter interface {
	Next(interface{}) bool
}

type Index struct {
	Key         []string
	Background  bool
	Sparse      bool
	ExpireAfter time.Duration
}

type ChangeInfo struct {
	Updated    int
	Removed    int
	UpsertedId interface{}
}
