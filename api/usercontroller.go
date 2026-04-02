package api

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	dbSchema "github.com/Adebusy/receiptonboardingsvc/dataaccess"
	inpuschema "github.com/Adebusy/receiptonboardingsvc/obj"
	psg "github.com/Adebusy/receiptonboardingsvc/postgresql"
	"github.com/Adebusy/receiptonboardingsvc/utilities"
	"github.com/EDDYCJY/go-gin-example/pkg/upload"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var (
	usww       = dbSchema.ConneectDeal(getdb)
	TAct       = dbSchema.ConnnectAction(getdb)
	ISig       = dbSchema.SinatureDeal(getdb)
	validateMe = validator.New()
)

// SignUp godoc
// @Summary		SignUp new user.
// @Description	SignUp new user.
// @Tags			user
// @Accept			*/*
// @User			json
// @Param user body inpuschema.SignUp true "SignUp new user"
// @Success		200	{object}	dbSchema.User
// @Router			/api/user/SignUp [post]
func SignUp(ctx *gin.Context) {
	currentTime := time.Now()
	usww := dbSchema.ConneectDeal(psg.GetDB())
	reqIn := &inpuschema.SignUp{}
	if err := ctx.ShouldBindJSON(reqIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}

	//validate request
	if validateObj := validateMe.Struct(reqIn); validateObj != nil {
		ctx.JSON(http.StatusBadRequest, validateObj.Error())
		return
	}

	enc := hex.EncodeToString([]byte(reqIn.Password))

	if CheckEmailExist := usww.GetUserByEmailAddress(reqIn.Email); CheckEmailExist.EmailAddress != "" {
		ctx.JSON(http.StatusBadRequest, "User with email address "+reqIn.Email+" already exist!!")
		return
	}

	if CheckMobile := usww.GetUserByMobileNumber(reqIn.MobileNumber); CheckMobile.MobileNumber != "" {
		ctx.JSON(http.StatusBadRequest, "User with mobile number "+reqIn.MobileNumber+" already exist!!")
		return
	}

	doCreate := usww.SignUp(reqIn.FirstName, reqIn.LastName, reqIn.Email, reqIn.MobileNumber, enc, currentTime.Format("01-02-2006"))
	logrus.Info(doCreate)
	ctx.JSON(http.StatusOK, doCreate)
}

// CompleteSignUp godoc
// @Summary		CompleteSignUp user signup.
// @Description	CompleteSignUp user signup.
// @Tags			user
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Param user body inpuschema.CompleteSignUp true "CompleteSignUp user signup"
// @Success		200	{object}	dbSchema.ResponseMessage
// @Router			/api/user/CompleteSignUp [post]
// func CompleteSignUp(ctx *gin.Context) {
// 	if !ValidateClient(ctx) {
// 		return
// 	}
// 	usww := dbSchema.ConneectDeal(psg.GetDB())
// 	reqIn := &inpuschema.CompleteSignUp{}
// 	if err := ctx.ShouldBindJSON(reqIn); err != nil {
// 		ctx.JSON(http.StatusBadRequest, err.Error())
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	//validate request
// 	if validateObj := validateMe.Struct(reqIn); validateObj != nil {
// 		ctx.JSON(http.StatusBadRequest, validateObj.Error())
// 		return
// 	}

// 	req := dbSchema.CompleteSignUpReq{EmailAddress: reqIn.EmailAddress,
// 		FirstName:    reqIn.FirstName,
// 		LastName:     reqIn.LastName,
// 		Status:       1,
// 		MobileNumber: reqIn.MobileNumber,
// 		CreatedAt:    time.Now().Format("01-02-2006")}

// 	if CheckEmailExist := usww.GetUserByEmailAddress(req.EmailAddress); CheckEmailExist.EmailAddress == "" {
// 		ctx.JSON(http.StatusBadRequest, "User with email address "+req.EmailAddress+" does not exist!!")
// 		return
// 	}

// 	if CheckMobile := usww.GetUserByMobileNumber(req.MobileNumber); CheckMobile.MobileNumber == "" {
// 		ctx.JSON(http.StatusBadRequest, "User with mobile number "+req.MobileNumber+" does not exist!!")
// 		return
// 	}

