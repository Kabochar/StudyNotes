# session&数据存储

文章整理自：https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/06.0.md

Web开发中一个很重要的议题：如何做好用户的整个浏览过程的控制？

因为 HTTP 协议是无状态的，所以用户的每一次请求都是无状态的，我们不知道在整个Web操作过程中哪些连接与该用户有关，我们应该如何来解决这个问题呢？

Web里面经典的解决方案是：cookie，session。

cookie 机制是一种 客户端 机制，把用户数据保存在客户端；session 机制是一种服务器端的机制，服务器使用一种类似于 散列表的结构 来保存信息，每一个网站 访客 都会被分配给一个 唯一 的 标志符 ，即 sessionID 。它的存放形式无非两种：一，经过url传递；二，保存在客户端的 cookies 里。

## 6.1 session 和 cookie

“登录” 过程中到底发生了什么？

当用户来到微博登录页面，输入用户名和密码之后点击 “登录 ” 后浏览器将认证信息  POST 给远端的服务器，服务器执行验证逻辑，如果验证通过，则浏览器会跳转到登录用户的微博首页。

在登录成功后，服务器如何验证我们对其他受限制页面的访问呢？因为 HTTP 协议是无状态的，所以很显然服务器不可能知道我们已经在上一次的 HTTP请求 中通过了验证。

当然，最简单的解决方案就是所有的请求里面都带上 用户名 和 密码，这样虽然可行，但大大加重了服务器的负担（对于每个request都需要到数据库验证），也大大降低了用户体验(每个页面都需要重新输入用户名密码，每个页面都带有登录表单)。

这时 COOKIE 为此而生--

cookie：在本地计算机保存一些用户操作的历史信息（当然包括登录信息），并在用户再次访问该站点时浏览器通过 HTTP 协议将 本地 cookie 内容 发送给服务器，从而完成验证，或继续上一步操作。

![1555076001634](pics\1555076001634.png)



session：在服务器上保存用户操作的历史信息。

服务器使用 session id 来标识 session，session id 由服务器负责产生，保证随机性与唯一性，相当于一个随机密钥，避免在握手或传输中暴露用户真实密码。

但该方式下，仍然需要将发送请求的 客户端 与 session 进行对应，所以可以借助cookie 机制 来获取 客户端的 标识（即session id），也可以通过 GET 方式将 id 提交给服务器。

![1555076059360](pics\1555076059360.png)

### Cookie

由 浏览器 维持的，存储 在客户端的一小段文本信息，伴随着用户请求和页面在 Web 服务器和浏览器之间传递。

cookie 是有时间限制的，根据生命期不同分成两种：会话 cookie，持久 cookie；

#### 会话 cookie

-   不设置过期时间。表示这个 cookie 的生命周期为从创建到浏览器关闭为止，只要关闭浏览器窗口，cookie 就消失了。
-   这种 生命期 为 浏览会话期 的 cookie 被称为 会话 cookie。
-   会话cookie 一般不保存在硬盘上而是保存在内存里。

####  持久 cookie

-   设置了过期时间 (setMaxAge(60*60*24)) ，浏览器就会把 cookie 保存到硬盘上，关闭后再次打开浏览器，这些 cookie 依然有效直到超过设定的过期时间。
-   存储在硬盘上的 cookie 可以在不同的浏览器进程间共享，比如两个 IE 窗口。
-   对于保存在内存的cookie，不同的浏览器有不同的处理方式。 

### Go 设置 cookie

通过 net/http 包中的 SetCookie 设置

```
http.SetCookie(w ResponseWriter, cookie *Cookie)
```

w 表示需要写入的 response，cookie 是一个struct

cookie 对象

```go
type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string

// MaxAge=0 means no 'Max-Age' attribute specified.
// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}
```

如何设置cookie

```go
expiration := time.Now()
expiration = expiration.AddDate(1, 0, 0)
cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
http.SetCookie(w, &cookie)
```

### Go 读取 cookie

如何读取cookie？通过 request 获取 cookie

```go
// 方式一
cookie, _ := r.Cookie("username")
fmt.Fprint(w, cookie)

// 方式二
for _, cookie := range r.Cookies() {
	fmt.Fprint(w, cookie.Name)
}
```

