package server

import (
	"crypto/rand"
	"encoding/base64"
)

type Tokener interface {
	Generate() string
}

type RandTokener struct{}

func (r RandTokener) Generate() string {
	b := make([]byte, 12)
	n, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	if n != 12 {
		panic("incorrect read")
	}

	return base64.StdEncoding.EncodeToString(b)
}
