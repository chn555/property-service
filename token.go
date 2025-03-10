package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type NextToken struct {
	Limit  int
	Offset int
}

func decodeNextToken(token string) (*NextToken, error) {
	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %v", err)
	}

	t := &NextToken{}
	if err := json.Unmarshal(b, t); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token %v", err)
	}
	return t, nil
}

func createNextToken(limit int, offset int) (string, error) {
	token := &NextToken{
		Limit:  limit,
		Offset: offset,
	}
	return token.encode()
}
func (t *NextToken) encode() (string, error) {
	if t.Offset == 0 && t.Limit == 0 {
		return "", nil
	}
	b, err := json.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("failed to json marshalling %v", err)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