// 	doCreate := usww.UpdateUserRecord(req)
// 	logrus.Info(doCreate)
// 	Response := &inpuschema.ResponseMessage{ResponseCode: "00",
// 		ResponseMessage: doCreate,
// 	}

// 	ctx.JSON(http.StatusOK, Response)
// }

// UpdateUserDetails godoc
// @Summary		Update User Details.
// @Description	Update User Details.
// @Tags			user
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Param user body inpuschema.CompleteSignUp true "Update User Details"
// @Success		200	{object}	dbSchema.ResponseMessage
// @Router			/api/user/UpdateUserDetails [post]
func UpdateUserDetails(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }
	usww := dbSchema.ConneectDeal(psg.GetDB())
	reqIn := &inpuschema.CompleteSignUp{}
	if err := ctx.ShouldBindJSON(reqIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}

	//validate request
	if validateObj := validateMe.Struct(reqIn); validateObj != nil {
		ctx.JSON(http.StatusBadRequest, validateObj.Error())
		return
	}

	req := dbSchema.CompleteSignUpReq{EmailAddress: reqIn.EmailAddress,
		FirstName:    reqIn.FirstName,
		Status:       1,
		MobileNumber: reqIn.MobileNumber,
		CreatedAt:    time.Now().Format("01-02-2006")}

	if CheckEmailExist := usww.GetUserByEmailAddress(req.EmailAddress); CheckEmailExist.EmailAddress == "" {
		ctx.JSON(http.StatusBadRequest, "User with email address "+req.EmailAddress+" does not exist!!")
		return
	}

	if CheckMobile := usww.GetUserByMobileNumber(req.MobileNumber); CheckMobile.MobileNumber == "" {
		ctx.JSON(http.StatusBadRequest, "User with mobile number "+req.MobileNumber+" does not exist!!")
		return
	}

	doCreate := usww.UpdateUserRecord(req)
	logrus.Info(doCreate)
	Response := &inpuschema.ResponseMessage{ResponseCode: "00",
		ResponseMessage: "records updated successfully!",
	}

	ctx.JSON(http.StatusOK, Response)
}

// ValidateAndSendTempPassword Validate email And Send Temp Password
// @Summary		Validate email and send temp password.
// @Description	Validate email and send temp password.
// @Tags			user
// @Param EmailAddress path string true "User email address"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.ResponseMessage
// @Router			/api/user/ValidateAndSendTempPassword/{EmailAddress} [get]
// func ValidateAndSendTempPassword(ctx *gin.Context) {
// 	res := inpuschema.ResponseMessage{}
// 	usww := dbSchema.ConnnectTemp(psg.GetDB())
// 	usw := dbSchema.ConneectDeal(psg.GetDB())
// 	requestEmail := ctx.Param("EmailAddress")
// 	getUSer := usw.GetUserByEmailAddress(requestEmail)
// 	logAction := fmt.Sprintf("ValidateAndSendTempPassword %v", requestEmail)
// 	logrus.Info(logAction)
// 	if getUSer.EmailAddress != "" {
// 		//generate temp password
// 		getToken := utilities.TempPassword(6, true, false, true)
// 		tempTable := dbSchema.TblTempPassword{
// 			EmailAddress: requestEmail,
// 			TempPassword: getToken,
// 			Status:       1,
// 			DateAdded:    time.Now(),
// 		}

// 		if creatTemp := usww.CreateTempPassword(tempTable); creatTemp != 0 {
// 			if sendNot := utilities.SendEmail(requestEmail, fmt.Sprintf("Here is a temporary password generate for your profile %s.", getToken)); sendNot == "00" {
// 				TAct.CreateAction(dbSchema.TblAction{EmailAddress: requestEmail, MobileNumber: "", RequestType: "Notification", Message: fmt.Sprintf("Here is a temporary password generate for your profile %s.", getToken), Status: 1, DateAdded: time.Now()})
// 				res.ResponseCode = "00"
// 				res.ResponseMessage = fmt.Sprintf("Here is a temporary password generate for your profile %s.", getToken)
// 				ctx.JSON(http.StatusOK, res)
// 				return
// 			} else {
// 				res.ResponseCode = "01"
// 				res.ResponseMessage = "service is unable to process request at the momment!!"
// 				ctx.JSON(http.StatusBadRequest, res)
// 				return
// 			}
// 		} else {
// 			res.ResponseCode = "01"
// 			res.ResponseMessage = "service is unable to process request at the momment!!"
// 			ctx.JSON(http.StatusBadRequest, res)
// 			return
// 		}
// 	} else {
// 		res.ResponseCode = "01"
// 		res.ResponseMessage = "service is unable to process request at the momment!!"
// 		ctx.JSON(http.StatusBadRequest, res)
// 		return
// 	}
// }

