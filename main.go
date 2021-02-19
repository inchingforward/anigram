package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/flosch/pongo2/v4"
	"github.com/gorilla/sessions"
	_ "github.com/iris-contrib/pongo2-addons/v4"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Renderer renders templates.
type Renderer struct {
	TemplateDir   string
	Reload        bool
	TemplateCache map[string]*pongo2.Template
}

var (
	debug    bool
	storeKey string
)

// GetTemplate returns a template, loading it every time if reload is true.
func (r *Renderer) GetTemplate(name string, reload bool) *pongo2.Template {
	filename := path.Join(r.TemplateDir, name)

	if r.Reload {
		return pongo2.Must(pongo2.FromFile(filename))
	}

	return pongo2.Must(pongo2.FromCache(filename))
}

// Render renders a pongo2 template.
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	template := r.GetTemplate(name, debug)
	pctx := data.(pongo2.Context)

	sess, err := session.Get("session", c)
	if err == nil {
		pctx["session"] = sess
	}

	return template.ExecuteWriter(pctx, w)
}

func init() {
	// db, err := sqlx.Connect("postgres", "user=postgres dbname=logbook sslmode=disable")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//setDB(db)

	storeKey := os.Getenv("ANIGRAM_STORE_KEY")
	if storeKey == "" {
		log.Fatal("The ANIGRAM_STORE_KEY is not set")
	}

	//gob.Register(SessionUser{})
}

func getIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", pongo2.Context{})
}

func getAbout(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented)
}

func getNewAnimation(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented)
}

func getAnimationEdit(c echo.Context) error {
	return c.Render(http.StatusOK, "animation_edit.html", pongo2.Context{})
}

func getAnimation(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented)
}

func postAnimation(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented)
}

func main() {
	flag.BoolVar(&debug, "debug", false, "true to enable debug")
	flag.Parse()

	log.Printf("debug: %v\n", debug)

	InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(storeKey))))

	e.Renderer = &Renderer{TemplateDir: "templates", Reload: debug, TemplateCache: make(map[string]*pongo2.Template)}

	e.Static("/static", "static")
	e.GET("/", getIndex)
	e.GET("/about", getAbout)
	e.GET("/animations/new", getNewAnimation)
	e.GET("/animations/:uuid/edit", getAnimationEdit)
	e.GET("/api/animations/:uuid", getAnimation)
	e.POST("/api/animations", postAnimation)

	e.Logger.Fatal(e.Start(":5000"))
}
