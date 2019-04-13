# Web服务

内容整理自：https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/08.1.md

[TOC]

## 8.1 Socket编程

我们每天打开浏览器浏览网页时，浏览器进程怎么和Web服务器进行通信的呢？

当你用QQ聊天时，QQ进程怎么和服务器或者是你的好友所在的QQ进程进行通信的呢？

当你打开PPstream观看视频时，PPstream进程如何与视频服务器进行通信的呢？ 

如此种种，都是靠Socket来进行通信的。

### 什么是Socket？

Socket起源于Unix，而Unix基本哲学之一就是“一切皆文件”，都可以用“打开open –> 读写write/read –> 关闭close”模式来操作。

Socket就是该模式的一个实现，网络的Socket数据传输是一种特殊的I/O，Socket也是一种文件描述符。

常用的Socket类型有两种：流式Socket（SOCK_STREAM）和数据报式Socket（SOCK_DGRAM）。

流式是一种面向连接的Socket，针对于面向连接的TCP服务应用；数据报式Socket是一种无连接的Socket，对应于无连接的UDP服务应用。

### Socket如何通信

络中的进程之间如何通过Socket通信呢？首要解决的问题：如何唯一标识一个进程，否则通信无从谈起！

网络层的 “ip地址” 可以唯一标识网络中的主机，而传输层的 “协议+端口” 可以唯一标识主机中的应用程序（进程）。

利用三元组（ip地址，协议，端口）就可以标识网络的进程了，网络中需要互相通信的进程，就可以利用这个标志在他们之间进行交互。

![1555126790145](D:\Documents\笔记本\offer学习复习\Web开发\1555126790145.png)

使用TCP/IP协议的应用程序通常采用应用编程接口：UNIX BSD的套接字（socket）和UNIX System V的TLI（已经被淘汰），来实现网络进程之间的通信。就目前而言，几乎所有的应用程序都是采用socket。

### Socket基础知识

Socket有两种：TCP Socket和UDP Socket，TCP和UDP是协议。

要确定一个进程的需要三元组，需要 IP地址 和 端口 。

#### IPv4地址

IPv4的地址位数为 32位

#### IPv6地址

为了解决IPv4在实施过程中遇到的各种问题而被提出的，IPv6采用 128位 地址长度，几乎可以不受限制地提供地址。

### Go支持的IP类型

Go 中 IP 的定义

```
type IP []byte
```

`ParseIP(s string) IP`函数会把一个IPv4或者IPv6的地址转化成IP类型

```go
package main
import (
	"net"
	"os"
	"fmt"
)
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}
	name := os.Args[1]
	addr := net.ParseIP(name)
	if addr == nil {
		fmt.Println("Invalid address")
	} else {
		fmt.Println("The address is ", addr.String())
	}
	os.Exit(0)
}
Input：
IPAddress
Ouput：
相应的IP格式
```

### 通过网络端口访问一个服务时，我们能够做什么呢？

作为客户端来说，我们可以通过向远端某台机器的的某个网络端口发送一个请求，然后得到在机器的此端口上监听的服务反馈的信息。

作为服务端，我们需要把服务绑定到某个指定端口，并且在此端口上监听，当有客户端来访问时能够读取信息并且写入反馈信息。

### TCP Socket

在Go语言的`net`包中有一个类型`TCPConn`，这个类型可以用来作为客户端和服务器端交互的通道，他有两个主要的函数：

```go
func (c *TCPConn) Write(b []byte) (int, error)
func (c *TCPConn) Read(b []byte) (int, error)
```

`TCPConn `可以用在客户端和服务器端来读写数据

`TCPAddr`类型，他表示一个TCP的地址信息，他的定义如下：

```go
type TCPAddr struct {
	IP IP
	Port int
	Zone string // IPv6 scoped addressing zone
}
```

在Go语言中通过 `ResolveTCPAddr` 获取一个`TCPAddr`

```go
func ResolveTCPAddr(net, addr string) (*TCPAddr, os.Error)
```