// ValidateTempToken Validate email And Send Temp Password
// @Summary		Validate email and send temp token.
// @Description	Validate email and send temp token.
// @Tags			user
// @Param EmailAddress path string true "User email address"
// @Param tempPassword path string true "temporary Password"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.ResponseMessage
// @Router			/api/user/ValidateTempToken/{EmailAddress}/{TempPassword}  [get]
// func ValidateTempToken(ctx *gin.Context) {
// 	res := inpuschema.ResponseMessage{}
// 	usww := dbSchema.ConnnectTemp(psg.GetDB())
// 	usw := dbSchema.ConneectDeal(psg.GetDB())
// 	requestEmail := ctx.Param("EmailAddress")
// 	tempPassword := ctx.Param("TempPassword")
// 	getUSer := usw.GetUserByEmailAddress(requestEmail)
// 	if getUSer.EmailAddress != "" {
// 		if creatTemp := usww.CheckTokenwithEmail(requestEmail, tempPassword); creatTemp != 0 {
// 			res.ResponseCode = "00"
// 			res.ResponseMessage = "Token validate successfully."
// 			ctx.JSON(http.StatusOK, res)
// 			return
// 		} else {
// 			res.ResponseCode = "01"
// 			res.ResponseMessage = "service is unable to process request at the momment!!"
// 			ctx.JSON(http.StatusBadRequest, res)
// 			return
// 		}
// 	} else {
// 		res.ResponseCode = "01"
// 		res.ResponseMessage = "Invalid email or token"
// 		ctx.JSON(http.StatusBadRequest, res)
// 		return
// 	}
// }

// GetUserByEmailAddress create new user
// @Summary		Get user by email address new cart user.
// @Description	Get user by email address new cart user.
// @Tags			user
// @Param EmailAddress path string true "User email address"
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/GetUserByEmailAddress/{EmailAddress} [get]
func GetUserByEmailAddress(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }
	requestEmail := ctx.Param("EmailAddress")
	getUSer := usww.GetUserByEmailAddress(requestEmail)
	logAction := fmt.Sprintf("GetUserByEmailAddress %v", requestEmail)
	logrus.Info(logAction)
	ctx.JSON(http.StatusOK, getUSer)
}

// GetUserByMobile existing user destails by mobile number
// @Summary		existing user destails by mobile number.
// @Description	existing user destails by mobile number.
// @Tags			user
// @Param MobileNumber path string true "User mobile number"
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	inpuschema.CartObj
// @Router			/api/user/GetUserByMobile/{MobileNumber} [get]
func GetUserByMobile(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
	userRespose := &inpuschema.UserResponse{}
	requestMobile := ctx.Param("MobileNumber")
	if getUSer := usww.GetUserByMobileNumber(requestMobile); getUSer.FirstName != "" {
		userRespose.FirstName = getUSer.FirstName
		userRespose.LastName = getUSer.LastName
		userRespose.Email = getUSer.EmailAddress
		userRespose.MobileNumber = getUSer.MobileNumber
		userRespose.Status = getUSer.Status
		userRespose.CreatedAt = getUSer.CreatedAt
		ctx.JSON(http.StatusOK, userRespose)
		return
	}

	logAction := fmt.Sprintf("GetUserByMobile %v", requestMobile)
	logrus.Info(logAction)
	ctx.JSON(http.StatusBadRequest, userRespose)
}

