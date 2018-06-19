package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"./handlers"
	"./models"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
	db        *gorm.DB
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var viewContext map[string]interface{}

	if vs, isMap := data.(map[string]interface{}); isMap {
		viewContext = vs
	} else {
		viewContext = map[string]interface{}{}
	}

	viewContext["reverse"] = c.Echo().Reverse

	user, err := handlers.GetUserFromContext(t.db, c)
	if err != nil {
		return err
	}
	if user != nil {
		viewContext["user"] = user
	}

	return t.templates.ExecuteTemplate(w, name, viewContext)
}

func main() {
	db, err := gorm.Open("mysql", "root:@/cabinet?charset=utf8&parseTime=True&loc=Local")
	// db, err := gorm.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.User{})

	e := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
		db:        db,
	}
	e.Renderer = renderer

	e.Static("/static", "./static")

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	authCtl := handlers.AuthController{
		DB: db,
	}

	// unauthorizedGroup := e.Group("/auth")
	unauthorizedGroup := e.Group("")

	unauthorizedGroup.Use(authCtl.CheckSessionForUnauthorized)

	unauthorizedGroup.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	unauthorizedGroup.GET("/auth/login", authCtl.GetLoginFrom)
	unauthorizedGroup.POST("/auth/login", authCtl.PostLoginForm)
	unauthorizedGroup.GET("/auth/registration", authCtl.GetRegistrationForm)
	unauthorizedGroup.POST("/auth/registration", authCtl.PostRegistrationForm)

	authorizedGroup := e.Group("")

	authorizedGroup.Use(authCtl.CheckSessionForAuthorized)

	// authorizedGroup.GET("/", func(c echo.Context) error {
	// 	return c.Render(http.StatusOK, "index.html", nil)
	// })
	authorizedGroup.POST("/logout", authCtl.Logout)

	profCtl := handlers.ProfileController{
		DB: db,
	}
	authorizedGroup.GET("/profile", profCtl.GetProfilePage)
	authorizedGroup.POST("/profile", profCtl.UpdateProfile)

	e.Logger.Fatal(e.Start(":8080"))
}
