# Golang wiki interface

目录

[TOC]

## Golang and Duck Typing 的关系

what is Duck Typing？

`Duck Typing`，鸭子类型，是动态编程语言的一种对象推断策略，它更关注对象能如何被使用，而不是对象的类型本身。Go 语言作为一门静态语言，它通过通过 接口的方式完美支持鸭子类型。

>   动态语言的特点：
>
>   变量绑定的类型是不确定的，在运行期间才能确定 函数和方法可以接收任何类型的参数，且调用时不检查参数类型 不需要实现接口

总结

鸭子类型是一种动态语言的风格，在这种风格中，一个对象有效的语义，不是由继承自特定的类或实现特定的接口，而是由它"当前方法和属性的集合"决定。

Go 作为一种静态语言，通过接口实现了 `鸭子类型`，实际上是 Go 的编译器在其中作了隐匿的转换工作。

## 值接收者 and 指针接收者 的区别

### 方法

方法有什么用？方法能给用户自定义的类型添加新的行为。

它和函数的区别在于方法有一个接收者，给一个函数添加一个接收者，那么它就变成了方法。接收者可以是`值接收者`，也可以是`指针接收者`。

KEY：在调用方法的时候，值类型既可以调用`值接收者`的方法，也可以调用`指针接收者`的方法；指针类型既可以调用`指针接收者`的方法，也可以调用`值接收者`的方法。（这是 Golang 的语法糖，编译器完成了转换工作）

### 值接收者 and 指针接收者

结论：实现了接收者是值类型的方法，相当于自动实现了接收者是指针类型的方法；而实现了接收者是指针类型的方法，不会自动生成对应接收者是值类型的方法。

```
// func (p *Gopher) debug()
// func (p Gopher) code()
func main() {
    var c coder = &Gopher{"Go"}
    c.code()
    c.debug()
}
output:
code()
debug()

// -----------
func main() {
    var c coder = Gopher{"Go"} // 这里更改为 值类型
    c.code()
    c.debug()
}
output:
./main.go:24:6: cannot use Programmer literal (type Programmer) as type coder in assignment:
    Programmer does not implement coder (debug method has pointer receiver)
```

解释：接收者是指针类型的方法，很可能在方法中会对接收者的属性进行更改操作，从而影响接收者；而对于接收者是值类型的方法，在方法中不会对接收者本身产生影响。

总结

如果实现了接收者是值类型的方法，会隐含地也实现了接收者是指针类型的方法。

### 两者分别在何时使用

如果需要修改接收者指向的值   or 降低值的复制成本，选用指针接收者；

否则，选用 值接收者；

如果类型具备非原始的本质，不能 被安全地 复制，这种类型总是应该被共享，那就定义指针接收者的方法。

## iface and eface 的区别是什么

略，建议阅读原文

## 接口的动态类型和动态值

### 接口类型和 `nil` 作比较

接口值的零值是指`动态类型`和`动态值`都为 `nil`。当仅且当这两部分的值都为 `nil` 的情况下，这个接口值就才会被认为 `接口值 == nil`。

```
type Coder interface {
    code()
}

type Gopher struct {
    name string
}

func (g Gopher) code() {
    fmt.Printf("%s is coding\n", g.name)
}

func main() {
    var c Coder
    fmt.Println(c == nil)
    fmt.Printf("c: %T, %v\n", c, c)

    var g *Gopher
    fmt.Println(g == nil)

    c = g
    fmt.Println(c == nil) 
    fmt.Printf("c: %T, %v\n", c, c)
}
output:
true
c: <nil>, <nil>
true
false // g 赋值给 c后，c 的动态类型为 *main.Gopher，所以输出为 false
c: *main.Gopher, <nil>
```

### 接口类型的隐形转换

```
type MyError struct {}

func (i MyError) Error() string {
    return "MyError"
}

func Process() error {
    var err *MyError = nil
    return err
}

func main() {
    err := Process()
    fmt.Println(err)

    fmt.Println(err == nil)
}
output:
<nil> // 值依旧为 nil
false // Process 函数返回了一个 error 接口，这块隐含了类型转换，它的类型转换为 *MyError
```

### 如何打印出接口的动态类型和值？

```
type iface struct {
    itab, data uintptr
}

func main() {
    var a interface{} = nil

    var b interface{} = (*int)(nil)

    x := 5
    var c interface{} = (*int)(&x)

    ia := *(*iface)(unsafe.Pointer(&a))
    ib := *(*iface)(unsafe.Pointer(&b))
    ic := *(*iface)(unsafe.Pointer(&c))

    fmt.Println(ia, ib, ic)

    fmt.Println(*(*int)(unsafe.Pointer(ic.data)))
}
output:
{0 0} {17426912 0} {17426912 842350714568}
5
```

a 的动态类型和动态值的地址均为 0，也就是 nil；b 的动态类型和 c 的动态类型一致，都是 `*int`；最后，c 的动态值为 5。

## 编译器自动检测类型是否实现接口

怎么操作

```
var _ io.Writer = (*myWriter)(nil)
```

这是什么操作？编译器会检查 `*myWriter` 类型是否实现了 `io.Writer` 接口。

总结

可添加以下代码，用来检测类型是否实现了接口

```
var _ io.Writer = (*myWriter)(nil)
var _ io.Writer = myWriter{}
```

