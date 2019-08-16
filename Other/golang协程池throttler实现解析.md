# golang协程池throttler实现解析

这次来介绍下,golang协程池的使用,以`throttler`实现为例。

首先介绍如何使用(拿作者github的例子为例)~

```
func ExampleThrottler() {
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}
	参数1:启动的协程数量
	参数2:需要执行的任务数
	t := New(2, len(urls))
	for _, url := range urls {

		// goroutine 启动
		go func(url string) {
			// 请求url
			err := http.Get(url)
			//让 throttler知道goroutines何时完成,然后throttler会新任命一个worker
		    t.Done(err)
		}(url)
		errorCount := t.Throttle()
		if errorCount > 0 {
			break
		}
	}
}
复制代码
```

虽然作者的`readme.md`没写,但是我们也可用这样用

```
package main

import (
	"github.com/throttler"
	"fmt"
)

func main() {

	p := throttler.New(10, 5)

	go func() {
		fmt.Println("hello world1")
		defer p.Done(nil)
	}()
	fmt.Println(1)
	p.Throttle()
	go func() {
		fmt.Println("hello world2")
		p.Done(nil)
	}()
	fmt.Println(2)
	p.Throttle()
	go func() {
		fmt.Println("hello world3")
		p.Done(nil)
	}()
	fmt.Println(3)
	p.Throttle()
	//fmt.Println(err + 3)
	go func() {
		fmt.Println("hello world4")
		p.Done(nil)
	}()
	fmt.Println(4)
	p.Throttle()
	//fmt.Println(err + 2)
	go func() {
		fmt.Println("hello world5")
		p.Done(nil)
	}()
	fmt.Println(5)
	p.Throttle()
}

复制代码
```

以上就是`Throttle`的使用例子,看起来非常简单,那么它是如何实现的呢？

首先我们看下`throttle`的主体结构,后续的操作都围绕着主体结构实现的

```
// Throttler stores all the information about the number of workers, the active workers and error information
type Throttler struct {
	maxWorkers    int32				// 最大的worker数
	workerCount   int32				// 正在工作的worker数量
	batchingTotal int32
	batchSize     int32				// 
	totalJobs     int32    			// 任务数量的和
	jobsStarted   int32  			// 任务开始的数量(初始值为0)
	jobsCompleted int32	 			// 任务完成的数量
	doneChan      chan struct{}		// 非缓冲队列,存储的一半是count(totalJobs)
	errsMutex     *sync.Mutex		// errMutex的并发
	errs          []error 			// 错误数组的集合,一般是业务处理返回的error
	errorCount    int32
}
复制代码
```

`New`操作创建一个协程池

```
func New(maxWorkers, totalJobs int) *Throttler {
	// 如果小于1 panic
	if maxWorkers < 1 {
		panic("maxWorkers has to be at least 1")
	}

	return &Throttler{
		// 最大协程数量
		maxWorkers: int32(maxWorkers),
		batchSize:  1,
		// 所有的任务数
		totalJobs:  int32(totalJobs),
		doneChan:   make(chan struct{}, totalJobs),
		errsMutex:  &sync.Mutex{},
	}
}
复制代码
```

当完成一个协程动作

```
func (t *Throttler) Done(err error) {
	if err != nil {
		// 如果出现错误,将错误追加到struct里面,因为struct非线程安全,所以需要加锁
		t.errsMutex.Lock()
		t.errs = append(t.errs, err)
		// errorCount ++
		atomic.AddInt32(&t.errorCount, 1)
		t.errsMutex.Unlock()
	} 
	// 每当一个goroutine进来,向struct写入一条数据
	t.doneChan <- struct{}{}
}
复制代码
```

等待协程完成的函数实现,可能稍微有点复杂

```
func (t *Throttler) Throttle() int {
	// 加载任务数  < 1 返回错误的数量
	if atomic.LoadInt32(&t.totalJobs) < 1 {
		return int(atomic.LoadInt32(&t.errorCount))
	}

	// jobStarted + 1 
	atomic.AddInt32(&t.jobsStarted, 1)
	// workerCount + 1
	atomic.AddInt32(&t.workerCount, 1)


	// 检查当前worker的数量是否和maxworker数量一致,等待这个workers完成

	// 实际上就是协程数量到达上限,需要等待运行中的协程释放资源
	if atomic.LoadInt32(&t.workerCount) == atomic.LoadInt32(&t.maxWorkers) {
		// 完成jobsCompleted - 1
		atomic.AddInt32(&t.jobsCompleted, 1)
		// workerCount - 1
		atomic.AddInt32(&t.workerCount, -1)
		<-t.doneChan
	}

	// check to see if all of the jobs have been started, and if so, wait until all
	// jobs have been completed before continuing

	// 如果任务开始的数量和总共的任务数一致
	if atomic.LoadInt32(&t.jobsStarted) == atomic.LoadInt32(&t.totalJobs) {
		// 如果完成的数量小于总job数 等待Job完成
		for atomic.LoadInt32(&t.jobsCompleted) < atomic.LoadInt32(&t.totalJobs) {
			// jobcomplete + 1
			atomic.AddInt32(&t.jobsCompleted, 1)
			<-t.doneChan
		}
	}

	return int(atomic.LoadInt32(&t.errorCount))
}
复制代码
```

简单枚举了下实现的流程:

假设有2个请求限制,3个请求,它的时序图是这样的

第一轮

```
totaljobs = 3
jobstarted = 1 workercount = 1   jobscompleted = 0 totaljobs = 3
复制代码
```

第二轮

```
jobstarted = 2 worker count = 2   jobscompleted = 0 totaljobs = 3
复制代码
```

第三轮

```
jobstarted = 3 worker count = 3 jobscompleted = 0 totaljobs = 3

// 操作1:因为goroutine限制为2,当前wokercount为3,需要阻塞,等待协程池释放

// 协程池释放:
jobstarted = 3 worker count = 2 jobscompleted = 1 totaljobs = 3

// 操作2:当前jobstarted与totaljobs相等,说明所有任务都已经池化了,则开始阻塞处理


//执行结束:

jobstarted = 3 worker count = 2 jobscompleted = 3 totaljobs = 3

复制代码
```

总的来说,该实现也是借用了`channel`的能力进行阻塞,实现起来还是非常简单的~

作者：WenruiGe

链接：https://juejin.im/post/5d4f9f396fb9a06b0703af47

来源：掘金

著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。