package handle

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"hiveify-core/database"
	"hiveify-core/response"
	"hiveify-core/util"
	"net/http"
)

const passwordSalt = "Yw#Lwpak#525178"

// RegisterHandler 用户注册
func RegisterHandler(c *gin.Context) {
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	email := c.Request.FormValue("email")
	fmt.Println(name, password, email)
	if len(name) == 0 || len(password) < 8 || !util.CheckEmail(email) {
		response.Fail(c, http.StatusInternalServerError, "参数不合法！")
		return
	}
	if signResult := database.Register(name, util.StringSha512(password+passwordSalt), email); !signResult {
		response.Fail(c, http.StatusInternalServerError, "注册失败！")
		return
	}
	response.Success(c, "注册成功！", database.User{
		Name:  name,
		Email: email,
		Phone: "",
	})
}

// LoginHandler 用户登录
func LoginHandler(c *gin.Context) {
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")

	if !util.CheckEmail(email) || len(password) < 8 {
		response.Fail(c, http.StatusInternalServerError, "参数不合法！")
		return
	}

	// 验证用户
	user, err := database.Login(email, util.StringSha512(password+passwordSalt))
	if err != nil {
		response.Fail(c, http.StatusNotFound, "email或手机或密码错误！")
		return
	}

	// 生成token
	token, err := util.GetToken(user)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "登陆失败！")
		return
	}
	response.Success(c, "登录成功！", token)
}

// GetUserInfoHandler 获取用户信息
func GetUserInfoHandler(c *gin.Context) {
	token, err := c.Cookie("Authorization")
	if err != nil {
		log.Warnf("Failed to get cookie %s", err.Error())
		response.Fail(c, http.StatusUnauthorized, "未登录，请先登录！")
		return
	}
	user, err := util.GetUserByString(token)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "未登录，请先登录！")
		return
	}
	response.Success(c, "获取用户信息成功！", user)
}
