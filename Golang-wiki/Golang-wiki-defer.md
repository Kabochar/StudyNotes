# Golang wiki defer

目录

[TOC]

## 它是什么？

defer：一种用于注册延迟调用的机制：让函数 or 语句可以在当前函数执行完毕后执行。（包括：return 正常结束 or panic 导致的异常结束）

使用场景：成对操作，打开/关闭连接；打开/关闭文件 etc

## 为什么需要它？

资源需要在用完之后需要释放掉，否则会造成内容泄露问题。欠钱不还，能有什么后果，你懂得啦！

它的出现正好解决了，打开资源的下一段，可以使用它来顺手关闭资源操作。

## 该如何合理操作它？

样例：

```go
f, err := os.Open(filename)
if err != nil {
	panic(err)
}
if f != nil {
    defer f.Close()
}
```

注意：defer 存在小小延迟问题，对时间操作要求高的程序，可以避免使用它。

## defer 进阶

### defer的底层原理

defer语句并不会马上执行，而是会进入一个栈，函数return前，会按先进后出（FILO）的顺序执行。

为什么是先进后出（FILO）？

后面定义的函数可能会依赖前面的资源，自然要先执行；否则，如果前面先执行，那后面函数的依赖就没有了；

defer如何操作对 “外部变量的引用” ？

在defer函数定义时，对外部变量的引用是有两种方式的，分别是作为函数参数和作为闭包引用：

-   作为函数参数，则在defer定义时就把值传递给defer，并被cache起来；
-   作为闭包引用的话，则会在defer函数真正调用时根据整个上下文确定当前的值；

注意：defer 后面语句执行时，函数调用的参数会自行 拷贝一份：

-   如果 变量是 “值”，与定义时一致；
-   如果 变量是 “引用”，可能与定义不一致；

```
func main() {
    var what [3]struct{}
    
    for i := range what {
        defer func() {
            fmt.Println(i)
        }
    }
}

Output：
2
2
2
```

为什么是 2？defer 后面跟的是一个闭包，i 是 引用类型，最后 i 的值为 2，所以变为 2

再一个

```
type number int

func (n number) print() { fmt.Println(n) }
func (n *number) pprint() { fmt.Println(*n) }

func main() {
    var n number
    
    defer n.print()  // 值
    defer n.pprint() // 引用 
    defer func() { n.print() } // 闭包
    defer func() { n.pprint() } // 闭包
}

Ouput:
3
3
3
0
```

### 利用 defer 先求值，再延迟调用

利用什么？defer 先求值，再延迟调用的性质；

```
func mergeFile() error {
    f, _ := os.Open("f1.txt")
    if f != nil { // 这样操作才是安全的
        defer func(f io.Closer) {
            if err := f.Close(); err != nil {
                fmt.Printf("defer close")
            }
        }(f) // 千万要加上括号
    }
    
    // 
    
    f, _ := os.Open("f2.txt")
    if f != nil { // 这样操作才是安全的
        defer func(f io.Closer) {
            if err := f.Close(); err != nil {
                fmt.Printf("defer close")
            }
        }(f) // 千万要加上括号
    }
    
    return nil
}
```

注意：先判断调用主体是否为空，否则会 panic

### defer & return xxx 联动

注意 defer  和  return xxx 的联动操作

经过  return xxx 的操作，变成三条指令

```
1，返回值 = xxx
2，调用 defer 函数
3，空的 return
```

1，3 才是return的真正命令，2 是 defer定义的语句，可能会操作返回指

案例 1

```
func fop() (res int) {
    t := 5
    defer func() {
        t = t + 5
    }()
    
    return t
}

func fop() (res int) {
    t := 5
    
    // 赋值指令
    res = t 
    
    // defer 被插入到 赋值与 return 之间执行，这个过程 res 的值没有被改变过
    func() { 
        t = t + 5
    }()
    
    // 空的 return
    return
}

Ouput:
5
```

案例 2

