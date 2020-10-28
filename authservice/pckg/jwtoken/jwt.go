package jwtoken

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

type JwtTokenMember struct {
	// uniqueSignature : уникальный ключ (подпись) сервиса.
	//
	// Используется для создания jwtToken-токена.
	// Читай здесь https://medium.com/@zhashkevych/jwt-%D0%B0%D0%B2%D1%82%D0%BE%D1%80%D0%B8%D0%B7%D0%B0%D1%86%D0%B8%D1%8F-%D0%B4%D0%BB%D1%8F-%D0%B2%D0%B0%D1%88%D0%B5%D0%B3%D0%BE-api-%D0%BD%D0%B0-go-80325de8691b.
	//
	uniqueSignature []byte

	// Время жизни токена.
	//
	accessTokenLifetime  time.Duration
	refreshTokenLifetime time.Duration
}

type Payload struct {
	FieldFirst  string `json:"first"`
	FieldSecond string `json:"second"`
}

type jwtClaim struct {
	Payload *Payload `json:"payload"`
	jwt.StandardClaims
}

func NewJwtTokenConstructor(uniqueSignature []byte, accessLifetime, refreshLifetime time.Duration) *JwtTokenMember {
	return &JwtTokenMember{
		uniqueSignature:      uniqueSignature,
		accessTokenLifetime:  accessLifetime,
		refreshTokenLifetime: refreshLifetime,
	}
}

func (attr *JwtTokenMember) CreateTokens(model *Payload) (accessToken string, refreshToken string, err error) {
	refreshToken = ""
	accessToken = ""
	err = errors.New("undefined error. ")
	if strings.TrimSpace(model.FieldSecond) == "" || strings.TrimSpace(model.FieldFirst) == "" {
		return accessToken, refreshToken, errors.New("invalid data. ")
	}
	refreshToken, err = attr.encode(model, attr.refreshTokenLifetime)
	if err != nil {
		return accessToken, refreshToken, err
	}
	accessToken, err = attr.encode(model, attr.accessTokenLifetime)
	if err != nil {
		return accessToken, refreshToken, err
	}
	return accessToken, refreshToken, nil
}

func (attr *JwtTokenMember) NewRefreshToken(model *Payload) (string, error) {
	if strings.TrimSpace(model.FieldSecond) == "" || strings.TrimSpace(model.FieldFirst) == "" {
		return "", errors.New("invalid data. ")
	}
	encodeTokenString, err := attr.encode(model, attr.refreshTokenLifetime)
	return encodeTokenString, err
}

func (attr *JwtTokenMember) NewAccessToken(model *Payload) (string, error) {
	if strings.TrimSpace(model.FieldSecond) == "" || strings.TrimSpace(model.FieldFirst) == "" {
		return "", errors.New("invalid data. ")
	}
	encodeTokenString, err := attr.encode(model, attr.accessTokenLifetime)
	return encodeTokenString, err
}

func (attr *JwtTokenMember) Encode(model *Payload, lifetime time.Duration) (string, error) {
	return attr.encode(model, lifetime)
}

func (attr *JwtTokenMember) encode(model *Payload, lifetime time.Duration) (string, error) {
	standardClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(lifetime).Unix(),
	}
	payload := &Payload{
		FieldFirst:  model.FieldFirst,
		FieldSecond: model.FieldSecond,
	}
	claims := &jwtClaim{
		Payload:        payload,
		StandardClaims: standardClaims,
	}
	encodeToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodeTokenString, err := encodeToken.SignedString(attr.uniqueSignature)
	if err != nil {
		return encodeTokenString, err
	}
	return encodeTokenString, err
}

func (attr *JwtTokenMember) Decode(encodeTokenString string) (*Payload, error) {
	//
	var keyFunc = func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			errorDesc := "Unexpected signing method: " + token.Header["alg"].(string)
			return nil, errors.New(errorDesc)
		}
		return attr.uniqueSignature, nil
	}

	decodeToken, err := jwt.ParseWithClaims(encodeTokenString, &jwtClaim{}, keyFunc)
	var payload *Payload
	claims, ok := decodeToken.Claims.(*jwtClaim)
	payload = claims.Payload
	if err != nil {
		return payload, err
	}
	if ok && decodeToken.Valid {
		if strings.TrimSpace(payload.FieldFirst) == "" || strings.TrimSpace(payload.FieldSecond) == "" {
			return nil, errors.New(" ")
		}
		return payload, nil
	}
	return payload, errors.New("invalid token. ")
}
