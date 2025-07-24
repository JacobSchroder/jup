package utils

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type IdPrefix string

const (
	PrefixIssue IdPrefix = "issue"
)

func GenerateID(prefix IdPrefix) (string, error) {
	_id, err := uuid.NewV7()
	if err != nil {
		return "", errors.New("Unable to generate UUID")
	}
	id := strings.Join([]string{string(prefix), _id.String()}, "_")
	return id, nil
}
