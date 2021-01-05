package util

import guuid "github.com/google/uuid"

func GetUUID() string {
	id, err := guuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return id.String()
}
