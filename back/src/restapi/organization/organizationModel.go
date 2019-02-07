package organization

import (
	// import to have gorm types
	_ "github.com/jinzhu/gorm"

	user "restapi/user"
)

// Organization database representation
type Organization struct {
	ID        uint           `gorm:"AUTO_INCREMENT" json:"id"`
	Title     string         `json:"name"`
	Employees []user.Profile `gorm:"many2many:organization_user" json:"employees"`
	is_active bool           `json:"is_active"`
}

// TableName sets table name of the struct
func (Organization) TableName() string {
	return "organization"
}
