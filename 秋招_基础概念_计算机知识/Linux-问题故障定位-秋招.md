# Linux 问题故障定位

目录

[TOC]

## 前言

套用5W2H方法，可以提出性能分析的几个问题

-   What-现象是什么样的
-   When-什么时候发生
-   Why-为什么会发生
-   Where-哪个地方发生的问题
-   How much-耗费了多少资源
-   How to do-怎么解决问题

## CPU

### 说明

针对应用程序，我们通常关注的是内核CPU调度器功能和性能。

线程的状态分析主要是分析线程的时间用在什么地方，而线程状态的分类一般分为：

-   on-CPU：执行中，执行中的时间通常又分为用户态时间user和系统态时间sys；
-   off-CPU：等待下一轮上CPU，或者等待I/O、锁、换页等等，其状态可以细分为可执行、匿名换页、睡眠、锁、空闲等状态；

如果大量时间花在CPU上，对CPU的剖析能够迅速解释原因；如果系统时间大量处于off-cpu状态，定位问题就会费时很多；

### 分析工具

| 工具   | 描述                          |
| ------ | ----------------------------- |
| uptime | 平均负载                      |
| vmstat | 包括系统范围内的 CPU 平均负载 |
| mpstat | 查看所有 CPU 核信息           |
| top    | 监控每个进程 CPU 用量         |
| sar -u | 查看 CPU 信息 |
| pidstat | 每个进程  CPU 用量分解 |
| pref | CPU 剖析和跟踪，性能计数分析 |

说明:

-   uptime, vmstat, mpstat, top, pidstat只能查询到 cpu 及负载的的使用情况。
-   perf可以跟着到进程内部具体函数耗时情况，并且可以指定内核函数进行统计，指哪打哪。

### 使用方式

```
//查看系统cpu使用情况
top

//查看所有cpu核信息
mpstat -P ALL 1

//查看cpu使用情况以及平均负载
vmstat 1

//进程cpu的统计信息
pidstat -u 1 -p pid

//跟踪进程内部函数级cpu使用情况
perf top -p pid -e cpu-clock
```

## 内存

### 说明

内存出现问题可能不只是影响性能，而是影响服务或者引起其他问题。

### 分析工具

| 工具     | 描述                               |
| -------- | ---------------------------------- |
| free     | 缓存容量统计                       |
| vmstat   | 虚拟内存统计信息                   |
| top      | 监视每个进程的内存使用情况         |
| pidstat  | 显示活动进程的内存使用统计         |
| pmap     | 查看进程内存映像信息               |
| sar -r   | 查看内存                           |
| dtrace   | 动态跟踪                           |
| valgrind | 分析程序性能以及程序中内存泄漏错误 |

说明：

-   free, vmstat, top, pidstat, pmap 只能统计内存信息以及进程的内存使用情况；
-   valgrind可以分析内存泄漏问题；
-   dtrace动态跟踪。需要对内核函数有很深入的了解，通过D语言编写脚本完成跟踪；

### 使用方式

```
//查看系统内存使用情况
free -m

//虚拟内存统计信息
vmstat 1

//查看系统内存情况
top

//1s采集周期，获取内存的统计信息
pidstat -p pid -r 1

//查看进程的内存映像信息
pmap -d pid

//检测程序内存问题
valgrind --tool=memcheck --leak-check=full --log-file=./log.txt  ./程序名
```

## 磁盘IO

### 说明

要监测 IO 性能，有必要了解一下基本原理和 Linux 是如何处理硬盘和内存之间的 IO 的。

### 分析工具

| 工具    | 描述                         |
| ------- | ---------------------------- |
| iostat  | 磁盘详细统计信息             |
| iotop   | 按进程查看磁盘 IO 的使用情况 |
| pidstat | 按进程查看磁盘 IO 的使用情况 |
| perf    | 动态跟踪工具 |

### 使用方式

```
//查看系统io信息
iotop

//统计 io 详细信息
iostat -d -x -k 1 10

//查看进程级 io 的信息
pidstat -d 1 -p  pid

//查看系统 IO 的请求，比如可以在发现系统IO异常时，可以使用该命令进行调查，就能指定到底是什么原因导致的IO异常
perf record -e block:block_rq_issue -ag
^C
perf report
```

## 网络

### 说明

网络的监测是所有 Linux 子系统里面最复杂的，有太多的因素在里面

### 分析工具

| 工具        | 描述                                                         |
| ----------- | ------------------------------------------------------------ |
| ping        | 主要通过 ICMP 封包来打印整个网络的状况报告                   |
| treceroute  | 检测发出数据包的主机到目标主机之间所经过的网络数量           |
| netstat     | 用于显示与 IP、TCP、UDP 和 ICMP 协议相关的统计数据，一般用于校验本机各端口的网络连接情况 |
| ss          | 获取 Socket 统计信息，比 netstat 更快更高效                  |
| host        | 查出某个主机名的 IP，跟 nslookup 作用一样                    |
| tcpdump     | 以包位单位进行输出，阅读起来不方便                           |
| tcpflow     | 面向 TCP 流，每个 TCP 传输保存成一个文件，很方便查看         |
| sar -n DEV  | 网卡流量情况                                                 |
| sar -n SOCK | 查询网络以及 TCP、UDP 状态信息                               |

### 使用方式

```
//显示网络统计信息
netstat -s

//显示当前UDP连接状况
netstat -nu

//显示UDP端口号的使用情况
netstat -apu

//统计机器中网络连接各个状态个数
netstat -a | awk '/^tcp/ {++S[$NF]} END {for(a in S) print a, S[a]}'

//显示TCP连接
ss -t -a

//显示sockets摘要信息
ss -s

//显示所有udp sockets
ss -u -a

//tcp,etcp状态
sar -n TCP,ETCP 1

//查看网络IO
sar -n DEV 1

//抓包以包为单位进行输出
tcpdump -i eth1 host 192.168.1.1 and port 80 

//抓包以流为单位显示数据内容
tcpflow -cp host 192.168.1.1
```

## 系统负载

### 说明

Load Average 就是一段时间（1分钟、5分钟、15分钟）内平均Load。

### 分析工具

| 工具   | 描述               |
| ------ | ------------------ |
| top    | 查看系统负载情况   |
| uptime | 查看系统负载情况   |
| strace | 统计跟踪内核态信息 |
| vmstat | 查看负载情况       |
| dmesg  | 查看内核日志信息   |

### 使用方式

```
//查看负载情况
uptime

top

vmstat

//统计系统调用耗时情况
strace -c -p pid

//跟踪指定的系统操作例如epoll_wait
strace -T -e epoll_wait -p pid

//查看内核日志信息
dmesg
```

## 火焰图

略

## 参考资料

Lucien_168@https://www.jianshu.com/p/0bbac570fa4c