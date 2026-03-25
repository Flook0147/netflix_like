package outbound

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/Flook0147/netflix_like/internal/auth/domain"
	userDomain "github.com/Flook0147/netflix_like/internal/user/domain"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	// Implementation details
	DB *gorm.DB
}

func NewRefreshTokenRepository(Db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		DB: Db,
	}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (r *RefreshTokenRepository) SaveRefreshToken(username, refreshToken string) error {

	var user userDomain.User

	// หา user จาก username
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}

	rt := domain.RefreshToken{
		UserID:    user.UserId,
		TokenHash: hashToken(refreshToken),
		CreatedAt: time.Now(),
		RevokedAt: time.Now(),
		ExpiredAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := r.DB.Create(&rt).Error; err != nil {
		return err
	}

	return nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	err := r.DB.Where("token_hash = ?", tokenHash).Delete(&domain.RefreshToken{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokenRepository) FindRefreshToken(refreshToken string) (*domain.RefreshToken, error) {
	tokenHash := hashToken(refreshToken)

	var rt domain.RefreshToken

	err := r.DB.Where("token_hash = ?", tokenHash).First(&rt).Error
	if err != nil {
		return nil, err
	}

	return &rt, nil
}
