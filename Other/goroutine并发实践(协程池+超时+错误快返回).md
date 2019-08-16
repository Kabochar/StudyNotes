# goroutine并发实践(协程池+超时+错误快返回)

<https://juejin.im/post/5d54fbeef265da03af19cc5c>

当我们使用`goroutine`的时候让函数并发执行的时候,可以借助着`sync.WaitGroup{}`的能力,其中代码如下:

```
func testGoroutine() {
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			fmt.Println("hello world")
		}()
	}
	wg.Wait()
}
复制代码
```

看完上述代码此时我们需要考虑的是,假设`goroutine`因为一些`rpc`请求过慢导致`hang`住,此时`goroutine`会一直卡住在`wg.Wait()`,最终导致请求失败

除非你使用的框架提供了一个超时的能力,或者你`go`出去的`rpc`请求存在超时断开的能力

#### 那么我们如何让代码不被`hang`住呢?

最简单的解法就是增加超时!

实际上超时也有很多解法

-   基于`ctx` 的`context.WithTimeOut()`实现
-   基于`select`实现

这里我选择基于`select`实现超时来给大家看下代码如何实现

```
func testWithGoroutineTimeOut() {
	var wg sync.WaitGroup
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
		}()
	}
	// wg.Wait()此时也要go出去,防止在wg.Wait()出堵住
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	// 正常结束完成
	case <-done:
	// 超时	
	case <-time.After(500 * time.Millisecond):
	}
}
复制代码
```

可以看到上述代码,已经基于`select`实现了超时,是不是非常简单呢~

但是我们对于这个接口会有更高的要求。

-   `goroutine`没有错误处理,
-   此时`go`出去的`goroutine`数量是依赖`for`循环的数量,假设`for`循环100w次,造成`goroutine`过多的问题

可以写一个协程池解决`goroutine`过多,,那么协程池如何实现呢?

我们可以使用`sync waitGroup`+ 非阻塞`channel`实现 代码如下:

```
package ezgopool

import "sync"

// goroutine pool
type GoroutinePool struct {
	c  chan struct{}
	wg *sync.WaitGroup
}

// 采用有缓冲channel实现,当channel满的时候阻塞
func NewGoroutinePool(maxSize int) *GoroutinePool {
	if maxSize <= 0 {
		panic("max size too small")
	}
	return &GoroutinePool{
		c:  make(chan struct{}, maxSize),
		wg: new(sync.WaitGroup),
	}
}

// add
func (g *GoroutinePool) Add(delta int) {
	g.wg.Add(delta)
	for i := 0; i < delta; i++ {
		g.c <- struct{}{}
	}

}

// done
func (g *GoroutinePool) Done() {
	<-g.c
	g.wg.Done()
}

// wait
func (g *GoroutinePool) Wait() {
	g.wg.Wait()
}

复制代码
```

以上就是协程池的实现,实际是非常简单的,我的博客也记录了另一个`golang`协程池的开源实现,具体见 [juejin.im/post/5d4f9f…](https://juejin.im/post/5d4f9f396fb9a06b0703af47)

然后最后我们的超时+错误快返回+协程池模型就完成了~

```
func testGoroutineWithTimeOut() {
	 wg :=sync.WaitGroup{}
	done := make(chan struct{})
	// 新增阻塞chan
	errChan := make(chan error)

	pool.NewGoroutinePool(10)
	for i := 0; i < 10; i++ {
		pool.Add(1)
		go func() {
			pool.Done()
			if err!=nil{
				errChan<-errors.New("error")
			}
		}()
	}

	go func() {
		pool.Wait()
		close(done)
	}()

	select {
	// 错误快返回,适用于get接口
	case err := <-errChan:
		return nil, err
	case <-done:
	case <-time.After(500 * time.Millisecond):
	}
}

复制代码
```

谢谢



作者：WenruiGe

链接：https://juejin.im/post/5d54fbeef265da03af19cc5c

来源：掘金

著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。