# 目标

1. 如何通过pprof查看当前程序的阻塞中的Goroutines数量，通过比对Goroutines数量判定是否有 Goroutines内存泄漏的情况
2. 如何查看一个请求的路径下，各个执行语句部分内存、CPU消耗的占比
3. Gin如何开启pprof，以及如何采集线上的执行指标信息
4. 如何对一个函数获取pprof(bench mark与pprof)


### 如何加载线上的pprof指标
```
#加载pprof文件然后分析
go tool pprof 127.0.0.1:8080/debug/pprof/profile

#加载pprof文件然后启用web服务器查看
go tool pprof -http=:9999 127.0.0.1:8080/debug/pprof/profile

#加载内存信息
go tool pprof 127.0.0.1:8080/debug/pprof/allocs
```

### 官方库提供的指标有哪些
```
"allocs":       "所有过去内存分配的采样"
"block":        "导致同步原语阻塞的堆栈跟踪"
"cmdline":      "当前程序的命令行调用"
"goroutine":    "所有当前goroutine的堆栈跟踪"
"heap":         "活动对象的内存分配采样。在获取堆样本之前，可以指定gcGET参数来运行gc"
"mutex":        "争用互斥锁持有者的堆栈跟踪"
"profile":      "CPU配置文件。可以在秒GET参数中指定持续时间。获取配置文件后，使用go工具pprof命令调查配置文件"
"threadcreate": "导致创建新操作系统线程的堆栈跟踪"
"trace":        "当前程序的执行轨迹。可以在秒GET参数中指定持续时间。获取跟踪文件后，使用go工具跟踪命令来调查跟踪"
```

### 命令行工具如何使用

1. top 查看占比
2. list 罗列一个具体函数的执行比如 {list main.main}
3. traces 执行路径
4. png 本地生成一个图片查看

### top的指标说明

``` 
flat	本函数的执行耗时
flat%	flat 占 CPU 总时间的比例。程序总耗时 16.22s, Eat 的 16.19s 占了 99.82%
sum%	前面每一行的 flat 占比总和
cum	    累计量。指该函数加上该函数调用的函数总耗时
cum%	cum 占 CPU 总时间的比例
```

### 参考博客

[Go 语言高性能编程pprof](https://geektutu.com/post/hpg-pprof.html)
[gin-contrib将pprof注入Gin路由工具](https://github.com/gin-contrib/pprof)
[http://127.0.0.1:8080/debug/pprof/](http://127.0.0.1:8080/debug/pprof/)
[深度解密Go语言之 pprof](https://www.cnblogs.com/qcrao-2018/p/11832732.html)
[实战Go内存泄露](https://segmentfault.com/a/1190000019222661)
[Go常用包(二十九):性能调试利器使用 - 如何分析一个函数执行的详细](http://liuqh.icu/2021/11/15/go/package/29-pprof-1/)

### 一个内存泄漏的例子

```
tick := time.Tick(time.Second / 100)
var buf []byte
for range tick {
    buf = append(buf, make([]byte, 1024*1024)...)
}
```

### 超大内存占用的例子
```
ch := make(chan bool)
go func() {
    var stringSlice []string
    for i := 0; i < 20; i++ {
        repeat := strings.Repeat("hello,world", 50000)
        stringSlice = append(stringSlice, repeat)
        time.Sleep(time.Millisecond * 500)
    }
    ch <- true
}()
<-ch
```

### bench
```
go test -bench="Fib$" -cpuprofile=cpu.pprof .
go tool pprof -text cpu.pprof
```

### 执行日志
```
➜  Documents go tool pprof 127.0.0.1:8080/debug/pprof/heap
Fetching profile over HTTP from http://127.0.0.1:8080/debug/pprof/heap
Saved profile in /Users/xuweiqiang/pprof/pprof.alloc_objects.alloc_space.inuse_objects.inuse_space.013.pb.gz
Type: inuse_space
Time: Dec 6, 2022 at 1:56pm (CST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) list GetUserList
Total: 9.43MB
ROUTINE ======================== /http/match.go
         0     4.06MB (flat, cum) 43.04% of Total
         .          .     45:	ch := make(chan bool)
         .          .     46:	go func() {
         .          .     47:		var stringSlice []string
         .          .     48:		for i := 0; i < 20; i++ {
         .          .     49:			// pprof 显示在这里占用2MB的内存开销
         .     4.06MB     50:			repeat := strings.Repeat("hello,world", 50000)
         .          .     51:			stringSlice = append(stringSlice, repeat)
         .          .     52:			time.Sleep(time.Millisecond * 500)
         .          .     53:		}
         .          .     54:		ch <- true
         .          .     55:	}()
(pprof) 
```