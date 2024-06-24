package auth

type Token interface {
	CreateToken(username *string, token *string) error
	VerifyToken(token *string) error
}
