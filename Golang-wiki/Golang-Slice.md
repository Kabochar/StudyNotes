# Golang Slice

目录

[TOC]

## 总结

-   切片是对底层数组的一个抽象，描述了它的一个片段；
-   切片实际上是一个结构体，它有三个字段：长度，容量，底层数据的地址；
-   多个切片可能共享同一个底层数组，这种情况下，对其中一个切片或者底层数组的更改，会影响到其他切片；
-   `append` 函数会在切片容量不够的情况下，调用 `growslice` 函数获取所需要的内存，这称为扩容，扩容会改变元素原来的位置；
-   扩容策略并不是简单的扩为原切片容量的 `2` 倍或 `1.25` 倍，还有内存对齐的操作。扩容后的容量 >= 原容量的 `2` 倍或 `1.25` 倍；
-   当直接用切片作为函数参数时，可以改变切片的元素，不能改变切片本身；想要改变切片本身，可以将改变后的切片返回，函数调用者接收改变后的切片或者将切片指针作为函数参数；

## 什么是 Slice？

上源码：

```
// runtime/slice.go
type slice struct {
    array unsafe.Pointer // 元素指针
    len   int // 长度 
    cap   int // 容量
}
```

主要三个内容：

-   array 指针，指向一个底层数组
-   len，slice 当前实际元素个数
-   cap，底层数组位置个数

slice 指向的底层数组可以被多个 slice 同时指向，因此 对某个 slice 的操作可以影响到其他 slice

## Slice 优雅创建姿势？

5 种方式 创建 slice

| 序号 | 方式               | 代码示例                                             |
| ---- | ------------------ | ---------------------------------------------------- |
| 1    | 直接声明           | `var slice []int`                                    |
| 2    | new                | `slice := *new([]int)`                               |
| 3    | 字面量             | `slice := []int{1,2,3,4,5}`                          |
| 4    | make               | `slice := make([]int, 5, 10)`                        |
| 5    | 从切片或数组“截取” | `slice := array[1:5]` 或 `slice := sourceSlice[1:5]` |

### 直接声明

注意：`nil slice` 和 `empty slice` 是不同的，`empty silce` 有底层的数组指针

| 创建方式      | nil切片              | 空切片                  |
| ------------- | -------------------- | ----------------------- |
| 方式一        | var s1 []int         | var s2 = []int{}        |
| 方式二        | var s4 = *new([]int) | var s3 = make([]int, 0) |
| 长度          | 0                    | 0                       |
| 容量          | 0                    | 0                       |
| 和 `nil` 比较 | `true`               | `false`                 |

### 字面量

```
func main() {
    s1 := []int{0, 1, 2, 3, 8: 100}
    fmt.Println(s1, len(s1), cap(s1))
}

output:
[0 1 2 3 0 0 0 0 100] 9 9
```

其他未注明的元素的初始值：对应类型零值；

### make

`make`函数需要传入三个参数：切片类型，长度，容量；容量可忽略，默认容量=长度；

秀操作

```
package main

import "fmt"

func main() {
    slice := make([]int, 5, 10) // 长度为5，容量为10
    slice[2] = 2 // 索引为2的元素赋值为2
    fmt.Println(slice)
}
```

编译查看

```
go tool compile -S main.go
```

具体编译内容忽略...

>   Go 语言汇编 `FUNCDATA` 和 `PCDATA` 是编译器产生的，用于保存一些和垃圾收集相关的信息，我们先不用 care；

关键几个函数


```
CALL    runtime.makeslice(SB) // 创建 slice
CALL    runtime.convT2Eslice(SB) // 类型转换
CALL    fmt.Println(SB) // 打印函数
CALL    runtime.morestack_noctxt(SB) // 栈空间扩容
```

### 截取

新老 slice 或者新 slice 老数组互相影响的前提是两者共用底层数组，如果因为执行 `append` 操作使得新 slice 底层数组扩容，移动到了新的位置，两者就不会相互影响了；

KEY：两者是否相互影响的问题关键在于两者是否会共用底层数组；

截取案例：

```
slice := data[2:4:6] // data[low, high, max]
```

切片相关操作

| **操作**            | **含义**                                                     |
| ------------------- | ------------------------------------------------------------ |
| **s[n]**            | 切片s中索引位置为 n 的项                                     |
| **s[:]**            | 从切片s的索引位置 0 到 len(s)-1 处所获得的切片               |
| **s[low:]**         | 从切片s的索引位置 low 到 len(s)-1 处所获得的切片             |
| **s[:high]**        | 从切片s的索引位置 0 到 high 处所获得的切片，len=high         |
| **s[low:high]**     | 从切片s的索引位置 low 到 high 处所获得的切片，len=high-low   |
| **s[low:high:max]** | 从切片s的索引位置 low 到 high 处所获得的切片，len=high-low，cap=max-low |
| **len(s)**          | 切片s的长度，总是 <=cap(s)                                   |
| **cap(s)**          | 切片s的容量，总是 >=len(s)                                   |

>   表格来源网络

细节：

-   max >= high >= low；
-   high == low 时，新 slice 为空；
-   high 和 max 必须在老数组或者老 slice 的容量（cap）范围内；

## slice VS array 恩怨情愁！

-   slice 的底层数据是数组，slice 是对数组的封装，它描述一个数组的片段；
-   两者都可以通过下标来访问单个元素；
-   数组是定长的，长度定义好之后，不能再更改；切片则非常灵活，它可以动态地扩容；
-   切片的类型和长度无关；

## append 你在干什么？

append 函数原型

```
func append(slice []Type, elems ...Type) []Type
```

append 操作，注意，必须有 变量 接收（slice变量是必要的）

```
slice = append(slice, elem1, elem2)     // 逐个添加操作
slice = append(slice, anotherSlice...) // 添加整个切片
```

为了应对未来的 append 操作，新的 slice 容量预留了一定的 buffer，但 buffer 的增长情况是怎样的呢？

-   在老 slice 容量小于1024的时候，新 slice 的容量的确是老 slice 的2倍；
-   后半部分还对 `newcap` 作了一个`内存对齐`，这个和内存分配策略相关。进行内存对齐之后，新 slice 的容量是要 `大于等于` 老 slice 容量的 `2倍`或者`1.25倍`；

## 为什么 nil slice 能和 append 好好相处？

`nil slice` or `empty slice`都是调用 `mallocgc` 来向 Go 的内存管理器申请到一块内存，然后再赋给原来的`nil slice` 或 `empty slice`，然后摇身一变，成为“真正”的 `slice` ；

## 传送 slice 和 slice 指针 有什么区别？

不管传的是 slice 还是 slice 指针，如果改变了 slice 底层数组的数据，会反应到实参 slice 的底层数据；

为什么能改变底层数组的数据？

底层数据在 slice 结构体里是一个指针，尽管 slice 结构体自身不会被改变，也就是说底层数据地址不会被改变；但是通过指向底层数据的指针，可以改变切片的底层数据，没有问题；

## 参考

饶全成@深度解密Go语言之Slice

<https://mp.weixin.qq.com/s?__biz=MjM5MDUwNTQwMQ==&mid=2257483743&idx=1&sn=af5059b90933bef5a7c9d491509d56d9&chksm=a5391709924e9e1f3839aef41d05820c52181ace22a43dc3177df98f9faa9edcfdfefe670d88&scene=27##>