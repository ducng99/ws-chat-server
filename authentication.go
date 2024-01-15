package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	math_rand "math/rand"
	"strings"
	"time"
)

type TokenInfo struct {
	Username string
	Expires  int64
}

var disallowedUsernames = map[string]bool{
	"admin": true, "administrator": true,
	"moderator": true, "mod": true,
	"system": true, "root": true,
	"guest": true, "anonymous": true,
	"support": true, "help": true, "staff": true, "customer_service": true, "info": true,
	"server": true, "bot": true, "service": true,
	"user": true, "member": true,
	"website": true, "site": true, "chat": true, "official": true, "webmaster": true, "team": true, "security": true,
	"verify": true, "verified": true,
	"error": true, "null": true, "undefined": true,
}

var userTokens = map[string]TokenInfo{}

func IsUsernameValid(username string) error {
	if len(username) < 3 || len(username) > 16 {
		return errors.New("Username must be between 3 and 16 characters")
	}

	if _, ok := disallowedUsernames[strings.ToLower(username)]; ok {
		return errors.New("Username is not allowed")
	}

	if _, ok := ChatServer.getUser(username); ok {
		return errors.New("Username is already in use")
	}

	return nil
}

func CreateAnonymousUsername() string {
	digits := 4
	counter := 0

	for {
		// Generate a random number with the specified number of digits
		number := math_rand.Intn(int(math.Pow10(digits)))
		// Create the username
		username := fmt.Sprintf("Anonymous#%0*d", digits, number)

		// Check if the username is valid and not already in use
		if _, ok := ChatServer.getUser(username); !ok {
			return username
		}

		counter++
		if counter >= 10 {
			digits++
			counter = 0
		}
	}
}

func CreateToken(username string) (string, error) {
	// Generate a random byte slice
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Convert the byte slice to a hexadecimal string
	token := fmt.Sprintf("%x", b)

	// Store the token
	userTokens[token] = TokenInfo{
		username,
		time.Now().Unix() + 5*60, // 5 minutes
	}

	return token, nil
}

func GetUsernameFromToken(token string) (string, bool) {
	info, ok := userTokens[token]
	if !ok || info.Expires < time.Now().Unix() {
		return "", false
	}

	return info.Username, ok
}
