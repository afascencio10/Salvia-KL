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
    _ "time/tzdata" // Embebe zonas horarias para que funcione en contenedores sin tzdata

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
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

    // Asegurar que usuarios sin town tengan Bogotá por defecto (evita error 500 en reasignación)
    gormDB.Exec(`UPDATE security.general_user_profile SET general_user_profile_town = '11001000' WHERE (general_user_profile_town IS NULL OR general_user_profile_town = '') AND general_user_profile_id IN (SELECT general_user_general_user_profile FROM security.general_user WHERE general_user_status = 'e')`)

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
        &models.Dupla{},
        &models.PsychosocialSupport{},
        &models.TeamContact{},
        &models.ContactAttempt{},
        &models.EconomicStabilization{},
        &models.CaseTimelineEvent{},
        &models.EntityLetter{},
        &models.CaseTask{},
        &models.RenderModification{},
        &models.MenTeamRemision{},
        &models.DiscapacidadRemision{},
        &models.BarrierV2{},
        &models.BarrierFollowUp{},
        &models.Directory{},
        &models.EntityCase{},
        &models.TeamMeeting{},
        &models.TeamMeetingAgent{},
    } {
        if err := gormDB.AutoMigrate(m); err != nil {
            log.Printf("[WARN] AutoMigrate %T: %v", m, err)
        }
    }

    // ─── AutoMigrate para tablas legacy (schemas security/salvia) ───
    // Agrega columnas faltantes sin ALTER TABLE manual.
    // Si se necesita un campo nuevo en una tabla legacy, agregarlo aquí.
    migrateLegacyTables(gormDB)

    // Fix: asignar sequence_number a seguimientos que lo tienen en 0 (bug de buildFollowUps).
    // Ordena por scheduled_date ASC dentro de cada caso para asignar 1, 2, 3...
    gormDB.Exec(`
        WITH numbered AS (
            SELECT id, ROW_NUMBER() OVER (PARTITION BY case_id ORDER BY scheduled_date ASC, created_at ASC) AS rn
            FROM salvia.follow_up_v2
            WHERE deleted_at IS NULL AND sequence_number = 0
        )
        UPDATE salvia.follow_up_v2 SET sequence_number = numbered.rn
        FROM numbered WHERE follow_up_v2.id = numbered.id
    `)

    // Repositories
    formRepo               := repository.NewFormRepository(gormDB)
    formSectionRepo        := repository.NewFormSectionRepository(gormDB)
    questionRepo           := repository.NewQuestionRepository(gormDB)
    repeaterGroupRepo      := repository.NewRepeaterGroupRepository(gormDB)
    visibilityCondRepo     := repository.NewVisibilityConditionRepository(gormDB)
    renderModificationRepo := repository.NewRenderModificationRepository(gormDB)
    formSubmissionRepo     := repository.NewFormSubmissionRepository(gormDB)
    repeaterEntryRepo      := repository.NewRepeaterEntryRepository(gormDB)
    answerRepo             := repository.NewAnswerRepository(gormDB)
    followUpRepo           := repository.NewFollowUpRepository(gormDB)
    optionRepo             := repository.NewOptionRepository(gormDB)
    barrierV2Repo          := repository.NewBarrierV2Repository(gormDB)
    barrierFollowUpRepo    := repository.NewBarrierFollowUpRepository(gormDB)
    victimCaseLightRepo    := repository.NewVictimCaseLightRepository(gormDB)
    townLightRepo          := repository.NewTownLightRepository(gormDB)
    attemptRepo            := repository.NewFollowUpAttemptRepository(gormDB)
    emRepo                 := repository.NewEmergencyMeasureRepository(gormDB)
    psRepo                 := repository.NewPsychosocialSupportRepository(gormDB)
    teamContactRepo        := repository.NewTeamContactRepository(gormDB)
    esRepo                 := repository.NewEconomicStabilizationRepository(gormDB)
    menTeamRemisionRepo       := repository.NewMenTeamRemisionRepository(gormDB)
    discapacidadRemisionRepo  := repository.NewDiscapacidadRemisionRepository(gormDB)
    agentLightRepo         := repository.NewAgentLightRepository(gormDB)
    caseTimelineRepo       := repository.NewCaseTimelineEventRepository(gormDB)
    caseDetailRepo         := repository.NewCaseDetailRepository(gormDB)
    caseInfoRepo           := repository.NewCaseInfoRepository(gormDB)
    reportRepo             := repository.NewReportRepository(gormDB)
    caseTaskRepo           := repository.NewCaseTaskRepository(gormDB)
    entityLetterRepo       := repository.NewEntityLetterRepository(gormDB)
    casesListRepo          := repository.NewCasesListRepository(gormDB)
    casesReassignRepo      := repository.NewCasesReassignRepository(gormDB)
    psychosocialListRepo   := repository.NewPsychosocialListRepository(gormDB)
    duplaRepo              := repository.NewDuplaRepository(gormDB)
    psychosocialReassignRepo := repository.NewPsychosocialReassignRepository(gormDB)
    entityCaseRepo         := repository.NewEntityCaseRepository(gormDB)
    psychosocialCalendarRepo := repository.NewPsychosocialCalendarRepository(gormDB)
    victimCaseFormRepo     := repository.NewVictimCaseFormRepository(gormDB)

    // Services
    casoCierreSvc := service.NewCasoCierreService(victimCaseLightRepo, caseTimelineRepo)

    formSectionSvc        := service.NewFormSectionService(formSectionRepo)
    questionSvc           := service.NewQuestionService(questionRepo)
    repeaterGroupSvc      := service.NewRepeaterGroupService(repeaterGroupRepo)
    visibilityCondSvc     := service.NewVisibilityConditionService(visibilityCondRepo)
    formSubmissionSvc     := service.NewFormSubmissionService(formSubmissionRepo)
    repeaterEntrySvc      := service.NewRepeaterEntryService(repeaterEntryRepo)
    answerSvc             := service.NewAnswerService(answerRepo)
    optionSvc             := service.NewOptionService(optionRepo)
    followUpV2Svc         := service.NewFollowUpV2Service(followUpRepo, formSubmissionRepo, barrierV2Repo, victimCaseLightRepo, townLightRepo, attemptRepo, emRepo, psRepo, esRepo, agentLightRepo, caseTimelineRepo)
    victimCaseFormSvc     := service.NewVictimCaseFormService(victimCaseFormRepo, victimCaseLightRepo, caseTimelineRepo, followUpV2Svc)

    formSvc := service.NewFormService(service.FormServiceDeps{
        FormRepo:                  formRepo,
        FormSectionRepo:           formSectionRepo,
        QuestionRepo:              questionRepo,
        RepeaterGroupRepo:         repeaterGroupRepo,
        OptionRepo:                optionRepo,
        VisibilityCondRepo:        visibilityCondRepo,
        RenderModificationRepo:    renderModificationRepo,
        FormSubmissionRepo:        formSubmissionRepo,
        RepeaterEntryRepo:         repeaterEntryRepo,
        AnswerRepo:                answerRepo,
        FollowUpRepo:              followUpRepo,
        EmergencyMeasureRepo:      emRepo,
        PsychosocialSupportRepo:   psRepo,
        EconomicStabilizationRepo: esRepo,
        MenTeamRemisionRepo:       menTeamRemisionRepo,
        DiscapacidadRemisionRepo:  discapacidadRemisionRepo,
        BarrierV2Repo:             barrierV2Repo,
        BarrierFollowUpRepo:       barrierFollowUpRepo,
        CaseTimelineEventRepo:     caseTimelineRepo,
        AgentLightRepo:            agentLightRepo,
        CasoCierreService:         casoCierreSvc,
        CaseRepo:                  victimCaseLightRepo,
        FollowUpV2Svc:             followUpV2Svc,
        CaseTaskRepo:              caseTaskRepo,
        EntityLetterRepo:          entityLetterRepo,
		TeamContactRepo:           teamContactRepo,
        DuplaRepo:                 duplaRepo,
        VictimCaseFormSvc:         victimCaseFormSvc,
    })
    caseDetailSvc         := service.NewCaseDetailService(caseDetailRepo, gormDB)
    caseInfoSvc           := service.NewCaseInfoService(caseInfoRepo)
    reportSvc             := service.NewReportService(reportRepo)
    casesListSvc          := service.NewCasesListService(casesListRepo)
    casesReassignSvc      := service.NewCasesReassignService(casesReassignRepo, caseTimelineRepo, gormDB)
    agentsSearchSvc       := service.NewAgentsSearchService(agentLightRepo)
    psychosocialListSvc   := service.NewPsychosocialListService(psychosocialListRepo, duplaRepo)
    psychosocialReassignSvc := service.NewPsychosocialReassignService(psychosocialReassignRepo, gormDB)
    duplaAdminSvc         := service.NewDuplaAdminService(duplaRepo, psychosocialReassignRepo)
    psychosocialCalendarSvc := service.NewPsychosocialCalendarService(psychosocialCalendarRepo)

    // Inyectar el servicio en el controller legacy para generación automática del calendario
    salvia_legacy.FollowUpSvc = followUpV2Svc
    salvia_legacy.CaseTimelineRepo = caseTimelineRepo
    salvia_legacy.VictimCaseLightRepo = victimCaseLightRepo

    // Controllers
    formCtrl               := salvia_ctrl.NewFormController(formSvc, victimCaseFormSvc)
    formSectionCtrl        := salvia_ctrl.NewFormSectionController(formSectionSvc)
    questionCtrl           := salvia_ctrl.NewQuestionController(questionSvc)
    repeaterGroupCtrl      := salvia_ctrl.NewRepeaterGroupController(repeaterGroupSvc)
    visibilityCondCtrl     := salvia_ctrl.NewVisibilityConditionController(visibilityCondSvc)
    formSubmissionCtrl     := salvia_ctrl.NewFormSubmissionController(formSubmissionSvc)
    repeaterEntryCtrl      := salvia_ctrl.NewRepeaterEntryController(repeaterEntrySvc)
    answerCtrl             := salvia_ctrl.NewAnswerController(answerSvc)
    followUpV2Ctrl         := salvia_ctrl.NewFollowUpV2Controller(followUpV2Svc)
    optionCtrl             := salvia_ctrl.NewOptionController(optionSvc)
    caseDetailCtrl         := salvia_ctrl.NewCaseDetailController(caseDetailSvc, caseTimelineRepo)
    caseInfoCtrl           := salvia_ctrl.NewCaseInfoController(caseInfoSvc, gormDB)
    reportCtrl             := salvia_ctrl.NewReportController(reportSvc)
    barrierV2Svc           := service.NewBarrierV2Service(barrierV2Repo, barrierFollowUpRepo, agentLightRepo)
    barrierV2GinCtrl       := salvia_ctrl.NewBarrierV2GinController(barrierV2Svc)

    entityLetterSvc        := service.NewEntityLetterService(entityLetterRepo, caseTimelineRepo, caseTaskRepo, barrierV2Repo)
    entityLetterCtrl       := salvia_ctrl.NewEntityLetterController(entityLetterSvc)

    entityCaseSvc          := service.NewEntityCaseService(entityCaseRepo)
    entityCaseCtrl         := salvia_ctrl.NewEntityCaseController(entityCaseSvc)
    entityAPICtrl          := salvia_ctrl.NewEntityAPIController(entityCaseSvc)

    caseTaskSvc             := service.NewCaseTaskService(service.CaseTaskServiceDeps{
        CaseTaskRepo:     caseTaskRepo,
        BarrierV2Repo:    barrierV2Repo,
        CaseTimelineRepo: caseTimelineRepo,
        EntityLetterSvc:  entityLetterSvc,
        EntityLetterRepo: entityLetterRepo,
        DB:               gormDB,
    })
    caseTaskCtrl            := salvia_ctrl.NewCaseTaskController(caseTaskSvc)

    entityBranchAPICtrl    := salvia_ctrl.NewEntityBranchAPIController(gormDB)
    casesListCtrl          := salvia_ctrl.NewCasesListController(casesListSvc)
    casesReassignCtrl      := salvia_ctrl.NewCasesReassignController(casesReassignSvc)
    agentsSearchCtrl       := salvia_ctrl.NewAgentsSearchController(agentsSearchSvc)
    psychosocialListCtrl   := salvia_ctrl.NewPsychosocialListController(psychosocialListSvc)
    psychosocialDetailCtrl := salvia_ctrl.NewPsychosocialDetailController(gormDB)
    psychosocialReassignCtrl := salvia_ctrl.NewPsychosocialReassignController(psychosocialReassignSvc)
    duplaAdminCtrl         := salvia_ctrl.NewDuplaAdminController(duplaAdminSvc)
    // Flujo 3x3 de Atención Psicosocial
    contactAttemptRepo      := repository.NewContactAttemptRepository(gormDB)
    psychosocial3x3Svc      := service.NewPsychosocial3x3Service(contactAttemptRepo, gormDB)
    psychosocialContactCtrl := salvia_ctrl.NewPsychosocialContactController(psychosocial3x3Svc)
    psychosocialCalendarCtrl := salvia_ctrl.NewPsychosocialCalendarController(psychosocialCalendarSvc)
    followUpV2Repo         := repository.NewFollowUpV2Repository(gormDB)
    assignCaseSvc          := service.NewAssignCaseService(victimCaseLightRepo, agentLightRepo, followUpV2Repo)
    assignCaseCtrl         := salvia_ctrl.NewAssignCaseController(assignCaseSvc)

    locationRepo := repository.NewLocationRepository(gormDB)
    locationCtrl := salvia_ctrl.NewLocationController(locationRepo)

    directoryRepo := repository.NewDirectoryRepository(gormDB)
    directorySvc := service.NewDirectoryService(directoryRepo, locationRepo)
    directoryCtrl := salvia_ctrl.NewDirectoryController(directorySvc)

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
    caseInfoCtrl.RegisterRoutes(api)
    reportCtrl.RegisterRoutes(api)
    entityLetterCtrl.RegisterRoutes(api)
    entityCaseCtrl.RegisterRoutes(api)
    entityAPICtrl.RegisterRoutes(api)
    entityBranchAPICtrl.RegisterRoutes(api)
    barrierV2GinCtrl.RegisterRoutes(api)
    caseTaskCtrl.RegisterRoutes(api)
    locationCtrl.RegisterRoutes(api)
    directoryCtrl.RegisterRoutes(api)
    casesListCtrl.RegisterRoutes(api)
    casesReassignCtrl.RegisterRoutes(api)
    agentsSearchCtrl.RegisterRoutes(api)
    psychosocialListCtrl.RegisterRoutes(api)
    psychosocialDetailCtrl.RegisterRoutes(api)
    psychosocialReassignCtrl.RegisterRoutes(api)
    duplaAdminCtrl.RegisterRoutes(api)
    psychosocialContactCtrl.RegisterRoutes(api)
    psychosocialCalendarCtrl.RegisterRoutes(api)
    assignCaseCtrl.RegisterRoutes(api)

    // Admin: endpoints de migración (protegidos por X-Security-Key)
    migrateCtrl := salvia_ctrl.NewMigrateController(gormDB)
    migrateCtrl.RegisterRoutes(api)

    // Admin: búsqueda de usuarios
    adminUsersCtrl := salvia_ctrl.NewAdminUsersController(gormDB)
    adminUsersCtrl.RegisterRoutes(api)

    // Admin: reporte Excel de usuarios
    adminUsersReportCtrl := salvia_ctrl.NewAdminUsersReportController(gormDB)
    adminUsersReportCtrl.RegisterRoutes(api)
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

