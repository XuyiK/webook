package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
	"webook/config"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middlewares"
	"webook/pkg/ginx/middleware/ratelimit"
)

func main() {
	//db := initDB()
	//server := initWebServer()
	//initUser(server, db)

	server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, it's started")
	})

	server.Run(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		// 用于允许前端访问你后端响应中的头部
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	redisClient := redis.NewClient(&redis.Options{Addr: config.Config.Redis.Addr})
	//useSession(server)
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	useJWT(server)
	return server
}

func initUser(server *gin.Engine, db *gorm.DB) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	c := web.NewUserHandler(us)
	c.RegisterRoutes(server)
}

func useJWT(server *gin.Engine) {
	login := &middlewares.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

func useSession(server *gin.Engine) {
	login := &middlewares.LoginMiddlewareBuilder{}
	// 基于cookie 实现
	store := cookie.NewStore([]byte("secret"))
	// 基于内存实现
	//store := memstore.NewStore([]byte("Upxnmo6PEdrbKRMBfCVOjUjjEoJY4D9e"), []byte("pBjVy5p318kMACbu84sjRTPpCOpV03HD"))
	// 基于Redis
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("Upxnmo6PEdrbKRMBfCVOjUjjEoJY4D9e"), []byte("pBjVy5p318kMACbu84sjRTPpCOpV03HD"))
	//if err != nil {
	//	panic(err)
	//}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
