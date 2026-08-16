package stringopr

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	str, _ := uuid.NewV7()

	return str.String()
}

func ParseUUID(strUUID string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(strUUID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse UUID: %v", err)
	}
	return parsedUUID, nil
}

func GenerateRandomStringFromString(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
