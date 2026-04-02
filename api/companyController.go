package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Adebusy/receiptonboardingsvc/dataaccess"
	dbSchema "github.com/Adebusy/receiptonboardingsvc/dataaccess"
	aiagents "github.com/Adebusy/receiptonboardingsvc/llmClientFunctions"
	inpuschema "github.com/Adebusy/receiptonboardingsvc/obj"
	psg "github.com/Adebusy/receiptonboardingsvc/postgresql"
	"github.com/Adebusy/receiptonboardingsvc/utilities"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var comp = dbSchema.CompanyDeal(psg.GetDB())
var userDeal = dbSchema.ConneectDeal(psg.GetDB())

func RegisterCompany(ctx *gin.Context) {
	// currentTime := time.Now()
	reqIn := &dbSchema.Company{}
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

	if CheckEmailExist := comp.GetCompanyByEmailAddress(reqIn.EmailAddress); CheckEmailExist.CompanyName != "" {
		ctx.JSON(http.StatusOK, "It appears that you already have a company name "+CheckEmailExist.CompanyName+" attached to this profile, please proceed to signature module if you do not have one.")
		return
	}

	if CheckEmailExist := comp.GetCompanyByMobileNumber(reqIn.MobileNumber); CheckEmailExist.CompanyName != "" {
		ctx.JSON(http.StatusOK, "It appears that you already have a company name "+CheckEmailExist.CompanyName+" attached to this profile, please proceed to signature module if you do not have one.")
		return
	}

	requestComp := dbSchema.TblCompany{
		CompanyName:     strings.Split(reqIn.CompanyName, ".")[1],
		CompanyAddress:  reqIn.CompanyAddress,
		EmailAddress:    reqIn.EmailAddress,
		MobileNumber:    reqIn.MobileNumber,
		RegNo:           reqIn.RegNo,
		CompanyLogoPath: reqIn.CompanyLogoPath,
		Status:          1,
		CreatedAt:       time.Now(),
		UserId:          reqIn.UserId,
	}

	doCreate := comp.CreateCompany(requestComp)

	//update user table
	usrUpdate := dbSchema.UpdateStatus{
		EmailAddress: reqIn.EmailAddress,
		Status:       1,
	}
	userDeal.UpdateUserStatusByUserEmail(usrUpdate)

	logrus.Info(doCreate)
	ctx.JSON(http.StatusOK, doCreate)
}

func GetCompanyByCompanyEmail(ctx *gin.Context) {
	usww := dbSchema.CompanyDeal(psg.GetDB())
	companyEmail := ctx.GetHeader("EmailAddress")
	if CheckEmailExist := usww.GetCompanyByEmailAddress(companyEmail); CheckEmailExist.CompanyName != "" {
		ctx.JSON(http.StatusOK, CheckEmailExist)
		return
	} else {
		ctx.JSON(http.StatusBadRequest, "This user has not completed his/her profile registration")
		return
	}
}

func GetCompanyByCompanyMobile(ctx *gin.Context) {
	companyMobile := ctx.GetHeader("MobileNumer")
	if CheckCompanyExist := comp.GetCompanyByMobileNumber(companyMobile); CheckCompanyExist.CompanyName != "" {
		ctx.JSON(http.StatusOK, CheckCompanyExist)
		return
	} else {
		ctx.JSON(http.StatusBadRequest, "This user has not completed his/her profile registration")
		return
	}
}

// SuggestCompNames suggest company names
// @Summary		Get available and top 3 company names.
// @Description	Get available and top 3 company names.
// @Tags			company
// @Param inpuschema.CompanyForSuggest path string true "company names"
// @Produce json
// @Accept			*/*
// @Company			json
// // @Param Authorization header string true "Authorization token"
// // @Param clientName header string true "registered client name"
// // @Security BearerAuth
// // @securityDefinitions.basic BearerAuth
// @Success		200	{object}	[]string
// @Router			/api/company/SuggestCompNames [post]
func SuggestCompNames(ctx *gin.Context) {
	reqIn := &inpuschema.CompanyForSuggest{}
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
	//check if names already suggested
	if checkNames := comp.FetchSuggestedCompanyNames(reqIn.UserId); checkNames[0].CompanyName != "" {
		ctx.JSON(http.StatusOK, checkNames)
		return
	}

	getNames := aiagents.SuggestCompanyNames(reqIn.CompanyName, reqIn.FieldsArea)

	//log suggested names
	for k := range getNames {
		sn := dataaccess.TblCompanyNameSuggested{
			UserId:      reqIn.UserId,
			CompanyName: getNames[k],
		}

		if retInsert := comp.CreateNameSuggested(sn); retInsert != "00" {
			logrus.Info(retInsert)
			ctx.JSON(http.StatusBadRequest, "Unable to complete process at the moment")
			return
		}
	}

	logAction := fmt.Sprintf("SuggestCompanyNames %v", reqIn.CompanyName)
	logrus.Info(logAction)
	ctx.JSON(http.StatusOK, getNames)
}

