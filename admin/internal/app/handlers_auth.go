package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) loginPage(c *gin.Context) {
	if _, ok := a.sessions.Get(c); ok {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}
	c.HTML(http.StatusOK, "auth/login", gin.H{
		"Title": "登录",
		"Error": c.Query("error"),
	})
}

func (a *App) apiLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if !bindJSON(c, &req) {
		return
	}
	user, err := a.store.Authenticate(req.Username, req.Password)
	if err != nil {
		jsonMessage(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	if err := a.sessions.Create(c, user.ID, user.Username); err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
		},
	})
}

func (a *App) apiLogout(c *gin.Context) {
	a.sessions.Destroy(c)
	c.JSON(http.StatusOK, gin.H{"message": "已退出"})
}

func (a *App) apiMe(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       currentUserID(c),
			"username": username,
		},
	})
}

func (a *App) apiNav(c *gin.Context) {
	nav, err := a.navMenus(currentUserID(c))
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"menus": nav})
}
