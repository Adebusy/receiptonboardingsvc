package dataaccess

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TConnect struct {
	DbGorm *gorm.DB
}

func ConnnectTemp(db *gorm.DB) ITemp {
	return &TConnect{db}
}

type TblTempPassword struct {
	Id           int       `gorm:"primary_key;auto_increment"`
	EmailAddress string    `gorm:"column:EmailAddress"`
	TempPassword string    `gorm:"column:TempPassword"`
	Status       int       `gorm:"column:Status"`
	DateAdded    time.Time `gorm:"column:DateAdded"`
}

type TempResp struct {
	Id        int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	Name      string    `json:"Name" validate:"omitempty"`
	Status    bool      `json:"Status" validate:"omitempty"`
	CreatedAt time.Time `json:"CreatedAt" validate:"omitempty"`
}

type ITemp interface {
	CreateTempPassword(prod TblTempPassword) int
	CheckTokenwithEmail(emailAddress, token string) int
}

func (cn TConnect) CreateTempPassword(temp TblTempPassword) int {

	if doInsert := cn.DbGorm.Table("TblTempPassword").Create(&temp).Error; doInsert == nil {
		return temp.Id
	} else {
		logrus.Error(doInsert)
		return 0
	}
}

func (cn TConnect) CheckTokenwithEmail(emailAddress, token string) int {
	tokentbl := TblTempPassword{}
	cn.DbGorm.Table("TblTempPassword").Select("Id", "EmailAddress", "TempPassword", "Status", "DateAdded").Where("\"EmailAddress\"=? and \"TempPassword\"=?", emailAddress, token).First(&tokentbl)
	return tokentbl.Id
}
