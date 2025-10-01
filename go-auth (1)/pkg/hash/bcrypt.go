package hash

import "golang.org/x/crypto/bcrypt"

func HashPassword(p string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
    return string(b), err
}

func CheckPassword(hashv, plain string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hashv), []byte(plain)) == nil
}
