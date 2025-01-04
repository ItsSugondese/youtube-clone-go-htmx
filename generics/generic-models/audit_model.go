package generic_models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AuditModel struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index" swaggerignore:"true"`
	CreatedBy  *string
	UpdatedBy  *string
	DeletedBy  *string
	IsDeleting bool `gorm:"-"`
}
