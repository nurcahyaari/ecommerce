package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

type User struct {
	Id       int64         `db:"id"`
	Name     string        `db:"name"`
	Email    string        `db:"email"`
	Phone    string        `db:"phone"`
	Password string        `db:"password"`
	Address  UserAddresses `db:"-"`
	Signature
}

func (a *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
}

func (a User) Auth(duration time.Duration) Auth {
	now := time.Now()
	exp := now.Add(duration)

	return Auth{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Id: fmt.Sprintf("%d", a.Id),
	}
}

type Users []User

// One returns the first warehouse on the list and the status (found or not found)
func (us Users) One() (User, bool) {
	if len(us) == 0 {
		return User{}, false
	}
	return us[0], true
}

type UserFilter struct {
	Or    *UserFilter
	Ids   []int64
	Phone null.String
	Email null.String
	Pagination
}

func (f UserFilter) composeFilter() ([]string, []interface{}) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	if len(f.Ids) > 0 {
		query = append(query, "id IN (?)")
		args = append(args, f.Ids)
	}

	if f.Phone.Valid {
		query = append(query, "phone = ?")
		args = append(args, f.Phone.ValueOrZero())
	}

	if f.Email.Valid {
		query = append(query, "email = ?")
		args = append(args, f.Email.ValueOrZero())
	}

	return query, args
}

func (f UserFilter) ComposeFilter() (string, []interface{}, error) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	query, args = f.composeFilter()
	if f.Or != nil {
		queryOr, argsOr := f.Or.composeFilter()

		args = append(args, argsOr...)
		whereClauseOr := strings.Join(queryOr, " OR ")
		query = append(query, "("+whereClauseOr+")")
	}

	// Combine query parts into a single WHERE clause
	whereClause := strings.Join(query, " AND ")
	if whereClause != "" {
		whereClause = "WHERE " + whereClause
	}

	return whereClause, args, nil
}

type UserAddress struct {
	Id          int64  `db:"id"`
	UserId      int64  `db:"user_id"`
	FullAddress string `db:"full_address"`
	Signature
}

type UserAddresses []UserAddress

// One returns the first warehouse on the list and the status (found or not found)
func (us UserAddresses) One() (UserAddress, bool) {
	if len(us) == 0 {
		return UserAddress{}, false
	}
	return us[0], true
}

func (us UserAddresses) MapById() MapUserAddress {
	mapById := make(MapUserAddress)
	for _, u := range us {
		mapById[u.Id] = u
	}
	return mapById
}

type UserAddressFilter struct {
	Or      *UserAddressFilter
	Ids     []int64
	UserIds []int64
	Pagination
}

func (f UserAddressFilter) composeFilter() ([]string, []interface{}) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	if len(f.Ids) > 0 {
		query = append(query, "id IN (?)")
		args = append(args, f.Ids)
	}

	if len(f.UserIds) > 0 {
		query = append(query, "user_id IN (?)")
		args = append(args, f.UserIds)
	}

	return query, args
}

func (f UserAddressFilter) ComposeFilter() (string, []interface{}, error) {
	var (
		query []string
		args  = make([]interface{}, 0)
	)

	query, args = f.composeFilter()
	if f.Or != nil {
		queryOr, argsOr := f.Or.composeFilter()

		args = append(args, argsOr...)
		whereClauseOr := strings.Join(queryOr, " OR ")
		query = append(query, "("+whereClauseOr+")")
	}

	// Combine query parts into a single WHERE clause
	whereClause := strings.Join(query, " AND ")
	if whereClause != "" {
		whereClause = "WHERE " + whereClause
	}

	return whereClause, args, nil
}

type MapUserAddress map[int64]UserAddress
