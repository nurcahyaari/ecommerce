package entity

import (
	"reflect"
	"testing"

	"gopkg.in/guregu/null.v4"
)

func TestUsers_One(t *testing.T) {
	tests := []struct {
		name  string
		us    Users
		want  User
		want1 bool
	}{
		{
			name:  "test1",
			us:    Users{},
			want:  User{},
			want1: false,
		},
		{
			name: "test1",
			us: Users{
				User{
					Id:   1,
					Name: "test",
				},
			},
			want: User{
				Id:   1,
				Name: "test",
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.us.One()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.One() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Users.One() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUserFilter_composeFilter(t *testing.T) {
	type fields struct {
		Or         *UserFilter
		Ids        []int64
		Phone      null.String
		Email      null.String
		Pagination Pagination
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
		want1  []interface{}
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := UserFilter{
				Or:         tt.fields.Or,
				Ids:        tt.fields.Ids,
				Phone:      tt.fields.Phone,
				Email:      tt.fields.Email,
				Pagination: tt.fields.Pagination,
			}
			got, got1 := f.composeFilter()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserFilter.composeFilter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("UserFilter.composeFilter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUserFilter_ComposeFilter(t *testing.T) {
	type fields struct {
		Or         *UserFilter
		Ids        []int64
		Phone      null.String
		Email      null.String
		Pagination Pagination
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			name: "filter1",
			fields: fields{
				Ids:   []int64{1, 2, 3},
				Email: null.StringFrom("test@test.com"),
			},
			want:    "WHERE id IN (?) AND email = ?",
			want1:   []interface{}{[]int64{1, 2, 3}, "test@test.com"},
			wantErr: false,
		},
		{
			name: "filter2 - Or",
			fields: fields{
				Ids: []int64{1, 2, 3},
				Or: &UserFilter{
					Email: null.StringFrom("test@test.com"),
					Phone: null.StringFrom("1"),
				},
			},
			want:    "WHERE id IN (?) AND (phone = ? OR email = ?)",
			want1:   []interface{}{[]int64{1, 2, 3}, "1", "test@test.com"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := UserFilter{
				Or:         tt.fields.Or,
				Ids:        tt.fields.Ids,
				Phone:      tt.fields.Phone,
				Email:      tt.fields.Email,
				Pagination: tt.fields.Pagination,
			}
			got, got1, err := f.ComposeFilter()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserFilter.ComposeFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserFilter.ComposeFilter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("UserFilter.ComposeFilter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
