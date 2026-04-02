package dataaccess

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TAction struct {
	DbGorm *gorm.DB
}

func ConnnectAction(db *gorm.DB) IAction {
	return &TAction{db}
}

type TblAction struct {
	Id           int       `gorm:"primary_key;auto_increment"`
	EmailAddress string    `gorm:"column:EmailAddress"`
	MobileNumber string    `gorm:"column:MobileNumber"`
	RequestType  string    `gorm:"column:RequestType"`
	Message      string    `gorm:"column:Message"`
	Status       int       `gorm:"column:Status"`
	DateAdded    time.Time `gorm:"column:DateAdded"`
}

type IAction interface {
	CreateAction(prod TblAction) int
	GetAction(email, requesttype string) []TblAction
}

func (cn TAction) CreateAction(temp TblAction) int {
	if doInsert := cn.DbGorm.Table("TblAction").Create(&temp).Error; doInsert == nil {
		return temp.Id
	} else {
		logrus.Error(doInsert)
		return 0
	}
}

func (cn TAction) GetAction(email, requesttype string) []TblAction {
	resp := []TblAction{}
	cn.DbGorm.Table("TblAction").Select("Id", "EmailAddress", "MobileNumber", "RequestType", "Message", "Status", "DateAdded").Where("\"EmailAddress\"=? and \"RequestType\"=? ", email, requesttype).Find(&resp)
	return resp
}
