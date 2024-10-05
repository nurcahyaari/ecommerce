package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type Signature struct {
	IsActived bool        `db:"is_actived"`
	IsRemoved bool        `db:"is_removed"`
	CreatedAt time.Time   `db:"created_at"`
	CreatedBy null.String `db:"created_by"`
	UpdatedAt time.Time   `db:"updated_at"`
	UpdatedBy null.String `db:"updated_by"`
	DeletedAt null.Time   `db:"deleted_at"`
	DeletedBy null.String `db:"deleted_by"`
}
