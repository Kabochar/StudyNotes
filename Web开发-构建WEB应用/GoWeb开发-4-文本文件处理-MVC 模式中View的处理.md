# 文本文件处理

文章整理自：https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/07.2.md

本文完成 MVC 模式中 V 的处理

## 7.2 JSON处理

JSON（Javascript Object Notation）是一种轻量级的数据交换语言，以文字为基础，具有自我描述性且易于让人阅读。

JSON 与 XML 最大的不同在于 XML  是一个完整的标记语言，而 JSON 不是。JSON 由于比 XML 更小、更快，更易解析,以及浏览器的内建快速解析支持，使得其更适用于网络数据传输领域。

### 解析JSON

#### 解析到结构体

使用场景：知道被解析的数据的格式

通过该函数

```
func Unmarshal(data []byte, v interface{}) error
```

解析例子

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Server struct {
	ServerName string
	ServerIP   string
}

type Serverslice struct {
	Servers []Server
}

func main() {
	var s Serverslice
	str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	json.Unmarshal([]byte(str), &s)
	fmt.Println(s)
}
```

如何将json数据与struct字段相匹配呢？例如JSON的key是`Foo`，那么怎么找对应的字段呢？

-   首先查找 tag 含有 `Foo`  的可导出的 struct 字段 (首字母大写)
-   其次查找字段名是 `Foo` 的导出字段
-   最后查找类似 `FOO` 或者 `FoO` 这样的除了首字母之外其他大小写不敏感的导出字段

PS：能够被赋值的字段必须是可导出字段(即首字母大写）。好处是：当你接收到一个很大的 JSON 数据结构而你却只想获取其中的部分数据的时候，你只需将你想要的数据对应的字段名大写，即可解决这个问题。

#### 解析到 interface

使用场景：不知道被解析的数据的格式

Go类型和JSON类型的对应关系如下：

-   bool 代表 JSON booleans
-   float64 代表 JSON numbers
-   string 代表 JSON strings
-   nil 代表 JSON null

解析到 interface{}

```
var f interface{}
err := json.Unmarshal(b, &f)
```

如何来访问这些数据呢？通过 断言 的方式。官方提供的解决方案

```go
for k, v := range m {
	switch vv := v.(type) {
	case string:
		fmt.Println(k, "is string", vv)
	case int:
		fmt.Println(k, "is int", vv)
	case float64:
		fmt.Println(k,"is float64",vv)
	case []interface{}:
		fmt.Println(k, "is an array:")
		for i, u := range vv {
			fmt.Println(i, u)
		}
	default:
		fmt.Println(k, "is of a type I don't know how to handle")
	}
}
```

bitly 公司开源了一个叫做  `simplejson` 的 包 ，在处理未知结构体的 JSON 时相当方便，详细例子如下所示：

```go
js, err := NewJson([]byte(`{
	"test": {
		"array": [1, "2", 3],
		"int": 10,
		"float": 5.150,
		"bignum": 9223372036854775807,
		"string": "simplejson",
		"bool": true
	}
}`))

arr, _ := js.Get("test").Get("array").Array()
i, _ := js.Get("test").Get("int").Int()
ms := js.Get("test").Get("string").MustString()
```

>   详细的请参考如下地址：<https://github.com/bitly/go-simplejson>

### 生成JSON

通过`Marshal`函数来处理

```
func Marshal(v interface{}) ([]byte, error)
```

下面的例子：

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Server struct {
	ServerName string
	ServerIP   string
}

type Serverslice struct {
	Servers []Server
}

func main() {
	var s Serverslice
	s.Servers = append(s.Servers, Server{ServerName: "Shanghai_VPN", ServerIP: "127.0.0.1"})
	s.Servers = append(s.Servers, Server{ServerName: "Beijing_VPN", ServerIP: "127.0.0.2"})
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Println(string(b))
}

Output:
{"Servers":[{"ServerName":"Shanghai_VPN","ServerIP":"127.0.0.1"},{"ServerName":"Beijing_VPN","ServerIP":"127.0.0.2"}]}
```