-   net参数是 "tcp4"、"tcp6"、"tcp" 中的任意一个，分别表示 TCP(IPv4-only), TCP(IPv6-only) 或者TCP(IPv4, IPv6的任意一个)。 
-   addr表示域名或者IP地址，例如"[www.google.com:80](http://www.google.com/)" 或者"127.0.0.1:22"。

#### TCP client

Go语言中通过 net 包中的 `DialTCP` 函数来建立一个 TCP 连接，并返回一个 `TCPConn` 类型的对象，当连接建立时服务器端也创建一个同类型的对象，此时客户端和服务器端通过各自拥有的 `TCPConn` 对象来进行数据交换。

一般而言，客户端通过`TCPConn`对象将请求信息发送到服务器端，读取服务器端响应的信息。服务器端读取并解析来自客户端的请求，并返回应答信息，这个连接只有当任一端关闭了连接之后才失效，不然这连接可以一直在使用。

建立连接的函数定义如下：

```go
func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)
```

-   network参数是"tcp4"、"tcp6"、"tcp"中的任意一个，分别表示TCP(IPv4-only)、TCP(IPv6-only)或者TCP(IPv4,IPv6的任意一个)
-   laddr表示本机地址，一般设置为 nil
-   raddr表示远程的服务地址

```go
package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	checkError(err)
	result, err := ioutil.ReadAll(conn)
	checkError(err)
	fmt.Println(string(result))
	os.Exit(0)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
```

首先，程序将用户的输入作为参数`service`传入`net.ResolveTCPAddr`获取一个tcpAddr，然后，把tcpAddr 传入 DialTCP 后创建了一个 TCP 连接 `conn`，通过 `conn` 来发送请求信息，最后通过`ioutil.ReadAll` 从 `conn`  中读取全部的文本，也就是服务端响应反馈的信息。

#### TCP server

在服务器端，我们需要绑定服务到指定的非激活端口，并监听此端口，当有客户端请求到达的时候可以接收到来自客户端连接的请求。

net包中有相应功能的函数，函数定义如下：

```
func ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error)
func (l *TCPListener) Accept() (Conn, error)
```

参数说明同DialTCP的参数一样。下面我们实现一个简单的时间同步服务，监听7777端口

```go
package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	service := ":7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		daytime := time.Now().String()
		conn.Write([]byte(daytime)) // don't care about return value
		conn.Close()                // we're finished with this client
	}
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
```

服务跑起来之后，它将会一直在那里等待，直到有新的客户端请求到达。当有新的客户端请求到达并同意接受`Accept`该请求的时候他会反馈当前的时间信息。

值得注意的是，在代码中 `for` 循环里，当有错误发生时，直接continue而不是退出，是因为在服务器端跑代码的时候，当有错误发生的情况下最好是由服务端记录错误，然后当前连接的客户端直接报错而退出，从而不会影响到当前服务端运行的整个服务。

有个缺点，执行的时候是单任务的，不能同时接收多个请求，改造以使它支持多并发

```go
package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn) // start goroutine handle server
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	daytime := time.Now().String()
	conn.Write([]byte(daytime)) // don't care about return value
	// we're finished with this client
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
```

如果我们需要通过从客户端发送不同的请求来获取不同的时间格式，而且需要一个长连接，该怎么做呢？

```go
package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
	"strings"
)

func main() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
	request := make([]byte, 128) // set maxium request length to 128B to prevent flood attack
	defer conn.Close()  // close connection before exit
	for {
		read_len, err := conn.Read(request)

		if err != nil {
			fmt.Println(err)
			break
		}

    		if read_len == 0 {
    			break // connection already closed by client
    		} else if strings.TrimSpace(string(request[:read_len])) == "timestamp" {
    			daytime := strconv.FormatInt(time.Now().Unix(), 10)
    			conn.Write([]byte(daytime))
    		} else {
    			daytime := time.Now().String()
    			conn.Write([]byte(daytime))
    		}

    		request = make([]byte, 128) // clear last read content
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
```

使用`conn.Read()`不断读取客户端发来的请求。由于我们需要保持与客户端的长连接，所以不能在读取完一次请求后就关闭连接。

由于`conn.SetReadDeadline()`设置了超时，当一定时间内客户端无请求发送，`conn`便会自动关闭，下面的for循环即会因为连接已关闭而跳出。

