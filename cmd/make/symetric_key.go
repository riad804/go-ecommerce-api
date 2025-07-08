package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		panic("failed to generate symmetric key: " + err.Error())
	}

	encodedKey := base64.StdEncoding.EncodeToString(key)
	fmt.Println("Symmetric Key:", encodedKey)
}
