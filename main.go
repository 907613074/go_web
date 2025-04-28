package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go_web/docs"
	"github.com/go_web/eth"
	"github.com/go_web/orm"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	eth.PrivateKey()
	// eth.TestQueryBlock()

	eth.Wallet()

	// eth.QueryTokenBalance()
	// eth.DigitalAssetToken()
	eth.LoadContract()
	// eth.QueryEvent()
	// eth.SubscribeNewHead()
	/*
		task3.InitDB()
		orm.InitDB()
		orm.QueryUsers(&orm.User{})

			orm.CURD()
			restfulAPI()

			task1.SingleNum()
			print()
			task1.Rob()
			print()
			task1.RverseString()
			print()

			task1.Sqrt(10)
			print()

			task1.RemoveDuplicates()
			print()
			task1.Merge()
			print()
			cal := task1.MyCalendar{}
			cal.Book(10, 20)
			cal.Book(15, 25)
			cal.Book(20, 30)
			fmt.Println("Bookings ", cal.Bookings)
			print()

			p := 11
			task2.Add10(&p)
			print()

			num := []int{1, 2, 13, 4, 1, 6, 5, 28, 6, 10}
			task2.DoubleSlice(&num)
			print()

			task2.GoRoutine()
			print()
			task2.TaskScheduler()
			print()

			task2.Inf()
			print()
			task2.Employee{Person: task2.Person{Name: "Alice", Age: 25}, EmployeeID: 1001}.PrintInfo()
			print()
	*/
}

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Data: data})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: msg})
}

func getInt(key string) uint {
	_id, _ := strconv.ParseUint(key, 10, 64)
	return uint(_id)
}

// @Summary Ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} string "pong"
// @Router /ping [get]
func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func restfulAPI() {
	r := gin.Default()

	r.GET("/ping", pingHandler)

	v1 := r.Group("/api/v1")
	{

		v1.GET("/users", func(c *gin.Context) {
			user := &orm.User{
				Name:  c.DefaultQuery("name", ""),
				Age:   int(getInt(c.Query("age"))),
				Email: c.DefaultQuery("email", ""),
			}
			user.ID = getInt(c.Query("id"))
			Success(c, orm.QueryUsers(user))
		})

		v1.GET("/users/:id", func(c *gin.Context) {
			var user orm.User
			user.ID = getInt(c.Param("id"))
			Success(c, orm.QueryUsers(&user))
		})

		r.POST("/users", func(c *gin.Context) {
			var newUser orm.User
			if err := c.ShouldBindJSON(&newUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			orm.CreateUser(&newUser)
			Success(c, newUser)
		})

		r.PUT("/users/:id", func(c *gin.Context) {
			var updatedUser orm.User
			if err := c.ShouldBindJSON(&updatedUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			updatedUser.ID = getInt(c.Param("id"))
			orm.UpdateUser(&updatedUser)
			Success(c, updatedUser)
		})

		r.DELETE("/users/:id", func(c *gin.Context) {
			orm.DeleteUser(getInt(c.Param("id")))
			Success(c, nil)
		})

	}

	// 添加docs路由（main函数内）
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 配置swagger.json路径（初始化路由前添加）
	r.StaticFile("/swagger.json", "./docs/swagger.json")

	// 添加文档访问中间件
	authMiddleware := gin.BasicAuth(gin.Accounts{
		"admin": "swagger123",
	})
	r.GET("/docs", authMiddleware, ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080" // 默认端口
	}
	log.Printf("http://localhost%s", port)

	// 启动服务
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
