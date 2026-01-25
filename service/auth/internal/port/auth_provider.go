package port

type AuthProvider interface {
	IssueToken(username string) (string, error)
	ValidateToken(token string) (string, error)
}
