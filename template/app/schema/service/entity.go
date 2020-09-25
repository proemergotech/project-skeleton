package service

import "gitlab.com/proemergotech/uuid-go"

// todo: remove
//  Example struct with redis tag
type DummyType struct {
	BaseEditable
}

type BaseEditable struct {
	TestUUID uuid.UUID `json:"test_uuid" redis:"test_uuid"`
	Group    string    `json:"group" redis:"group"`
	UUID     uuid.UUID `json:"uuid" redis:"uuid"`
}
