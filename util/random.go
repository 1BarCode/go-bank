package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	db "github.com/1BarCode/go-bank/db/sqlc"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1) // this is equiv to min + (0 -> max - min) => min -> max
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte((c))
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// RandomUser generates a random user
func RandomUser() (user db.User, password string, err error) {
	password = RandomString(6)

	hashedPassword, err := HashPassword(password)

	user = db.User{
		Username: RandomOwner(),
		HashedPassword: hashedPassword,
		FullName: RandomOwner(),
		Email: RandomEmail(),
	}

	return
}