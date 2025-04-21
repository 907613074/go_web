package task3

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go_web/orm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
题目1：基本CRUD操作
假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
要求 ：
编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
*/
var _db *sql.DB
var db *gorm.DB

func InitDB() {
	var err error
	_db, err = sql.Open("mysql", orm.DBUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer _db.Close()

	db, err = gorm.Open(mysql.Open(orm.DBUrl), &gorm.Config{})
	if err != nil {
		panic("failed to connect ain")
	}

	// 测试插入操作
	InsertStudent("张三", 20, "三年级")

	// 测试查询操作
	QueryStudents()
	UpdateStudentGrade("张三", "四年级")
	DeleteStudent()

	err = db.AutoMigrate(&Employee{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	QueryEmployeeBySalary()
	QueryEmployees()
}

func InsertStudent(name string, age int, grade string) {
	// INSERT INTO students (name, age, grade) VALUES ('张三', 20, '三年级')
	sql := "INSERT INTO students (name, age, grade) VALUES (?, ?, ?)"
	_, err := _db.Exec(sql, name, age, grade)
	if err != nil {
		fmt.Println(err)
	}
}

func QueryStudents() {
	// SELECT * FROM students WHERE age > 18
	sql := "SELECT * FROM students WHERE age > ? order by id desc limit 2"
	rows, err := _db.Query(sql, 18)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var age int
		var grade string
		err = rows.Scan(&id, &name, &age, &grade)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(id, name, age, grade)
	}
}

func UpdateStudentGrade(name string, grade string) {
	// UPDATE students SET grade = '四年级' WHERE name = '张三'
	sql := "UPDATE students SET grade = ? WHERE name = ?"
	_, err := _db.Exec(sql, grade, name)
	if err != nil {
		fmt.Println(err)
	}
}

func DeleteStudent() {
	// DELETE FROM students WHERE age < 15
	sql := "DELETE FROM students WHERE age < ?"
	_, err := _db.Exec(sql, 15)
	if err != nil {
		fmt.Println(err)
	}
}

/*
假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表
（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
要求 ：
编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，
如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
*/
func TransferMoney(from_account_id int, to_account_id int, amount int) {
	tx, err := _db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tx.Rollback()

	// 检查账户 A 的余额是否足够
	var balance int
	sql := "SELECT balance FROM accounts WHERE id = ?"
	err = tx.QueryRow(sql, from_account_id).Scan(&balance)
	if err != nil {
		fmt.Println(err)
		return
	}
	if balance < amount {
		fmt.Println("账户余额不足")
		return
	}

	// 扣除 100 元，向账户 B 增加 100 元
	sql = "UPDATE accounts SET balance = balance - ? WHERE id = ?"
	_, err = tx.Exec(sql, amount, from_account_id)
	if err != nil {
		fmt.Println(err)
		return
	}
	sql = "UPDATE accounts SET balance = balance + ? WHERE id = ?"
	_, err = tx.Exec(sql, amount, to_account_id)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 在 transactions 表中记录该笔转账信息
	sql = "INSERT INTO transactions (from_account_id, to_account_id, amount) VALUES (?, ?, ?)"
	_, err = tx.Exec(sql, from_account_id, to_account_id, amount)
	if err != nil {
		fmt.Println(err)
		return
	}

	tx.Commit()
}

/*
题目1：使用SQL扩展库进行查询
假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
要求 ：
编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
*/
type Employee struct {
	Id         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

func QueryEmployees() {
	var employees []Employee
	err := db.Model(&Employee{}).Debug().Where("department = ?", "技术部").Find(&employees)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(employees)
}
func QueryEmployeeBySalary() {
	var employee Employee
	err := db.Model(&Employee{}).Debug().Order("salary DESC").Limit(1).Find(&employee)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(employee)
}

/*
实现类型安全映射
假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
要求 ：
定义一个 Book 结构体，包含与 books 表对应的字段。
编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
*/

type Book struct {
	Id     int    `db:"id"`
	Title  string `db:"title"`
	Author string `db:"author"`
	Price  int    `db:"price"`
}

func QueryBooks() {
	var books []Book
	err := db.Model(&Book{}).Debug().Where("price > ?", 50).Find(&books)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(books)
}

/*
假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
要求 ：
使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
编写Go代码，使用Gorm创建这些模型对应的数据库表。

编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
编写Go代码，使用Gorm查询评论数量最多的文章信息。

为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
*/
type User struct {
	Id        int    `db:"id"`
	Name      string `db:"name"`
	Password  string `db:"password"`
	PostCount int    `db:"post_count"`
}

type Post struct {
	Id           int    `db:"id"`
	Title        string `db:"title"`
	Content      string `db:"content"`
	AuthorId     int    `db:"author_id"`
	CommentCount int    `db:"comment_count"`
	Status       string `db:"status"`
}

type Comment struct {
	Id       int    `db:"id"`
	Content  string `db:"content"`
	PostId   int    `db:"post_id"`
	AuthorId int    `db:"author_id"`
}

func CreateTables() {
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

func QueryUserPosts(userId int) {
	var posts []Post
	err := db.Model(&Post{}).Debug().Where("author_id = ?", userId).Find(&posts)
	if err != nil {
		fmt.Println(err)
	}
	for _, post := range posts {
		var comments []Comment
		err = db.Model(&Comment{}).Debug().Where("post_id = ?", post.Id).Find(&comments)
		if err != nil {
			fmt.Println(err)
		}
		post.CommentCount = len(comments)
	}
	fmt.Println(posts)
}

func QueryMostCommentedPost() {
	var post Post
	err := db.Model(&Post{}).Debug().Order("comment_count DESC").Limit(1).Find(&post)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(post)
}

func UpdatePostCommentCount(postId int) {
	var post Post
	db.Model(&Post{}).Debug().Where("id = ?", postId).First(&post)

	post.CommentCount = post.CommentCount - 1
	if post.CommentCount == 0 {
		post.Status = "无评论"
	}
	err := db.Save(&post).Error
	if err != nil {
		fmt.Println(err)
	}
}

func DeleteComment(commentId int) {
	var comment Comment
	db.Model(&Comment{}).Debug().Where("id = ?", commentId).First(&comment)

	err := db.Delete(&comment).Error
	if err != nil {
		fmt.Println(err)
	}
	UpdatePostCommentCount(comment.PostId)
}
