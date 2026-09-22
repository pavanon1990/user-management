package port

type TokenProvider interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (string, error) // return userID, ใช้ตอนทำ middleware ทีหลัง
}
