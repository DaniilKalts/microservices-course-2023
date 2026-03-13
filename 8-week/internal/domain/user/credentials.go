package user

type Credentials struct {
	ID           string
	PasswordHash string
	Role         Role
}

func ValidatePassword(password string) error {
	var hasUpper, hasLower, hasDigit bool

	for _, r := range password {
		switch {
		case 'A' <= r && r <= 'Z':
			hasUpper = true
		case 'a' <= r && r <= 'z':
			hasLower = true
		case '0' <= r && r <= '9':
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return ErrWeakPassword
	}

	return nil
}