## session

session：当 session 一词与网络协议相关联时，它隐含了 “面向连接 ” 和（或） “保持状态” 这样两个含义。

含义（使用场景）：指一类用来在 客户端 与 服务器端 之间保持状态的解决方案。有时候Session也用来指这种解决方案的存储结构。

session机制：一种服务器端的机制，服务器使用一种类似于 散列表的结构 (也可能就是使用散列表) 来保存信息。

程序需要为某个客户端的请求创建一个 session 的时候，服务器首先检查这个客户端的请求里是否包含了一个 session 标识－称为 session id，如果已经包含一个 session id 则说明以前已经为此客户创建过 session，服务器就按照 session id 把这个 session 检索出来使用 ( 如果检索不到，可能会新建一个，这种情况可能出现在服务端已经删除了该用户对应的session对象，但用户人为地在请求的URL后面附加上一个JSESSION的参数 )。如果客户请求 不包含 session id，则为此客户 创建一个 session 并且 同时生成一个与此 session相关联的 session id，这个session id 将在本次响应中返回给客户端保存。

### 总结

session 和c ookie 的目的相同，都是为了克服 http 协议无状态的缺陷，但完成的方法不同。session 通过 cookie，在客户端保存 session id，而将用户的其他会话消息保存在服务端的 session 对象中，与此相对的，cookie 需要将所有信息都保存在客户端。

cookie 存在着一定的安全隐患，例如本地 cookie 中保存的用户名密码被破译，或cookie被其他网站收集（例如：1. appA 主动设置域 B cookie，让域 B cookie 获取；2. XSS，在appA上通过 javascript 获取 document.cookie，并传递给自己的 appB）。

## 6.2 Go如何使用 session

### session 创建过程

session 的基本原理：由服务器为每个会话维护一份信息数据，客户端和服务端依靠一个 全局唯一 的 标识 来访问这份数据，以达到交互的目的。

当用户访问 Web 应用时，服务端程序会随需要创建 session，这个过程可以概括为三个步骤：

-   生成全局唯一标识符（sessionid）；
-   开辟数据存储空间。写到，内存 or 文件里或存储在数据库中。
-   将 session 的全局唯一标示符发送给客户端。

最关键的是如何发送这个 session 的唯一标识这一步上。一般来说会有两种常用的方式：cookie 和 URL 重写。

1.  Cookie 服务端通过设置 Set-cookie 头就可以将 session 的标识符传送到客户端。

    而客户端此后的每一次请求都会带上这个标识符，另外一般包含 session 信息的 cookie 会将 失效时间 设置为 0 (会话cookie)，即浏览器进程有效时间。

    至于浏览器怎么处理这个 0，每个浏览器都有自己的方案，但差别都不会太大(一般体现在新建浏览器窗口的时候)；

2.  URL重写。所谓 URL 重写，就是在返回给用户的页面里的所有的 URL 后面追加session 标识符。

    这样用户在收到响应之后，无论点击响应页面里的哪个链接或提交表单，都会自动带上 session 标识符，从而就实现了会话的保持。

    虽然这种做法比较麻烦，但是，如果客户端禁用了cookie的话，此种方案将会是首选。

### Go 实现 session 管理

结合 session 的生命周期（lifecycle），来实现 go 语言版本的 session 管理。

### session 管理设计

session 管理涉及到如下几个因素：

-   全局 session 管理器
-   保证 session id 的全局唯一性
-   为每个 客户端 关联一个 session
-   session 的存储(可以存储到内存、文件、数据库等)
-   session 过期处理

### session 管理器

定义一个全局的session管理器

```go
type Manager struct {
	cookieName  string     // private cookiename
	lock        sync.Mutex // protects session
	provider    Provider
	maxLifeTime int64
}

func NewManager(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxLifeTime: maxLifeTime}, nil
}
```

抽象出一个Provider接口，用以表征session管理器底层存储结构。

```go
type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}
```

-   SessionInit()  实现 Session 的初始化，操作成功，则返回此新的 Session 变量
-   SessionRead()  返回 sid 所代表的 Session 变量，如果不存在，那么将以 sid为参数调用 SessionInit 函数创建并返回一个新的 Session 变量
-   SessionDestroy()  用来销毁 sid 对应的 Session 变量
-   SessionGC()  根据 maxLifeTime 来删除过期的数据

