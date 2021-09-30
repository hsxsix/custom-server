/**
 * @File: user.go
 * @Author: hsien
 * @Description:
 * @Date: 9/22/21 6:17 PM
 */

package user

import (
	"context"
	"custom_server/internal/consts"
	"custom_server/internal/model"
	"custom_server/internal/model/postgres"
	"custom_server/internal/model/request"
	"custom_server/internal/model/response"
	"custom_server/internal/repository"
	"custom_server/pkg/util"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type Service struct {
	repository *repository.User
}

func NewService(db *gorm.DB) *Service {
	return &Service{repository: repository.NewUser(db)}
}

func (s *Service) GetUserList(ctx context.Context, query *request.UserQuery) (*response.UserListRes, error) {
	count, users, err := s.repository.QueryUser(ctx, query)
	if err != nil {
		return nil, err
	}

	var res = make([]*response.UserInfoRes, len(users))

	for i, user := range users {
		res[i] = &response.UserInfoRes{
			UserID:      user.UserID,
			UserName:    user.UserName,
			PhoneNumber: user.PhoneNumber,
			UserAvatar:  user.UserAvatar,
			UserDesc:    user.UserDesc,
			UserSex:     s.getSexDesc(user.UserSex),
			CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return &response.UserListRes{
		Total: int(count),
		Page: response.Page{
			PageNo:   query.PageNo,
			PageSize: query.PageSize,
		},
		List: res,
	}, nil
}

func (s *Service) GetUserInfoByID(ctx context.Context, userId string) (*response.UserInfoRes, error) {
	user, err := s.repository.QueryUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return &response.UserInfoRes{
		UserID:      user.UserID,
		UserName:    user.UserName,
		PhoneNumber: user.PhoneNumber,
		UserAvatar:  user.UserAvatar,
		UserDesc:    user.UserDesc,
		UserSex:     s.getSexDesc(user.UserSex),
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil

}

func (s *Service) Login(ctx context.Context, login *request.UserLogin) (string, error) {
	user, err := s.repository.QueryUserByPhone(ctx, login.PhoneNumber)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Shadow), []byte(login.Password))
	if err != nil {
		return "", consts.IncorrectPasswordWithError(err)
	}

	claims := &model.UserClaims{
		UserId: user.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * time.Hour).Unix(),
			Issuer:    consts.JwtIssuer,
		},
	}

	return util.GenerateToken(consts.JwtSecret, claims)
}

func (s *Service) Register(ctx context.Context, register *request.UserRegister) error {
	if register == nil {
		return errors.New("register user is nil object")
	}

	exist, err := s.repository.QueryUserByPhone(ctx, register.PhoneNumber)
	if err != nil {
		return err
	}

	if exist.PhoneNumber != "" {
		return consts.PhoneRegistered
	}

	shadow, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &postgres.User{
		UserID:      util.StrMD5(uuid.New().String()),
		UserName:    register.UserName,
		PhoneNumber: register.PhoneNumber,
		Shadow:      string(shadow[:]),
	}

	return s.repository.AddUser(ctx, user)
}

func (s *Service) getSexDesc(sex int) string {
	switch sex {
	case 1:
		return "男"
	case 2:
		return "女"
	default:
		return "其他"
	}
}
