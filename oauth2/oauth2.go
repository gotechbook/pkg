package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gotechbook/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

var CtxUserTokenKey = &ctxUserToken{}

type ctxUserToken struct{}

type Generate interface {
	Token(context.Context, *Oauth2) (string, string, error)
}

type Oauth2 struct {
	Client   Client
	Token    Token
	UserID   string
	DeviceNo string
	CreateAt time.Time
}

type UserToken struct {
	UserID   uint64 `json:"userID"`
	ClientID uint64 `json:"clientID"`
	DeviceNo string `json:"deviceNo"`
}

func (u *UserToken) Marshal() ([]byte, error) {
	return json.Marshal(u)
}
func (u *UserToken) Unmarshal(b []byte) error {
	return json.Unmarshal(b, u)
}
func ParseToken(token string, publicKeyFile string) (*UserToken, error) {
	claims := Claims{}
	if parse, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		b, err := utils.ReadFile(publicKeyFile)
		if err != nil {
			return nil, err
		}
		pk, err := jwt.ParseECPublicKeyFromPEM(b)
		if err != nil {
			return nil, err
		}
		return pk, err
	}); err != nil {
		return nil, err
	} else {
		if parse.Valid {
			ut := &UserToken{}
			err := json.Unmarshal([]byte(claims.Subject), ut)
			if err != nil {
				return nil, err
			}
			return ut, nil
		}
	}
	return nil, errors.New("invalid access token")
}
func ExtractUserToken(ctx context.Context) (*UserToken, error) {
	value := ctx.Value(CtxUserTokenKey)
	if v, ok := value.(*UserToken); ok {
		return v, nil
	}
	return nil, status.Errorf(codes.Unauthenticated, "user token not found in context")
}
func toUserToken(data *Oauth2) (*UserToken, error) {
	uid, err := strconv.Atoi(data.UserID)
	if err != nil {
		return nil, err
	}

	cid, err := strconv.Atoi(data.Client.GetID())
	if err != nil {
		return nil, err
	}

	return &UserToken{UserID: uint64(uid), ClientID: uint64(cid), DeviceNo: data.DeviceNo}, nil
}