// SuggestCompanyLogo suggests company logo
// @Summary		Get available and top 3 company logo.
// @Description	Get available and top 3 company logo.
// @Tags			company
// @Param inpuschema.CompanyForSuggest path string true "company logo"
// @Produce json
// @Accept			*/*
// @Company			json
// // @Param Authorization header string true "Authorization token"
// // @Param clientName header string true "registered client name"
// // @Security BearerAuth
// // @securityDefinitions.basic BearerAuth
// @Success		200	{object}	[]string
// @Router			/api/company/SuggestCompanyLogo [post]
func SuggestCompanyLogo(ctx *gin.Context) {
	currentTime := time.Now()
	reqIn := &inpuschema.CompanyForSuggest{}
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

	getNames := aiagents.SuggestCompanyLogo(reqIn.CompanyName, reqIn.FieldsArea)
	logAction := fmt.Sprintf("SuggestCompanyNames %v", reqIn.CompanyName)
	logrus.Info(logAction)

	jsonString := utilities.ReconstructJSON(getNames)

	var logos []utilities.LogoConcept
	if err := json.Unmarshal([]byte(jsonString), &logos); err != nil {
		panic(err)
	}

	for _, logo := range logos {
		err := utilities.SaveBase64Image(
			logo.LogoBase64.Data,
			"./logos",
			logo.LogoID,
			logo.LogoBase64.Extension,
		)

		if err != nil {
			log.Printf("Failed to save %s: %v", logo.LogoID, err)
		} else {
			//save image path and name inside database/table
			lgsPath := dataaccess.TblCompanylogosPrev{CompanyName: reqIn.CompanyName, EmailAddress: reqIn.EmailAddress, CompanyLogoPath: logo.LogoID + "." + logo.LogoBase64.Extension, Status: 0, CreatedAt: string(currentTime.String())}
			if savePrev := comp.CreateCompanylogosPrev(lgsPath); savePrev != "00" {
				log.Fatal("unable to log prev image to db")
			}
		}
	}
	ctx.JSON(http.StatusOK, getNames)
}

// SuggestManagerSign suggests company logo
// @Summary		Get available and top 3 company logo.
// @Description	Get available and top 3 company logo.
// @Tags			company
// @Param inpuschema.CompanyForSuggest path string true "company logo"
// @Produce json
// @Accept			*/*
// @Company			json
// // @Param Authorization header string true "Authorization token"
// // @Param clientName header string true "registered client name"
// // @Security BearerAuth
// // @securityDefinitions.basic BearerAuth
// @Success		200	{object}	[]string
// @Router			/api/company/SuggestManagerSign [post]
func SuggestManagerSign(ctx *gin.Context) {
	currentTime := time.Now()
	reqIn := &inpuschema.CompanyForSuggest{}
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

	getNames := aiagents.SuggestCompanyLogo(reqIn.CompanyName, reqIn.FieldsArea)
	logAction := fmt.Sprintf("SuggestCompanyNames %v", reqIn.CompanyName)
	logrus.Info(logAction)

	jsonString := utilities.ReconstructJSON(getNames)

	var logos []utilities.LogoConcept
	if err := json.Unmarshal([]byte(jsonString), &logos); err != nil {
		panic(err)
	}

	for _, logo := range logos {
		err := utilities.SaveBase64Image(
			logo.LogoBase64.Data,
			"./logos",
			logo.LogoID,
			logo.LogoBase64.Extension,
		)

		if err != nil {
			log.Printf("Failed to save %s: %v", logo.LogoID, err)
		} else {
			//save image path and name inside database/table
			lgsPath := dataaccess.TblCompanylogosPrev{CompanyName: reqIn.CompanyName, EmailAddress: reqIn.EmailAddress, CompanyLogoPath: logo.LogoID + "." + logo.LogoBase64.Extension, Status: 0, CreatedAt: string(currentTime.String())}
			if savePrev := comp.CreateCompanylogosPrev(lgsPath); savePrev != "00" {
				log.Fatal("unable to log prev image to db")
			}
		}
	}
	ctx.JSON(http.StatusOK, getNames)
}