// LogIn exiting user In
// @Summary		Log user In with username and password.
// @Description	Log user In with username and password.
// @Tags			user
// @Param user body inpuschema.LogInReq true "LogInReq"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/LogIn/ [post]
func LogIn(ctx *gin.Context) {
	reqIn := &inpuschema.LogInReq{}
	if err := ctx.ShouldBindJSON(reqIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}

	getUSerobj := dbSchema.User{}
	userRespose := &inpuschema.UserResponse{}
	enc := hex.EncodeToString([]byte(reqIn.Password))

	if utilities.IsEmailValid(reqIn.Email) {
		getUSerobj = usww.GetUserByEmailAddress(reqIn.Email)
	} else if utilities.IsNumberValid(reqIn.Email) {
		getUSerobj = usww.GetUserByMobileNumber(reqIn.Email)
	} else {
		logAction := fmt.Sprintf("Incorrect username %s", reqIn.Email)
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}

	if getUSerobj.EmailAddress != "" || getUSerobj.MobileNumber != "" {
		if getUSerobj.Password == enc {
			userRespose.FirstName = getUSerobj.FirstName
			userRespose.LastName = getUSerobj.LastName
			userRespose.Email = getUSerobj.EmailAddress
			userRespose.MobileNumber = getUSerobj.MobileNumber
			userRespose.Status = getUSerobj.Status
			userRespose.CreatedAt = getUSerobj.CreatedAt
			userRespose.Id = uint(getUSerobj.Id)
			if valcompany, _ := strconv.Atoi(getUSerobj.Status); valcompany >= 2 {
				//get companyName
				getCompany := comp.GetCompanyByEmailAddress(getUSerobj.EmailAddress)
				userRespose.CompanyName = getCompany.CompanyName
			}
			if val, _ := strconv.Atoi(getUSerobj.Status); val >= 3 {
				//get signature
				getSignature := ISig.GetSignatureByEmailAddress(getUSerobj.EmailAddress)
				userRespose.Signature = getSignature.CompanySignaturePathOnDevice
				userRespose.SignatureType = getSignature.SignatureType
				userRespose.SignatureText = getSignature.SignatureText
				userRespose.SignatureImage = getSignature.SignatureImage
			}

			//get logo
			if val, _ := strconv.Atoi(getUSerobj.Status); val >= 4 {

			}

			logrus.Info(fmt.Sprintf("LogIn for user %s", reqIn.Email))
			ctx.JSON(http.StatusOK, userRespose)
			return
		} else {
			logAction := fmt.Sprintf("Incorrect password %s", reqIn.Email)
			logrus.Info(logAction)
			ctx.JSON(http.StatusBadRequest, logAction)
			return
		}
	} else {
		logAction := fmt.Sprintf("Incorrect username %s", reqIn.Email)
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}
}

// LogOut exiting user In
// @Summary		Log user Out with username and password.
// @Description	Log user Out with username and password.
// @Tags			user
// @Param UserName path string true "Username"
// @Param Password path string true "Password"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	string
// @Router			/api/user/LogOut/{UserName}/{Password} [get]
func LogOut(ctx *gin.Context) {
	getUSerobj := dbSchema.User{}
	userRespose := &inpuschema.UserResponse{}
	UserName := ctx.Param("UserName")
	Password := ctx.Param("Password")
	enc := hex.EncodeToString([]byte(Password))

	if utilities.IsEmailValid(UserName) {
		getUSerobj = usww.GetUserByEmailAddress(UserName)
	} else if utilities.IsNumberValid(UserName) {
		getUSerobj = usww.GetUserByMobileNumber(UserName)
	} else {
		logAction := fmt.Sprintf("Incorrect username %s", UserName)
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}

	if getUSerobj.EmailAddress != "" || getUSerobj.MobileNumber != "" {
		if getUSerobj.Password == enc {
			userRespose.FirstName = getUSerobj.FirstName
			userRespose.LastName = getUSerobj.LastName
			userRespose.Email = getUSerobj.EmailAddress
			userRespose.MobileNumber = getUSerobj.MobileNumber
			userRespose.Status = getUSerobj.Status
			userRespose.CreatedAt = getUSerobj.CreatedAt
			userRespose.Id = uint(getUSerobj.Id)
			// if newToken := CreateOrGetToken(userRespose.Email); newToken != "" {
			// 	userRespose.Token = newToken
			// }
			logrus.Info(fmt.Sprintf("LogIn for user %s", UserName))
			ctx.JSON(http.StatusOK, userRespose)
			return
		} else {
			logAction := fmt.Sprintf("Incorrect password %s", UserName)
			logrus.Info(logAction)
			ctx.JSON(http.StatusBadRequest, logAction)
			return
		}
	} else {
		logAction := fmt.Sprintf("Incorrect username %s", UserName)
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}
}