需要注意的是，`request`在创建时需要指定一个最大长度以防止 flood attack ；每次读取到请求处理完毕后，需要清理request，因为`conn.Read()`会将新读取到的内容 append 到原内容之后。

### 控制TCP连接

平常用到比较多 TCP有很多连接控制函数

```go
func DialTimeout(net, addr string, timeout time.Duration) (Conn, error)
```

设置建立连接的超时时间，客户端和服务器端都适用，当超过设置时间时，连接自动关闭。

```go
func (c *TCPConn) SetReadDeadline(t time.Time) error
func (c *TCPConn) SetWriteDeadline(t time.Time) error
```

设置写入/读取一个连接的超时时间。当超过设置时间时，连接自动关闭。

```go
func (c *TCPConn) SetKeepAlive(keepalive bool) os.Error
```

设置keepAlive属性，是操作系统层在tcp上没有数据和ACK的时候，会间隔性的发送keepalive包，操作系统可以通过该包来判断一个tcp连接是否已经断开。案例：windows上默认2个小时没有收到数据和keepalive包的时候人为tcp连接已经断开，

### UDP Socket

处理UDP Socket和TCP Socket不同的地方就是在服务器端处理多个客户端请求数据包的方式不同，UDP缺少了对客户端连接请求的Accept函数。UDP的几个主要函数如下所示：

```go
func ResolveUDPAddr(net, addr string) (*UDPAddr, os.Error)
func DialUDP(net string, laddr, raddr *UDPAddr) (c *UDPConn, err os.Error)
func ListenUDP(net string, laddr *UDPAddr) (c *UDPConn, err os.Error)
func (c *UDPConn) ReadFromUDP(b []byte) (n int, addr *UDPAddr, err os.Error)
func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (n int, err os.Error)
```

一个UDP的客户端代码

```go
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)
	_, err = conn.Write([]byte("anything"))
	checkError(err)
	var buf [512]byte
	n, err := conn.Read(buf[0:])
	checkError(err)
	fmt.Println(string(buf[0:n]))
	os.Exit(0)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s", err.Error())
		os.Exit(1)
	}
}
```

UDP服务器端如何来处理：

```go
package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	for {
		handleClient(conn)
	}
}
func handleClient(conn *net.UDPConn) {
	var buf [512]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	daytime := time.Now().String()
	conn.WriteToUDP([]byte(daytime), addr)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s", err.Error())
		os.Exit(1)
	}
}
```

## 8.2 WebSocket

WebSocket 它实现了 基于 浏览器的 远程 socket，它使浏览器和服务器可以进行  全双工通信 ，许多浏览器（Firefox、Google Chrome和Safari）都已对此做了支持。

在 WebSocket 出现之前，为了实现即时通信，采用的技术都是 “轮询” 。

即，在特定的时间间隔内，由浏览器对服务器发出HTTP Request，服务器在收到请求后，返回最新的数据给浏览器刷新，“轮询”使得浏览器需要对服务器不断发出请求，这样会占用大量带宽。

WebSocket采用了一些特殊的报头，使得浏览器和服务器只需要做一个握手的动作，就可以在浏览器和服务器之间建立一条连接通道。且此连接会保持在活动状态。它解决了Web实时化的问题，相比传统HTTP有如下好处：

-   一个Web客户端只建立一个TCP连接
-   Websocket服务端可以推送(push)数据到web客户端.
-   有更加轻量级的头，减少数据传送量

WebSocket URL的起始输入是ws://或是wss://（在SSL上）。

下图展示了 WebSocket 的通信过程，一个带有特定报头的HTTP握手被发送到了服务器端，接着在服务器端或是客户端就可以通过 JavaScript 来使用某种 套接口（socket），这一套接口可被用来通过事件句柄异步地接收数据。

![1555128584892](D:\Documents\笔记本\offer学习复习\Web开发\1555128584892.png)

### WebSocket原理

在第一次handshake通过以后，连接便建立成功，其后的通讯数据都是以”\x00″开头，以”\xFF”结尾。

在客户端，这个是透明的，WebSocket组件会自动将原始数据“掐头去尾”。

