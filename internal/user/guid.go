package user

import "github.com/google/uuid"

func VerifyGUID(guid string) (uuid.UUID, error) {
	uuid, err := uuid.Parse(guid)
	if err != nil {
		return uuid, err
	}
	return uuid, nil
}
