# govendor 管理你的 go 项目包版本

## Why?

>   例1：项目依赖一个 `github.com/foo 1.0.0` 的包，如果不使用包版本管理工具，他人在本地部署安装你的项目时，安装的包版本可能是最新的 `github.com/foo 2.0.0` ，如果两个版本存在兼容问题，就会出现 `crashed` 。

>   例2：使用 `go get` 安装的项目依赖包的存放位置为 `$GOPATH/src`，即与你的项目路径同级，我们使用 `git` 时就没办法管理这些依赖包，不能也不应该将它们也提交到`git`仓库，如果提供一个包版本说明说，将说明书提交仓库，他人则根据此说明书安装依赖包。

govendor 解决：

-   避免出现他人使用项目时出现版本管理问题。
-   方便 git 管理依赖包

`govendor` 会将项目依赖的包版本记录到  `your_proj/vendor/vendor.json`  中，后期将此文件提交至   `git`   仓库，则其他人  `pull or clone`  你的项目到本地后，即可使用  `govendor`  安装稳定的包依赖。

## govendor 工作流

### Major

如何使用`govendor`管理你的项目包版本，及如何安装他人的`govendor`管理的`go`项目

### 项目初始化

使用  `govendor`  初始化你的项目，将会在工程目录下自动创建 `vendor` 目录及 `vendor/vendor.json` 文件。

>   如果是已有项目，也没关系，`govendor`允许你在项目开发的任何阶段去使用它，它总能将你的项目包版本管理起来。

```
mkdir go_proj && cd go_proj
# init proj
govendor init
#查看目录结构
tree
.
└── vendor
    └── vendor.json

1 directory, 1 file
```

### 安装 & 管理包

#### `govendor get`

通过 `go get`安装包到 `$GOPATH/src`下，但使用 `govendor get` 可以在安装包的同时将包纳入版本管理，而且会将包安装在`$GOPATH/src/your_proj/vendor`，更符合我们的要求。

```
#安装在 $GOPATH/src 下
go get github.com/go-sql-driver/mysql

#安装在$GOPATH/src/your_proj/vendor下
govendor get github.com/go-sql-driver/mysql
```

#### `govendor list & govendor add`

`govendor list ` 可以帮助查看项目中引入的包的状态，即哪些是 没有纳入版本管理的外部包，哪些是 纳入版本管理的包，哪些是 标准包，哪些是 本地包 等。

```
govendor list

Status Types

        +local    (l) packages in your project 你自己在项目中定义的包
        +external (e) referenced packages in GOPATH but not in current project 使用 go get 安装的项目外部包
        +vendor   (v) packages in the vendor folder 使用 govendor get 安装的纳入版本管理的包
        +std      (s) packages in the standard library 标准包 fmt/time/runtime 等
        
        +excluded (x) external packages explicitly excluded from vendoring 排除的外部包
        +unused   (u) packages in the vendor folder, but unused 安装但没引用的包
        +missing  (m) referenced packages but not found 引用但没安装的包 缺失了

        +program  (p) package is a main package 你的项目主包,它总会同 l 一起出现 pl 这个很好理解吧

        +outside  +external +missing
        +all      +all packages
```

`govendor add`  则是方便我们在任何时间将项目包纳入版本管理。

比如我们前期一直使用或现在偶然使用  `go get  `安装了一个项目的依赖包，此包是不会被记录在`vendor/vendor.json`  中的，即没有纳入版本管理，那该如何将其纳入呢？

```
go add +external
```

执行上方命令即可，这样项目依赖的包都纳入了版本管理。

### 提交git仓库

在提交源码至 `git` 仓库时，我们没有必要将依赖包源文件也一并提交至仓库，所以 `.gitignore` 的编排要加上如下规则：

```
# vi .gitignore
vendor/*
!vendor/vendor.json
```

即排除  `vendor`  下的除  `vendor/vendor.json`  外的所有文件（这些文件其实就是依赖包），将  `vendor/vendor.json`  提交至  `git` 仓库即可。

### 安装/部署 govendor 项目

当我们从`git`仓库下载好`govendor`管理的`golang`项目时，需要安装好项目的包依赖，才可以正常的运行程序，类似 `composer install`的作用，这里则是使用`govendor sync`。

这里使用我的一个 `govendor` 管理的[基于`Gin`的`MVC`简易框架](https://github.com/sqrtcat/easy-gin)给大家演示一下：

```
cd $GOPATH/src && git clone git@github.com:sqrtcat/easy-gin.git && cd easy-gin

govendor sync
```

运行程序即可，简单！

## 备注

整理自：https://segmentfault.com/a/1190000018528053#articleHeader2

Author：[**big_cat**](https://segmentfault.com/u/big_cat)