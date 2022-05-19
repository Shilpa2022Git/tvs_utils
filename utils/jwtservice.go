package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	//GenerateToken(UserID string) string
	GenerateServiceValidationToken(chasisno, vehicleregno string) string
	ValidateToken(token string) (*jwt.Token, error)
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	ChasisNumber string `json:"chasisno"`
	Vehicleregno string `json:"vehicleregno"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService(secretkey string) JWTService {
	return &jwtService{
		issuer:    "trafficviolationsystem",
		secretKey: secretkey,
	}
}

func (j *jwtService) GenerateServiceValidationToken(chasisno, vehicleregno string) string {

	claims := &Payload{
		chasisno,
		vehicleregno,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 45).Unix(), //time.Now().AddDate(1, 0, 0).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(j.secretKey))

	if err != nil {
		panic(err)
	}
	return t
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return verifyToken(token, j.secretKey)
}

func verifyToken(token, secretKey string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
}

func (j *jwtService) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Algo used ", token.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing method %s", token.Header["alg"])
		}

		str, _ := base64.StdEncoding.DecodeString(j.secretKey)
		return str, nil
	})

	fmt.Println("verify Error : ", err, err.(*jwt.ValidationError))

	// if err != nil {
	// 	verr, ok := err.(*jwt.ValidationError)
	// 	if ok && errors.Is(verr.Inner, errors.New("token is expired")) {
	// 		return nil, errors.New("token is invalid")
	// 	}
	// 	return nil, errors.New("token is invalid")
	// }

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errors.New("token is invalid, could not parse claims")
	}

	if payload.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token is expired")
	}

	return payload, nil
}