Session接口需要实现什么样的功能呢？

对 Session 的处理基本就 设置值、读取值、删除值 以及 获取当前 sessionID 这四个操作，所以我们的 Session 接口也就实现这四个操作。

```
type Session interface {
	Set(key, value interface{}) error // set session value
	Get(key interface{}) interface{}  // get session value
	Delete(key interface{}) error     // delete session value
	SessionID() string                // back current sessionID
}
```

>   以上设计思路来源于database/sql/driver，先定义好接口，然后具体的存储 session 的结构实现相应的接口并注册后，相应功能这样就可以使用了.

Register()  随时注册存储 session 的结构

```go
var provides = make(map[string]Provider)

// Register makes a session provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	provides[name] = provider
}
```

### 全局唯一的 Session ID

Session ID 是用来识别访问 Web 应用的每一个用户，因此必须保证它是全局唯一的（GUID）

```go
func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
```

### session 创建

需要为每个来访用户 分配 或 获取 与他相关连的 Session，以便后面根据 Session信息来验证操作。

SessionStart()：用来检测是否已经有某个 Session 与 当前来访用户 发生了 关联，如果没有，则创建之。

```go
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxLifeTime)}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}
```

Login 操作演示 session 的应用：

```go
func login(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, sess.Get("username"))
	} else {
		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/", 302)
	}
}
```

### 操作值：设置、读取和删除

SessionStart 函数 返回的是一个满足 Session 接口的 变量，那么我们该如何用他来对 session 数据进行操作呢？

```go
func count(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	createtime := sess.Get("createtime")
	if createtime == nil {
		sess.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 360) < (time.Now().Unix()) {
		globalSessions.SessionDestroy(w, r)
		sess = globalSessions.SessionStart(w, r)
	}
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	t, _ := template.ParseFiles("count.gtpl")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("countnum"))
}
```

因为Session有过期的概念，所以我们定义了GC操作，当访问过期时间满足GC的触发条件后将会引起GC，但是当我们进行了任意一个session操作，都会对Session实体进行更新，都会触发对最后访问时间的修改，这样当GC的时候就不会误删除还在使用的Session实体。

### session 重置

Web应用中有用户退出这个操作，那么当用户退出应用的时候，我们需要对该用户的session数据进行销毁操作。

```go
//Destroy sessionid
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}
```

### session 销毁

Session管理器如何来管理销毁，只要我们在 Main 启动的时候启动：

```
func init() {
	go globalSessions.GC()
}
```

GC充分利用了 time 包中的定时器功能，当超时 `maxLifeTime` 之后调用GC函数，这样就可以保证 `maxLifeTime` 时间内的 session 都是可用的，类似的方案也可以用于统计在线用户数之类的。

```go
func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() { manager.GC() })
}
```

## 6.3 session 存储

基于内存的 session 存储接口的实现

```go
package memory

import (
	"container/list"
	"github.com/astaxie/session"
	"sync"
	"time"
)

var pder = &Provider{list: list.New()}

type SessionStore struct {
	sid          string                      //session id唯一标示
	timeAccessed time.Time                   //最后访问时间
	value        map[interface{}]interface{} //session里面存储的值
}

func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) Get(key interface{}) interface{} {
	pder.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	} else {
		return nil
	}
}

func (st *SessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) SessionID() string {
	return st.sid
}

type Provider struct {
	lock     sync.Mutex               //用来锁
	sessions map[string]*list.Element //用来存储在内存
	list     *list.List               //用来做gc
}

func (pder *Provider) SessionInit(sid string) (session.Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := pder.list.PushBack(newsess)
	pder.sessions[sid] = element
	return newsess, nil
}

func (pder *Provider) SessionRead(sid string) (session.Session, error) {
	if element, ok := pder.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		sess, err := pder.SessionInit(sid)
		return sess, err
	}
	return nil, nil
}

func (pder *Provider) SessionDestroy(sid string) error {
	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		pder.list.Remove(element)
		return nil
	}
	return nil
}

func (pder *Provider) SessionGC(maxlifetime int64) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (pder *Provider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}

func init() {
	pder.sessions = make(map[string]*list.Element, 0)
	session.Register("memory", pder)
}
```