// LogInWithMobileNumber for exiting user
// @Summary		Log user In with mobile number and password.
// @Description	Log user In with mobile number and password.
// @Tags			user
// @Param MobileNumber path string true "MobileNumber"
// @Param Password path string true "Password"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/LogInWithMobileNumber/{MobileNumber}/{Password} [get]
func LogInWithMobileNumber(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }
	userRespose := &inpuschema.UserResponse{}
	MobileNumber := ctx.Param("MobileNumber")
	Password := ctx.Param("Password")
	password, _ := utilities.HashPassword(Password)

	if getUSer := usww.GetUserByMobileNumber(MobileNumber); getUSer.EmailAddress != "" {
		if utilities.CheckPasswordHash(Password, password) {
			userRespose.FirstName = getUSer.FirstName
			userRespose.LastName = getUSer.LastName
			userRespose.Email = getUSer.EmailAddress
			userRespose.MobileNumber = getUSer.MobileNumber
			userRespose.Status = getUSer.Status
			userRespose.CreatedAt = getUSer.CreatedAt
			logrus.Info(fmt.Sprintf("LogIn for user with LogInWithMobileNumber %s", MobileNumber))
			ctx.JSON(http.StatusOK, userRespose)
			return
		} else {
			logAction := fmt.Sprintf("Incorrect password %s", MobileNumber)
			logrus.Info(logAction)
			ctx.JSON(http.StatusBadRequest, logAction)
			return
		}
	} else {
		logAction := fmt.Sprintf("Incorrect username %s", MobileNumber)
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}
}

// LogInWithEmailAddress for exiting user
// @Summary		Log user In with email address and password.
// @Description	Log user In with email address and password.
// @Tags			user
// @Param EmailAddress path string true "EmailAddress"
// @Param Password path string true "Password"
// @Produce json
// @Accept			*/*
// @User			json
// @Success		200	{object}	inpuschema.UserResponse
// @Router			/api/user/LogInWithEmailAddress/{EmailAddress}/{Password} [get]
func LogInWithEmailAddress(ctx *gin.Context) {
	userRespose := &inpuschema.UserResponse{}
	EmailAddress := ctx.Param("EmailAddress")
	Password := ctx.Param("Password")
	password, _ := utilities.HashPassword(Password)

	if getUSer := usww.GetUserByEmailAddress(EmailAddress); getUSer.EmailAddress != "" {
		if utilities.CheckPasswordHash(Password, password) {
			userRespose.FirstName = getUSer.FirstName
			userRespose.LastName = getUSer.LastName
			userRespose.Email = getUSer.EmailAddress
			userRespose.MobileNumber = getUSer.MobileNumber
			userRespose.Status = getUSer.Status
			userRespose.CreatedAt = getUSer.CreatedAt
			logrus.Info(fmt.Sprintf("LogIn for user with LogInWithEmailAddress %s", EmailAddress))
			ctx.JSON(http.StatusOK, userRespose)
			return
		} else {
			logAction := fmt.Sprintf("Incorrect password %s", EmailAddress)
			logrus.Info(logAction)
			ctx.JSON(http.StatusBadRequest, logAction)
			return
		}
	} else {
		logAction := fmt.Sprintf("Incorrect username %s", EmailAddress)
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}
}

// SendEmail godoc
// @Summary		Send Email.
// @Description	Send Email.
// @Tags			user
// @Accept			*/*
// @User			json
// @Param user body inpuschema.EmailObj true "Send Email"
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{string}	string "Email sent successfully!!"
// @Failure		400		{string} string	"Unable to send email at the monent!!"
// @Router			/api/user/SendEmail [post]
func SendEmail(ctx *gin.Context) {
	if !ValidateClient(ctx) {
		return
	}
	reqIn := &inpuschema.EmailObj{}
	if err := ctx.ShouldBindJSON(reqIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	sender := os.Getenv("SMTP_SENDER")
	recipient := reqIn.ToEmail

	from := "From: " + sender + "\n"
	to := "To: " + recipient + "\n"
	subject := "Subject: Digital company update\n"
	body := reqIn.MailBody
	message := []byte(from + to + subject + "\n" + body)
	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, []string{recipient}, message)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
		resp := fmt.Sprintf("Error sending email: %v", err)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	logAction := fmt.Sprintf("SendEmail to %v", recipient)
	logrus.Info(logAction)
	ctx.JSON(http.StatusOK, "Email sent successfully!!!")
}

