package main

import (
	"OttBot2/models"
	"OttBot2/settings"
	"OttBot2/telegram"
	"OttBot2/webroutes"
	"github.com/kataras/go-template/html"
	"github.com/kataras/iris"
)

func hello(ctx *iris.Context) {
	ctx.Write("Hello world!")
}
func main() {
	models.MakeDB("./userdatabase.db")
	settings.LoadSettings()
	telegram.Register("\\/ping", settings.GetChannelID(), telegram.TestCmd)
	telegram.Register(".*", settings.GetChannelID(), telegram.HandleUsers)
	telegram.Register("^\\/info @.+", settings.GetControlID(), telegram.FindUserByUsername)
	telegram.Register("^\\/warn @.+", settings.GetControlID(), telegram.WarnUserByUsername)
	telegram.Register("^\\/find .+", settings.GetControlID(), telegram.LookupAlias)
	go telegram.InitBot(settings.GetBotToken())
	api := iris.New()
	api.StaticServe("./static/")
	api.UseTemplate(html.New(html.Config{
		Layout: "layout.html",
	})).Directory("./templates", ".html")
	api.Config.IsDevelopment = true
	api.Get("/", webroutes.ServeIndex)
	api.Post("/register", webroutes.AddUser)
	api.Get("/login", webroutes.GetRenderIndex)
	api.Get("/logout", webroutes.GetLogout)
	api.Post("/login", webroutes.GetLogin)
	api.Get("/user/:id", webroutes.ShowUser)
	api.Get("/users", webroutes.RenderSearchPage)
	api.Post("/users/uname", webroutes.SearchByUName)
	api.Post("/users/alias", webroutes.SearchByAlias)
	//api.ListenLETSENCRYPT("127.0.0.1:443")
	api.Listen(":8080")
}
