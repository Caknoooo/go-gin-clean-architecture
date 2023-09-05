package services

import (
	"context"

	"github.com/Caknoooo/golang-clean_template/dto"
	"github.com/Caknoooo/golang-clean_template/entities"
	"github.com/Caknoooo/golang-clean_template/helpers"
	"github.com/Caknoooo/golang-clean_template/repository"
)

type UserService interface {
	RegisterUser(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error)
	GetAllUser(ctx context.Context) ([]dto.UserResponse, error)
	GetUserByID(ctx context.Context, userID string) (dto.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error)
	CheckUser(ctx context.Context, email string) (bool, error)
	UpdateUser(ctx context.Context, req dto.UserUpdateRequest, userId string) error
	DeleteUser(ctx context.Context, userId string) error 
	Verify(ctx context.Context, email string, password string) (bool, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepository: ur,
	}
}

func (us *userService) RegisterUser(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error) {
	user := entities.User{
		Nama:     req.Nama,
		NoTelp:   req.NoTelp,
		Role: 	 "user",
		Email:    req.Email,
		Password: req.Password,
	}

	userResponse, err := us.userRepository.RegisterUser(ctx, user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	return dto.UserResponse{
		ID:       userResponse.ID.String(),
		Nama:     userResponse.Nama,
		NoTelp:   userResponse.NoTelp,
		Role: 	  userResponse.Role,
		Email:    userResponse.Email,
		Password: userResponse.Password,
	}, nil
}

func (us *userService) GetAllUser(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := us.userRepository.GetAllUser(ctx)
	if err != nil {
		return nil, dto.ErrGetAllUser
	}

	var userResponse []dto.UserResponse
	for _, user := range users {
		userResponse = append(userResponse, dto.UserResponse{
			ID:       user.ID.String(),
			Nama:     user.Nama,
			NoTelp:   user.NoTelp,
			Role: 	  user.Role,
			Email:    user.Email,
			Password: user.Password,
		})
	}

	return userResponse, nil
}

func (us *userService) GetUserByID(ctx context.Context, userID string) (dto.UserResponse, error) {
	user, err := us.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserById
	}

	return dto.UserResponse{
		ID:       user.ID.String(),
		Nama:     user.Nama,
		NoTelp:   user.NoTelp,
		Role: 	  user.Role,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (us *userService) GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error) {
	emails, err := us.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserByEmail
	}

	return dto.UserResponse{
		ID:       emails.ID.String(),
		Nama:     emails.Nama,
		NoTelp:   emails.NoTelp,
		Role: 	  emails.Role,
		Email:    emails.Email,
		Password: emails.Password,
	}, nil
}

func (us *userService) CheckUser(ctx context.Context, email string) (bool, error) {
	res, err := us.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if res.Email == "" {
		return false, err
	}
	return true, nil
}

func (us *userService) UpdateUser(ctx context.Context, req dto.UserUpdateRequest, userId string) error {
	user, err := us.userRepository.GetUserByID(ctx, userId)
	if err != nil {
		return dto.ErrUserNotFound
	}

	userUpdate := entities.User{
		ID:       user.ID,
		Nama:     req.Nama,
		NoTelp:   req.NoTelp,
		Role: 	  user.Role,
		Email:    req.Email,
		Password: req.Password,
	}

	err = us.userRepository.UpdateUser(ctx, userUpdate)
	if err != nil {
		return dto.ErrUpdateUser
	}

	return nil
}

func (us *userService) DeleteUser(ctx context.Context, userId string) error {
	user, err := us.userRepository.GetUserByID(ctx, userId)
	if err != nil {
		return dto.ErrUserNotFound
	}

	err = us.userRepository.DeleteUser(ctx, user.ID.String())
	if err != nil {
		return dto.ErrDeleteUser
	}

	return nil
}

func (us *userService) Verify(ctx context.Context, email string, password string) (bool, error) {
	res, err := us.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	checkPassword, err := helpers.CheckPassword(res.Password, []byte(password))
	if err != nil {
		return false, err
	}

	if res.Email == email && checkPassword {
		return true, nil
	}
	return false, nil
}
