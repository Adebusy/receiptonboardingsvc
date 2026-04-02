package dataaccess

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TblCompanyNameSuggested struct {
	Id          int       `gorm:"primary_key;auto_increment"`
	UserId      int       `gorm:"column:UserId"`
	CompanyName string    `gorm:"column:CompanyName"`
	Status      int       `gorm:"column:Status;default:1"`
	CreatedAt   time.Time `gorm:"column:CreatedAt;autoCreateTime"`
}

type TblCompanylogosPrev struct {
	Id              int    `gorm:"primary_key;auto_increment"`
	CompanyId       int    `gorm:"column:CompanyId"`
	CompanyName     string `gorm:"column:CompanyName"`
	EmailAddress    string `gorm:"column:EmailAddress"`
	CompanyLogoPath string `gorm:"column:CompanyLogoPath"`
	Status          int    `gorm:"column:Status"`
	CreatedAt       string `gorm:"column:CreatedAt"`
}

// type TblCompanySignature struct {
// 	Id                           int       `gorm:"primary_key;auto_increment"`
// 	CompanyId                    int       `gorm:"column:CompanyId"`
// 	EmailAddress                 string    `gorm:"column:EmailAddress"`
// 	CompanySignaturePath         string    `gorm:"column:CompanySignaturePath"`
// 	CompanySignaturePathOnDevice string    `gorm:"column:CompanySignaturePathOnDevice"`
// 	Status                       int       `gorm:"column:Status"`
// 	CreatedAt                    time.Time `gorm:"column:CreatedAt;autoCreateTime"`
// }

type CompanySignature struct {
	CompanyId                    int    `gorm:"column:CompanyId"`
	EmailAddress                 string `gorm:"column:EmailAddress"`
	CompanySignaturePath         string `gorm:"column:CompanySignaturePath"`
	CompanySignaturePathOnDevice string `gorm:"column:CompanySignaturePathOnDevice"`
	SignatureType                string `gorm:"column:SignatureType"`
	SignatureText                string `gorm:"column:SignatureText"`
	SignatureImage               string `gorm:"column:SignatureImage"`
}

type TblCompany struct {
	Id                 int                 `gorm:"primary_key;auto_increment"`
	CompanyName        string              `gorm:"column:CompanyName" validate:"required, min=8"`
	CompanyAddress     string              `gorm:"column:CompanyAddress"`
	EmailAddress       string              `gorm:"column:EmailAddress" validate:"required,email"`
	MobileNumber       string              `gorm:"column:MobileNumber" validate:"required,min=8"`
	RegNo              string              `gorm:"column:RegNo"`
	CompanyLogoPath    string              `gorm:"column:CompanyLogoPath"`
	CompanyReceiptPath string              `gorm:"column:CompanyReceiptPath"`
	Status             int                 `gorm:"column:Status;default:1"`
	CreatedAt          time.Time           `gorm:"column:CreatedAt;autoCreateTime"`
	UserId             int                 `gorm:"column:UserId;not null"`             // foreign key column
	User               TblUser             `gorm:"foreignKey:UserId;references:Id"`    // association
	LogoPrev           TblCompanylogosPrev `gorm:"foreignKey:CompanyId;references:Id"` // association
	SignatureRequest   TblSignatureRequest `gorm:"foreignKey:CompanyId;references:Id"` // association
}

//CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP

type Company struct {
	UserId             int    `json:"UserId"`
	CompanyName        string `json:"CompanyName" validate:"required,min=4"`
	CompanyAddress     string `json:"CompanyAddress"`
	EmailAddress       string `json:"EmailAddress" validate:"required,email"`
	MobileNumber       string `json:"MobileNumber" validate:"required,min=8"`
	RegNo              string `json:"RegNo"`
	CompanyLogoPath    string `json:"CompanyLogoPath"`
	CompanyReceiptPath string `json:"CompanyReceiptPath"`
}

type SuggestedCompanyNames struct {
	CompanyName string `json:"CompanyName"`
	UserId      int    `json:"UserId"`
	Status      int    `json:"Status"`
}

type Icompany interface {
	CreateCompany(companydetails TblCompany) string
	GetCompanyById(UserId int) Company
	GetCompanyByEmailAddress(EmailAddress string) Company
	GetCompanyByMobileNumber(mobile string) Company
	CheckCompanyByMobileNumber(MobileNumber string) (string, error)
	GetAllCompanies() []Company
	CreateCompanylogosPrev(logosPrev TblCompanylogosPrev) string
	GetCompanyDetailsByCompanyName(companyName string) Company
	CreateNameSuggested(companydetails TblCompanyNameSuggested) string
	FetchSuggestedCompanyNames(UserId int) []SuggestedCompanyNames
}

