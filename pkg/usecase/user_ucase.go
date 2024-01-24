package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/taroxii/vote-api/pkg/config"
	"github.com/taroxii/vote-api/pkg/constants"
	"github.com/taroxii/vote-api/pkg/entity"
	"github.com/taroxii/vote-api/pkg/utils/logger"
	"go.uber.org/zap"
)

type userUsecase struct {
	userRepository entity.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(u entity.UserRepository, timeout time.Duration) entity.UserUsecase {
	return &userUsecase{
		userRepository: u,
		contextTimeout: timeout,
	}
}

func (uc *userUsecase) SignIn(ctx context.Context, username string) (tokens *string, userDM *entity.User, err error) {
	u, err := uc.userRepository.GetByUsername(ctx, username)
	if err != nil {
		return nil, nil, err
	}
	c := entity.JWTClaims{
		ID:       u.ID,
		Username: u.Username,
		Issuer:   constants.ACCOUNT_CENTER,
	}
	var claimsDat jwt.MapClaims
	claimsBytes, err := json.Marshal(c)
	if err != nil {
		logger.Logger.Error("Error to create map string ", zap.Error(err))
		return nil, nil, err
	}
	err = json.Unmarshal(claimsBytes, &claimsDat)
	if err != nil {
		logger.Logger.Error("Failed to unmarshal []bytes(claims) ", zap.Error(err))
		return nil, nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsDat)

	tokenString, err := token.SignedString([]byte(config.Config.JWTSecret))
	if err != nil {
		logger.Logger.Error("Failed to signedString) ", zap.Error(err))
	}
	return &tokenString, &u, nil

}
