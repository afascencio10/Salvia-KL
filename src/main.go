package main

import (
	common_routers "bitsflow/common/facades"
	"bitsflow/common/utils"
	_ "bitsflow/docs" // Swagger docs generados por swag init
	internaldb "bitsflow/internal/db"
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	salvia_ctrl "bitsflow/salvia/controller"
	salvia_legacy "bitsflow/salvia/controllers"
	salvia_facades "bitsflow/salvia/facades"
	"bitsflow/salvia/service"
	security_routers "bitsflow/security/facades"
	"embed"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed config/*
var configAssets embed.FS

//go:embed frontend/* frontend/*/* frontend/*/*/* frontend/*/*/*/* frontend/*/*/*/*/*
var frontendAssets embed.FS

func main() {

	//Definimos todos los archivos estáticos
	utils.ConfigAssets = configAssets
	utils.FrontendAssets = frontendAssets

	var taskFuncs map[string]func() = map[string]func(){"UpdateVictimCasesOwnersAndRoles": salvia_facades.UpdateVictimCasesOwnersAndRoles}
	go utils.Load(taskFuncs)

	var router *gin.Engine = common_routers.InitRouter()
	salvia_facades.StartRouter(router)
	security_routers.StartRouter(router)

	// ── Nuevo patrón: GORM + Repository + Service + Controller ──────────────
	gormDB, err := internaldb.NewGormDB(utils.LoadDBCLientConfig())
	if err != nil {
		log.Fatalf("Error conectando GORM: %v", err)
	}
	// AutoMigrate por tabla — warning en lugar de fatal para tablas ya existentes
	for _, m := range []interface{}{
		&models.Form{},
		&models.FormSection{},
		&models.RepeaterGroup{},
		&models.Question{},
		&models.VisibilityCondition{},
		&models.FormSubmission{},
		&models.RepeaterEntry{},
		&models.Answer{},
		&models.FollowUpV2{},
		&models.EmergencyMeasure{},
		&models.PsychosocialSupport{},
		&models.EconomicStabilization{},
	} {
		if err := gormDB.AutoMigrate(m); err != nil {
			log.Printf("[WARN] AutoMigrate %T: %v", m, err)
		}
	}

	// Repositories
	formRepo               := repository.NewFormRepository(gormDB)
	formSectionRepo        := repository.NewFormSectionRepository(gormDB)
	questionRepo           := repository.NewQuestionRepository(gormDB)
	repeaterGroupRepo      := repository.NewRepeaterGroupRepository(gormDB)
	visibilityCondRepo     := repository.NewVisibilityConditionRepository(gormDB)
	formSubmissionRepo     := repository.NewFormSubmissionRepository(gormDB)
	repeaterEntryRepo      := repository.NewRepeaterEntryRepository(gormDB)
	answerRepo             := repository.NewAnswerRepository(gormDB)
	followUpRepo           := repository.NewFollowUpRepository(gormDB)
	barrierV2Repo          := repository.NewBarrierV2Repository(gormDB)
	victimCaseLightRepo    := repository.NewVictimCaseLightRepository(gormDB)
	emRepo                 := repository.NewEmergencyMeasureRepository(gormDB)
	psRepo                 := repository.NewPsychosocialSupportRepository(gormDB)
	esRepo                 := repository.NewEconomicStabilizationRepository(gormDB)

	// Services
	formSvc               := service.NewFormService(formRepo)
	formSectionSvc        := service.NewFormSectionService(formSectionRepo)
	questionSvc           := service.NewQuestionService(questionRepo)
	repeaterGroupSvc      := service.NewRepeaterGroupService(repeaterGroupRepo)
	visibilityCondSvc     := service.NewVisibilityConditionService(visibilityCondRepo)
	formSubmissionSvc     := service.NewFormSubmissionService(formSubmissionRepo)
	repeaterEntrySvc      := service.NewRepeaterEntryService(repeaterEntryRepo)
	answerSvc             := service.NewAnswerService(answerRepo)
	followUpV2Svc         := service.NewFollowUpV2Service(followUpRepo, barrierV2Repo, victimCaseLightRepo, emRepo, psRepo, esRepo)

	// Inyectar el servicio en el controller legacy para generación automática del calendario
	salvia_legacy.FollowUpSvc = followUpV2Svc

	// Controllers
	formCtrl               := salvia_ctrl.NewFormController(formSvc)
	formSectionCtrl        := salvia_ctrl.NewFormSectionController(formSectionSvc)
	questionCtrl           := salvia_ctrl.NewQuestionController(questionSvc)
	repeaterGroupCtrl      := salvia_ctrl.NewRepeaterGroupController(repeaterGroupSvc)
	visibilityCondCtrl     := salvia_ctrl.NewVisibilityConditionController(visibilityCondSvc)
	formSubmissionCtrl     := salvia_ctrl.NewFormSubmissionController(formSubmissionSvc)
	repeaterEntryCtrl      := salvia_ctrl.NewRepeaterEntryController(repeaterEntrySvc)
	answerCtrl             := salvia_ctrl.NewAnswerController(answerSvc)
	followUpV2Ctrl         := salvia_ctrl.NewFollowUpV2Controller(followUpV2Svc)

	// Routes
	api := router.Group("/api/v1")
	formCtrl.RegisterRoutes(api)
	formSubmissionCtrl.RegisterRoutes(api)
	formSectionCtrl.RegisterRoutes(api)
	questionCtrl.RegisterRoutes(api)
	repeaterGroupCtrl.RegisterRoutes(api)
	visibilityCondCtrl.RegisterRoutes(api)
	repeaterEntryCtrl.RegisterRoutes(api)
	answerCtrl.RegisterRoutes(api)
	followUpV2Ctrl.RegisterRoutes(api)

	// Swagger UI — accesible en https://localhost/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// ────────────────────────────────────────────────────────────────────────

	common_routers.StartRouter()
}