JSON 输出的时候必须注意，只有导出的字段才会被输出，如果修改字段名，那么就会发现什么都不会输出，所以必须通过 struct tag 定义来实现：

```
type Server struct {
	ServerName string `json:"serverName"`
	ServerIP   string `json:"serverIP"`
}

type Serverslice struct {
	Servers []Server `json:"servers"`
}
```

需要注意的几点是:

-   字段的tag是 `"-"` ，那么这个字段不会输出到JSON
-   tag 中带有自定义名称，那么这个自定义名称会出现在 JSON 的字段名中，例如上面例子中 serverName
-   tag中如果带有 `"omitempty"` 选项，那么如果该 字段值为空，就 不会输出到JSON串中
-   如果字段类型是 bool, string, int, int64等，而tag中带有 `",string"` 选项，那么这个字段在输出到 JSON 的时候会把该字段对应的值 转换 成 JSON字符串 

案例

```go
type Server struct {
	// ID 不会导出到JSON中
	ID int `json:"-"`

	// ServerName2 的值会进行二次JSON编码
	ServerName  string `json:"serverName"`
	ServerName2 string `json:"serverName2,string"`

	// 如果 ServerIP 为空，则不输出到JSON串中
	ServerIP   string `json:"serverIP,omitempty"`
}

s := Server {
	ID:         3,
	ServerName:  `Go "1.0" `,
	ServerName2: `Go "1.0" `,
	ServerIP:   ``,
}
b, _ := json.Marshal(s)
os.Stdout.Write(b)

Output:
{"serverName":"Go \"1.0\" ","serverName2":"\"Go \\\"1.0\\\" \""}
```

Marshal函 数只有在 转换成功 的时候 才会 返回数据，在转换的过程中我们需要注意几点：

-   JSON 对象 **只支持 string 作为 key**，所以要 编码一个 map，那么必须是`map[string]T`  这种类型  (T是Go语言中任意的类型)
-   Channel, complex 和 function 是不能被编码成 JSON 的
-   嵌套的数据 是 不能编码的，不然会让 JSON 编码进入死循环
-   指针在编码的时候会 输出 指针 指向的内容，而 空指针会输出 null

## 7.3 正则处理

更多：https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/07.3.md

## 7.4 模板处理

MVC 的设计模式，Model 处理数据，View 展现结果，Controller 控制用户的请求，至于 View 层的处理，在很多动态语言里面都是通过在 静态 HTML 中插入动态语言生成的数据。

![1555123036633](D:\Documents\笔记本\offer学习复习\Web开发\1555123036633.png)

Web应用反馈给客户端的信息中的大部分内容是静态的，不变的，而另外少部分是根据用户的请求来动态生成的，例如要显示用户的访问记录列表。

用户之间只有记录数据是不同的，而 列表的样式 则是 固定 的，此时采用模板可以复用很多静态代码。

### Go模板使用

使用 `html/template` 包来进行模板处理，使用类似 `Parse`、 `ParseFile`、 `Execute`等方法从文件或者字符串加载模板，然后执行类似上面图片展示的模板的merge 操作。

```go
func handler(w http.ResponseWriter, r *http.Request) {
	t := template.New("some template") //创建一个模板
	t, _ = t.ParseFiles("tmpl/welcome.html")  //解析模板文件
	user := GetUser() //获取当前用户信息
	t.Execute(w, user)  //执行模板的merger操作
}
```

### 模板中如何插入数据？

模板通过 `{{}}` 来包含需要在渲染时被 替换的字段，`{{.}}`表示当前的对象。

注意一点：这个字段必须是导出的(字段首字母必须是大写的)，否则，在渲染的时候 无法显示。

```go
package main

import (
	"html/template"
	"os"
)

type Person struct {
	UserName string
}

func main() {
	t := template.New("fieldname example")
	t, _ = t.Parse("hello {{.UserName}}!")
	p := Person{UserName: "Astaxie"}
	t.Execute(os.Stdout, p)
}
Output:
hello Astaxie
```