// ChangePassword godoc
// @Summary		ChangePassword user password.
// @Description	ChangePassword user password.
// @Tags			user
// @Accept			*/*
// @User			json
// @Param user body inpuschema.ChangePassword true "Update password"
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{string}	string "Password updated successfully!!"
// @Failure		400		{string} string	"Unable to update password at the monent!!"
// @Router			/api/user/ChangePassword [post]
func ChangePassword(ctx *gin.Context) {
	reqIn := &inpuschema.ChangePassword{}
	if err := ctx.ShouldBindJSON(reqIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}

	getUSerobj := dbSchema.User{}
	userRespose := &inpuschema.UserResponse{}
	enc := hex.EncodeToString([]byte(reqIn.CurrentPassword))
	newPasswordenc := hex.EncodeToString([]byte(reqIn.NewPassword))

	if reqIn.CurrentPassword == reqIn.NewPassword {
		logAction := "Current password is equal to new password"
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}

	if utilities.IsEmailValid(reqIn.UserName) {
		getUSerobj = usww.GetUserByEmailAddress(reqIn.UserName)
	} else if utilities.IsNumberValid(reqIn.UserName) {
		getUSerobj = usww.GetUserByMobileNumber(reqIn.UserName)
	} else {
		logAction := fmt.Sprintf("Incorrect username %s", reqIn.UserName)
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}

	if getUSerobj.EmailAddress != "" || getUSerobj.MobileNumber != "" {
		if getUSerobj.Password == enc {
			userRespose.Email = getUSerobj.EmailAddress
			userRespose.MobileNumber = getUSerobj.MobileNumber
			if doupdate := usww.ChangePassword(userRespose.Email, userRespose.MobileNumber, newPasswordenc); doupdate != 0 {
				ctx.JSON(http.StatusOK, "Password updated successfully!!")
				return
			} else {
				ctx.JSON(http.StatusBadRequest, "Unable to update password at the moment")
			}
		} else {
			logAction := fmt.Sprintf("Incorrect current password %s", reqIn.UserName)
			logrus.Info(logAction)
			ctx.JSON(http.StatusBadRequest, logAction)
			return
		}
	} else {
		logAction := fmt.Sprintf("Incorrect username or password for  %s", reqIn.UserName)
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}
}

// ChangePasswordWithoutValidation godoc
// @Summary		ChangePasswordWithoutValidation user password.
// @Description	ChangePasswordWithoutValidation user password.
// @Tags			user
// @Accept			*/*
// @User			json
// @Param user body inpuschema.ChangePassword true "Change password"
// @Success		200	{string}	string "Password updated successfully!!"
// @Failure		400		{string} string	"Unable to change password at the monent!!"
// @Router			/api/user/ChangePasswordWithoutValidation [post]
func ChangePasswordWithoutValidation(ctx *gin.Context) {
	reqIn := &inpuschema.ChangePassword{}
	if err := ctx.ShouldBindJSON(reqIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}

	getUSerobj := dbSchema.User{}
	userRespose := &inpuschema.UserResponse{}
	newPasswordenc := hex.EncodeToString([]byte(reqIn.NewPassword))

	if reqIn.CurrentPassword == reqIn.NewPassword {
		logAction := "Current password is equal to new password"
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}

	if utilities.IsEmailValid(reqIn.UserName) {
		getUSerobj = usww.GetUserByEmailAddress(reqIn.UserName)
	} else {
		logAction := fmt.Sprintf("Incorrect username %s", reqIn.UserName)
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}

	if getUSerobj.EmailAddress != "" {
		userRespose.Email = getUSerobj.EmailAddress
		// userRespose.MobileNumber = getUSerobj.MobileNumber
		if doupdate := usww.ChangePassword(userRespose.Email, "", newPasswordenc); doupdate != 0 {
			ctx.JSON(http.StatusOK, "Password changed successfully!!")
			return
		} else {
			ctx.JSON(http.StatusBadRequest, "Unable to change password at the moment")
		}
	} else {
		logAction := fmt.Sprintf("Incorrect email address  %s", reqIn.UserName)
		logrus.Info(logAction)
		ctx.JSON(http.StatusBadRequest, logAction)
		return
	}
}

