package dataaccess

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type User struct {
	Id           int    `gorm:"column:Id"`
	FirstName    string `gorm:"column:FirstName"`
	LastName     string `gorm:"column:LastName"`
	EmailAddress string `gorm:"column:EmailAddress"`
	MobileNumber string `gorm:"column:MobileNumber"`
	Password     string `gorm:"column:Password"`
	Status       string `gorm:"column:Status"`
	CreatedAt    string `gorm:"column:CreatedAt"`
}

type TblSignatureRequest struct {
	Id              int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	CompanyId       int       `gorm:"column:CompanyId"`
	EmailAddress    string    `gorm:"column:EmailAddress"`
	SignatureBase64 string    `json:"signature_base64"`
	Status          int       `gorm:"column:Status"`
	CreatedAt       time.Time `gorm:"column:CreatedAt"`
}

type TblStatus struct {
	Id         int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	StatusName string    `json:"StatusName" validate:"omitempty"`
	CreatedAt  time.Time `json:"CreatedAt" validate:"omitempty"`
}

type TblRole struct {
	Id        int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	RoleName  string    `json:"RoleName" validate:"omitempty"`
	Status    bool      `json:"Status" validate:"omitempty"`
	CreatedAt time.Time `json:"CreatedAt" validate:"omitempty"`
}

type CompleteSignUpReq struct {
	EmailAddress string `gorm:"column:EmailAddress"`
	FirstName    string `gorm:"column:FirstName"`
	LastName     string `gorm:"column:LastName"`
	MobileNumber string `gorm:"column:MobileNumber"`
	Status       int    `gorm:"column:Status"`
	CreatedAt    string `gorm:"column:CreatedAt"`
}

type UpdateStatus struct {
	EmailAddress string `gorm:"column:EmailAddress"`
	Status       int    `gorm:"column:Status"`
	CreatedAt    string `gorm:"column:CreatedAt"`
}

type TblUser struct {
	Id           int       `json:"Id" gorm:"unique;primaryKey;autoIncrement"`
	FirstName    string    `gorm:"column:FirstName"`
	LastName     string    `gorm:"column:LastName"`
	EmailAddress string    `gorm:"column:EmailAddress" validate:"required,email"`
	MobileNumber string    `gorm:"column:MobileNumber" validate:"required,min=8"`
	Password     string    `gorm:"column:Password"`
	Status       int       `gorm:"column:Status;default:1"`
	CreatedAt    time.Time `gorm:"column:CreatedAt;autoCreateTime"`

	Companies []TblCompany `gorm:"foreignKey:UserId"`
}

type ResponseMessage struct {
	ResponseCode    string
	ResponseMessage string
}

type DbConnect struct {
	DbGorm *gorm.DB
}

func ConneectDeal(db *gorm.DB) Iuser {
	return &DbConnect{db}
}

type Iuser interface {
	CreateUser(usr *User) string
	UpdateUserRecord(usr CompleteSignUpReq) string
	SignUp(firstName, lastName, emailAddress, mobileNumber, password, createdAt string) string
	GetUserByEmailAddress(EmailAddress string) User
	GetUserByUsername(EmailAddress string) User
	GetUserByMobileNumber(MobileNumber string) User
	LoginUser(UserName, Password string) User
	GetUserByUserId(UserId int) User
	UpdateUserStatusByUserEmail(usr UpdateStatus) string
	GetClientByName(clientName string) TblClient
	RegisterNewClient(req TblClient) string

	LogOut(token, username string) string

	ChangePassword(emailAddress, mobileNumber, password string) int
}

func (cn DbConnect) LogOut(token, username string) string {
	if doinssert := cn.DbGorm.Table("TblBlacklisted").Create("&usr").Error; doinssert != nil {
		logrus.Error(doinssert)
		return "Unable to create user at the moment!!"
	} else {
		return "User created successfully!!"
	}
}

func (cn DbConnect) CreateUser(usr *User) string {
	if doinssert := cn.DbGorm.Table("TblUser").Create(&usr).Error; doinssert != nil {
		logrus.Error(doinssert)
		return "Unable to create user at the moment!!"
	} else {
		return "User created successfully!!"
	}
}