如果字段的首字母没有大写。无法输出该字段

```go
type Person struct {
	UserName string
	email	string  //未导出的字段，首字母是小写的
}

t, _ = t.Parse("hello {{.UserName}}! {{.email}}")
```

### 输出嵌套字段内容

如果字段里面还有对象，如何循环的输出这些内容呢？可以使用`{{with …}}…{{end}}`和`{{range …}}{{end}}`来进行数据的输出。

-   {{range}} 这个和 Go语法里面的 range 类似，循环操作数据
-   {{with}} 操作是 指当前对象的值，类似上下文的概念

```go
package main

import (
	"html/template"
	"os"
)

type Friend struct {
	Fname string
}

type Person struct {
	UserName string
	Emails   []string
	Friends  []*Friend
}

func main() {
	f1 := Friend{Fname: "minux.ma"}
	f2 := Friend{Fname: "xushiwei"}
	t := template.New("fieldname example")
	t, _ = t.Parse(`hello {{.UserName}}!
			{{range .Emails}}
				an email {{.}}
			{{end}}
			{{with .Friends}}
			{{range .}}
				my friend name is {{.Fname}}
			{{end}}
			{{end}}
			`)
	p := Person{UserName: "Astaxie",
		Emails:  []string{"astaxie@beego.me", "astaxie@gmail.com"},
		Friends: []*Friend{&f1, &f2}}
	t.Execute(os.Stdout, p)
}
Output:
hello Astaxie!
			
				an email astaxie@beego.me
			
				an email astaxie@gmail.com
			
			
			
				my friend name is minux.ma
			
				my friend name is xushiwei
```

### 条件处理

如果 pipeline 为空，那么if就认为是false/如何使用`if-else`语法

```go
package main

import (
	"os"
	"text/template"
)

func main() {
	tEmpty := template.New("template test")
	tEmpty = template.Must(tEmpty.Parse("空 pipeline if demo: {{if ``}} 不会输出. {{end}}\n"))
	tEmpty.Execute(os.Stdout, nil)

	tWithValue := template.New("template test")
	tWithValue = template.Must(tWithValue.Parse("不为空的 pipeline if demo: {{if `anything`}} 我有内容，我会输出. {{end}}\n"))
	tWithValue.Execute(os.Stdout, nil)

	tIfElse := template.New("template test")
	tIfElse = template.Must(tIfElse.Parse("if-else demo: {{if `anything`}} if部分 {{else}} else部分.{{end}}\n"))
	tIfElse.Execute(os.Stdout, nil)
}
```

>   注意：if里面无法使用条件判断，例如.Mail=="[astaxie@gmail.com](mailto:astaxie@gmail.com)"，这样的判断是不正确的，if里面只能是bool值

### pipelines

Go语言模板最强大的一点就是支持 pipe 数据，在Go语言里面任何`{{}}`里面的都是pipelines数据。（`ls | grep "beego"`类似这样的语法）

输出的email里面如果还有一些可能引起XSS注入的，那么我们如何来进行转化呢？

```
{{. | html}}
```

在email输出的地方我们可以采用如上方式可以把输出全部转化html的实体

### 模板变量

Go语言通过申明的局部变量格式如下所示

```go
$variable := pipeline
```

详细的例子

```go
{{with $x := "output" | printf "%q"}}{{$x}}{{end}}
{{with $x := "output"}}{{printf "%q" $x}}{{end}}
{{with $x := "output"}}{{$x | printf "%q"}}{{end}}
```

### 模板函数

模板在输出对象的字段值时，采用了`fmt`包把对象转化成了字符串。

案例背景：有时候我们为了防止垃圾邮件发送者通过采集网页的方式来发送给我们的邮箱信息，我们希望把 `@` 替换成 `at` 例如：`astaxie at beego.me`

每一个模板函数都有一个唯一值的名字，然后与一个Go函数关联，通过如下的方式来关联。

```go
type FuncMap map[string]interface{}
```

