package encrypt

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("no input value")
	} else {
		fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Sprint("cannot generate encrypted password"), err
		}
		hashedString := string(fromPassword)
		return hashedString, nil
	}

}

func Verify(password, hashedPassword string) (bool, error) {
	if password == "" || hashedPassword == "" {
		return false, nil
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, fmt.Errorf("invalid string comparision : %v", err)
		}
		return false, err
	}
	return true, nil
}
