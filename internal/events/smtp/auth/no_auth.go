package auth

type NoAuth struct {
}

func NewNoAuth() (*NoAuth, error) {
	return &NoAuth{}, nil
}

func (a *NoAuth) Authenticate(username, password string) bool {
	return true
}
