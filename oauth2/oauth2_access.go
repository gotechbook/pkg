package oauth2

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

type Access struct{}

func NewAccess() *Access {
	return &Access{}
}

func (a *Access) Token(_ context.Context, data *Oauth2) (access, refresh string, err error) {
	buf := bytes.NewBufferString(data.Client.GetID())
	buf.WriteString(data.UserID)
	buf.WriteString(strconv.FormatInt(data.CreateAt.UnixNano(), 10))

	access = base64.URLEncoding.EncodeToString([]byte(uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes()).String()))
	access = strings.ToUpper(strings.TrimRight(access, "="))

	refresh = base64.URLEncoding.EncodeToString([]byte(uuid.NewSHA1(uuid.Must(uuid.NewRandom()), buf.Bytes()).String()))
	refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	return access, refresh, nil
}