// ─── migrateLegacyTables ─────────────────────────────────────────────────────
// AutoMigrate para tablas legacy que no tienen modelos GORM propios.
// Cuando se necesite agregar una columna nueva a una tabla legacy,
// solo hay que agregarla al struct correspondiente aquí.
// GORM detecta columnas faltantes y las crea automáticamente (ADD COLUMN).
func migrateLegacyTables(db *gorm.DB) {
    // ─── security.general_user ───
    type GeneralUserSync struct {
        ID                 uint    `gorm:"column:general_user_id;primaryKey"`
        Team               string  `gorm:"column:general_user_team;type:varchar(50)"`
        AssignedDepartment string  `gorm:"column:general_user_assigned_department;type:varchar(20)"`
        // EntityBranchId: sede (salvia.entity_branch) del usuario rol et. Nullable.
        EntityBranchId *int64 `gorm:"column:entity_branch_id;type:bigint;index"`
    }
    if err := db.Table("security.general_user").AutoMigrate(&GeneralUserSync{}); err != nil {
        log.Printf("[WARN] AutoMigrate security.general_user: %v", err)
    }

    // ─── salvia.victim_case ───
    type VictimCaseSync struct {
        ID       uint   `gorm:"column:victim_case_id;primaryKey"`
        Team     string `gorm:"column:victim_case_team;type:varchar(64)"`
        AgentID  string `gorm:"column:agent_id;type:varchar(64)"`
    }
    if err := db.Table("salvia.victim_case").AutoMigrate(&VictimCaseSync{}); err != nil {
        log.Printf("[WARN] AutoMigrate salvia.victim_case: %v", err)
    }
}
