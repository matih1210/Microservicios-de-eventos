package token

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    SID      string `json:"sid"`
    UserID   string `json:"uid"`
    Username string `json:"usr"`
    jwt.RegisteredClaims
}

func NewJWT(secret, sid, userID, username string, expMin int) (string, error) {
    now := time.Now()
    claims := Claims{
        SID:      sid,
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expMin) * time.Minute)),
        },
    }
    t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return t.SignedString([]byte(secret))
}
