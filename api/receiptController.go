package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Adebusy/receiptonboardingsvc/obj"
	inpuschema "github.com/Adebusy/receiptonboardingsvc/obj"
	"github.com/Adebusy/receiptonboardingsvc/utilities"
	"github.com/gin-gonic/gin"
)

// SuggestCompanyReceiptTemplate suggests company logo
// @Summary		Get available and top 3 company logo.
// @Description	Get available and top 3 company logo.
// @Tags			receipt
// @Param inpuschema.ReceiptForSuggest path string true "company logo"
// @Produce json
// @Accept			*/*
// @Receipt			json
// // @Param Authorization header string true "Authorization token"
// // @Param clientName header string true "registered client name"
// // @Security BearerAuth
// // @securityDefinitions.basic BearerAuth
// @Success		200	{object}	[]string
// @Router			/api/receipt/SuggestCompanyReceiptTemplate [post]
func SuggestCompanyReceiptTemplate(ctx *gin.Context) {

	reqIn := &inpuschema.ReceiptForSuggest{}
	if err := ctx.ShouldBindJSON(reqIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}

	//retrieve templates
	templatesPath := "./receiptTemplates"

	data := obj.TemplateData{
		CompanyName: reqIn.CompanyName,
		Signature:   reqIn.ManagerSignaturePath,
		LogoURL:     reqIn.CompanyLogoPath,
	}

	htmlContents, err := utilities.LoadAndUpdateHTMLTemplates(templatesPath, data)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(htmlContents); i++ {
		resp := htmlContents[i]
		os.WriteFile(fmt.Sprintf("suggestedReceipt/receipt%s%v.html", reqIn.CompanyName, i), []byte(resp), 0644)
	}

	jsonBytes, _ := json.Marshal(htmlContents)
	ctx.JSON(http.StatusOK, fmt.Sprintf(string(jsonBytes)))
}

// CompanyReceiptSelect select company receipt
// @Summary		select company receipt.
// @Description	select company receipt.
// @Tags			receipt
// @Param inpuschema.ReceiptSelect path string true "company receipt"
// @Produce json
// @Accept			*/*
// @Receipt			json
// // @Param Authorization header string true "Authorization token"
// // @Param clientName header string true "registered client name"
// // @Security BearerAuth
// // @securityDefinitions.basic BearerAuth
// @Success		200	{object}	[]string
// @Router			/api/receipt/CompanyReceiptSelect [post]
func CompanyReceiptSelect(ctx *gin.Context) {
	reqIn := &inpuschema.ReceiptSelect{}
	if err := ctx.ShouldBindJSON(reqIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}

	if utilities.PickSelectedFile(reqIn.ReceiptName) == "00" {
		for i := 0; i <= 3; i++ {
			if err := os.Remove(fmt.Sprintf("suggestedReceipt/receipt%s%v.html", reqIn.CompanyName, i)); err != nil {
				fmt.Print(err.Error())
			}
		}
		ctx.JSON(http.StatusOK, "Receipt selected successfully!")
		return
	}
	ctx.JSON(http.StatusBadRequest, "Unable to select receipt at the moment!")
}

// GenerateReceipt generate company receipt
// @Summary		generate company receipt.
// @Description	generate company receipt.
// @Tags			receipt
// @Param inpuschema.ReceiptRequest path string true "receipt request"
// @Produce json
// @Accept			*/*
// @Receipt			json
// // @Param Authorization header string true "Authorization token"
// // @Param clientName header string true "registered client name"
// // @Security BearerAuth
// // @securityDefinitions.basic BearerAuth
// @Success		200	{object}	[]string
// @Router			/api/receipt/GenerateReceipt [post]
func GenerateReceipt(ctx *gin.Context) {
	reqIn := &inpuschema.ReceiptRequest{}
	if err := ctx.ShouldBindJSON(reqIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}

	type Item struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}

	parseRequest, TotalAmount := utilities.ParseItems(reqIn.SalesDetails)

	//generate receipttableandContent
	var allRole string
	for i := 0; i <= len(parseRequest)-1; i++ {
		troll := fmt.Sprintf(`<tr><td>%s</td><td>%d</td><td>%.2f</td></tr>`, parseRequest[i].Name, parseRequest[i].Quantity, parseRequest[i].Price)
		allRole += troll
	}

	//get template name
	template := comp.GetCompanyDetailsByCompanyName(reqIn.CompanyName)

	utilities.GenerateNewReceipt(template.CompanyReceiptPath, allRole, reqIn.CustumerName, strconv.FormatFloat(TotalAmount, 'f', 2, 64))

	ctx.JSON(http.StatusBadRequest, "Receipt generated successfully!")
}
