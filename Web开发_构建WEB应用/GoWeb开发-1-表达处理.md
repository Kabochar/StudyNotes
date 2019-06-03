# 表单处理

文章整理自：https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/04.0.md

## 4.2 验证表单处理

开发Web的一个原则就是，不能信任用户输入的任何信息。这其中存在一定恶意行为。

验证信息有两个方式：前端 JS 验证，服务端验证

### 通用正则检测

匹配数字

```go 
if m, _ := regexp.MatchString("^[0-9]+$", r.Form.Get("age")); !m {
	return false
}
```

### 4.3 预防跨站脚本

>   现在的网站包含大量的动态内容以提高用户体验，比过去要复杂得多。所谓动态内容，就是根据用户环境和需要，Web应用程序能够输出相应的内容。动态站点会受到一种名为“跨站脚本攻击”（Cross Site Scripting, 安全专家们通常将其缩写成 XSS）的威胁，而静态站点则完全不受其影响。

对 XSS 最佳的防护：一是 验证所有输入数据，有效检测攻击（正则表达式验证）。一个是 对所有输出数据进行适当的处理，以防止任何已成功注入的脚本在浏览器端运行。

Go 里面是怎么做这个有效防护的呢？Go 的 `html/template` 里面带有下面几个函数可以帮你转义。

-   func HTMLEscape(w io.Writer, b []byte) // 把 b 进行转义之后写到 w
-   func HTMLEscapeString(s string) string // 转义 s 之后返回结果字符串
-   func HTMLEscaper(args ...interface{}) string // 支持多个参数一起转义，返回结果字符串

案例：

```go
import "text/template"
...
t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
err = t.ExecuteTemplate(out, "T", "<script>alert('you have been pwned')</script>")
Output:
Hello, <script>alert('you have been pwned')</script>!


import "html/template"
...
t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
err = t.ExecuteTemplate(out, "T", template.HTML("<script>alert('you have been pwned')</script>"))
Output:
Hello, <script>alert('you have been pwned')</script>!
```

## 4.4 防止多次递交表单

如何有效的防止用户多次递交相同的表单？

解决方案：在表单中添加一个带有唯一值的隐藏字段。在验证表单时，先检查带有该唯一值的表单是否已经递交过了。如果是，拒绝再次递交；如果不是，则处理表单进行逻辑处理。另外，如果是采用了Ajax模式递交表单的话，当表单递交后，通过javascript来禁用表单的递交按钮。

```html
<input type="hidden" name="token" value="{{.}}">
<input type="submit" value="登陆">
```

增加了一个隐藏字段`token`，这个值我们通过MD5(时间戳)来获取唯一值，然后我们把这个值存储到服务器端 (session来控制)。

>   更多：https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/04.4.md

### 4.5 处理文件上传

添加form的`enctype`属性

```
application/x-www-form-urlencoded   表示在发送前编码所有字符（默认）
multipart/form-data	  不对字符编码。在使用包含文件上传控件的表单时，必须使用该值。
text/plain	  空格转换为 "+" 加号，但不对特殊字符编码。
```

案例

```html
<form enctype="multipart/form-data" action="/upload" method="post">
    <input type="file" name="uploadfile" />
    <input type="hidden" name="token" value="{{.}}"/>
    <input type="submit" value="upload" />
</form>
```

#### 服务端操作

```go
http.HandleFunc("/upload", upload)

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix() // 防止多次递交表单
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20) // ！！！
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)  // 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
```

处理文件上传我们需要调用 `r.ParseMultipartForm` ，里面的参数表示 `maxMemory` ，调用 `ParseMultipartForm` 之后，上传的文件存储在 `maxMemory` 大小的内存里面，如果文件大小超过了 `maxMemory` ，那么剩下的部分将存储在系统的临时文件中。

我们可以通过 `r.FormFile` 获取上面的文件句柄，然后实例中使用了 `io.Copy` 来存储文件。

>   获取其他非文件字段信息的时，不需要调用`r.ParseForm`。因为，在需要的时候Go自动会去调用。
>
>   而且`ParseMultipartForm`调用一次之后，后面再次调用不会再有效果。

上传文件主要三步处理：

1.  表单中增加 enctype="multipart/form-data"
2.  服务端调用 `r.ParseMultipartForm`，把上传的文件存储在内存和临时文件中
3.  使用 `r.FormFile` 获取文件句柄，然后对文件进行存储等处理。

#### 客户端操作

```go
package main
import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func postFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// 关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// 打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()
	
	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}

// sample usage
func main() {
	target_url := "http://localhost:9090/upload"
	filename := "./astaxie.pdf"
	postFile(filename, target_url)
}
```

客户端通过 multipart.Write 把文件的文本流写入一个缓存中，然后调用http的Post方法把缓存传到服务器。

