package oauth2

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
}

func (c *Claims) Valid() error {
	if c.ExpiresAt.Unix() == c.IssuedAt.Unix() {
		return nil
	}

	if time.Unix(c.ExpiresAt.Unix(), 0).Before(time.Now()) {
		return errors.New("invalid access token")
	}
	return nil
}

type JWT struct {
	issuer string
	prk    *ecdsa.PrivateKey
	puk    *ecdsa.PublicKey
}

func NewJWT(issuer string, prk []byte, puk []byte) *JWT {
	a := &JWT{issuer: issuer}
	if err := a.setPrivateKey(prk); err != nil {
		return nil
	}
	if err := a.setPublicKey(puk); err != nil {
		return nil
	}
	return a
}
func (j *JWT) Token(_ context.Context, data *Oauth2) (access, refresh string, err error) {
	ut, err := toUserToken(data)
	if err != nil {
		return
	}
	utb, err := ut.Marshal()
	if err != nil {
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   string(utb),
			Audience:  []string{data.Client.GetID()},
			ExpiresAt: jwt.NewNumericDate(data.Token.GetCreateAt().Add(data.Token.GetExpiresIn())),
			IssuedAt:  jwt.NewNumericDate(data.Token.GetCreateAt()),
		},
	})
	access, err = token.SignedString(j.prk)
	if err != nil {
		return
	}
	t := uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).String()
	refresh = base64.URLEncoding.EncodeToString([]byte(t))
	refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	return access, refresh, nil
}
func (j *JWT) setPrivateKey(b []byte) error {
	pk, err := jwt.ParseECPrivateKeyFromPEM(b)
	if err != nil {
		return errors.WithStack(err)
	}
	j.prk = pk
	return nil
}
func (j *JWT) setPublicKey(b []byte) error {
	pk, err := jwt.ParseECPublicKeyFromPEM(b)
	if err != nil {
		return errors.Wrap(err, "parse public key error")
	}
	j.puk = pk
	return nil
}
