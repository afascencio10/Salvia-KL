package common_facades

import (
	common_config "bitsflow/common/config"
	"bitsflow/common/utils"
	security_config "bitsflow/security/config"
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var modules map[string]string = map[string]string{}

var router *gin.Engine
var staticURI string = "/static"
var publicURI string = "/public"
var sem = make(chan struct{}, 9999) // Limitar a 9999 peticiones concurrentes
func limitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sem <- struct{}{}
		defer func() {
			<-sem
		}()
		c.Next()
	}
}

func InitRouter() *gin.Engine {
	modules["security"] = "security"
	modules["salvia"] = "salvia"

	//loginModule = modules["security"]

	router = gin.Default()
	store := cookie.NewStore([]byte("secret"))

	store.Options(sessions.Options{
		MaxAge:   14400,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	})
	// Middleware para restringir la cantidad de threads concurrentes
	router.Use(limitMiddleware())
	
	if os.Getenv("PORT") == "" {
		// Middleware para redirigir HTTP a HTTPS solo en local
		router.Use(ForceHTTPS())
	}

	router.Use(sessions.Sessions("default_session", store))

	router.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("refreshKey", utils.GetRandomString(3))
		session.Save()
		c.Next()
	})

	router.Use(AuthMiddleware())

	router.MaxMultipartMemory = 50 << 20

	publicFS, _ := fs.Sub(utils.FrontendAssets, "frontend")
	// ---------------------------
	// 1. Cargar plantillas embed
	// ---------------------------

	var tmplFiles []string
	fs.WalkDir(utils.FrontendAssets, "frontend/templates", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && (strings.HasSuffix(path, ".tmpl") || strings.HasSuffix(path, ".html")) {
			tmplFiles = append(tmplFiles, path)
			fmt.Println("Loaded file:", path)
		}
		return nil
	})

	tmpl := template.Must(template.ParseFS(utils.FrontendAssets, tmplFiles...))

	router.SetHTMLTemplate(tmpl)

	for _, t := range tmpl.Templates() {
		fmt.Println("Loaded template:", t.Name())
	}

	// ---------------------------
	// 2. Servir archivos estáticos
	// ---------------------------
	router.StaticFS("/static", http.FS(publicFS))

	router.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("favicon.ico", http.FS(publicFS))
	})

	//router.StaticFile("/favicon.ico", "./frontend/favicon.ico")

	//router.LoadHTMLGlob("frontend/templates/**/*")

	//Templates, _ = template.ParseGlob("frontend/templates/**/*")

	router.GET("/", home)

	return router
}

// NewHTTPServer retorna un *http.Server configurado pero sin arrancar.
// Úsalo en main.go para implementar graceful shutdown.
func NewHTTPServer() *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "443"
	}
	return &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
}

func StartRouter() {
	port := os.Getenv("PORT")
	if port != "" {
		err := router.Run(":" + port)
		if err != nil {
			println("Error iniciando servidor en puerto " + port + ": ", err.Error())
		}
	} else {
		err := router.RunTLS(":443", "certs/salvia.crt", "certs/salvia.key")
		if err != nil {
			println("Error iniciando servidor TLS local: ", err.Error())
		}
	}
}

func ForceHTTPS() gin.HandlerFunc {
	return func(c *gin.Context) {
		go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "https://"+req.Host+req.URL.String(), http.StatusMovedPermanently)
		}))
		c.Next()
	}
}

func home(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, staticURI+"/landing.html")
	c.Abort()
}

/*
*

	Esta función se encarga de renderizar la plantilla HTML.

	@c 				Contexto gin
	@entity			Nombre de la entidad que representa el request
	@htmlFolder 	Directorio donde se encuentra el html con el contenido principal
	@configTemplate Nombre de la constante para ubicar el archivo html principal con el contenido
	@htmlScripts	Ruta completa del html que importa los scripts javascript
	@viewTemplate	Ruta completa con la vista que representa toda la página a renderizar
	@errorTemplate 	Ruta completa del html de error en caso de presentarse, se renderiza

*
*/
func RenderTemplate(c *gin.Context, entity string, module string, htmlFolder string, configTemplate map[string]string, configTemplateName string, extraTemplates []string, viewTemplate string, errorTemplate string, templateFields map[string]interface{}, funcMap template.FuncMap) bool {
	var res string
	var code int = http.StatusOK

	var errorMap map[string]map[string]string = map[string]map[string]string{}

	var b bytes.Buffer
	t := template.New(configTemplate[configTemplateName]).Funcs(funcMap)
	t, err := t.ParseFS(utils.FrontendAssets, append([]string{security_config.HTML_Templates_folder + modules[module] + "/" + htmlFolder + configTemplate[configTemplateName]}, extraTemplates...)...)

	if err != nil {
		utils.SetError(errorMap, entity, "default", common_config.Enums.GLOBAL_ERROR, "", common_config.Locale)
		res = utils.CommMsgGetJSONErrors(errorMap)
		code = 400
	} else {
		err = t.Execute(&b, templateFields)
		if err != nil {
			utils.SetError(errorMap, entity, "default", common_config.Enums.GLOBAL_ERROR, "", common_config.Locale)
			res = utils.CommMsgGetJSONErrors(errorMap)
			code = 400
		}
	}

	if code == http.StatusOK {
		c.HTML(
			code,
			viewTemplate,
			gin.H{
				"content": template.HTML(b.String()),
			},
		)
	} else {
		c.HTML(
			code,
			errorTemplate,
			gin.H{
				"content": res,
			},
		)
	}
	return code == http.StatusOK
}

func IsLoggedIn(c *gin.Context) bool {

	session := sessions.Default(c)
	var error error = nil
	v := session.Get("userData")
	if v != nil {
		_, error = utils.GetCommonSession(v.(string))
	}

	return v != nil && error == nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		location := c.Request.RequestURI

		if !strings.Contains(location, publicURI) && !strings.Contains(location, staticURI) && !IsLoggedIn(c) && !strings.Contains(location, "login") && !strings.HasPrefix(location, "/api/") && !strings.HasPrefix(location, "/swagger") {
			//c.Redirect(302, "/"+security_config.Locale["sp"][loginModule]+"/login")
			c.Redirect(http.StatusTemporaryRedirect, staticURI+"/landing.html")
			c.Abort()
			return
		}

		// El usuario está autenticado, continuar con la siguiente ruta
		c.Next()
	}
}

func SetHeaderNoCache(c *gin.Context) {
	// Agregar encabezados para evitar el almacenamiento en caché
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}
