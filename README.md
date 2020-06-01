# go pprof

```
Go语言工具链中的 go pprof 可以帮助开发者快速分析及定位各种性能问题，如 CPU消耗、内存分配及阻塞分析
```

# Go语言项目中的性能优化主要有以下几个方面
+ CPU profile: 报告程序的CPU使用情况，按照一定频率去采集应用程序在CPU寄存器上面的数据
+ Memory profile: 报告程序的内存使用情况
+ Block profile: 报告 goroutines 不在运行状态的情况，可以用来分析和查找死锁等性能瓶颈
+ Goroutine profile: 报告 goroutines 的使用情况，有哪些 goroutine，它们的调用关系是怎样的

------------------------------------------------------------

## CPU性能分析

+ 在cpu文件夹下执行如下命令，加 -cpu ture 是为了开启CPU profile分析，生成cpu.pprof文件
```bash
go run main.go -cpu true
```
+ 在cpu文件夹下执行如下命令，加 -cpu ture 是为了开启CPU profile分析，生成cpu.pprof文件
```bash
go run main.go -cpu true

```
+ 对生成cpu.pprof文件进行分析，执行如下命令
```bash
go tool pprof cpu.pprof
```

+ 进入交互交互界面，输入top，获取分析结果
```bash
Type: cpu
Time: May 24, 2020 at 2:02pm (CST)
Duration: 607.77ms, Total samples = 520ms (85.56%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 520ms, 100% of 520ms total
Showing top 10 nodes out of 44
      flat  flat%   sum%        cum   cum%
     240ms 46.15% 46.15%      240ms 46.15%  runtime.memmove
     100ms 19.23% 65.38%      100ms 19.23%  runtime.usleep
      60ms 11.54% 76.92%       60ms 11.54%  runtime.memclrNoHeapPointers
      60ms 11.54% 88.46%       60ms 11.54%  runtime.nanotime
      20ms  3.85% 92.31%      320ms 61.54%  main.joinSlice
      20ms  3.85% 96.15%       20ms  3.85%  runtime.madvise
      10ms  1.92% 98.08%       10ms  1.92%  runtime.procyield
      10ms  1.92%   100%       10ms  1.92%  runtime.pthread_cond_wait
         0     0%   100%      320ms 61.54%  main.main
         0     0%   100%       20ms  3.85%  runtime.(*mheap).alloc

```
```
其中参数说明：
flat: 当前函数占用CPU的耗时
flat%: 当前函数占用CPU的耗时百分比
sum%: 函数占用CPU的耗时累计百分比
cum：当前函数加上调用当前函数的函数占用CPU的总耗时
cum%：当前函数加上调用当前函数的函数占用CPU的总耗时百分比
最后一列：函数名称
```

+ 现在对耗时长的函数详细分析，执行命令 list joinSlice
```bash
(pprof) list joinSlice
Total: 520ms
ROUTINE ======================== main.joinSlice in /Users/wutianxiang/uniontech/golang_pprof/cpu/main.go
      20ms      320ms (flat, cum) 61.54% of Total
         .          .     30:func joinSlice() []int {
         .          .     31:   var slices []int
         .          .     32:   num := 10000000
         .          .     33:   for i := 0; i < num; i++ {
         .          .     34:           // 故意造成多次切片添加，由于每次操作可能会有内存重新分配和移动，性能较低
         .      300ms     35:           slices = append(slices, i+1)
         .          .     36:   }
      20ms       20ms     37:   return slices
         .          .     38:}
(pprof) 

```

+ 通过上面分析发现大部分CPU资源被35行占用，我们分析出slices在定义时，没有提前给容量，每次操作可能内存分配和移动

+ 如果将joinSlice函数代码重新优化，那么如何优化？可以对已知切片元素数量的情况下直接分配内存，代码如下
```
func joinSlice() []int {
	num := 10000000
	slices := make([]int, 0, num)
	for i := 0; i < num; i++ {
		slices = append(slices, i+1)
	}
	return slices
}
```


