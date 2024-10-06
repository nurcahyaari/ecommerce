package transferobject

import (
	"strconv"
	"strings"
	"time"

	"github.com/nurcahyaari/ecommerce/src/domain/entity"
	"gopkg.in/guregu/null.v4"
)

type RequestSearchUser struct {
	Or    *RequestSearchUser
	Ids   string `json:"ids"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Pagination
}

func (r RequestSearchUser) composeUserFilter() (entity.UserFilter, error) {
	userFilter := entity.UserFilter{}
	if r.Ids != "" {
		ids := strings.Split(r.Ids, ",")
		for _, id := range ids {
			idInt, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return userFilter, err
			}

			userFilter.Ids = append(userFilter.Ids, idInt)
		}
	}

	if r.Phone != "" {
		userFilter.Phone = null.StringFrom(r.Phone)
	}

	if r.Email != "" {
		userFilter.Email = null.StringFrom(r.Email)
	}

	return userFilter, nil
}

func (r RequestSearchUser) UserFilter() (entity.UserFilter, error) {
	userFilter, err := r.composeUserFilter()
	if err != nil {
		return userFilter, err
	}
	if r.Or != nil {
		userFilterOr, err := r.Or.composeUserFilter()
		if err != nil {
			return userFilter, err
		}
		userFilter.Or = &userFilterOr
	}

	if r.Pagination.Page == 0 && r.Pagination.Size == 0 {
		r.Pagination.Default()
	}

	if r.Pagination.Page != 0 {
		userFilter.Pagination.Page = r.Pagination.Page
	}

	if r.Pagination.Size != 0 {
		userFilter.Pagination.Size = r.Pagination.Size
	}

	return userFilter, nil
}

type User struct {
	Id        int64       `json:"id"`
	Name      string      `json:"name"`
	Phone     string      `json:"phone"`
	Email     string      `json:"email"`
	Password  string      `json:"-"`
	IsActived bool        `json:"is_actived"`
	IsRemoved bool        `json:"is_removed"`
	CreatedAt time.Time   `json:"created_at"`
	CreatedBy null.String `json:"created_by"`
	UpdatedAt time.Time   `json:"updated_at"`
	UpdatedBy null.String `json:"updated_by"`
	DeletedAt null.Time   `json:"deleted_at"`
	DeletedBy null.String `json:"deleted_by"`
}

func (u User) Entity() entity.User {
	return entity.User{
		Id:       u.Id,
		Name:     u.Name,
		Phone:    u.Phone,
		Email:    u.Email,
		Password: u.Password,
		Signature: entity.Signature{
			IsActived: u.IsActived,
			IsRemoved: u.IsRemoved,
			CreatedAt: u.CreatedAt,
			CreatedBy: u.CreatedBy,
			UpdatedAt: u.UpdatedAt,
			UpdatedBy: u.UpdatedBy,
			DeletedAt: u.DeletedAt,
			DeletedBy: u.DeletedBy,
		},
	}
}

func NewUser(user entity.User) User {
	return User{
		Id:        user.Id,
		Name:      user.Name,
		Phone:     user.Phone,
		Email:     user.Email,
		Password:  user.Password,
		IsActived: user.IsActived,
		IsRemoved: user.IsRemoved,
		CreatedAt: user.CreatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
		DeletedAt: user.DeletedAt,
		DeletedBy: user.DeletedBy,
	}
}

type Users []User

func NewUsers(users entity.Users) Users {
	respUser := make(Users, 0)
	for _, user := range users {
		respUser = append(respUser, NewUser(user))
	}

	return respUser
}

type ResponseSearchUser struct {
	Users      Users      `json:"users"`
	Pagination Pagination `json:"pagination"`
}

func NewResponseSearchUser(users entity.Users, pagination entity.Pagination) ResponseSearchUser {
	return ResponseSearchUser{
		Users:      NewUsers(users),
		Pagination: NewPagination(pagination),
	}
}

type ResponseGetUser struct {
	User User `json:"user"`
}

func NewResponseGetUser(user entity.User) ResponseGetUser {
	return ResponseGetUser{
		User: NewUser(user),
	}
}

type UserAddress struct {
	Id          int64       `json:"id"`
	UserId      int64       `json:"user_id"`
	FullAddress string      `json:"full_address"`
	IsActived   bool        `json:"is_actived"`
	IsRemoved   bool        `json:"is_removed"`
	CreatedAt   time.Time   `json:"created_at"`
	CreatedBy   null.String `json:"created_by"`
	UpdatedAt   time.Time   `json:"updated_at"`
	UpdatedBy   null.String `json:"updated_by"`
	DeletedAt   null.Time   `json:"deleted_at"`
	DeletedBy   null.String `json:"deleted_by"`
}

func (ua UserAddress) Entity() entity.UserAddress {
	return entity.UserAddress{
		Id:          ua.Id,
		UserId:      ua.UserId,
		FullAddress: ua.FullAddress,
		Signature: entity.Signature{
			IsActived: ua.IsActived,
			IsRemoved: ua.IsRemoved,
			CreatedAt: ua.CreatedAt,
			CreatedBy: ua.CreatedBy,
			UpdatedAt: ua.UpdatedAt,
			UpdatedBy: ua.UpdatedBy,
			DeletedAt: ua.DeletedAt,
			DeletedBy: ua.DeletedBy,
		},
	}
}

func NewUserAddress(userAddress entity.UserAddress) UserAddress {
	return UserAddress{
		Id:          userAddress.Id,
		UserId:      userAddress.UserId,
		FullAddress: userAddress.FullAddress,
		IsActived:   userAddress.IsActived,
		IsRemoved:   userAddress.IsRemoved,
		CreatedAt:   userAddress.CreatedAt,
		CreatedBy:   userAddress.CreatedBy,
		UpdatedAt:   userAddress.UpdatedAt,
		UpdatedBy:   userAddress.UpdatedBy,
		DeletedAt:   userAddress.DeletedAt,
		DeletedBy:   userAddress.DeletedBy,
	}
}

type UserAddresses []UserAddress

func (uas UserAddresses) Entity() entity.UserAddresses {
	ent := entity.UserAddresses{}

	for _, ua := range uas {
		ent = append(ent, ua.Entity())
	}

	return ent
}

func NewUserAddresses(userAddresses entity.UserAddresses) UserAddresses {
	resp := UserAddresses{}

	for _, userAddress := range userAddresses {
		resp = append(resp, NewUserAddress(userAddress))
	}

	return resp
}

type RequestSearchUserAddress struct {
	Or      *RequestSearchUserAddress
	Ids     string `json:"ids"`
	UserIds string `json:"user_ids"`
	Pagination
}

func (r RequestSearchUserAddress) composeUserFilter() (entity.UserAddressFilter, error) {
	userFilter := entity.UserAddressFilter{}
	if r.Ids != "" {
		ids := strings.Split(r.Ids, ",")
		for _, id := range ids {
			idInt, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return userFilter, err
			}

			userFilter.Ids = append(userFilter.Ids, idInt)
		}
	}

	if r.UserIds != "" {
		ids := strings.Split(r.UserIds, ",")
		for _, id := range ids {
			idInt, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return userFilter, err
			}

			userFilter.UserIds = append(userFilter.UserIds, idInt)
		}
	}

	return userFilter, nil
}

func (r RequestSearchUserAddress) UserFilter() (entity.UserAddressFilter, error) {
	userFilter, err := r.composeUserFilter()
	if err != nil {
		return userFilter, err
	}
	if r.Or != nil {
		userFilterOr, err := r.Or.composeUserFilter()
		if err != nil {
			return userFilter, err
		}
		userFilter.Or = &userFilterOr
	}

	if r.Pagination.Page == 0 && r.Pagination.Size == 0 {
		r.Pagination.Default()
	}

	if r.Pagination.Page != 0 {
		userFilter.Pagination.Page = r.Pagination.Page
	}

	if r.Pagination.Size != 0 {
		userFilter.Pagination.Size = r.Pagination.Size
	}

	return userFilter, nil
}

type ResponseSearchUserAddress struct {
	UserAddresses UserAddresses `json:"userAddresses"`
	Pagination    Pagination    `json:"pagination"`
}

func NewResponseSearchUserAddress(userAddresses entity.UserAddresses, pagination entity.Pagination) ResponseSearchUserAddress {
	return ResponseSearchUserAddress{
		UserAddresses: NewUserAddresses(userAddresses),
		Pagination:    NewPagination(pagination),
	}
}

type ResponseGetUserAddress struct {
	UserAddress UserAddress `json:"user_address"`
}

func NewResponseGetUserAddress(userAddress entity.UserAddress) ResponseGetUserAddress {
	return ResponseGetUserAddress{
		UserAddress: NewUserAddress(userAddress),
	}
}
