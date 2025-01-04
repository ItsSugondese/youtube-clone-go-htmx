package model

type Role struct {
	ID string `json:"name" gorm:"primarykey"`
}

func (b *Role) HasAuditModel() bool {
	return false
}