```
func fop() (res int) {
    defer func() {
        t = t + 5
    }(res)
    
    return 1
}

func fop() (res int) {
	res = 1
	
	// 这里改的 res 是之前传值传进去的 res，不会改变要 return 的 res 值
    defer func(res int) {
        res = res + 5
    }(res)
    
    return
}
Output:
1
```

## defer 参数

defer语句表达式的值在定义时就已经确定！！

```
func fop1() {
	var err error
	
    defer fmt.Println(err)
    
    err = errors.New("defer error")
    return
}

func fop2() {
	var err error
	
    defer func() {
       fmt.Println(err)
    }()
    
    err = errors.New("defer error")
    return
}

func fop3() {
	var err error
	
    defer func() {
       fmt.Println(err)
    }(err)
    
    err = errors.New("defer error")
    return
}

func main() {
    fop1()
    fop2()
    fop3()
}
```

第 1，3个函数是因为作为函数参数，定义的时候就会求值，定义的时候err变量的值都是nil；

第 2个例子中是一个闭包，它引用的变量err在执行的时候最终变成 `defer error`了；

## 闭包 闭包 闭包？

闭包是由函数及其相关引用环境组合而成的实体

```
闭包=函数+引用环境
```

一般的函数都有函数名，但是匿名函数就没有。匿名函数不能独立存在，但可以直接调用或者赋值于某个变量。匿名函数也被称为闭包，一个闭包继承了函数声明时的作用域。在Golang中，所有的匿名函数都是闭包。

有个不太恰当的例子，可以把闭包看成是一个类，一个闭包函数调用就是实例化一个类。闭包在运行时可以有多个实例，它会将同一个作用域里的变量和常量捕获下来，无论闭包在什么地方被调用（实例化）时，都可以使用这些变量和常量。而且，闭包捕获的变量和常量是引用传递，不是值传递。

```
func main() {
    var a = Accumulator()
    
    fmt.Printf("%d'n", a(1))
    fmt.Printf("%d'n", a(10))
    fmt.Printf("%d'n", a(100))
    
    var b = Accumulator()
    fmt.Printf("%d'n", b(1))
    fmt.Printf("%d'n", b(10))
    fmt.Printf("%d'n", b(100))
}

func Accumulator() func(int) int {
    var x int
    
    return func(delta int) int {
        fmt.Printf("(%+v, %+v) - ", &x, x)
        x += delta
        return x
    }
}
output:
(0xc420014070, 0) - 1
(0xc420014070, 1) - 11
(0xc420014070, 11) - 111

(0xc4200140b8, 0) - 1
(0xc4200140b8, 1) - 11
(0xc4200140b8, 11) - 111
```

闭包引用了 x 变量，a, b可看作 2个不同的实例，实例之间互不影响。实例内部，x 变量是同一个地址，因此具有“累加效应”。

## defer recover

panic会停掉当前正在执行的程序，不只是当前协程。在这之前，它会有序地执行完当前协程defer列表里的语句，其它协程里挂的defer语句不作保证。因此，我们经常在 defer 里挂一个 recover 语句，防止程序直接挂掉，这起到了  `try...catch` 的效果。

注意，recover() 函数只在defer的上下文中才有效（且只有通过在defer中用匿名函数调用才有效），直接调用的话，只会返回 `nil`.

```
func main() {
	defer fmt.Println("defer main")
	var user=os.Getenv
	
	go func() {
        defer func() { // 优雅处理协程的 panic
            fmt.Println("defer caller")
            if err := recover(); err != nil {
                fmt.Println("recover success. err: ", err)
            }
        }()
        
        func() {
            defer func() {
                fmt.Println("defer here")
            }()
            
            if user == "" {
                panic("should set user env.")
            }
            
            fmt.Println("after panic")
        }()
	}()
	
	time.Sleep(100)
	fmt.Println("end of main function")
}
```

上面的 panic 最终会被 recover 捕获到。这样的处理方式在一个 http server 的主流程常常会被用到。

一次偶然的请求可能会触发某个 bug，这时用 recover 捕获 panic，稳住主流程，不影响其他请求。

## 参考

 [饶全成@Golang之轻松化解defer的温柔陷阱](https://zhuanlan.zhihu.com/p/56557423)

