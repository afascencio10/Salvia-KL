package salvia_facades

import (
	"fmt"

	"bitsflow/common/db"
	"bitsflow/common/utils"
	salvia_config "bitsflow/salvia/config"
	salvia_ctrl "bitsflow/salvia/controllers"
	salvia_daos "bitsflow/salvia/dao"

	"github.com/gin-gonic/gin"
)

//var dbClientConfig db.DBClientConfig = db.DBClientConfig{Hostname: "localhost", Port: "5432", DatabaseName: "salvia", UserName: "postgres", Password: "asd876.!@asdSDS5436"}

// var dbClientConfig db.DBClientConfig = db.DBClientConfig{Hostname: "localhost", Port: "5432", DatabaseName: "salvia", UserName: "salvia_admin", Password: "asd876.!@asdSDS5a36Z"}
var dbClientConfig db.DBClientConfig

// var dbClientConfig db.DBClientConfig = db.DBClientConfig{Hostname: "localhost", Port: "5432", DatabaseName: "salvia", UserName: "postgres", Password: "123456"}
var dbServerConfig db.DBServerConfig = db.DBServerConfig{PoolSize: db.DefaultPoolSize}

const module string = "salvia"

//var Templates *template.Template//

func StartRouter(router *gin.Engine) {
	dbClientConfig = utils.LoadDBCLientConfig()

	// Sincroniza todas las secuencias del schema salvia al arrancar el servidor.
	// Previene errores de duplicate key causados por secuencias desincronizadas
	// (p.ej. tras migraciones o inserts con IDs explícitos).
	go func() {
		if err := db.SyncAllSequences(&dbClientConfig, "salvia"); err != nil {
			fmt.Printf("[StartRouter] Advertencia: error al sincronizar secuencias: %v\n", err)
		}
	}()

	var translatedModule string = salvia_config.Locale["sp"][module]
	var translatedEntity string
	var translatedNew string = salvia_config.Locale["sp"]["new"]
	var translatedUpdate string = salvia_config.Locale["sp"]["update"]
	var translatedInvalidate string = salvia_config.Locale["sp"]["invalidate"]
	var translatedReport string = salvia_config.Locale["sp"]["report"]
	var translatedDocument string = salvia_config.Locale["sp"]["documentType"]

	secRouter := router.Group("/" + translatedModule)
	{

		/*
			VictimContact
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.VictimContactEntityName]

		secRouter.POST("/public", VictimContactPOST_Public)
		secRouter.GET("/public/"+translatedNew, VictimContactPOST_GET_Public)

		secRouter.POST("/"+translatedEntity, VictimContactPOST)
		secRouter.PUT("/"+translatedEntity+"/:id/"+translatedInvalidate, VictimContactPUT) // Invalidar
		secRouter.GET("/"+translatedEntity, VictimContactGET)
		secRouter.GET("/"+translatedEntity+"/p/:p", VictimContactGET)
		secRouter.GET("/"+translatedEntity+"/f/:f/p/:p", VictimContactGET)
		secRouter.GET("/"+translatedEntity+"/"+translatedNew, VictimContactPOST_GET)
		secRouter.GET("/"+translatedEntity+"/:id", VictimContactGET)
		/*
			----------------------------------------------------------
		*/

		/*
			VictimCase
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.VictimCaseEntityName]

		secRouter.POST("/"+translatedEntity, VictimCasePOST)
		secRouter.POST("/"+translatedEntity+"/simular", SimulateCasePOST)
		secRouter.POST("/"+translatedEntity+"/:id", VictimCasePOST)

		// Rutas estáticas primero (antes de /:id para evitar conflictos en Gin)
		secRouter.GET("/"+translatedEntity+"/"+translatedNew, VictimCasePOST_GET)
		secRouter.GET("/"+translatedEntity+"/"+translatedReport, VictimCaseReportGET)
		secRouter.POST("/"+translatedEntity+"/"+translatedReport, VictimCaseReportPOST)
		secRouter.GET("/"+translatedEntity+"/p/:p", VictimCaseGET)
		secRouter.GET("/"+translatedEntity+"/f/:f/p/:p", VictimCaseGET)
		secRouter.GET("/"+translatedEntity, VictimCaseGET)

		// Rutas con parámetro :id
		secRouter.GET("/"+translatedEntity+"/:id"+"/"+translatedNew, VictimCasePOST_GET)
		secRouter.GET("/"+translatedEntity+"/:id"+"/"+translatedUpdate, VictimCasePUT_GET)
		secRouter.GET("/"+translatedEntity+"/:id/detalle", CaseDetailGET)
		secRouter.GET("/"+translatedEntity+"/:id/"+translatedDocument+"/:docType", VictimCaseGET)
		secRouter.GET("/"+translatedEntity+"/:id/"+translatedDocument+"/:docType"+"/p/:p", VictimCaseGET)
		secRouter.GET("/"+translatedEntity+"/:id", VictimCaseGET)
		secRouter.GET("/"+translatedEntity+"/:id/", VictimCaseGET)

		secRouter.PUT("/"+translatedEntity, VictimCasePUT)
		secRouter.PUT("/"+translatedEntity+"/:id/:by", VictimCasePUT)
		secRouter.PUT("/"+translatedEntity+"/:id/:by/:form", VictimCasePUT)

		/*
			EntityBranch
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.EntityBranchName]

		secRouter.GET("/"+translatedEntity, EntityBranchGET)

		secRouter.GET("/"+translatedEntity+"/e/:entityICode/tc/:townCode/", EntityBranchGET)
		secRouter.GET("/"+translatedEntity+"/e/:entityICode/tc/:townCode", EntityBranchGET)

		secRouter.GET("/"+translatedEntity+"/:townCode/by/:by/", EntityBranchGET)
		secRouter.GET("/"+translatedEntity+"/:townCode/by/:by", EntityBranchGET)

		secRouter.POST("/"+translatedEntity, EntityBranchPOST)
		secRouter.PUT("/"+translatedEntity, EntityBranchPUT)

		/*
			Entity
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.EntityEntityName]

		secRouter.GET("/"+translatedEntity, EntityGET)
		secRouter.GET("/"+translatedEntity+"/:sectorCode", EntityGET)

		/*
			Moment
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.MomentEntityName]

		secRouter.PUT("/"+translatedEntity+"/:id/:momentCode/:entityBranchIcode", MomentPUT)

		/*
			Alert
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.AlertEntityName]

		secRouter.GET("/"+translatedEntity, AlertGET)
		secRouter.GET("/"+translatedEntity+"/:id/:by/", AlertGET)
		secRouter.GET("/"+translatedEntity+"/:id/:by", AlertGET)

		/*
			CaseLog
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.CaseLogEntityName]

		secRouter.POST("/"+translatedEntity+"/:id", CaseLogPOST)

		/*
			LoadPlainFile
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.LoadPlainFilesEntityName]

		secRouter.POST("/"+translatedEntity+"/:srv", LoadPlainFilePOST)
		secRouter.POST("/"+translatedEntity+"/:srv/:id", LoadPlainFilePOST)
		secRouter.GET("/"+translatedEntity+"/"+translatedNew, LoadPlainFilePOST_GET)

		/*
			AssignOperators
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.AssignOperatorsEntityName]

		secRouter.POST("/"+translatedEntity+"/:by/:p1/:p2", AssignOperatorsPOST)
		secRouter.GET("/"+translatedEntity, AssignOperatorsPOST_GET)

		/*
			FollowUp
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.FollowUpEntityName]

		secRouter.POST("/"+translatedEntity+"/:id", FollowUpPOST)
		secRouter.PUT("/"+translatedEntity+"/:id/", FollowUpPUT)
		secRouter.PUT("/"+translatedEntity+"/:id", FollowUpPUT)
		secRouter.GET("/"+translatedEntity+"/:id", FollowUpGET)

		/*
			FollowUpEntry
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.FollowUpEntryEntityName]

		secRouter.PUT("/"+translatedEntity+"/:id/", FollowUpEntryPUT)
		secRouter.PUT("/"+translatedEntity+"/:id", FollowUpEntryPUT)

		/*
			Barrier
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.BarrierEntityName]

		secRouter.GET("/"+translatedEntity+"/:id/", BarrierGET)
		secRouter.GET("/"+translatedEntity+"/:id", BarrierGET)
		secRouter.GET("/"+translatedEntity+"/:id/:by/", BarrierGET)
		secRouter.GET("/"+translatedEntity+"/:id/:by", BarrierGET)

		/*
			FollowUpEntryActing
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.FollowUpEntryActingEntityName]

		secRouter.POST("/"+translatedEntity+"/e/:e/b/:b", FollowUpEntryActingPOST)

		/*
			Mis Seguimientos (My Follow Ups)
		*/
		secRouter.GET("/mis-seguimientos", MyFollowUpsGET)

		/*
			Mis Casos
		*/
		secRouter.GET("/mis-casos", MyCasesGET)

		/*
			Lista de Casos
		*/
		secRouter.GET("/lista-casos", ListCasesGET)

		/*
			Historial de Remisiones (Atención Psicosocial — solo sv)
		*/
		secRouter.GET("/historial-remisiones", HistorialRemisionesGET)

		/*
			Mis remisiones Psicosocial (roles ps, ts)
		*/
		secRouter.GET("/mis-remisiones-psicosocial", MisRemisionesPsicosocialGET)

		/*
			Remisión Temporal (host del card 3x3 — pantalla puente, roles ps, ts)
		*/
		secRouter.GET("/remision-temporal/:id", RemisionTemporalGET)

		/*
			KPIs Dashboard (embudo analítico)
		*/
		secRouter.GET("/kpis", KpiDashboardGET)
		secRouter.GET("/kpis/funnel", KpiFunnelGET)

		/*
			Feminicide
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.FeminicideEntityName]

		secRouter.GET("/"+translatedEntity, FeminicideGET)
		secRouter.GET("/"+translatedEntity+"/"+translatedNew, FeminicidePOST_GET)
		secRouter.POST("/"+translatedEntity, FeminicidePOST)

		/*
			HacerSeguimiento
		*/
		secRouter.GET("/hacer-seguimiento/:id", HacerSeguimientoGET)

		/*
			SeguimientosArea
		*/
		secRouter.GET("/seguimientos/area", SeguimientosAreaGET)

		/*
			Psicosocial — Registrar sesión
			Rutas estáticas antes de la dinámica para evitar conflictos en Gin.
		*/
		secRouter.GET("/psicosocial/test", TestCasosPsicosocialGET)
		secRouter.GET("/psicosocial/registrar/:id", RegistrarSesionPsicosocialGET)

		/*
			Notificaciones
		*/
		secRouter.GET("/notificaciones", NotificacionesGET)

		/*
			Mis Barreras
		*/
		secRouter.GET("/mis-barreras", MisBarrerasGET)

		/*
			Barrera Detalle — pantalla interna de una barrera
			Nota: usa /barreras (plural) para evitar conflicto con las rutas JSON /barrera/:id
		*/
		secRouter.GET("/barreras/:id", BarreraDetalleGET)

		/*
			Remisión Psicosocial Detalle — pantalla interna de una remisión
		*/
		secRouter.GET("/remision-psicosocial/:id", RemisionPsicosocialDetalleGET)

		/*
			FeminicideRisk
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.FeminicideRiskEntityName]

		secRouter.GET("/"+translatedEntity, FeminicideRiskGET)
		secRouter.GET("/"+translatedEntity+"/:id", FeminicideRiskGET)
		secRouter.GET("/"+translatedEntity+"/:id/"+translatedNew, FeminicideRiskPOST_GET)
		secRouter.POST("/"+translatedEntity, FeminicideRiskPOST)

		/*
			----------------------------------------------------------
		*/

	}

	secRouter = router.Group("/public")
	{

		/*
			VictimContact
		*/
		translatedEntity = salvia_config.Locale["sp"][salvia_daos.VictimContactEntityName]

		secRouter.POST("/"+translatedEntity+"/"+translatedNew, VictimContactPOST_Public)
		secRouter.GET("/"+translatedEntity+"/"+translatedNew, VictimContactPOST_GET_Public)

		/*
			SimularCaso — sin autenticación, solo para desarrollo/pruebas
		*/
		secRouter.POST("/simular/caso", SimulateCasePOST)

		/*
			SyncSequences — sin autenticación, solo para desarrollo/pruebas
		*/
		secRouter.POST("/sync/sequences", SyncSequencesPOST)

		/*
			----------------------------------------------------------
		*/

	}

	salvia_ctrl.GetVictimCaseForm2EnumsByAll(&db.ConnData{}, dbClientConfig, dbServerConfig)

}
