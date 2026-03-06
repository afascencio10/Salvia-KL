package main

import (
	common_routers "bitsflow/common/facades"
	"bitsflow/common/utils"
	salvia_facades "bitsflow/salvia/facades"
	security_routers "bitsflow/security/facades"
	"embed"

	"github.com/gin-gonic/gin"
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

	common_routers.StartRouter()
}