如果我们想要的email函数的模板函数名是 `emailDeal`，它关联的Go函数名称是 `EmailDealWith` ，那么我们可以通过下面的方式来注册这个函数

```go
t = t.Funcs(template.FuncMap{"emailDeal": EmailDealWith})
```

`EmailDealWith`  参数和返回值定义如下：

```go
func EmailDealWith(args …interface{}) string
```

实现例子：

```go
package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

type Friend struct {
	Fname string
}

type Person struct {
	UserName string
	Emails   []string
	Friends  []*Friend
}

// 把 @ 符号 转换成 at 显示
func EmailDealWith(args ...interface{}) string {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}
	// find the @ symbol
	substrs := strings.Split(s, "@")
	if len(substrs) != 2 {
		return s
	}
	// replace the @ by " at "
	return (substrs[0] + " at " + substrs[1])
}

func main() {
	f1 := Friend{Fname: "minux.ma"}
	f2 := Friend{Fname: "xushiwei"}
	t := template.New("fieldname example")
	t = t.Funcs(template.FuncMap{"emailDeal": EmailDealWith})
	t, _ = t.Parse(`hello {{.UserName}}!
				{{range .Emails}}
					an emails {{.|emailDeal}}
				{{end}}
				{{with .Friends}}
				{{range .}}
					my friend name is {{.Fname}}
				{{end}}
				{{end}}
				`)
	p := Person{UserName: "Astaxie",
		Emails:  []string{"astaxie@beego.me", "astaxie@gmail.com"},
		Friends: []*Friend{&f1, &f2}}
	t.Execute(os.Stdout, p)
}
```

在模板包内部已经有内置的实现函数，下面代码截取自模板包里面

```go
var builtins = FuncMap{
	"and":      and,
	"call":     call,
	"html":     HTMLEscaper,
	"index":    index,
	"js":       JSEscaper,
	"len":      length,
	"not":      not,
	"or":       or,
	"print":    fmt.Sprint,
	"printf":   fmt.Sprintf,
	"println":  fmt.Sprintln,
	"urlquery": URLQueryEscaper,
}
```

### Must操作

作用：检测模板是否正确，例如大括号是否匹配，注释是否正确的关闭，变量是否正确的书写。

```go
package main

import (
	"fmt"
	"text/template"
)

func main() {
	tOk := template.New("first")
	template.Must(tOk.Parse(" some static text /* and a comment */"))
	fmt.Println("The first one parsed OK.")

	template.Must(template.New("second").Parse("some static text {{ .Name }}"))
	fmt.Println("The second one parsed OK.")

	fmt.Println("The next one ought to fail.")
	tErr := template.New("check parse error with Must")
	template.Must(tErr.Parse(" some static text {{ .Name }"))
}
Output：
The first one parsed OK.
The second one parsed OK.
The next one ought to fail.
panic: template: check parse error with Must:1: unexpected "}" in command
```

## 嵌套模板

出现背景：开发Web应用的时候，经常会遇到一些模板有些部分是固定不变的，然后可以抽取出来作为一个独立的部分。我们可以定义成 `header`、 `content`、 `footer`三个部分。Go语言中通过如下的语法来申明

通过如下的语法来 申明

```
{{define "子模板名称"}}
	内容
{{end}}
```

通过如下方式来 调用：

```
{{template "子模板名称"}}
```

定义三个文件，`header.tmpl`、`content.tmpl`、`footer.tmpl`文件，里面的内容如下

```html
//header.tmpl
{{define "header"}}
<html>
    <head>
        <title>演示信息</title>
    </head>
<body>
{{end}}

//content.tmpl
{{define "content"}}
{{template "header"}}
    <h1>演示嵌套</h1>
    <ul>
        <li>嵌套使用define定义子模板</li>
        <li>调用使用template</li>
    </ul>
        {{template "footer"}}
{{end}}

//footer.tmpl
{{define "footer"}}
	</body>
</html>
{{end}}
```

演示代码如下：

