package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/riad804/go_ecommerce_api/internals/models"
	"github.com/riad804/go_ecommerce_api/internals/repositories"
	"github.com/riad804/go_ecommerce_api/token"
	"github.com/riad804/go_ecommerce_api/workers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	tokenMaker  token.Maker
	userRepo    repositories.UserRepository
	cfg         *config.Config
	distributor workers.TaskDistributor
}

func NewAuthService(tokenMaker token.Maker, userRepo repositories.UserRepository, cfg *config.Config, distributor workers.TaskDistributor) *AuthService {
	return &AuthService{
		tokenMaker:  tokenMaker,
		userRepo:    userRepo,
		cfg:         cfg,
		distributor: distributor,
	}
}

func (s *AuthService) CreateUser(req models.RegisterRequest) (*models.User, int, error) {
	hashPassword, err := config.HashPassword(req.Password)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashPassword,
		Phone:    req.Phone,
		IsAdmin:  false,
	}
	r, err := s.userRepo.Create(user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, fiber.StatusConflict, fmt.Errorf("user already existed")
		}
		return nil, fiber.StatusInternalServerError, err
	}
	oid, ok := r.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("InsertedID is not an ObjectID")
	}
	result, err := s.userRepo.FindByID(oid)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	return result, fiber.StatusAccepted, nil
}

func (s *AuthService) FindUser(req models.LoginRequest) (*map[string]any, int, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fiber.StatusNotFound, fmt.Errorf("user not found")
	}

	err = config.CheckPassword(req.Password, user.Password)
	if err != nil {
		return nil, fiber.StatusBadRequest, fmt.Errorf("incorrect password")
	}

	accessTokenPayload, err := s.tokenMaker.CreateAccessToken(user.Email, user.IsAdmin, s.cfg.Token.AccDuration)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("access err: %s", err.Error())
	}

	refreshTokenPayload, err := s.tokenMaker.CreateRefreshToken(user.Email, user.IsAdmin, s.cfg.Token.RefDuration)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("refresh err: %s", err.Error())
	}

	user.Password = ""
	result := map[string]any{
		"access_token":  accessTokenPayload.Token,
		"refresh_token": refreshTokenPayload.Token,
		"user":          user,
	}
	return &result, fiber.StatusOK, nil
}

func (s *AuthService) ForgotPassword(req models.ForgotPassRequest) (*models.User, int, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fiber.StatusNotFound, fmt.Errorf("user not found")
	}

	otp := config.RandomInt(1000, 9000)

	user.ResetPasswordOtp = &otp
	expiry := time.Now().Add(10 * time.Minute)
	user.ResetPasswordOtpExpires = &expiry

	_, err = s.userRepo.Update(*user)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("User updating failed")
	}

	payload := workers.PayloadSendVerifyEmail{
		Name:  user.Name,
		Email: user.Email,
		OTP:   fmt.Sprintf("%d", otp),
	}
	err = s.distributor.DistributeTaskSendVerifyEmail(&payload)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return user, fiber.StatusOK, nil
}

func (s *AuthService) OtpVerify(req models.VerifyOtpRequest) (*models.User, int, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fiber.StatusNotFound, fmt.Errorf("user not found")
	}

	otp, err := strconv.ParseInt(req.OTP, 10, 64)
	if err != nil {
		return nil, fiber.StatusBadRequest, fmt.Errorf("Invalid otp")
	}

	if *user.ResetPasswordOtp != otp || time.Now().Unix() > user.ResetPasswordOtpExpires.Unix() {
		return nil, fiber.StatusUnauthorized, fmt.Errorf("Expired otp")
	}
	ptr := new(int64)
	*ptr = 1
	user.ResetPasswordOtp = ptr
	user.ResetPasswordOtpExpires = nil

	_, err = s.userRepo.Update(*user)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("User updating failed")
	}
	return user, fiber.StatusOK, nil
}

func (s *AuthService) ResetPassword(req models.ResetPasswordRequest) (*models.User, int, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fiber.StatusNotFound, fmt.Errorf("user not found")
	}
	ptr := new(int64)
	*ptr = 1
	if user.ResetPasswordOtp != nil && *user.ResetPasswordOtp != *ptr {
		return nil, fiber.StatusUnauthorized, fmt.Errorf("Confirm OTP before reseting password")
	}
	hashPass, err := config.HashPassword(req.Password)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	user.Password = hashPass
	user.ResetPasswordOtp = nil
	_, err = s.userRepo.Update(*user)
	if err != nil {
		return nil, fiber.StatusInternalServerError, fmt.Errorf("User updating failed")
	}
	return user, fiber.StatusOK, nil
}
