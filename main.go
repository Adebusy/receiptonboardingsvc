package main

import (
	"context"
	"net/http"

	"github.com/Adebusy/receiptonboardingsvc/api"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

func main() {
	shutdown := initTracer()
	defer shutdown(context.Background())
	svc := gin.Default()
	svc.Use(otelgin.Middleware("receipt-onboarding-service"))

	svc.POST("/api/user/SignUp", api.SignUp)                                           //done
	svc.GET("api/user/GetUserByEmailAddress/:EmailAddress", api.GetUserByEmailAddress) //done
	svc.POST("api/user/LogIn", api.LogIn)                                              //done
	svc.GET("api/user/GetUserByMobile/:MobileNumber", api.GetUserByMobile)             //done
	svc.POST("/api/user/UpdateUserDetails", api.UpdateUserDetails)                     //done
	svc.POST("/api/user/UploadSignatureHandler", api.UploadSignatureHandler)           //done

	svc.POST("api/company/SuggestCompanyNames", api.SuggestCompNames)       //done
	svc.POST("api/company/SuggestCompanyLogo", api.SuggestCompanyLogo)      //done
	svc.POST("api/company/CompleteSignUp", api.RegisterCompany)             //done
	svc.POST("api/company/GetCompanyEmail", api.GetCompanyByCompanyEmail)   //done
	svc.POST("api/company/GetCompanyMobile", api.GetCompanyByCompanyMobile) //done

	svc.POST("/api/receipt/SuggestCompanyReceiptTemplate", api.SuggestCompanyReceiptTemplate) //done
	svc.POST("api/receipt/CompanyReceiptSelect", api.CompanyReceiptSelect)                    //done

	svc.POST("api/receipt/GenerateReceipt", api.GenerateReceipt) //done
	svc.GET("/", CheckServiceStatus)
	svc.Run(":8098")
}

func CheckServiceStatus(ctx *gin.Context) {
	tracer := otel.Tracer("gin-service")

	// Start custom span
	_, span := tracer.Start(ctx.Request.Context(), "CheckService")
	defer span.End()
	// Simulate business logic
	//processData(ctx) sample method call

	ctx.JSON(http.StatusOK, "I am up and running!!!")
}

// func processData(ctx *gin.Context) {
// 	tracer := otel.Tracer("gin-service")

// 	_, span := tracer.Start(ctx.Request.Context(), "processData")
// 	defer span.End()

// 	// Simulate work
// 	time.Sleep(100 * time.Millisecond)

// 	calculateSomething(ctx)
// }

// func calculateSomething(ctx *gin.Context) {
// 	tracer := otel.Tracer("gin-service")

// 	_, span := tracer.Start(ctx.Request.Context(), "calculateSomething")
// 	defer span.End()

// 	time.Sleep(50 * time.Millisecond)
// }
