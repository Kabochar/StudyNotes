# 数据库

文章整理自：https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/05.0.md

[TOC]

## 5.1 database/sql接口

Go官方没有提供数据库驱动，而是为开发数据库驱动定义了一些标准接口，开发者可以根据定义的接口来开发相应的数据库驱动。这样做有一个好处，只要是按照标准接口开发的代码， 以后需要迁移数据库时，不需要任何修改。

### sql.Register

用来注册数据库驱动的，当第三方开发者开发数据库驱动时，都会实现init函数，在init里面会调用这个`Register(name string, driver driver.Driver)`完成本驱动的注册。

```go
//https://github.com/mattn/go-sqlite3驱动
func init() {
	sql.Register("sqlite3", &SQLiteDriver{})
}

//https://github.com/mikespook/mymysql驱动
// Driver automatically registered in database/sql
var d = Driver{proto: "tcp", raddr: "127.0.0.1:3306"}
func init() {
	Register("SET NAMES utf8")
	sql.Register("mymysql", &d)
}
```

在 `database/sql` 内部通过一个 map 来存储用户定义的相应驱动。

```go
var drivers = make(map[string]driver.Driver)

drivers[name] = driver
```

通过 database/sql 的注册函数可以同时注册多个数据库驱动，只要不重复。

### driver.Driver

Driver 是一个数据库驱动的接口。Open(name string)，这个方法返回一个数据库的Conn接口。

```go
type Driver interface {
	Open(name string) (Conn, error)
}
```

返回的Conn只能用来进行一次goroutine的操作。多个 goroutine 执行不同操作，Go不知道某个操作究竟是由哪个goroutine发起的,从而导致数据混乱。

### driver.Conn

Conn 是一个数据库连接的接口定义。这个 Conn 只能应用在一个 goroutine 里面

```go
type Conn interface {
	Prepare(query string) (Stmt, error)
	Close() error
	Begin() (Tx, error)
}
```

>   更多
>
>   https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/05.1.md

## 5.2 使用MySQL数据库

相关库 ：[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql ) 

基本操作

```go
package main

import (
	"database/sql"
	"fmt"
	//"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "astaxie:astaxie@/test?charset=utf8")
	checkErr(err)

	//插入数据
	stmt, err := db.Prepare("INSERT INTO userinfo SET username=?,department=?,created=?")
	checkErr(err)

	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
	//更新数据
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("astaxieupdate", id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	//查询数据
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	//删除数据
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
```

go-sql-driver 中注册 mysql 这个数据库驱动。

第二个参数是DSN(Data Source Name)，它是go-sql-driver定义的一些数据库链接和配置信息。它支持如下格式：

```
user@unix(/path/to/socket)/dbname?charset=utf8
user:password@tcp(localhost:5555)/dbname?charset=utf8
user:password@/dbname
user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname
```

`db.Prepare()` 用来返回准备要执行的 sql 操作，然后返回准备完毕的执行状态。

`db.Query()` 用来直接执行 Sql 返回Rows结果。

`stmt.Exec()` 用来执行 stmt 准备好的 SQL 语句

我们可以看到我们传入的参数都是=?对应的数据，这样做的方式可以一定程度上防止 SQL 注入。

## 5.5 使用Beego orm库进行ORM开发

### 安装

```
go get github.com/astaxie/beego
```

### 初始化

```go
import (
	"database/sql"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// 注册驱动
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	// 设置默认数据库
	orm.RegisterDataBase("default", "mysql", "root:root@/my_db?charset=utf8", 30)
	// 注册定义的model
    orm.RegisterModel(new(User))

   	// 创建table
    orm.RunSyncdb("default", false, true)
}
```

### 简单示例

```go
package main

import (
    "fmt"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql" // 导入数据库驱动
)

// Model Struct
type User struct {
    Id   int
    Name string `orm:"size(100)"`
}

func init() {
    // 设置默认数据库
    orm.RegisterDataBase("default", "mysql", "root:root@/my_db?charset=utf8", 30)
    
    // 注册定义的 model
    orm.RegisterModel(new(User))
//RegisterModel 也可以同时注册多个 model
//orm.RegisterModel(new(User), new(Profile), new(Post))

    // 创建 table
    orm.RunSyncdb("default", false, true)
}

func main() {
    // 创建一个beego orm对象
    o := orm.NewOrm()

    user := User{Name: "slene"}

    // 插入表
    id, err := o.Insert(&user)
    fmt.Printf("ID: %d, ERR: %v\n", id, err)

    // 更新表
    user.Name = "astaxie"
    num, err := o.Update(&user)
    fmt.Printf("NUM: %d, ERR: %v\n", num, err)

    // 读取 one
    u := User{Id: user.Id}
    err = o.Read(&u)
    fmt.Printf("ERR: %v\n", err)

    // 删除表
    num, err = o.Delete(&u)
    fmt.Printf("NUM: %d, ERR: %v\n", num, err)
}
```

