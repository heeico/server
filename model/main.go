package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt time.Time      `json:"createdAt,omitempty"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type URLShortener struct {
	BaseModel
	RedirectURL string `json:"redirectURL"`
	Alias       string `json:"alias"`
	VisitCount  int    `json:"visitCount" gorm:"default:0"`

	URLTrack []URLTrack `json:"-"`
}

type URLTrack struct {
	BaseModel
	IPAddress      string `json:"ipAddress"`
	URLShortenerID uint   `json:"-"`
}

type TotalVisit struct {
	CreatedAt time.Time `json:"createdAt"`
	Total     int       `json:"totalVisit"`
}

type User struct {
	BaseModel
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty" gorm:"unique"`
	Role     string `json:"role,omitempty"`
}

func (u *User) HashPassword() {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 2)
	if err != nil {
		panic(err)
	}
	u.Password = string(bytes)
}

func (u *User) CheckPasswordHash(userPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(u.Password)); err != nil {
		return errors.New("password does not match")
	}
	return nil
}

var hmacSampleSecret []byte = []byte("secretKey")

func (u *User) GetAuthToken() string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["ID"] = u.ID
	claims["user"] = u.Email
	claims["exp"] = time.Now().Add(time.Minute * 100).Unix()
	claims["role"] = u.Role

	tokenString, err := token.SignedString(hmacSampleSecret)

	if err != nil {
		panic(err)
	}

	return tokenString

}

func ValidateAuthToken(tokenString string) (userId int64, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["ID"]
		if userId, ok := userId.(float64); ok {
			return int64(userId), nil
		}
		return 0, errors.New("invalid Token, user id not valid")
	} else {
		return 0, err
	}
}
