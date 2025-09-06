package service

import (
	"github.com/aborgas90/expense-tracker-api/internal/model"
	"github.com/aborgas90/expense-tracker-api/internal/repo"
	auth "github.com/aborgas90/expense-tracker-api/internal/dto/auth"
	token "github.com/aborgas90/expense-tracker-api/internal/auth"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"strconv"
)

type UserService struct {
	repo *repo.UserRepo
}

func NewUserService(r *repo.UserRepo) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) RegisterUser(req *auth.RegisterUserRequest) (*auth.RegisterUserResponse, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("username, email, dan password wajib diisi")
	}
	if len(req.Password) < 6 {
		return nil, errors.New("password minimal 6 karakter")
	}
	
	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashed),
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}

	res := &auth.RegisterUserResponse{
		ID:        strconv.FormatUint(uint64(u.Id), 10),
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}

	return res, nil
}

func (s *UserService) LoginUser(req *auth.LoginUserRequest) (*auth.LoginUserResponse, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return  nil, errors.New("username atau password salah")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, errors.New("username atau password salah")
	}

	accessToken, refreshToken, err := token.GenerateToken(user.Id)
	if err != nil {
		return  nil, err 
	}

	return &auth.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
