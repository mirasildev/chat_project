package usecase

import (
	"github.com/google/uuid"
	"github.com/mirasildev/chat_task/domain"
)

type UserRepository interface {
	Store(user *domain.User) (*domain.User, error)
	Get(id string) (*domain.User, error)
	GetAll(limit, page int64, search string) (*domain.GetAllUsersReponse, error)
	GetByEmail(email string) (*domain.User, error)
	UpdatePassword(userID, password string) error
	Delete(id string) error
	Update(user *domain.User) (*domain.User, error)
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(g UserRepository) *UserService {
	return &UserService{
		userRepo: g,
	}
}

func (c *UserService) CreateUser(user *domain.User) (*domain.User, error) {
	user.ID = uuid.New().String()
	return c.userRepo.Store(user)
}

func (c *UserService) GetUser(id string) (*domain.User, error) {
	return c.userRepo.Get(id)
}

func (c *UserService) GetAllUsers(limit, page int64, search string) (*domain.GetAllUsersReponse, error) {
	return c.userRepo.GetAll(limit, page, search)
}

func (c *UserService) GetUserByEmail(email string) (*domain.User, error){
	return c.userRepo.GetByEmail(email)
}

func (c *UserService) UpdateUserPassword(userID, password string) error {
	return c.userRepo.UpdatePassword(userID, password)
}

func (c *UserService) DeleteUser(id string) error {
	return c.userRepo.Delete(id)
}

func (c *UserService) Update(user *domain.User) (*domain.User, error) {
	return c.userRepo.Update(user)
}
