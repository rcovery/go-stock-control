package part

import "github.com/google/uuid"

type ID string

func NewID() (ID, error) {
	newuuid, err := uuid.NewV7()
	return ID(newuuid.String()), err
}

func ParseID(s string) (ID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return "", err
	}
	return ID(parsed.String()), nil
}

func (id ID) String() string {
	return string(id)
}
