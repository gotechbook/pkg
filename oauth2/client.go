package oauth2

import "time"

type Client interface {
	GetID() string
	GetName() string
	GetDescribe() string
	GetSecret() string
	GetVersion() string
	GetDomain() string
	GetScope() string
	GetGrant() string
	GetCodeExpiresIn() time.Duration
	GetAccessExpiresIn() time.Duration
	GetRefreshExpiresIn() time.Duration
}
