package orm

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go_web/model"
	"github.com/go_web/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DBUrl = "root:root@tcp(127.0.0.1:13306)/ain?charset=utf8mb4&parseTime=True&loc=Local"
)

// 定义数据库连接
var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open(mysql.Open(DBUrl), &gorm.Config{})
	if err != nil {
		panic("failed to connect ain")
	}

	fmt.Println("Database connected!")

	generatedModel()
}

func generatedModel() {
	// 初始化数据库连接（确保配置正确）
	mysql := sqlx.NewMysql(DBUrl)
	m := model.NewTLambPointModel(mysql)

	// 插入数据
	result, err := m.Insert(context.Background(), &model.TLambPoint{
		Person:   "张三",
		Page:     1,
		Stroke:   1,
		Lng:      123.456,
		Lat:      123.456,
		Ts:       sql.NullInt64{Int64: 1, Valid: true},
		Pressure: sql.NullInt64{Int64: 25, Valid: true},
	})
	if err != nil {
		log.Fatalf("插入数据失败: %v", err)
	}
	id, _ := result.LastInsertId()
	util.Log("插入结果: ", id)

	// 查询数据
	user, err := m.FindOne(context.Background(), id) // 假设主键是 ID
	if err != nil {
		log.Fatalf("查询数据失败: %v", err)
	}
	Json(user)

}

type MyTime struct {
	time.Time
}

// MarshalJSON 实现了 json.Marshaler 接口，以自定义时间格式
func (t MyTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:00:05"))
	return []byte(formatted), nil
}

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null;" json:"name"`
	Password string `gorm:"type:varchar(100);not null;" json:"password"`
	Email    string `gorm:"type:varchar(100);index;" json:"email"`
	Status   int8   `gorm:"default:1;" json:"status"`
	Age      int    `gorm:"default:18;" json:"age"`
}

type PageData struct {
	PageSize int
	Page     int
	Total    int64
	Array    interface{}
}

func CURD() {
	// 自动迁移，根据结构体创建或更新表
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 创建记录并选择字段
	db.Select("Name", "Email", "Age").Create(&User{Name: "张三", Email: "123@qq.com", Age: 25})

	// 批量插入
	users := []User{{Email: "456@qq.com", Name: "李四", Age: 15}, {Email: "789@qq.com", Name: "王五"}}
	db.CreateInBatches(users, 100)

	// 链式条件查询

	var results []struct {
		Age   int
		Count int
	}
	db.Model(&User{}).Select("age, COUNT(*) as count").Group("age").Find(&results)
	Json(results)

	// 条件更新
	updates := map[string]interface{}{
		"Name":   "New Name",
		"Status": 0, // 零值也会被更新
	}
	db.Model(&User{}).Where("age < ?", 18).Updates(updates)

	db.Delete(&User{}, "email = ?", "789@qq.com")
	// 批量硬删除
	db.Unscoped().Where("created_at <= ?", time.Now().Add(-30*time.Minute)).Delete(&User{})

	util.Separator()
}

func structToMapReflect(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	// 确保传入的是结构体
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		result[field.Name] = value.Interface()
	}

	return result
}

func CreateUser(user *User) {
	db.Create(user)
	Json(user)
}

func UpdateUser(user *User) {
	updates := structToMapReflect(user)
	db.Model(&User{}).Where("id = ?", user.ID).Updates(updates)
	Json(user)
}

func DeleteUser(id uint) {
	db.Delete(&User{}, "id = ?", id)
	Json(id)
}

func QueryUsers(user *User) PageData {
	page := 1
	pageSize := 5
	offset := (page - 1) * pageSize
	var userList []User
	var count int64

	where := db.Debug().Unscoped().Model(&User{})
	if user.ID != 0 {
		where = where.Where("id = ?", user.ID)
	}
	if user.Status != 0 {
		where = where.Where("status = ?", user.Status)
	}
	if user.Age != 0 {
		where = where.Where("age = ?", user.Age)
	}
	if user.Name != "" {
		where = where.Where("name like ?", "%"+user.Name+"%")
	}
	if user.Email != "" {
		where = where.Where("email like ?", "%"+user.Email+"%")
	}
	where.Count(&count)

	where.Order("created_at desc").Limit(pageSize).Offset(offset).Find(&userList)
	result := PageData{
		Page:     page,
		PageSize: pageSize,
		Total:    count,
		Array:    userList,
	}
	Json(result)
	return result
}

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var storedUser User
	if err := db.Where("username = ?", user.Name).First(&storedUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       storedUser.ID,
		"username": storedUser.Name,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	util.Log(tokenString)
}

func Json(v any) string {
	jsonData, err := json.MarshalIndent(v, "", "  ") //json格式化输出
	jsonData, err = json.Marshal(v)
	if err != nil {
		log.Fatalf("failed to serialize page data to JSON: %v", err)
	}
	r := string(jsonData)
	util.Log(r)
	return r
}