func (cn DbConnect) UpdateUserRecord(usr CompleteSignUpReq) string {

	if doinssertupdate := cn.DbGorm.Table("TblUser").Debug().Where("\"EmailAddress\"=? or \"MobileNumber\"=?", usr.EmailAddress, usr.MobileNumber).Updates(&usr).Error; doinssertupdate != nil {
		logrus.Error(doinssertupdate)
		return "Unable to create user at the moment!!"
	} else {
		logrus.Error(fmt.Sprintf("UpdateUserRecord for %s", usr.EmailAddress))
		return "User created successfully!!"
	}
}

func (cn DbConnect) UpdateUserStatusByUserEmail(usr UpdateStatus) string {
	if doinssertupdate := cn.DbGorm.Table("TblUser").Debug().Where("\"EmailAddress\"=?", usr.EmailAddress).Update("Status", usr.Status).Error; doinssertupdate != nil {
		logrus.Error(doinssertupdate)
		return "Unable to create user at the moment!!"
	} else {
		logrus.Error(fmt.Sprintf("UpdateUserRecord for %s", usr.EmailAddress))
		return "User created successfully!!"
	}
}

func (cn DbConnect) ChangePassword(emailAddress, mobileNumber, password string) int {

	retval := 0
	if emailAddress != "" {
		if doinssertupdate := cn.DbGorm.Table("TblUser").Debug().Where("\"EmailAddress\"=? ", emailAddress).Update("Password", password).Error; doinssertupdate == nil {
			logrus.Error(doinssertupdate)
			retval = 1
			return retval
		} else {
			logrus.Error(fmt.Sprintf("UpdateUserRecord for %s", emailAddress))
			return retval
		}
	}

	if mobileNumber != "" {
		if doinssertupdate := cn.DbGorm.Table("TblUser").Debug().Where("\"MobileNumber\"=? ", mobileNumber).Update("Password", password).Error; doinssertupdate != nil {
			logrus.Error(doinssertupdate)
			retval = 1
			return retval
		} else {
			logrus.Error(fmt.Sprintf("UpdateUserRecord for %s", emailAddress))
			return retval
		}
	}
	return retval
}

func (cn DbConnect) SignUp(firstName, lastName, emailAddress, mobileNumber, password, createdAt string) string {

	if doinssert := cn.DbGorm.Table("TblUser").Select("FirstName", "LastName", "EmailAddress", "MobileNumber", "Password", "Status", "CreatedAt").Create(map[string]interface{}{"FirstName": firstName, "LastName": lastName, "EmailAddress": emailAddress, "MobileNumber": mobileNumber, "Password": password, "Status": "0", "CreatedAt": createdAt}).Error; doinssert != nil {
		logrus.Error(doinssert)
		return "Unable to create create sign up at the moment!!"
	} else {
		return "User signed up successfully!!"
	}
}

func (cn DbConnect) GetUserByUserId(UserId int) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("UserName", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Password", "Status", "CreatedAt").Where("\"Id\"=?", UserId).First(&res)
	return res
}

func (cn DbConnect) GetUserByEmailAddress(EmailAddress string) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Debug().Select("Id", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Password", "Status", "CreatedAt").Where("\"EmailAddress\"=?", EmailAddress).First(&res)
	return res
}

func (cn DbConnect) GetUserByMobileNumber(MobileNumber string) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("Id", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Password", "Status", "CreatedAt").Where("\"MobileNumber\"=?", MobileNumber).First(&res)
	return res
}

func (cn DbConnect) GetUserByUsername(username string) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("Id", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Password", "Status", "CreatedAt", "Password").Where("\"UserName\"=?", username).First(&res)
	return res
}

func (cn DbConnect) LoginUser(UserName, Password string) User {
	res := User{}
	cn.DbGorm.Table("TblUser").Select("Id", "FirstName", "LastName", "EmailAddress", "MobileNumber", "Status", "CreatedAt", "Password").Where("\"UserName\"=? and \"Password\"=?", UserName, Password).First(&res)
	return res
}
