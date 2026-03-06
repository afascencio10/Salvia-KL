package security_routers

import (
	"bitsflow/common/db"
	"bitsflow/common/utils"
	security_config "bitsflow/security/config"
	security_daos "bitsflow/security/dao"

	"github.com/gin-gonic/gin"
)

//var dbClientConfig db.DBClientConfig = db.DBClientConfig{Hostname: "localhost", Port: "5432", DatabaseName: "salvia", UserName: "postgres", Password: "asd876.!@asdSDS5436"}

//var dbClientConfig db.DBClientConfig = db.DBClientConfig{Hostname: "localhost", Port: "5432", DatabaseName: "salvia", UserName: "postgres", Password: "123456"}

// var dbClientConfig db.DBClientConfig = db.DBClientConfig{Hostname: "localhost", Port: "5432", DatabaseName: "salvia", UserName: "salvia_admin", Password: "asd876.!@asdSDS5a36Z"}
var dbClientConfig db.DBClientConfig

var dbServerConfig db.DBServerConfig = db.DBServerConfig{PoolSize: 80}

const module string = "security"

//var Templates *template.Template

func StartRouter(router *gin.Engine) {
	dbClientConfig = utils.LoadDBCLientConfig()

	var translatedModule string = security_config.Locale["sp"][module]
	var translatedEntity string
	var translatedNew string = security_config.Locale["sp"]["new"]
	var translatedUpdate string = security_config.Locale["sp"]["update"]
	//var translatedDisable string = security_config.Locale["sp"]["disable"]
	var translatedRemove string = security_config.Locale["sp"]["remove"]

	secRouter := router.Group("/" + translatedModule)
	{

		/*
			GeneralUser
		*/
		translatedEntity = security_config.Locale["sp"][security_daos.GeneralUserEntityName]
		secRouter.POST("/"+translatedEntity, GeneralUserPOST)
		secRouter.POST("/login", GeneralUserLOGIN_POST)
		secRouter.POST("/logout", GeneralUserLOGOUT_POST)

		secRouter.GET("/login", GeneralUserLOGIN_GET)
		secRouter.GET("/login/:id", GeneralUserLOGIN_GET)
		secRouter.GET("/login/:id/", GeneralUserLOGIN_GET)
		secRouter.GET("/login/captcha/:id", GeneralUserLOGIN_CAPTCHA)
		secRouter.GET("/login/captcha/:id/:media", GeneralUserLOGIN_CAPTCHA)
		secRouter.GET("/login/captcha", GeneralUserLOGIN_CAPTCHA)

		secRouter.POST("/login/:by", GeneralUserLOGIN_POST)
		secRouter.POST("/login/:by/", GeneralUserLOGIN_POST)

		secRouter.GET("/"+translatedEntity, GeneralUserGET)
		secRouter.GET("/"+translatedEntity+"/f/:f/p/:p", GeneralUserGET)
		secRouter.GET("/"+translatedEntity+"/"+translatedNew, GeneralUserPOST_GET)

		secRouter.GET("/"+translatedEntity+"/:id", GeneralUserGET)
		secRouter.GET("/"+translatedEntity+"/:id/", GeneralUserGET)

		secRouter.DELETE("/"+translatedEntity+"/:id", GeneralUserDELETE)
		secRouter.GET("/"+translatedEntity+"/:id/"+translatedRemove, GeneralUserDELETE_GET)
		secRouter.PUT("/"+translatedEntity+"/:id", GeneralUserPUT)
		secRouter.GET("/"+translatedEntity+"/:id/"+translatedUpdate, GeneralUserPUT_GET)
		secRouter.GET("/"+translatedEntity+"/:id/:by", GeneralUserPUT_GET)
		secRouter.PUT("/"+translatedEntity+"/:id/:by", GeneralUserPUT)
		/*
			----------------------------------------------------------
		*/

		/*
			City
		*/

		translatedEntity = security_config.Locale["sp"][security_daos.CityEntityName]

		secRouter.GET("/"+translatedEntity+"/:id/:by"+"/", CityGET)
		secRouter.GET("/"+translatedEntity+"/:id/:by", CityGET)

		/*
			Town
		*/

		translatedEntity = security_config.Locale["sp"][security_daos.TownEntityName]

		secRouter.GET("/"+translatedEntity+"/:id/:by"+"/", TownGET)
		secRouter.GET("/"+translatedEntity+"/:id/:by", TownGET)

		/*
			----------------------------------------------------------
		*/

	}

	secRouter = router.Group("/public")
	{

		/*
			City
		*/
		translatedEntity = security_config.Locale["sp"][security_daos.CityEntityName]
		secRouter.GET("/"+translatedEntity+"/:id/:by", CityGET_Public)

		/*
			----------------------------------------------------------
		*/
		/*
			Town
		*/
		translatedEntity = security_config.Locale["sp"][security_daos.TownEntityName]
		secRouter.GET("/"+translatedEntity+"/:id/:by", TownGET_Public)

		/*
			----------------------------------------------------------
		*/

	}

}
