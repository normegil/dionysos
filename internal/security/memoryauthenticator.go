package security

type MemoryAuthenticator struct {
	Username string
	Password string
}

func (a MemoryAuthenticator) Authenticate(username string, password string) (bool, error) {
	if username != a.Username {
		return false, nil
	}
	if password != a.Password {
		return false, nil
	}
	return true, nil
}