浏览器发出WebSocket连接请求，然后服务器发出回应，然后连接建立成功，这个过程通常称为“握手” (handshaking)。

![1555131254044](D:\Documents\笔记本\offer学习复习\Web开发\1555131254044.png)

在请求中的 "Sec-WebSocket-Key" 是随机的，这个是一个经过base64编码后的数据。服务器端接收到这个请求之后需要把这个字符串连接上一个固定的字符串：

```
258EAFA5-E914-47DA-95CA-C5AB0DC85B11
```

即：`f7cb4ezEAl6C3wRaU6JORA==` 连接上那一串固定字符串，生成一个这样的字符串：

```
f7cb4ezEAl6C3wRaU6JORA==258EAFA5-E914-47DA-95CA-C5AB0DC85B11
```

对该字符串先用 sha1安全散列算法计算出二进制的值，然后用base64对其进行编码，即可以得到握手后的字符串：

```
rE91AJhfC+6JdVcVXOGJEADEJdQ=
```

将之作为响应头`Sec-WebSocket-Accept`的值反馈给客户端。

### Go实现WebSocket

实现一个简单的例子：用户输入信息，客户端通过 WebSocket 将信息发送给服务器端，服务器端收到信息之后主动 Push 信息到客户端，然后客户端将输出其收到的信息，客户端的代码如下：

```go
<html>
<head></head>
<body>
	<script type="text/javascript">
		var sock = null;
		var wsuri = "ws://127.0.0.1:1234";

		window.onload = function() {

			console.log("onload");

			sock = new WebSocket(wsuri);

			sock.onopen = function() {
				console.log("connected to " + wsuri);
			}

			sock.onclose = function(e) {
				console.log("connection closed (" + e.code + ")");
			}

			sock.onmessage = function(e) {
				console.log("message received: " + e.data);
			}
		};

		function send() {
			var msg = document.getElementById('message').value;
			sock.send(msg);
		};
	</script>
	<h1>WebSocket Echo Test</h1>
	<form>
		<p>
			Message: <input id="message" type="text" value="Hello, world!">
		</p>
	</form>
	<button onclick="send();">Send Message</button>
</body>
</html>
```

客户端JS，很容易的就通过WebSocket函数建立了一个与服务器的连接sock，当握手成功后，会触发WebScoket对象的onopen事件，告诉客户端连接已经成功建立。

客户端一共绑定了四个事件。

-   1）onopen 建立连接后触发
-   2）onmessage 收到消息后触发
-   3）onerror 发生错误时触发
-   4）onclose 关闭连接时触发

服务器端的实现

```go
package main

import (
	"golang.org/x/net/websocket"
	"fmt"
	"log"
	"net/http"
)

func Echo(ws *websocket.Conn) {
	var err error

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

func main() {
	http.Handle("/", websocket.Handler(Echo))

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
```

当 客户端 将用户输入的信息 Send 之后，服务器端 通过 Receive 接收到了相应信息，然后通过 Send 发送了应答信息。

## 8.3 REST

RESTful，是目前最为流行的一种互联网软件架构。因为它结构清晰、符合标准、易于理解、扩展方便。

REST是"资源表现层状态转化"。

### 什么是REST

REST(REpresentational State Transfer)：它指的是一组架构约束条件和原则。满足这些约束条件和原则的应用程序或设计就是 RESTful 的。

理解下面几个概念：

-   资源（Resources） REST是"表现层状态转化"，其实它省略了主语。"表现层"其实指的是"资源"的"表现层"。

    什么是资源呢？

    就是我们平常上网访问的一张图片、一个文档、一个视频等。这些资源我们通过URI来定位，也就是一个URI表示一个资源。

-   表现层（Representation）资源是做一个具体的实体信息，它可以有多种的展现方式。而把实体展现出来就是表现层。例如：一个txt文本信息，他可以输出成html、json、xml等格式

    URI 确定一个资源，但是如何确定它的具体表现形式呢？

    应该在HTTP请求的头信息中用 `Accept` 和 `Content-Type` 字段指定，这两个字段才是对"表现层"的描述。

