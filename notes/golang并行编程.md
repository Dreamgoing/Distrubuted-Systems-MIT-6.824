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

## Reference

[go pattern](https://pdos.csail.mit.edu/6.824/notes/gopattern.pdf)