#### SetMaxIdleConns

根据数据库的别名，设置数据库的最大空闲连接

```go
orm.SetMaxIdleConns("default", 30)
```

#### SetMaxOpenConns

根据数据库的别名，设置数据库的最大数据库连接 (goversion >= 1.2)

```go
orm.SetMaxOpenConns("default", 30)
```

#### Debug

beego orm支持打印调试，通过如下的代码实现调试

```go
 orm.Debug = true
```

### 插入数据

操作的是struct对象，而不是原生的sql语句，最后通过调用Insert接口将数据保存到数据库。

```go
o := orm.NewOrm()
var user User
user.Name = "zxxx"
user.Departname = "zxxx"

id, err := o.Insert(&user)
if err == nil {
	fmt.Println(id)
}
```

### 更新数据

user 的 主键 已经有值了，此时调用 Insert 接口，beego orm 内部会自动调用 update以进行数据的更新而非插入操作。

```go
o := orm.NewOrm()
user := User{Uid: 1}
if o.Read(&user) == nil {
	user.Name = "MyName"
	if num, err := o.Update(&user); err == nil {
		fmt.Println(num)
	}
}
```

Update 默认更新所有的字段，可以更新指定的字段：

```go
// 只更新 Name
o.Update(&user, "Name")
// 指定多个字段
// o.Update(&user, "Field1", "Field2", ...)

// Where:用来设置条件，支持多个参数，第一个参数如果为整数，相当于调用了Where("主键=?",值)。
```

### 查询数据

例子1，根据主键获取数据：

```go
o := orm.NewOrm()
var user User

user := User{Id: 1}

err = o.Read(&user)

if err == orm.ErrNoRows {
	fmt.Println("查询不到")
} else if err == orm.ErrMissPK {
	fmt.Println("找不到主键")
} else {
	fmt.Println(user.Id, user.Name)
}
```

例子2：

```go
o := orm.NewOrm()
var user User

qs := o.QueryTable(user) // 返回 QuerySeter
qs.Filter("id", 1) // WHERE id = 1
qs.Filter("profile__age", 18) // WHERE profile.age = 18
```

例子3，WHERE IN 查询条件：

```go
qs.Filter("profile__age__in", 18, 20) 
// WHERE profile.age IN (18, 20)
```

例子4，更加复杂的条件：

```go
qs.Filter("profile__age__in", 18, 20).Exclude("profile__lt", 1000)
// WHERE profile.age IN (18, 20) AND NOT profile_id < 1000
```

可以通过如下接口获取多条数据，请看示例

获取多条数据，请看示例

例子1，根据条件age>17，获取20位置开始的10条数据的数据

```go
var allusers []User
qs.Filter("profile__age__gt", 17)
// WHERE profile.age > 17
```

例子2，limit默认从10开始，获取10条数据

```go
qs.Limit(10, 20)
// LIMIT 10 OFFSET 20 注意跟SQL反过来的
```

### 删除数据

例子1，删除单条数据

```go
o := orm.NewOrm()
if num, err := o.Delete(&User{Id: 1}); err == nil {
	fmt.Println(num)
}
```

Delete 操作会对反向关系进行操作，此例中 Post 拥有一个到 User 的外键。删除 User 的时候。如果 on_delete 设置为默认的级联操作，将删除对应的 Post

### 关联查询

```go
type Post struct {
	Id    int    `orm:"auto"`
	Title string `orm:"size(100)"`
	User  *User  `orm:"rel(fk)"`
}

var posts []*Post
qs := o.QueryTable("post")
num, err := qs.Filter("User__Name", "slene").All(&posts)
```

### Group By和Having

```go
qs.OrderBy("id", "-profile__age")
// ORDER BY id ASC, profile.age DESC

qs.OrderBy("-profile__age", "profile")
// ORDER BY profile.age DESC, profile_id ASC
```

GroupBy ：用来指定进行groupby的字段

Having ：用来指定having执行的时候的条件

### 使用原生sql

简单示例：

```
o := orm.NewOrm()
var r orm.RawSeter
r = o.Raw("UPDATE user SET name = ? WHERE name = ?", "testing", "slene")
```

复杂原生sql使用：

```go
func (m *User) Query(name string) user []User {
	var o orm.Ormer
	var rs orm.RawSeter
	o = orm.NewOrm()
	rs = o.Raw("SELECT * FROM user "+
		"WHERE name=? AND uid>10 "+
		"ORDER BY uid DESC "+
		"LIMIT 100", name)
	//var user []User
	num, err := rs.QueryRows(&user)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(num)
		//return user
	}
	return
}	
```

## 5.6 NOSQL数据库操作

>   更多：https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/05.6.md

### redis

### mongoDB

