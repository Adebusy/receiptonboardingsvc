package dataaccess

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TblCompanySignature struct {
	Id                           int       `gorm:"primary_key;auto_increment"`
	CompanyId                    int       `gorm:"column:CompanyId"`
	EmailAddress                 string    `gorm:"column:EmailAddress"`
	SignatureText                string    `gorm:"column:SignatureText"`
	SignatureImage               string    `gorm:"column:SignatureImage"`
	CompanySignaturePath         string    `gorm:"column:CompanySignaturePath"`
	CompanySignaturePathOnDevice string    `gorm:"column:CompanySignaturePathOnDevice"`
	SignatureType                string    `gorm:"column:SignatureType"`
	Status                       int       `gorm:"column:Status;default:1"`
	CreatedAt                    time.Time `gorm:"column:CreatedAt;autoCreateTime"`
}

type SignatureRequest struct {
	CompanyName  string `json:"CompanyName"`
	EmailAddress string `json:"EmailAddress"`
}

type ISignature interface {
	CreateCompanySignature(companydetails CompanySignature) string
	GetSignatureByEmailAddress(EmailAddress string) CompanySignature
}

func SinatureDeal(db *gorm.DB) ISignature {
	return &DbConnect{db}
}

func (cn DbConnect) CreateCompanySignature(companydetails CompanySignature) string {
	if doinssert := cn.DbGorm.Table("TblCompanySignature").Create(&companydetails).Error; doinssert != nil {
		logrus.Error(doinssert)
		return "01"
	} else {
		return "00"
	}
}

func (cn DbConnect) GetSignatureByEmailAddress(EmailAddress string) CompanySignature {
	res := CompanySignature{}
	cn.DbGorm.Table("TblCompanySignature").Debug().Select("CompanyId", "CompanySignaturePath", "CompanySignaturePathOnDevice", "SignatureType").Where("\"EmailAddress\"=?", EmailAddress).First(&res)
	return res
}
