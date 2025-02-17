package services

import (
	"context"
	db "github.com/LuccChagas/my-chat-app/db/sqlc"
	"github.com/LuccChagas/my-chat-app/internal/models"
	"github.com/LuccChagas/my-chat-app/internal/repository"
	"github.com/LuccChagas/my-chat-app/utils"
	"github.com/google/uuid"
)

type UserService struct {
	repository repository.RepositoryInterface
}

func NewUserService(repository repository.RepositoryInterface) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user models.UserRequest) (models.UserResponse, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return models.UserResponse{}, err
	}

	arg := db.CreateUsersParams{
		ID:        uuid.New(),
		Password:  hashedPassword,
		Cpf:       user.Cpf,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      user.Name,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		NickName:  user.NickName,
	}

	response, err := s.repository.CreateUser(ctx, arg)
	if err != nil {
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:        response.ID,
		Cpf:       response.Cpf,
		Email:     response.Email,
		Phone:     response.Phone,
		Name:      response.Name,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		NickName:  response.NickName,
		CreatedAt: response.CreatedAt,
		UpdatedAt: response.UpdatedAt,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, ID uuid.UUID) (models.UserResponse, error) {
	user, err := s.repository.GetUser(ctx, ID)
	if err != nil {
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:        user.ID,
		Password:  user.Password,
		Cpf:       user.Cpf,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      user.Name,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		NickName:  user.NickName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
func (s *UserService) GetAllUsers(ctx context.Context) (response []models.UserResponse, err error) {
	allUsers, err := s.repository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range allUsers {
		response = append(response, models.UserResponse{
			ID:        user.ID,
			Cpf:       user.Cpf,
			Email:     user.Email,
			Phone:     user.Phone,
			Name:      user.Name,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			NickName:  user.NickName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return response, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (models.UserResponse, error) {

	user, err := s.repository.GetUserByNickname(ctx, username)
	if err != nil {
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:        user.ID,
		Password:  user.Password,
		Cpf:       user.Cpf,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      user.Name,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		NickName:  user.NickName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

//if !crypt.CheckPasswordHash(user.Password, dbUser.Password) {
//return model.UserLoginResponse{}, errors.New("invalid password")
//}
