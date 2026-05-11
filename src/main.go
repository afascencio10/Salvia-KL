package main

import (
    common_routers "bitsflow/common/facades"
    commondb "bitsflow/common/db"
    "bitsflow/common/utils"
    internaldb "bitsflow/internal/db"
    "bitsflow/internal/models"
    "bitsflow/internal/repository"
    salvia_ctrl "bitsflow/salvia/controller"
    salvia_legacy "bitsflow/salvia/controllers"
    salvia_facades "bitsflow/salvia/facades"
    "bitsflow/salvia/service"
    security_routers "bitsflow/security/facades"
    "context"
    "embed"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
)

//go:embed config/*
var configAssets embed.FS

//go:embed frontend/* frontend/*/* frontend/*/*/* frontend/*/*/*/* frontend/*/*/*/*/*
var frontendAssets embed.FS

func main() {

    // Definimos todos los archivos estáticos
    utils.ConfigAssets = configAssets
    utils.FrontendAssets = frontendAssets

    var taskFuncs map[string]func() = map[string]func(){"UpdateVictimCasesOwnersAndRoles": salvia_facades.UpdateVictimCasesOwnersAndRoles}
    go utils.Load(taskFuncs)

    var router *gin.Engine = common_routers.InitRouter()
    salvia_facades.StartRouter(router)
    security_routers.StartRouter(router)

    // ── Nuevo patrón: GORM + Repository + Service + Controller ──────────────
    gormDB, err := internaldb.NewGormDB(utils.LoadDBCLientConfig().AsGormConfig())
    if err != nil {
        log.Fatalf("Error conectando GORM: %v", err)
    }

    // Asegurar que columnas requeridas existan en la BD (compatibilidad con QA/producción)
    if err := gormDB.Exec("ALTER TABLE security.general_user ADD COLUMN IF NOT EXISTS general_user_team VARCHAR(64) DEFAULT NULL").Error; err != nil {
        log.Printf("[WARN] No se pudo crear columna general_user_team (puede que ya exista o falten permisos): %v", err)
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
        &models.Option{},
        &models.FollowUpV2{},
        &models.EmergencyMeasure{},
        &models.PsychosocialSupport{},
        &models.EconomicStabilization{},
        &models.CaseTimelineEvent{},
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
    optionRepo             := repository.NewOptionRepository(gormDB)
    barrierV2Repo          := repository.NewBarrierV2Repository(gormDB)
    victimCaseLightRepo    := repository.NewVictimCaseLightRepository(gormDB)
    townLightRepo          := repository.NewTownLightRepository(gormDB)
    attemptRepo            := repository.NewFollowUpAttemptRepository(gormDB)
    emRepo                 := repository.NewEmergencyMeasureRepository(gormDB)
    psRepo                 := repository.NewPsychosocialSupportRepository(gormDB)
    esRepo                 := repository.NewEconomicStabilizationRepository(gormDB)
    agentLightRepo         := repository.NewAgentLightRepository(gormDB)
    caseDetailRepo         := repository.NewCaseDetailRepository(gormDB)
    reportRepo             := repository.NewReportRepository(gormDB)

    // Services
    formSvc := service.NewFormService(service.FormServiceDeps{
        FormRepo:           formRepo,
        FormSectionRepo:    formSectionRepo,
        QuestionRepo:       questionRepo,
        RepeaterGroupRepo:  repeaterGroupRepo,
        OptionRepo:         optionRepo,
        VisibilityCondRepo: visibilityCondRepo,
        FormSubmissionRepo: formSubmissionRepo,
        RepeaterEntryRepo:  repeaterEntryRepo,
        AnswerRepo:         answerRepo,
    })
    formSectionSvc        := service.NewFormSectionService(formSectionRepo)
    questionSvc           := service.NewQuestionService(questionRepo)
    repeaterGroupSvc      := service.NewRepeaterGroupService(repeaterGroupRepo)
    visibilityCondSvc     := service.NewVisibilityConditionService(visibilityCondRepo)
    formSubmissionSvc     := service.NewFormSubmissionService(formSubmissionRepo)
    repeaterEntrySvc      := service.NewRepeaterEntryService(repeaterEntryRepo)
    answerSvc             := service.NewAnswerService(answerRepo)
    optionSvc             := service.NewOptionService(optionRepo)
    followUpV2Svc         := service.NewFollowUpV2Service(followUpRepo, formSubmissionRepo, barrierV2Repo, victimCaseLightRepo, townLightRepo, attemptRepo, emRepo, psRepo, esRepo, agentLightRepo)
    caseDetailSvc         := service.NewCaseDetailService(caseDetailRepo, gormDB)
    reportSvc             := service.NewReportService(reportRepo)

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
    optionCtrl             := salvia_ctrl.NewOptionController(optionSvc)
    caseDetailCtrl         := salvia_ctrl.NewCaseDetailController(caseDetailSvc)
    reportCtrl             := salvia_ctrl.NewReportController(reportSvc)

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
    optionCtrl.RegisterRoutes(api)
    caseDetailCtrl.RegisterRoutes(api)
    reportCtrl.RegisterRoutes(api)
    // ────────────────────────────────────────────────────────────────────────

    // ── Graceful shutdown ────────────────────────────────────────────────────
    srv := common_routers.NewHTTPServer()

    go func() {
        port := os.Getenv("PORT")
        if port != "" {
            log.Printf("Servidor corriendo en :%s", port)
            if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
                log.Fatalf("Error iniciando servidor: %v", err)
            }
        } else {
            log.Println("Servidor corriendo en :443 (TLS)")
            if err := srv.ListenAndServeTLS("certs/salvia.crt", "certs/salvia.key"); err != nil && err.Error() != "http: Server closed" {
                log.Fatalf("Error iniciando servidor TLS: %v", err)
            }
        }
    }()

    // Esperar señal de apagado (Ctrl+C o SIGTERM)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Apagando servidor — cerrando conexiones...")

    // 1. Detener el servidor HTTP (esperar hasta 10s a que terminen requests activos)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("[HTTP] Error en shutdown: %v", err)
    }

    // 2. Cerrar pool GORM
    if sqlDB, err := gormDB.DB(); err == nil {
        sqlDB.Close()
        log.Println("[GORM] Conexiones cerradas")
    }

    // 3. Cerrar pool legacy pgx
    commondb.CloseAllConnections()

    log.Println("Servidor apagado correctamente ✓")
    // ─────────────────────────────────────────────────────────────────────────
}