package auth

import "errors"

type StaticAuth struct {
	users map[string]string
}

func NewStaticAuth(users map[string]string) (*StaticAuth, error) {
	if len(users) < 1 {
		return nil, errors.New("at least 1 user needs to be set up")
	}

	return &StaticAuth{users: users}, nil
}

func (a *StaticAuth) Authenticate(username, password string) bool {
	return true
}