-   状态转化（State Transfer）访问一个网站，就代表了客户端和服务器的一个互动过程。在这个过程中，肯定涉及到数据和状态的变化。而 HTTP协议 是 无状态 的，那么这些状态肯定保存在服务器端，所以如果客户端想要通知服务器端改变数据和状态的变化，肯定要通过某种方式来通知它。

    客户端能通知服务器端的手段，只能是HTTP协议。

    具体来说，就是HTTP协议里面，四个表示操作方式的动词：GET、POST、PUT、DELETE。它们分别对应四种基本操作：GET用来获取资源，POST用来新建资源（也可以用于更新资源），PUT用来更新资源，DELETE用来删除资源。

总结一下什么是RESTful架构：

-   （1）每一个 URI 代表一种资源；
-   （2）客户端 和 服务器之间，传递 这种资源的 某种表现层；
-   （3）客户端通过 四个 HTTP 动词，对服务器端 资源进行操作，实现 "表现层状态转化"。

Web应用要满足 REST 最重要的原则：客户端和服务器之间的交互在请求之间是无状态的。

即，从客户端到服务器的每个请求都必须包含理解请求所必需的信息。如果服务器在请求之间的任何时间点重启，客户端不会得到通知。

此外，此请求可以由任何可用服务器回答，这十分适合云计算之类的环境。因为是无状态的，所以客户端可以缓存数据以改进性能。

另一个重要的REST原则：系统分层。

这表示组件无法了解除了与它直接交互的层次以外的组件。通过将系统知识限制在单个层，可以限制整个系统的复杂性，从而促进了底层的独立性。

REST的架构图：

![1555133012763](D:\Documents\笔记本\offer学习复习\Web开发\1555133012763.png)

当REST架构的约束条件作为一个整体应用时：

-   将生成一个可以扩展到大量客户端的应用程序。
-   它还降低了客户端和服务器之间的交互延迟。
-   统一界面简化了整个系统架构，改进了子系统之间交互的可见性。
-   REST 简化了客户端和服务器的实现，而且对于使用REST开发的应用程序更加容易扩展。

展示 REST 的扩展性：

![1555133104743](D:\Documents\笔记本\offer学习复习\Web开发\1555133104743.png)

### RESTful的实现

RESTful 是基于 HTTP协议 实现的，所以我们可以利用 `net/http` 包来自己实现，当然需要针对 REST 做一些改造，REST 是根据不同的 method 来处理相应的资源。REST 的 level 分级：

![1555133164490](D:\Documents\笔记本\offer学习复习\Web开发\1555133164490.png)

在应用开发的时候也不一定全部按照 RESTful 的规则全部实现他的方式。

因为有些时候完全按照RESTful的方式未必是可行的，RESTful服务充分利用每一个HTTP方法，包括`DELETE`和`PUT`。可有时，HTTP客户端只能发出 `GET` 和 `POST` 请求：

-   HTML 标准只能通过 链接 和 表单 支持 `GET` 和 `POST` 。

    在没有 Ajax 支持的网页浏览器中不能发出 `PUT` 或 `DELETE` 命令

-   有些防火墙会挡住 HTTP  `PUT` 和 `DELETE` 请求，要绕过这个限制，客户端需要把实际的`PUT`和`DELETE`请求通过 POST 请求穿透过来。

    RESTful 服务则要负责在收到的 POST 请求中 找到原始的 HTTP 方法并还原。

我们现在可以通过`POST`里面增加隐藏字段`_method`这种方式可以来模拟`PUT`、`DELETE`等方式，但是服务器端需要做转换。我现在的项目里面就按照这种方式来做的REST接口。当然Go语言里面完全按照RESTful来实现是很容易的，我们通过下面的例子来说明如何实现RESTful的应用设计。

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func getuser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid := ps.ByName("uid")
	fmt.Fprintf(w, "you are get user %s", uid)
}

func modifyuser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid := ps.ByName("uid")
	fmt.Fprintf(w, "you are modify user %s", uid)
}

func deleteuser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid := ps.ByName("uid")
	fmt.Fprintf(w, "you are delete user %s", uid)
}

