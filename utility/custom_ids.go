package utility

import (
	"crypto/rand"
	"strconv"
)

func generateRandomNumber(max int) int {
	var b [1]byte
	rand.Read(b[:])
	return int(b[0]) % max
}

func GenerateSellerID() string {
	// generated seller id will be 10 letters long
	sellerID := ""
	for range 10 {
		// generating a random choice between 0 and 2 inclusive
		// 0 means capital letters
		// 1 means numbers
		// 2 means special chars
		specialChars := "!@#$%&"

		randomChoice := generateRandomNumber(3)
		switch randomChoice {
		case 0:
			randomLetterIndex := generateRandomNumber(26)
			sellerID = sellerID + string(rune(65+randomLetterIndex))
			break
		case 1:
			randomNumber := generateRandomNumber(10)
			sellerID = sellerID + strconv.Itoa(randomNumber)
			break
		case 2:
			randomSpecialCharIndex := generateRandomNumber(6)
			sellerID = sellerID + string(specialChars[randomSpecialCharIndex])
			break
		}
	}

	return sellerID
}
