package oauth2

import "time"

type Token interface {
	New() Token
	GetClientID() string
	SetClientID(string)
	GetUserID() string
	SetUserID(string)
	GetRedirectURI() string
	SetRedirectURI(string)
	GetScope() string
	SetScope(string)
	GetCode() string
	SetCode(string)
	GetCodeType() string
	SetCodeType(string)
	GetExpiresIn() time.Duration
	SetExpiresIn(time.Duration)
	GetCreateAt() time.Time
	SetCreateAt(createAt time.Time)
}
