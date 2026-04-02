package obj

type UserObj struct {
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	Email        string `json:"email" validate:"required,email"`
	MobileNumber string `json:"MobileNumber" validate:"required,min=8"`
	Status       string `json:"Status" validate:"required,min=1"`
	Password     string `json:"Password" validate:"required,min=8"`
}

type CompanyForSuggest struct {
	UserId       int    `json:"UserId"`
	EmailAddress string `json:"EmailAddress"`
	CompanyName  string `json:"CompanyName" validate:"required,min=4"`
	FieldsArea   string `json:"FieldsArea"`
}

type CompanyForSignature struct {
	CompanyName string `json:"CompanyName" validate:"required,min=4"`
}

type ReceiptForSuggest struct {
	EmailAddress         string `json:"EmailAddress"`
	CompanyName          string `json:"CompanyName"`
	FieldsArea           string `json:"FieldsArea"`
	CompanyLogoPath      string `json:"CompanyLogoPath"`
	ManagerSignaturePath string `json:"ManagerSignaturePath"`
}

type ReceiptSelect struct {
	CompanyName string `json:"CompanyName"`
	ReceiptName string `json:"ReceiptName"`
}

type ReceiptRequest struct {
	EmailAddress string `json:"EmailAddress"`
	CompanyName  string `json:"CompanyName"`
	CustumerName string `json:"CustumerName"`
	SalesDetails string `json:"SalesDetails"`
}

type SignUp struct {
	FirstName    string `json:"FirstName" validate:"required"`
	LastName     string `json:"LastName" validate:"required"`
	Email        string `json:"Email" validate:"required,email"`
	MobileNumber string `json:"MobileNumber" validate:"required,min=8"`
	Password     string `json:"Password" validate:"required,min=8"`
	SignUpBy     string `json:"SignUpBy" validate:"required"`
}

type LogOut struct {
	Email        string `json:"Email" validate:"required,email"`
	MobileNumber string `json:"MobileNumber" validate:"required,min=8"`
	Password     string `json:"Password" validate:"required,min=8"`
}

type LogInReq struct {
	Email    string `json:"Email" validate:"required,email"`
	Password string `json:"Password" validate:"required,min=8"`
}

type CompleteSignUp struct {
	EmailAddress string `gorm:"column:EmailAddress"`
	MobileNumber string `gorm:"column:MobileNumber"`
	FirstName    string `gorm:"column:FirstName"`
	LastName     string `gorm:"column:LastName"`
}

type UserResponse struct {
	Id             uint
	FirstName      string
	LastName       string
	Email          string
	MobileNumber   string
	Status         string
	CreatedAt      string
	CompanyName    string
	SignatureType  string
	Signature      string
	LogoType       string
	LogoPath       string
	SignatureText  string
	SignatureImage string
}

type EmailObj struct {
	ToEmail  string `gorm:"column:ToEmail"`
	Subject  string `gorm:"column:Subject"`
	MailBody string `gorm:"column:MailBody"`
}

type ChangePassword struct {
	UserName        string `gorm:"column:UserName"`
	CurrentPassword string `gorm:"column:CurrentPassword"`
	NewPassword     string `gorm:"column:NewPassword"`
}

type SignatureRequest struct {
	EmailAddress   string `json:"EmailAddress" validate:"required,email"`
	PathOnDevice   string `json:"PathOnDevice"`
	SignatureType  string `json:"SignatureType"`
	SignatureText  string `json:"SignatureText"`
	SignatureImage string `json:"SignatureImage"`
}

type ResponseMessage struct {
	ResponseCode    string
	ResponseMessage string
}

type ConfigStruct struct {
	CreateTable          bool
	IsDropExistingTables bool
}

type TokenResp struct {
	Token string
}

type TemplateData struct {
	CompanyName string
	Signature   string
	LogoURL     string
}

type LogoBase64 struct {
	Extension string `json:"extension"`
	Data      string `json:"data"`
}

type LogoConcept struct {
	LogoID     string     `json:"logo_id"`
	Concept    string     `json:"concept_name"`
	LogoBase64 LogoBase64 `json:"logo_base64"`
}