通过 init 函数注册到 session 管理器中。这样就可以方便的调用了。

如何来调用该引擎呢？

```
import (
	"github.com/astaxie/session"
	_ "github.com/astaxie/session/providers/memory"
)
```

当 import 的时候已经执行了 memory 函数里面的 init 函数，这样就已经注册到session 管理器中，我们就可以使用了，通过如下方式就可以初始化一个session管理器：

```
var globalSessions *session.Manager

//然后在init函数中初始化
func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
}
```

## 6.4 预防 session 劫持

在session技术中，客户端和服务端通过session的标识符来维护会话， 但这个标识符很容易就能被嗅探到，从而被其他人利用。它是中间人攻击的一种类型。

## session劫持过程

测试案例

```
func count(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	t, _ := template.ParseFiles("count.gtpl")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("countnum"))
}
```

count.gtpl的代码如下所示：

```
Hi. Now count:{{.}}
```

 打开另一个浏览器(这里我打开了firefox浏览器)，复制 chrome 地址栏里的地址到新打开的浏览器的地址栏中。然后打开 firefox 的 cookie 模拟插件，新建一个cookie ，把按上图中 cookie 内容原样在 firefox 中重建一份。

虽然换了浏览器，但是却获得了 sessionID，然后模拟了 cookie 存储的过程。

如果交替点击两个浏览器里的链接你会发现它们其实操纵的是同一个计数器。不必惊讶，此处 firefox 盗用了 chrome 和 goserver 之间的维持会话的钥匙，即gosessionid，这是一种类型的“会话劫持”。

在 goserver 看来，它从 http 请求中得到了一个 gosessionid，由于 HTTP 协议的无状态性，它无法得知这个 gosessionid 是从 chrome 那里“劫持”来的，它依然会去查找对应的 session，并执行相关计算。与此同时 chrome 也无法得知自己保持的会话已经被 “劫持”。

## session劫持防范

如何有效的防止session劫持呢？

#### cookieonly 和 token

其中一个解决方案：

sessionID 的值只允许 cookie 设置，而不是通过 URL 重置方式设置，同时设置cookie 的 httponly 为 true。

**httponly** ：设置是否可通过客户端脚本访问这个设置的 cookie。

第一，这个可以防止这个 cookie 被 XSS 读取从而引起 session 劫持；第二，cookie 设置不会像 URL 重置方式那么容易获取 sessionID。

第二步，在每个请求里面加上token（防止 form 重复递交类似），我们在每个请求里面加上一个隐藏的 token，然后每次验证这个token，从而保证用户的请求都是唯一性。

```
h := md5.New()
salt:="astaxie%^7&8888"
io.WriteString(h,salt+time.Now().String())
token:=fmt.Sprintf("%x",h.Sum(nil))
if r.Form["token"]!=token{
	//提示登录
}
sess.Set("token",token)
```

#### 间隔生成新的SID

给session额外设置一个创建时间的值，一旦过了一定的时间，我们销毁这个sessionID，重新生成新的session，这样可以一定程度上防止session劫持的问题。

````go
createtime := sess.Get("createtime")
if createtime == nil {
	sess.Set("createtime", time.Now().Unix())
} else if (createtime.(int64) + 60) < (time.Now().Unix()) {
	globalSessions.SessionDestroy(w, r)
	sess = globalSessions.SessionStart(w, r)
}
````

session 启动后，我们设置了一个值，用于记录生成sessionID的时间。通过判断每次请求是否过期(这里设置了60秒)定期生成新的ID，这样使得攻击者获取有效sessionID的机会大大降低。

### 总结

两个手段的组合可以在实践中消除 session 劫持的风险，一方面，由于 sessionID 频繁改变，使攻击者难有机会获取有效的 sessionID；另一方面，因为sessionID只能在 cookie  中传递，然后设置了 httponly，所以基于 URL 攻击的可能性为零，同时被 XSS 获取 sessionID 也不可能。最后，由于我们还设置了 MaxAge=0，这样就相当于 session cookie 不会留在浏览器的历史记录里面。