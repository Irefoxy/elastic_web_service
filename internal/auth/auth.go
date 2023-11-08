package auth

type Auth interface {
	GenerateToken() (string, error)
	VerifyToken(tokenString string) (bool, error)
}
