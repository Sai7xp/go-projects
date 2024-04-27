/*
* Created on 27 April 2024
* @author Sai Sumanth
 */

package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between given min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money between 0 and 7k
func RandomMoney() int64 {
	return RandomInt(0, 7000)
}

func RandomCurrency() string {
	supportedCurrencies := []string{"INR", "EUR", "USD"}
	return supportedCurrencies[rand.Intn(len(supportedCurrencies))]
}
