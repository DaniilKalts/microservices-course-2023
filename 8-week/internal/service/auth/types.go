package auth

type TokenPair struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresIn  int64
	RefreshTokenExpiresIn int64
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type LogoutInput struct {
	RefreshToken string
}

type RefreshInput struct {
	RefreshToken string
}
