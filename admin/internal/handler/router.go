package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puoxiu/discron/admin/internal/middlerware"
	"github.com/puoxiu/discron/admin/internal/model/resp"
)


func RegisterRouters(r *gin.Engine) {
	r.Use(middlerware.Cors())

	configRoute(r)

	configNoRoute(r)
}


func configRoute(r *gin.Engine) {
	hello := r.Group("/ping")
	{
		hello.GET("", func(c *gin.Context) {
			c.JSON(200, "pong")
		})
		hello.POST("", func(c *gin.Context) {
			type Hello struct {
				Name string `json:"name" form:"name"`
			}
			var h Hello
			var err error
			err = c.ShouldBindJSON(&h)
			if err != nil {
				c.JSON(resp.ERROR, err.Error())
			}
			c.JSON(200, "hello,"+h.Name)
		})
	}

	base := r.Group("")
	{
		base.POST("register", defaultUserRouter.Register)
		base.POST("login", defaultUserRouter.Login)
	}

	stat := r.Group("/statis")
	stat.Use(middlerware.JWTAuth())
	{
		stat.GET("today", defaultStatRouter.GetTodayStatistics)
		stat.GET("week", defaultStatRouter.GetWeekStatistics)
		stat.GET("system", defaultStatRouter.GetSystemInfo)

	}

	job := r.Group("/job")
	job.Use(middlerware.JWTAuth())
	{
		job.POST("add", defaultJobRouter.CreateOrUpdate)
		job.POST("del", defaultJobRouter.Delete)
		job.GET("find", defaultJobRouter.FindById)
		job.POST("search", defaultJobRouter.Search)
		job.POST("log", defaultJobRouter.SearchLog)
		job.POST("once", defaultJobRouter.Once)
		// job.POST("kill", defaultJobRouter.Kill)
	}

	user := r.Group("/user")
	user.Use(middlerware.JWTAuth())
	{
		user.POST("del", defaultUserRouter.Delete)
		user.POST("update", defaultUserRouter.Update)
		user.POST("change_pw", defaultUserRouter.ChangePassword)
		user.GET("find", defaultUserRouter.FindById)
		user.POST("search", defaultUserRouter.Search)
	}
	node := r.Group("/node")
	node.Use(middlerware.JWTAuth())
	{
		node.POST("search", defaultNodeRouter.Search)
		node.POST("del", defaultNodeRouter.Delete)
	}
	script := r.Group("/script")
	script.Use(middlerware.JWTAuth())
	{
		script.POST("add", defaultScriptRouter.CreateOrUpdate)
		script.POST("del", defaultScriptRouter.Delete)
		script.GET("find", defaultScriptRouter.FindById)
		script.POST("search", defaultScriptRouter.Search)
	}
}

// func configNoRoute(r *gin.Engine) {
// 	r.LoadHTMLGlob("./dist/*.html") // npm打包成dist的路径
// 	r.StaticFile("favicon.ico", "./dist/favicon.ico")
// 	r.Static("/css", "./dist/css")         // dist里面的静态资源
// 	r.Static("/fonts", "./dist/fonts")     // dist里面的静态资源
// 	r.Static("/js", "./dist/js")           // dist里面的静态资源
// 	r.Static("/img", "./dist/img")         // dist里面的静态资源
// 	r.StaticFile("/", "./dist/index.html") // 前端网页入口页面
// }

func configNoRoute(r *gin.Engine) {
    // 加载前端打包的 HTML 文件
    r.LoadHTMLGlob("./dist/*.html")
    
    // 映射所有静态资源（关键修改：将 assets 目录整体暴露）
    r.Static("/assets", "./dist/assets")  // 匹配 dist/assets 下的所有资源
    
    // 映射根目录的静态文件（如 vite.svg）
    r.StaticFile("/vite.svg", "./dist/vite.svg")
    
    // 单页应用入口：所有未匹配路由指向 index.html
    r.NoRoute(func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })
}