# 14丨什么是事务处理，如何使用COMMIT和ROLLBACK进行操作？

[TOC]

## 前情

在 MySQL 5.5 版本之前，默认的存储引擎是 MyISAM，在 5.5 版本之后默认存储引擎是 InnoDB。InnoDB 和 MyISAM 区别之一就是 InnoDB 支持事务，也可以说这是 InnoDB 取代 MyISAM 的重要原因。

什么是事务呢？

事务的英文是 transaction，从英文中你也能看出来它是进行一次处理的基本单元，要么完全执行，要么都不执行。

事务保证了一次处理的完整性，也保证了数据库中的数据一致性。

## 疑问

事务的特性是什么？如何理解它们？

如何对事务进行控制？控制的命令都有哪些？

为什么我们执行 COMMIT、ROLLBACK 这些命令的时候，有时会成功，有时会失败？

## 事务的特性：ACID

从事务的 4 个特性说起，这 4 个特性用英文字母来表达就是 ACID

-   A，也就是原子性（Atomicity）。原子的概念就是不可分割
-   C，就是一致性（Consistency）。当事务提交后，或者当事务发生回滚后，数据库的完整性约束不能被破坏。
-   I，就是隔离性（Isolation）。一个事务在提交之前，对其他事务都是不可见的。
-   D，指的是持久性（Durability）。事务提交之后对数据的修改是持久性的

在这四个特性中，原子性是基础，隔离性是手段，一致性是约束条件，而持久性是我们的目的。

## 事务的控制

MySQL，可以通过 SHOW ENGINES 命令来查看当前 MySQL 支持的存储引擎都有哪些，以及这些存储引擎是否支持事务。

在 MySQL 中，InnoDB 是支持事务的，而 MyISAM 存储引擎不支持事务。

事务的常用控制语句都有哪些：

-   START TRANSACTION 或者 BEGIN，作用是显式开启一个事务。
-   COMMIT：提交事务。当提交事务后，对数据库的修改是永久性的。
-   ROLLBACK 或者 ROLLBACK TO [SAVEPOINT]，意为回滚事务。意思是撤销正在进行的所有没有提交的修改，或者将事务回滚到某个保存点。
-   SAVEPOINT：在事务中创建保存点，方便后续针对保存点进行回滚。一个事务中可以存在多个保存点。
-   RELEASE SAVEPOINT：删除某个保存点。
-   SET TRANSACTION，设置事务的隔离级别。

使用事务有两种方式，分别为隐式事务和显式事务。隐式事务实际上就是自动提交，Oracle 默认不自动提交，需要手写 COMMIT 命令，而 MySQL 默认自动提交，当然我们可以配置 MySQL 的参数：

```

mysql> set autocommit =0;  //关闭自动提交


mysql> set autocommit =1;  //开启自动提交
```

在 MySQL 的默认状态下，认真观察下面几个事务最后的处理结果是什么

第一个

解析：name 设置为唯一，第二次插入张飞会产生错误，执行 ROLLBACK 相当于对事务进行回滚。所以只能看到第一次提交的事务。

```

CREATE TABLE test(name varchar(255), PRIMARY KEY (name)) ENGINE=InnoDB;
BEGIN;
INSERT INTO test SELECT '关羽';
COMMIT;
BEGIN;
INSERT INTO test SELECT '张飞';
INSERT INTO test SELECT '张飞';
ROLLBACK;
SELECT * FROM test;

RESULT：
name
----
关羽
```

第二个

这里 MySQL 默认插入张飞是两个事务，因为在 autocommit=1 的情况下，MySQL 会进行隐式事务，也就是自动提交，因此在进行第一次插入“张飞”后，数据表里就存在了两行数据，而第二次插入“张飞”就会报错。

最后在执行 ROLLBACK 时候，实际上事务已经自动提交，就没法进行回滚。

```

CREATE TABLE test(name varchar(255), PRIMARY KEY (name)) ENGINE=InnoDB;
BEGIN;
INSERT INTO test SELECT '关羽';
COMMIT;
INSERT INTO test SELECT '张飞'; // 隐式COMMIT
INSERT INTO test SELECT '张飞';
ROLLBACK;
SELECT * FROM test;

RESULT：
name
----
关羽
张飞
```

第三个

MySQL 中 completion_type 参数的作用，实际上这个参数有 3 种可能：

-   completion=0，这是默认情况。也就是说当我们执行 COMMIT 的时候会提交事务，在执行下一个事务时，还需要我们使用 START TRANSACTION 或者 BEGIN 来开启。
-   completion=1，这种情况下，当我们提交事务后，相当于执行了 COMMIT AND CHAIN，也就是开启一个链式事务，即当我们提交事务之后会开启一个相同隔离级别的事务（隔离级别会在下一节中进行介绍）。
-   completion=2，这种情况下 COMMIT=COMMIT AND RELEASE，也就是当我们提交后，会自动与服务器断开连接。

回来下面的代码

这里使用了 completion=1，也就是说当我提交之后，相当于在下一行写了一个 START TRANSACTION 或 BEGIN。这时两次插入“张飞”会被认为是在同一个事务之内的操作，那么第二次插入“张飞”就会导致事务失败，而回滚也将这次事务进行了撤销，所以你能看到的结果就只有一个“关羽”。

```

CREATE TABLE test(name varchar(255), PRIMARY KEY (name)) ENGINE=InnoDB;
SET @@completion_type = 1;
BEGIN;
INSERT INTO test SELECT '关羽';
COMMIT;
INSERT INTO test SELECT '张飞'; // 这两句被当成同一个事务
INSERT INTO test SELECT '张飞';
ROLLBACK;
SELECT * FROM test;

RESULT：
name
----
关羽
```

当我们设置 autocommit=0 时，不论是否采用 START TRANSACTION 或者 BEGIN 的方式来开启事务，都需要用 COMMIT 进行提交，让事务生效，使用 ROLLBACK 对事务进行回滚。

当我们设置 autocommit=1 时，每条 SQL 语句都会自动进行提交。

不过这时，如果你采用 START TRANSACTION 或者 BEGIN 的方式来显式地开启事务，那么这个事务只有在 COMMIT 时才会生效，在 ROLLBACK 时才会回滚。

## 总结

正是因为有事务的存在，即使在数据库操作失败的情况下，也能保证数据的一致性。同样，多个应用程序访问数据库的时候，事务可以提供隔离，保证事务之间不被干扰。最后，事务一旦提交，结果就会是永久性的，这就意味着，即使系统崩溃了，数据库也可以对数据进行恢复。

事务是数据库区别于文件系统的重要特性之一，当我们有了事务就会让数据库始终保持一致性，同时我们还能通过事务的机制恢复到某个时间点，这样可以保证已提交到数据库的修改不会因为系统崩溃而丢失。