func CompanyDeal(db *gorm.DB) Icompany {
	return &DbConnect{db}
}

func (cn DbConnect) CreateNameSuggested(companydetails TblCompanyNameSuggested) string {
	if doinssert := cn.DbGorm.Table("TblCompanyNameSuggested").Create(&companydetails).Error; doinssert != nil {
		logrus.Error(doinssert)
		return "01"
	} else {
		return "00"
	}
}

func (cn DbConnect) FetchSuggestedCompanyNames(userId int) []SuggestedCompanyNames {
	res := []SuggestedCompanyNames{}
	cn.DbGorm.Table("TblCompanyNameSuggested").Find(&res).Where("UserId=?", userId)
	return res
}

func (cn DbConnect) CreateCompany(companydetails TblCompany) string {
	if doinssert := cn.DbGorm.Table("TblCompany").Create(&companydetails).Error; doinssert != nil {
		logrus.Error(doinssert)
		return "Unable to create companu at the moment!!"
	} else {
		return "Company created successfully!!"
	}
}

func (cn DbConnect) GetCompanyById(UserId int) Company {
	res := Company{}
	cn.DbGorm.Table("TblCompany").Select("Id", "CompanyName", "CompanyAddress", "EmailAddress", "MobileNumber", "RegNo", "CompanyLogoPath", "CompanyReceiptPath", "Status", "CreatedAt").Where("\"Id\"=?", UserId).First(&res)
	return res
}

func (cn DbConnect) GetCompanyDetailsByCompanyName(companyName string) Company {
	res := Company{}
	cn.DbGorm.Table("TblCompany").Select("Id", "CompanyName", "CompanyAddress", "EmailAddress", "MobileNumber", "RegNo", "CompanyLogoPath", "CompanyReceiptPath", "Status", "CreatedAt").Where("\"CompanyName\"=?", companyName).First(&res)
	return res
}

func (cn DbConnect) GetCompanyByEmailAddress(EmailAddress string) Company {
	res := Company{}
	cn.DbGorm.Table("TblCompany").Debug().Select("CompanyName", "CompanyAddress", "EmailAddress", "MobileNumber", "RegNo", "CompanyLogoPath", "CompanyReceiptPath", "'Status'", "'CreatedAt'").Where("\"EmailAddress\"=?", EmailAddress).First(&res)
	return res
}

func (cn DbConnect) GetCompanyByMobileNumber(MobileNumber string) Company {
	res := Company{}
	cn.DbGorm.Table("TblCompany").Debug().Select("CompanyName", "CompanyAddress", "EmailAddress", "MobileNumber", "RegNo", "CompanyLogoPath", "CompanyReceiptPath", "'Status'", "'CreatedAt'").Where("\"MobileNumber\"=?", MobileNumber).First(&res)
	return res
}

func (cn DbConnect) CheckCompanyByMobileNumber(MobileNumber string) (string, error) {
	res := Company{}
	result := cn.DbGorm.Table("TblCompany").Select("CompanyName", "CompanyAddress", "EmailAddress", "MobileNumber", "RegNo", "CompanyLogoPath", "CompanyReceiptPath", "'Status'", "'CreatedAt'").Where("\"MobileNumber\"=?", MobileNumber).First(&res)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", nil // not found
		}
		fmt.Printf(result.Error.Error())
		return result.Error.Error(), result.Error // real DB error
	}
	return "00", nil
}

func (cn DbConnect) GetCompanyByMobileNumbers(mobile string) (Company, error) {
	var res Company
	result := cn.DbGorm.
		Table("TblCompany").
		Where(`"MobileNumber" = ?`, mobile).
		First(&res)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return res, nil // not found
		}
		return res, result.Error // real DB error
	}
	return res, nil
}

func (cn DbConnect) GetAllCompanies() []Company {
	res := []Company{}
	cn.DbGorm.Table("TblCompany").Find(&res)
	return res
}

func (cn DbConnect) CreateCompanylogosPrev(logosPrev TblCompanylogosPrev) string {
	if doinssert := cn.DbGorm.Table("TblCompanylogosPrev").Create(logosPrev).Error; doinssert != nil {
		logrus.Error(doinssert)
		return "01"
	} else {
		return "00"
	}
}
