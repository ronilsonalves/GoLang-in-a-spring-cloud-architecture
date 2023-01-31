package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/cmd/server/handler"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/docs"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/internal/appointment"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/internal/dentist"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/internal/patient"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/pkg/middleware"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/pkg/sd"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/pkg/store"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"os/signal"
	"time"
	_ "time/tzdata"
)

func init() {
	time.Local, _ = time.LoadLocation("America/Fortaleza")
}

// @title Dental Clinic API
// @version 1.0
// @description This API handle appointments, dentists and patients for dental clinic system.
// @termsOfService https://www.linkedin.com/in/ronilsonalves

// @contact.name Ronilson Alves
// @contact.url https://github.com/ronilsonalves

// @license.name Apache 2.0
// @license.url https://www.apache.org/license/LICENSE-2.0.html

// @securityDefinitions.apikey SECRET_TOKEN
// @in header
// @name SECRET_TOKEN
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file", err.Error())
	}

	eurekaRegister := sd.BuildFargoInstance()
	eurekaRegister.Register()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	// DB INIT
	sqlStore := store.NewSQLStore()
	apStore := store.NewSQLAp()
	//Handlers INIT
	appRepo := appointment.NewRepository(apStore)
	appService := appointment.NewService(appRepo)
	appHandler := handler.NewAppointmentHandler(appService)

	dentistRepo := dentist.NewRepository(sqlStore)
	dentistService := dentist.NewService(dentistRepo)
	dentistHandler := handler.NewDentistHandler(dentistService)

	patientRepo := patient.NewRepository(sqlStore)
	patientService := patient.NewService(patientRepo)
	patientHandler := handler.NewPatientHandler(patientService)

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	docs.SwaggerInfo.Host = os.Getenv("HOST") + ":" + os.Getenv("PORT")
	docs.SwaggerInfo.BasePath = os.Getenv("BASE_PATH")
	r.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(middleware.Authentication())

	api := r.Group("/api/v1")
	{
		appointments := api.Group("/appointments")
		{
			appointments.GET("", appHandler.GetAll())
			appointments.GET(":id", appHandler.GetByID())
			appointments.GET("/patient/:identity_number", appHandler.GetAllByIdentityNumber())
			appointments.GET("/dentist/:license_number", appHandler.GetAllByLicenseNumber())
			appointments.POST("", appHandler.Post())
			appointments.PUT(":id", appHandler.Put())
			appointments.PATCH(":id", appHandler.Patch())
			appointments.DELETE(":id", appHandler.Delete())
		}
		dentists := api.Group("/dentists")
		{
			dentists.GET("", dentistHandler.GetAll())
			dentists.GET(":id", dentistHandler.GetByID())
			dentists.POST("", dentistHandler.Post())
			dentists.PUT(":id", dentistHandler.Put())
			dentists.PATCH(":id", dentistHandler.Patch())
			dentists.DELETE(":id", dentistHandler.Delete())
		}
		patients := api.Group("/patients")
		{
			patients.GET("", patientHandler.GetAll())
			patients.GET(":id", patientHandler.GetByID())
			patients.POST("", patientHandler.Post())
			patients.PUT(":id", patientHandler.Put())
			patients.PATCH(":id", patientHandler.Patch())
			patients.DELETE(":id", patientHandler.Delete())
		}
	}

	go func() {
		select {
		case signal := <-c:
			_ = signal
			time.Sleep(4 * time.Second)
			eurekaRegister.Deregister()
			os.Exit(1)
		}
	}()

	r.Run(fmt.Sprintf("%s", os.Getenv("HOST")) + fmt.Sprintf(":%s", os.Getenv("PORT")))
}
