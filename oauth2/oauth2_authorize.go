package oauth2

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/google/uuid"
	"strings"
)

type Authorize struct{}

func NewAuthorize() *Authorize {
	return &Authorize{}
}

func (a *Authorize) Token(_ context.Context, data *Oauth2) (string, string, error) {
	buf := bytes.NewBufferString(data.Client.GetID())
	buf.WriteString(data.UserID)
	token := uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes())
	code := base64.URLEncoding.EncodeToString([]byte(token.String()))
	code = strings.ToUpper(strings.TrimRight(code, "="))
	return code, "", nil
}