## 接口的构造过程是怎样的

过程略，建议阅读原文

如何打印出接口类型的 `Hash` 值？

```
type iface struct {
    tab  *itab
    data unsafe.Pointer
}
type itab struct {
    inter uintptr
    _type uintptr
    link uintptr
    hash  uint32
    _     [4]byte
    fun   [1]uintptr
}

func main() {
    var qcrao = Person(Student{age: 18})

    iface := (*iface)(unsafe.Pointer(&qcrao))
    fmt.Printf("iface.tab.hash = %#x\n", iface.tab.hash)
}
output:
iface.tab.hash = 0xd4209fda
```

`hash` 值只和他的字段、方法相关。

## 类型转换 and 断言 的区别

`类型转换`、`类型断言`本质都是把一个类型转换成另外一个类型。不同之处在于，类型断言是对接口变量进行的操作。

### 类型转换

案例

```
func main() {
    var i int = 9

    var f float64
    f = float64(i)
    fmt.Printf("%T, %v\n", f, f)

    f = 10.8
    a := int(f)
    fmt.Printf("%T, %v\n", a, a)
}
```

语法

```
<结果类型> := <目标类型> (<表达式>)
```

### 断言

案例

```
	switch v := v.(type) {
    case nil:
        ...
    case Student:
        ...
```

语法

```
// 安全类型断言
<目标类型的值>，<布尔参数> := <表达式>.( 目标类型 )  

//非安全类型断言
<目标类型的值> := <表达式>.( 目标类型 )
```

类型不同声明形式的底层状况

```
// 三次输出的地址不一致：
// 调用函数时，实际上是复制了一份参数；
// 断言之后，又生成了一份新的拷贝；

// --- var i interface{} = new(Student)
0xc4200701b0 [Name: ], [Age: 0]
0xc4200701d0 [Name: ], [Age: 0]
0xc420080020 [Name: ], [Age: 0]
*Student type[*main.Student] [Name: ], [Age: 0] // 类型 *main.Student，值 [Name: ], [Age: 0]


// --- var i interface{} = (*Student)(nil)
0xc42000e1d0 <nil>
0xc42000e1f0 <nil>
0xc42000c030 <nil>
*Student type[*main.Student] <nil> // 类型 *main.Student，值 nil

// --- var i interface{}
0xc42000e1d0 <nil>
0xc42000e1e0 <nil>
0xc42000e1f0 <nil>
nil type[<nil>] <nil> // 类型 nil，值 nil
```

拓展

`fmt.Println` 函数的参数是 `interface`。

对于内置类型，函数内部会用穷举法，得出它的真实类型，然后转换为字符串打印。而对于自定义类型，首先确定该类型是否实现了 `String()` 方法，如果实现了，则直接打印输出 `String()` 方法的结果；否则，会通过反射来遍历对象的成员进行打印。

>   类型 `T` 只有接受者是 `T` 的方法；而类型 `*T` 拥有接受者是 `T` 和 `*T` 的方法。语法上 `T` 能直接调 `*T` 的方法仅仅是 `Go` 的语法糖。

## 接口转换 的原理

当判定一种类型是否满足某个接口时，Go 使用类型的方法集和接口所需要的方法集进行匹配，如果类型的方法集完全包含接口的方法集，则可认为该类型实现了该接口。当然，能转换的原因必然是类型兼容。

探索将一个接口转换给另外一个接口背后的原理。。。建议阅读原文

## 如何用 interface 实现多态

多态的特点

1.一种类型具有多种类型的能力

2.允许不同的对象对同一消息做出灵活的反应

3.以一种通用的方式对待个使用的对象

4.非动态语言必须通过继承和接口的方式来实现

案例

```

func main() {
    qcrao := Student{age: 18}
    whatJob(&qcrao)

    growUp(&qcrao)
    fmt.Println(qcrao)

    stefno := Programmer{age: 100}
    whatJob(stefno)

    growUp(stefno)
    fmt.Println(stefno)
}

func whatJob(p Person) {
    p.job()
}

func growUp(p Person) {
    p.growUp()
}

type Person interface {
    job()
    growUp()
}

type Student struct {
    age int
}

func (p Student) job() {
    fmt.Println("I am a student.")
    return
}

func (p *Student) growUp() {
    p.age += 1
    return
}

type Programmer struct {
    age int
}

func (p Programmer) job() {
    fmt.Println("I am a programmer.")
    return
}

func (p Programmer) growUp() {
    // 程序员老得太快 ^_^
    p.age += 10
    return
}
output:
I am a student.
{19}
I am a programmer.
{100}
```

## Golang 接口 与 C++ 接口 异同

C++ 略，具体研究原文

接口定义了一种规范，描述了类的行为和功能，而不做具体实现。

Golang 采用的是 “非侵入式”，不需要显示声明，更方便后续拓展，只需要实现接口定义的函数，编译器自动会识别。

## 参考

饶全成@深度解密Go语言之关于 interface 的 10 个问题

<https://mp.weixin.qq.com/s?__biz=MjM5MDUwNTQwMQ==&mid=2257483749&idx=1&sn=b6bca6ac5afab7ac6963871d41a51473&chksm=a5391733924e9e2595174b0f8f354710a9bfc84ab5d1725119cea6030ccf3444708a2f98ff6a&scene=27#wechat_redirect&cpage=0>