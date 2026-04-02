package dataaccess

import (
	"time"

	"gorm.io/gorm"
)

func ConnectClient(db *gorm.DB) IClient {
	return &DbConnect{db}
}

type IClient interface {
	GetClientByName(clientName string) TblClient
	RegisterNewClient(req TblClient) string
}

type TblClient struct {
	Id          int       `json:"Id" validate:"omitempty"`
	Name        string    `json:"Name" validate:"omitempty"`
	Status      int       `json:"Status" validate:"omitempty"`
	Description string    `json:"Description" validate:"omitempty"`
	CreatedAt   time.Time `gorm:"column:CreatedAt;autoCreateTime"`
}

type ClientRequest struct {
	Name        string `json:"Name" validate:"omitempty"`
	Description string `json:"Description" validate:"omitempty"`
}

type ClientResp struct {
	Id        int    `json:"Id" validate:"omitempty"`
	Name      string `json:"Name" validate:"omitempty"`
	RespToken string `json:"respToken" validate:"omitempty"`
}

func (cn DbConnect) GetClientByName(clientName string) TblClient {
	res := TblClient{}
	cn.DbGorm.Table("TblClient").Debug().Select("Id", "Name", "Status", "Description", "DateAdded").Where("\"Name\"=? and \"Status\"=1", clientName).First(&res)
	return res
}

func (cn DbConnect) RegisterNewClient(req TblClient) string {
	if doinssert := cn.DbGorm.Table("TblClient").Create(&req).Error; doinssert != nil {
		return "01"
	} else {
		return "00"
	}
}