func adduser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// uid := r.FormValue("uid")
	uid := ps.ByName("uid")
	fmt.Fprintf(w, "you are add user %s", uid)
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	router.GET("/user/:uid", getuser)
	router.POST("/adduser/:uid", adduser)
	router.DELETE("/deluser/:uid", deleteuser)
	router.PUT("/moduser/:uid", modifyuser)

	log.Fatal(http.ListenAndServe(":8080", router))
}
```

代码演示了如何编写一个REST的应用，我们访问的资源是用户，我们通过不同的method来访问不同的函数，这里使用了第三方库 [github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) 。这个库实现了自定义路由和方便的路由规则映射，通过它，我们可以很方便的实现REST的架构。

通过上面的代码可知，REST 就是根据不同的 method 访问 同一个资源 的时候 实现不同的 逻辑处理。

### 总结

REST是一种架构风格，汲取了WWW的成功经验：无状态，以资源为中心，充分利用HTTP协议和URI协议，提供统一的接口定义，使得它作为一种设计Web服务的方法而变得流行。

## 8.4 RPC

Socket 和 HTTP 采用的是类似 "信息交换" 模式，即客户端发送一条信息到服务端，然后 (一般来说)服务器端都会返回一定的信息以表示响应。客户端和服务端之间约定了交互信息的格式，以便双方都能够解析交互所产生的信息。

但是很多独立的应用并没有采用这种模式，而是采用类似常规的函数调用的方式来完成想要的功能。

RPC 就是 想实现 函数 调用模式的 网络化。客户端就像调用本地函数一样，然后客户端把这些参数打包之后通过网络传递到服务端，服务端解包到处理过程中执行，然后执行的结果反馈给客户端。

RPC（Remote Procedure Call Protocol）——远程过程调用协议，是一种通过网络从远程计算机程序上请求服务，而不需要了解底层网络技术的协议。

它假定某些传输协议的存在，如TCP或UDP，以便为通信程序之间携带信息数据。通过它可以使函数调用模式网络化。在OSI网络通信模型中，RPC跨越了传输层和应用层。RPC使得开发包括网络分布式多程序在内的应用程序更加容易。

### RPC工作原理

![1555135204758](D:\Documents\笔记本\offer学习复习\Web开发\1555135204758.png)

运行时,一次客户机对服务器的RPC调用,其内部操作大致有如下十步：

-   1，调用客户端句柄；执行传送参数
-   2，调用本地系统内核发送网络消息
-   3，消息传送到远程主机
-   4，服务器句柄得到消息并取得参数
-   5，执行远程过程
-   6，执行的过程将结果返回服务器句柄
-   7，服务器句柄返回结果，调用远程系统内核
-   8，消息传回本地主机
-   9，客户句柄由内核接收消息
-   10，客户接收句柄返回的数据

### Go RPC

Go标准包中已经提供了对RPC的支持，而且支持三个级别的RPC：TCP、HTTP、JSONRPC。

但，Golang 的 RPC包 是独一无二的 RPC，它和传统的 RPC 系统不同，它 只支持 Go 开发的服务器与客户端之间的交互，因为在内部，它们采用了 Gob 来 编码。

Go RPC 的函数只有符合下面的条件才能被远程访问，不然会被忽略，详细的要求如下：

-   函数必须是 导出 的(首字母大写)
-   必须有两个导出类型的参数，
-   第一个参数是接收的参数，第二个参数是返回给客户端的参数，第二个参数必须是指针类型的
-   函数还要有一个返回值error

正确的RPC函数格式如下：

```
func (t *T) MethodName(argType T1, replyType *T2) error
```

T、T1 和 T2 类型必须能被 `encoding/gob` 包编解码。

任何的 RPC都需要通过网络来传递数据，Go RPC可以利用HTTP和TCP来传递数据，利用 HTTP 的好处是可以直接复用 `net/http` 里面的一些函数。

### HTTP RPC

http的服务端代码实现如下：

```
package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main() {

	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
```

可以看到，我们 注册了一个 Arith 的RPC服务，然后通 过`rpc.HandleHTTP `函数把该服务注册到了 HTTP 协议上，然后我们就可以利用http的方式来传递数据了。

请看下面的客户端代码：

```
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "server")
		os.Exit(1)
	}
	serverAddress := os.Args[1]

	client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	args := Args{17, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	var quot Quotient
	err = client.Call("Arith.Divide", args, &quot)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d/%d=%d remainder %d\n", args.A, args.B, quot.Quo, quot.Rem)

}
```

我们把上面的服务端和客户端的代码分别编译，然后先把服务端开启，然后开启客户端，输入代码，就会输出如下信息：

```
$ ./http_c localhost
Arith: 17*8=136
Arith: 17/8=2 remainder 1
```

通过 调用，可以看到，参数 和 返回值 是我们定义的 struct 类型，在服务端，我们把它们当做调用函数的参数的类型，在客户端作为`client.Call`的第2，3两个参数的类型。

客户端最重要的就是这个 `Call()` ，它有  3 个参数，第  1 个要调用的函数的名字，第  2 个是要传递的参数，第  3个要返回的参数 (注意是 指针类型)。

### TCP RPC

实现基于TCP协议的RPC，服务端的实现代码如下所示：

```
package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main() {

	arith := new(Arith)
	rpc.Register(arith)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		// rpc.ServeConn(conn)
		go rpc.ServeConn(conn) // 启动 goroutine 提高并发性
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
```

上面这个 代码 和 http 的服务器相比，不同在于：在此处我们采用了 TCP 协议，然后需要自己控制连接，当有客户端连接上来后，我们需要把这个连接交给 rpc 来处理。

下面展现了TCP实现的RPC客户端：

```
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "server:port")
		os.Exit(1)
	}
	service := os.Args[1]

	client, err := rpc.Dial("tcp", service)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	args := Args{17, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	var quot Quotient
	err = client.Call("Arith.Divide", args, &quot)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d/%d=%d remainder %d\n", args.A, args.B, quot.Quo, quot.Rem)

}
```

这个客户端代码和http的客户端代码对比，唯一的区别一个是 `DialHTTP()` ，一个是 `Dial(tcp)`，其他处理一模一样。

### JSON RPC

JSON RPC 是数据编码采用了 JSON，而不是  gob 编码，其他和上面介绍的RPC概念一模一样。使用Go提供的json-rpc标准包，请看服务端代码的实现：

```
package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main() {

	arith := new(Arith)
	rpc.Register(arith)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		jsonrpc.ServeConn(conn)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
```

通过示例我们可以看出  json-rpc 是基于TCP协议实现的，**目前它还不支持HTTP方式。**

请看客户端的实现代码：

```
package main

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
	"os"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "server:port")
		log.Fatal(1)
	}
	service := os.Args[1]

	client, err := jsonrpc.Dial("tcp", service)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	args := Args{17, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	var quot Quotient
	err = client.Call("Arith.Divide", args, &quot)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d/%d=%d remainder %d\n", args.A, args.B, quot.Quo, quot.Rem)

}
```

### 总结

Go 已经提供了对 RPC 的良好支持，通过上面 HTTP、TCP、JSON RPC 的实现，我们就可以很方便的开发很多分布式的Web应用。但遗憾的是目前 Go 尚未提供对 SOAP RPC 的支持，欣慰的是现在已经有第三方的开源实现了。

## 回望

介绍了目前流行的几种主要的网络应用开发方式：

第一小节介绍了网络编程中的基础：Socket 编程，因为现在网络正在朝云的方向快速进化，作为这一技术演进的基石的的 socket 知识，作为开发者的你，是必须要掌握的。

第二小节介绍了正愈发流行的 HTML5 中一个重要的特性 WebSocket ，通过它,服务器可以实现主动的push消息，以简化以前ajax轮询的模式。

第三小节介绍了 REST编写模式，这种模式特别适合来开发网络应用API，目前移动应用的快速发展，我觉得将来会是一个潮流。

第四小节介绍了Go实现的RPC相关知识。

对于上面四种开发方式，Go都已经提供了良好的支持，net包及其子包,是所有涉及到网络编程的工具的所在地。如果你想更加深入的了解相关实现细节，可以尝试阅读这个包下面的源码。