// @Summary Import Image
// @Produce  json
// @Param image formData file true "Image File"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /api/user/UploadImage [post]
func UploadImage(ctx *gin.Context) {
	file, image, err := ctx.Request.FormFile("image")
	if err != nil {
		logrus.Warn(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if image == nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		ctx.JSON(http.StatusBadRequest, "ERROR_UPLOAD_CHECK_IMAGE_FORMAT")
		return
	}

	err = upload.CheckImage(fullPath)
	if err != nil {
		logrus.Warn(err)
		ctx.JSON(http.StatusInternalServerError, "ERROR_UPLOAD_CHECK_IMAGE_FAIL")
		return
	}
}

// GetAllNotificationsByEmail  user notifications by email
// @Summary		get user notification by email.
// @Description	get user notification by email.
// @Tags			user
// @Param EmailAddress path string true "User email address"
// @Param NotificationType path string true "Notification type"
// @Produce json
// @Accept			*/*
// @User			json
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{object}	inpuschema.CartObj
// @Router			/api/user/GetAllNotificationsByEmail/{EmailAddress}/{NotificationType} [get]
func GetAllNotificationsByEmail(ctx *gin.Context) {
	// if !ValidateClient(ctx) {
	// 	return
	// }
	EmailAddress := ctx.Param("EmailAddress")
	NotificationType := ctx.Param("NotificationType")
	logAction := fmt.Sprintf("GetAllNotificationsByEmail %s and %s", EmailAddress, NotificationType)
	logrus.Info(logAction)
	ctx.JSON(http.StatusBadRequest, TAct.GetAction(EmailAddress, NotificationType))
}

// UploadSignatureHandler godoc
// @Summary		Upload user Signature.
// @Description	Upload user Signature.
// @Tags			user
// @Accept			*/*
// @User			json
// @Param user body inpuschema.SignatureRequest true "upload signature"
// @Param Authorization header string true "Authorization token"
// @Param clientName header string true "registered client name"
// @Security BearerAuth
// @securityDefinitions.basic BearerAuth
// @Success		200	{string}	string "Signature saved successfully!"
// @Failure		400		{string} string	"Unable to save signature at the monent!!"
// @Router			/api/user/UploadSignatureHandler [post]
func UploadSignatureHandler(ctx *gin.Context) {
	reqIn := &inpuschema.SignatureRequest{}

	if err := ctx.ShouldBindJSON(reqIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}

	getDetailwithEmail := usww.GetUserByEmailAddress(reqIn.EmailAddress)
	if getDetailwithEmail.EmailAddress == "" {
		ctx.JSON(http.StatusBadRequest, "Unable to perform action at the moment!!")
		return
	}

	requestIn := &dbSchema.CompanySignature{
		CompanyId:                    getDetailwithEmail.Id,
		EmailAddress:                 getDetailwithEmail.EmailAddress,
		CompanySignaturePath:         reqIn.PathOnDevice,
		CompanySignaturePathOnDevice: reqIn.PathOnDevice,
		SignatureType:                reqIn.SignatureType,
		SignatureText:                reqIn.SignatureText,
		SignatureImage:               reqIn.SignatureImage,
	}

	//save request
	if doInsert := ISig.CreateCompanySignature(*requestIn); doInsert != "00" {
		ctx.JSON(http.StatusBadRequest, "Unable to save Signature at the moment!!")
		return
	}
	//update user records
	usrUpdate := dbSchema.UpdateStatus{
		EmailAddress: reqIn.EmailAddress,
		Status:       3,
	}

	userDeal.UpdateUserStatusByUserEmail(usrUpdate)

	ctx.JSON(http.StatusOK, "Signature saved successfully")
}
