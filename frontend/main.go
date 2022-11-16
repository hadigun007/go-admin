package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	_ "github.com/hadigun007/go-admin/adapter/gin" // web framework adapter
	"github.com/hadigun007/go-admin/frontend/pages"

	_ "github.com/GoAdminGroup/themes/adminlte"                 // ui theme
	_ "github.com/hadigun007/go-admin/modules/db/drivers/mysql" // sql driver

	"github.com/hadigun007/go-admin/engine"

	"github.com/gin-gonic/gin"
	"github.com/hadigun007/go-admin/template"
	"github.com/hadigun007/go-admin/template/chartjs"
)

func main() {
	startServer()
}

func startServer() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()

	template.AddComp(chartjs.NewChart())

	eng := engine.Default()

	if err := eng.AddConfigFromYAML("./config.yml").
		// AddGenerators(tables.Generators).
		Use(r); err != nil {
		panic(err)
	}

	r.Static("/uploads", "./uploads")

	eng.HTML("GET", "/hadi", pages.GetDashBoard)
	eng.HTMLFile("GET", "/hadi/hello", "./html/hello.tmpl", map[string]interface{}{
		"msg": "Hello world",
	})

	_ = r.Run(":8800")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	eng.MysqlConnection().Close()
}
