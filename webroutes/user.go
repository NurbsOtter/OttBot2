package webroutes

import (
	"OttBot2/models"
	"fmt"
	"github.com/kataras/iris"
	"strings"
)

type BareMinRender struct {
	LoggedIn bool
}
type UserOutStruct struct {
	Name        string
	LoggedIn    bool
	UserID      int64
	PingAllowed bool
	Warnings    []models.Warning
	Aliases     []models.UserAlias
}
type AliasStruct struct {
	LoggedIn bool
	NotFound bool
	Aliases  []models.UserAlias
}
type SearchStruct struct {
	UserName string `json:"UserName"`
}
type WarnStruct struct {
	WarningText string
}

func AddUser(ctx *iris.Context) {
	newUser := &models.User{}
	err := ctx.ReadJSON(&newUser)
	if err != nil {
		panic(err)
	}
	if (newUser.UserName == "") || (newUser.Password == "") {
		ctx.Write("You suck ass!")
	} else {
		if models.InsertUser(strings.Trim(strings.ToLower(newUser.UserName), " "), newUser.Password) {
			ctx.Write("Worked!")
		} else {
			ctx.Write("Username Taken")
		}
	}
}
func ServeIndex(ctx *iris.Context) {
	ctx.ServeFile("./static/index.html", false)
}
func GetLogin(ctx *iris.Context) {
	newUser := models.User{}
	err := ctx.ReadForm(&newUser)
	if err != nil {
		panic(err)
	}
	foundUser := models.VerifyUser(strings.Trim(strings.ToLower(newUser.UserName), " "), newUser.Password)
	if foundUser != nil {
		ctx.Session().Set("userID", foundUser.ID)
		ctx.Redirect("/home")
	} else {
		ctx.Redirect("/login")
	}
}
func GetLogout(ctx *iris.Context) {
	ctx.Session().Clear()
	ctx.Redirect("/login")
}
func GetRenderIndex(ctx *iris.Context) {
	sessUser := models.GetUserFromContext(ctx)
	if sessUser != nil {
		ctx.Redirect("/users")
	}
	outStruct := BareMinRender{LoggedIn: sessUser != nil}
	ctx.Render("login.html", outStruct)
}
func RenderSearchPage(ctx *iris.Context) {
	sessUser := models.GetUserFromContext(ctx)
	userOutStruct := BareMinRender{LoggedIn: sessUser != nil}
	if sessUser != nil {
		ctx.Render("searchusers.html", userOutStruct)
	} else {
		ctx.Redirect("/")
	}
}
func SearchByUName(ctx *iris.Context) {
	sessUser := models.GetUserFromContext(ctx)
	if sessUser != nil {
		foundSearch := SearchStruct{}
		ctx.ReadForm(&foundSearch)
		foundSearch.UserName = strings.ToLower(foundSearch.UserName)
		foundSearch.UserName = strings.Trim(foundSearch.UserName, " ")
		user := models.SearchUserByUsername(foundSearch.UserName)
		if user == nil {
			ctx.Redirect("/users")
		} else {
			outString := fmt.Sprintf("/user/%d", user.ID)
			ctx.Redirect(outString)
		}

	}
}
func SearchByAlias(ctx *iris.Context) {
	sessUser := models.GetUserFromContext(ctx)
	if sessUser != nil {
		foundSearch := AliasStruct{}
		foundSearch.LoggedIn = true
		aliasSearch := SearchStruct{}
		ctx.ReadForm(&aliasSearch)
		aliases := models.SearchAliases(aliasSearch.UserName)
		if len(aliases) == 0 {
			foundSearch.NotFound = true
		} else {
			foundSearch.Aliases = aliases
			foundSearch.NotFound = false
		}
		ctx.Render("searchalias.html", foundSearch)
	}
}
func ShowUser(ctx *iris.Context) {
	sessUser := models.GetUserFromContext(ctx)
	if sessUser == nil {
		ctx.Redirect("/index.html")
		return
	} else {
		userID, err := ctx.ParamInt64("id")
		if err != nil {
			fmt.Println("Failed to parse a route")
			ctx.Redirect("/")
			return
		}
		user := models.ChatUserFromID(userID)
		if user == nil {
			ctx.WriteString("Not found") //Todo proper handler.
		} else {
			fmt.Println(user)
			outData := UserOutStruct{}
			outData.Name = user.UserName
			outData.LoggedIn = true
			outData.Warnings = models.GetUsersWarnings(user)
			outData.Aliases = user.GetAliases()
			outData.PingAllowed = user.PingAllowed
			outData.UserID = user.ID
			ctx.Render("userdisplay.html", outData)
		}
	}
}
func WarnUser(ctx *iris.Context) {
	sessUser := models.GetUserFromContext(ctx)
	if sessUser != nil {
		userID, err := ctx.ParamInt64("userID")
		if err != nil {
			fmt.Println("Failed to parse a userID on the warning web route.")
			ctx.Redirect("/home")
		}
		inWarn := WarnStruct{}
		err = ctx.ReadForm(&inWarn)
		if err != nil {
			fmt.Println("Failed to parse a form on the warning web route.")
			ctx.Redirect("/home")
			return
		}
		models.AddWarningToID(userID, inWarn.WarningText)
		redir := fmt.Sprintf("/user/%d", userID)
		ctx.Redirect(redir)
	} else {
		ctx.Redirect("/login")
	}
}
func ToggleAllowedPing(ctx *iris.Context) {
	sessUser := models.GetUserFromContext(ctx)
	if sessUser == nil {
		ctx.Redirect("/login")
	} else {
		userID, err := ctx.ParamInt64("id")
		if err != nil {
			fmt.Println("Failed to parse a route")
			ctx.Redirect("/")
			return
		}
		user := models.ChatUserFromID(userID)
		if user == nil {
			ctx.WriteString("This is a false user. \nWhy for you do this. :(")
		} else {
			user.ToggleModPing()
			redir := fmt.Sprintf("/user/%d", user.ID)
			ctx.Redirect(redir)
		}
	}
}
