package util

import uuid "github.com/satori/go.uuid"

func GetUUID() string {
	id := uuid.NewV4()
	ids := id.String()
	return ids
}
