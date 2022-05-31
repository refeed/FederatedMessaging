package message

import (
	"crypto/sha1"
	"fmt"
)

// This is essentially a function to convert from string to SHA1
// as the user identity is just a hash string
func GetUserIdentity(passphrase string) string {
	sha := sha1.New()
	sha.Write([]byte(passphrase))
	return fmt.Sprintf("%x", sha.Sum(nil))
}
