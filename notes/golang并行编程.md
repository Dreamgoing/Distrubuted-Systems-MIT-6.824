# Golang 并行编程

> 介绍golang使用goroutine, 以及几种常用的并行模式.

## 概述
#### 并行!=并发 `Concurrency` not equal `Parallelism`

+ 并行: 同时处理多个事情, 并行执行.
+ 并发: 多个事情同时发生, 不一定同时处理, 可以轮询处理.

## 重要的点

#### 编程思想
+ 在生产和开发环境中,使用`race detector`, 来检测竞争条件. (有多个`goroutine`同时写某个共享的变量)
+ 使用`goroutine`运行独立的`golang代码片段`. `goroutine`更重要的是独立性.
+ 最好在单个的`goroutine`使用`mutex`, 这样会使得代码更加简单清晰.
+ 最好将变量的状态转化为代码的状态. (减少对`mutex`的使用, 即使用`channel`来实现对状态的转换)
+ 可以另开`goroutine`来保存程序中的状态转换主要是通过`channel`来实现.

#### 考虑
+ 考虑`slow goroutine`可能产生的影响
+ 考虑每个`goroutine`的运行和结束
+ 使用HTTP服务的`/debug/pprof/goroutine`来检查活跃`goroutine`的栈情况

#### 使用`channel`
+ 使用`buffered channel`来作为并行`goroutine`的队列.
+ 使用`sync channel`即`unbounded channel`会阻塞`goroutine`, 需要考虑何时阻塞.
+ 使用`close`来关闭`channel`, 即表明不再有更多的值放入`channel`中.

#### 并行编程
+ 使用`mutex`,`goroutine`, `channel`尽可能的简单, 清晰.

## 编程模式

#### 1. 将数据的状态转化为程序状态
> 可以理解为使用`goroutine`和`channel`来并行程序, `channel`可以对`goroutine`进行控制.

对于正则匹配`"([^"//]|//.)*"` 有限状态机,如下为一般思路下的代码. 代码中没有突出状态的转换.

``` go
type quoteReader struct {
 	state int
}
func (q *quoteReader) Init() {
 	q.state = 0
}
func (q *quoteReader) ProcessChar(c rune) Status {
 	switch q.state {
 	case 0:
 	if c != '"' {
 		return BadInput
 	}
 	q.state = 1
 	case 1:
 	if c == '"' {
 		return Success
 	}
 	if c == '\\' {
 		q.state = 2
 	} else {
 		q.state = 1
 	}
 	case 2:
 	q.state = 1
 	}
 	return NeedMoreInput
}

```

如下代码, 将数据的状态转换体现在了代码中, 使用额外其他的`goroutine`和`channel`
来控制状态转换.

其代码核心使用`goroutine`来读取`char`, 使用`channel`来控制状态.

``` go
func readString(readChar func() rune) bool {
	// state==0
	if readChar() != '"' { 
	 return false
	}
	var c rune
	for c != '"' {
		// state==1
		c := readChar()
	 	if c == '\\' {
	 	// state==2
 			readChar()
 		}
	}
	return true
}
type quoteReader struct {
	char chan rune
	status chan Status
}
func (q *quoteReader) Init() {
 	q.char = make(chan rune)
	q.status = make(chan Status)
 	go readString(q.readChar)
 	<-q.status // always NeedMoreInput
}
func (q *quoteReader) ProcessChar(c rune) Status {
 	q.char <- c
 	return <-q.status
}
func (q *quoteReader) readChar() int {
 	q.status <- NeedMoreInput
 	return <-q.char
}
```
#### 2. 发布者订阅者模式(Publish/subscribe server)

``` go
type PubSub interface {
	// Publish publishes the event e to
	// all current subscriptions.
	Publish(e Event)
	// Subscribe registers c to receive future events.
	// All subscribers receive events in the same order,
	// and that order respects program order:
	// if Publish(e1) happens before Publish(e2),
	// subscribers receive e1 before e2.
	Subscribe(c chan<- Event)
	// Cancel cancels the prior subscription of channel c.
	// After any pending already-published events
	// have been sent on c, the server will signal that the
	// subscription is cancelled by closing c.
	Cancel(c chan<- Event)
}
```



#### 3. `for-range channel` => `for-select-exit`
> 优点在于通过select监听多个`channel`可以实现对`servers`的监听控制, 比如退出, 超时等等.

``` go
go func() {
 	for _, srv := range servers {
 		go runTasks(srv)
 	}
}()

go func() {
	for {
		select {
		case srv := <-servers:
			go runTasks(srv)
		case <-exit:
			return
		case <-timer.C:
			//timeout
		}
	}
}()
```

#### 4. 使用`mu.Lock defer mu.Lock`
> 推荐在独立的`goroutine`中使用`mutex`.

``` go
go func(server *Server){
	server.mu.Lock()
	defer server.mu.Unlock()
	//do sometings
	
}
```


## Reference

[go pattern](https://pdos.csail.mit.edu/6.824/notes/gopattern.pdf)