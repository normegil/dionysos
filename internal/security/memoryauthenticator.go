package security

type MemoryAuthenticator struct {
	Username string
	Password string
}

func (a MemoryAuthenticator) Authenticate(username string, password string) bool {
	if username != a.Username {
		return false
	}
	if password != a.Password {
		return false
	}
	return true
}
