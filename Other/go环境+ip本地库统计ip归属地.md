# go环境+ip本地库统计ip归属地

来源：<https://blog.51cto.com/wujianwei/2122667>

# 一、服务器环境介绍：

dell服务器PowerEdge R420

```
[root@git-server ~]# cat /etc/redhat-release 
CentOS release 6.9 (Final)
[root@git-server ~]# uname -a
Linux git-server 2.6.32-696.23.1.el6.x86_64 #1 SMP Tue Mar 13 22:44:18 UTC 2018 x86_64 x86_64 x86_64 GNU/Linux
```

系统是最小化安装
由于我的系统是CentOS6.9 x89_64位的，所以下载安装64位的go tar包

# 二、安装go环境：

参考文档：
<https://www.cnblogs.com/1111zhiping-tian/p/8086736.html>

安装过程：

```
cd /root
wget https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz
tar xf go1.9.2.linux-amd64.tar.gz
mkdir -p workspace/src
```

go环境变量设置：
[root@git-server ~]# cat ~/.bashrc

```
#.bashrc

#User specific aliases and functions

alias rm='rm -i'
alias cp='cp -i'
alias mv='mv -i'

#Source global definitions
if [ -f /etc/bashrc ]; then
    . /etc/bashrc
fi
export GOROOT=$HOME/go
export GOPATH=$HOME/workspace
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
```

source ~/.bashrc

```
[root@git-server ~]# go version
go version go1.9.2 linux/amd64
```

到此go环境安装成功

# 二、下载开源的本地ip库文件以及查询ip归属地：

**以下的操作过程要严格按照以下的顺序执行**

```
下载main.go文件：
https://github.com/miraclesu/ipip
wget https://github.com/miraclesu/ipip/archive/master.zip
unzip matser.zip

[root@git-server ipip-master]# ls
main.go  README.md
```

[root@git-server ipip-master]# pwd
/root/ipip-master

**main.go文件脚本内容：**

```
[root@localhost ipip-master]# cat /root/ipip-master/main.go
package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "os"
    "strings"
    "github.com/wangtuanjie/ip17mon"
)

var (
    path string
    file string

    batch = 100
)

func init() {
    flag.StringVar(&path, "path", "17monipdb.dat", "/root/ipip-master") ###17monipdb 数据文件路径
    flag.StringVar(&file, "ips", "ips.txt", "/root/ipip-master") ###要查询 ip 数据文件路径
}

func main() {
    flag.Parse()

    f, err := os.Open(file)
    if err != nil {
        fmt.Printf("打开日志[%s]出错: %s\n", file, err.Error())
        return
    }

    if err = ip17mon.Init(path); err != nil {
        fmt.Printf("初始化 ip 库失败: %s\n", err.Error())
        return
    }

    buf, output := bufio.NewReader(f), make([]string, 0, batch)
    for {
        ip, err := buf.ReadString('\n')
        if err != nil {
            if err != io.EOF {
                fmt.Printf("读取日志[%s]出错: %s\n", file, err.Error())
                return
            }
            break
        }

        ip = strings.TrimSpace(ip)
        loc, err := ip17mon.Find(ip)
        if err != nil {
            output = append(output, fmt.Sprintf("%s,%s,%s,%s,%s", ip, err.Error(), "", "", ""))
            continue
        }

        output = append(output, fmt.Sprintf("%s,%s,%s,%s,%s", ip, loc.Country, loc.Region, loc.City, loc.Isp))
        if len(output) < batch {
            continue
        }

        fmt.Println(strings.Join(output, "\n"))
        output = output[:0]
    }

    if len(output) > 0 {
        fmt.Println(strings.Join(output, "\n"))
    }
}
[root@localhost ipip-master]# 
```

**在数据文件路径下执行如下命令：**

```
go get -u github.com/wangtuanjie/ip17mon
[root@git-server ipip-master]# pwd
/root/ipip-master
[root@git-server ipip-master]# go get -u github.com/wangtuanjie/ip17mon
[root@git-server ipip-master]# cd /root/workspace/src/github.com/wangtuanjie/ip17mon
[root@git-server ip17mon]# ls
17monipdb.dat  circle.yml  example  ip17mon.go  ip17mon_test.go  LICENSE  Makefile  README.md  tools
[root@git-server ipip-master]# ls
main.go  README.md
```

**下载datIP本地库文件库:**

```
[root@git-server ipip-master]# wget https://github.com/wangtuanjie/ip17mon/blob/master/17monipdb.dat
[root@git-server ipip-master]# ls
17monipdb.dat  main.go  README.md
```

**上传存放IP的文件文件，要求每行一个单独的ip地址，上传一个存放1000多万的唯一的ip地址文件**

```
cd /root/ipip-master
[root@git-server ipip-master]# cat /root/ipip-master/unip_guishu |wc -l
10178478
```

**然后查询IP归属地址**

执行一下命令：

```
[root@git-server ipip-master]# time go run main.go >newfile

real    0m34.420s
user    0m33.978s
sys 0m2.100s

[root@cacti ipip]# tail -10 newfile 
99.97.22.46,美国,美国,N/A,N/A
99.97.26.40,美国,美国,N/A,N/A
99.97.29.64,美国,美国,N/A,N/A
99.97.30.121,美国,美国,N/A,N/A
99.98.238.52,美国,美国,N/A,N/A
99.98.254.160,美国,美国,N/A,N/A
99.99.218.167,美国,美国,N/A,N/A
99.99.232.239,美国,美国,N/A,N/A
99.99.37.198,美国,美国,N/A,N/A
99.99.55.172,美国,美国,N/A,N/A
```

**到此处已经完成**
小插曲：
最开始业务需求是这样的：部门领导让我统计业务机器一天日志中所有访问ip的数量，而且要求是查询出来IP的归属地，并且对ip的归属地进行大陆和国外以及港澳台的占比计算。一开始我先写shell脚本批量获取所有业务机器上的日志中的ip，并且排序去重后，发现ip量有1000多万条。后面我有采用shell脚本来一个ip一个ip的来查询，本以为很快地，结果shell脚本执行了24小时了，结果才获取到不到10万个不通ip的归属地。这时，我犯怵了，这样下去，这查到什么时候。网上各种找资料，请教社区网友，而然一个论坛网友凯强回复了我，而且提供了我方法，在此非常的感谢他。写此博文，希望能帮助更多的遇到类似问题的朋友们。