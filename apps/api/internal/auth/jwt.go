package auth

import (
"errors"
"time"

"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
secret []byte
}

func NewJWT(secret string) *JWT {
return &JWT{secret: []byte(secret)}
}

type Claims struct {
UserID int64 `json:"user_id"`
jwt.RegisteredClaims
}

func (j *JWT) Sign(userID int64) (string, error) {
now := time.Now()
claims := Claims{
UserID: userID,
RegisteredClaims: jwt.RegisteredClaims{
IssuedAt:  jwt.NewNumericDate(now),
ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
},
}
t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return t.SignedString(j.secret)
}

func (j *JWT) Parse(tokenString string) (Claims, error) {
var claims Claims
t, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
return j.secret, nil
})
if err != nil {
return Claims{}, err
}
if !t.Valid {
return Claims{}, errors.New("invalid token")
}
return claims, nil
}