```go
package main

import (
	"fmt"
	"os"
	"text/template"
)

func main() {
	s1, _ := template.ParseFiles("header.tmpl", "content.tmpl", "footer.tmpl")
	s1.ExecuteTemplate(os.Stdout, "header", nil)
	fmt.Println()
	s1.ExecuteTemplate(os.Stdout, "content", nil)
	fmt.Println()
	s1.ExecuteTemplate(os.Stdout, "footer", nil)
	fmt.Println()
	s1.Execute(os.Stdout, nil)
}
```

通过 `template.ParseFiles` 把所有的嵌套模板全部解析到模板里面，其实每一个定义的  {{define}}  都是一个独立的模板，他们相互独立，是并行存在的关系，内部其实存储的是类似 map 的一种关系 (key 是 模板的名称，value 是 模板的内容)。

然后我们通过 `ExecuteTemplate` 来执行相应的子模板内容，我们可以看到 header、footer 都是相对独立的，都能输出内容，content 中因为嵌套了 header 和 footer 的内容，就会同时输出三个的内容。

但是当我们 执行`s1.Execute`，没有任何的输出，因为在默认的情况下没有默认的子模板，所以不会输出任何的东西。

>   同一个集合类的模板是互相知晓的，如果同一模板被多个集合使用，则它需要在多个集合中分别解析

## 7.5 文件操作

### 目录操作

-   func Mkdir(name string, perm FileMode) error

    创建名称为name的目录，权限设置是perm，例如0777

-   func MkdirAll(path string, perm FileMode) error

    根据path创建多级子目录，例如astaxie/test1/test2。

-   func Remove(name string) error

    删除名称为name的目录，当目录下有文件或者其他目录时会出错

-   func RemoveAll(path string) error

    根据path删除多级子目录，如果path是单个名称，那么该目录下的子目录全部删除。

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	os.Mkdir("astaxie", 0777)
	os.MkdirAll("astaxie/test1/test2", 0777)
	err := os.Remove("astaxie")
	if err != nil {
		fmt.Println(err)
	}
	os.RemoveAll("astaxie")
}
```

### 文件操作

#### 建立与打开文件

-   func Create(name string) (file *File, err Error)

    根据提供的文件名创建新的文件，返回一个文件对象，默认权限是0666的文件，返回的文件对象是可读写的。

-   func NewFile(fd uintptr, name string) *File

    根据文件描述符创建相应的文件，返回一个文件对象

通过如下两个方法来打开文件：

-   func Open(name string) (file *File, err Error)

    该方法打开一个名称为name的文件，但是是只读方式，内部实现其实调用了OpenFile。

-   func OpenFile(name string, flag int, perm uint32) (file *File, err Error)

    打开名称为name的文件，flag是打开的方式，只读、读写等，perm是权限

### 写文件

-   func (file *File) Write(b []byte) (n int, err Error)

    写入byte类型的信息到文件

-   func (file *File) WriteAt(b []byte, off int64) (n int, err Error)

    在指定位置开始写入byte类型的信息

-   func (file *File) WriteString(s string) (ret int, err Error)

    写入string信息到文件

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	userFile := "astaxie.txt"
	fout, err := os.Create(userFile)		
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	defer fout.Close()
	for i := 0; i < 10; i++ {
		fout.WriteString("Just a test!\r\n")
		fout.Write([]byte("Just a test!\r\n"))
	}
}
```

### 读文件

读文件函数：

-   func (file *File) Read(b []byte) (n int, err Error)

    读取数据到b中

-   func (file *File) ReadAt(b []byte, off int64) (n int, err Error)

    从off开始读取数据到b中

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	userFile := "asatxie.txt"
	fl, err := os.Open(userFile)		
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	defer fl.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := fl.Read(buf)
		if 0 == n {
			break
		}
		os.Stdout.Write(buf[:n])
	}
}
```

### 删除文件

Go语言里面删除文件和删除文件夹是同一个函数

-   func Remove(name string) Error

    调用该函数就可以删除文件名为name的文件

# 7.6 字符串处理

更多：https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/07.6.md