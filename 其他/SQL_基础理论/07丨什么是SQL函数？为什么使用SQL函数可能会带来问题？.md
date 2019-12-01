# 07丨什么是SQL函数？为什么使用SQL函数可能会带来问题？

[TOC]

## 什么是SQL函数

函数的作用是什么呢？它可以把我们经常使用的代码封装起来，需要的时候直接调用即可。这样既提高了代码效率，又提高了可维护性。

封装在 SQL 中的内部的函数

## 常用的SQL函数有哪些

### 算术函数

![1575203408590](D:\Documents\笔记本\offer学习复习\其他\SQL_基础理论\1575203408590.png)

案例

SELECT ABS(-2)，运行结果为 2。

SELECT MOD(101,3)，运行结果 2。

SELECT ROUND(37.25,1)，运行结果 37.3。

### 字符串函数

![1575203469895](D:\Documents\笔记本\offer学习复习\其他\SQL_基础理论\1575203469895.png)

案例

SELECT CONCAT('abc', 123)，运行结果为 abc123。

SELECT LENGTH('你好')，运行结果为 6。

SELECT CHAR_LENGTH('你好')，运行结果为 2。

SELECT LOWER('ABC')，运行结果为 abc。

SELECT UPPER('abc')，运行结果 ABC。

SELECT REPLACE('fabcd', 'abc', 123)，运行结果为 f123d。

SELECT SUBSTRING('fabcd', 1,3)，运行结果为 fab。

### 日期函数

![1575203513231](D:\Documents\笔记本\offer学习复习\其他\SQL_基础理论\1575203513231.png)

案例

SELECT CURRENT_DATE()，运行结果为 2019-04-03。

ELECT CURRENT_TIME()，运行结果为 21:26:34。

SELECT CURRENT_TIMESTAMP()，运行结果为 2019-04-03 21:26:34。

SELECT EXTRACT(YEAR FROM '2019-04-03')，运行结果为 2019。

SELECT DATE('2019-04-01 12:00:05')，运行结果为 2019-04-01。

### 转换函数

![1575203545175](D:\Documents\笔记本\offer学习复习\其他\SQL_基础理论\1575203545175.png)

案例

SELECT CAST(123.123 AS INT)，运行结果会报错。

SELECT CAST(123.123 AS DECIMAL(8,2))，运行结果为 123.12。

SELECT COALESCE(null,1,2)，运行结果为 1。

CAST函数在转换数据类型的时候，不会四舍五入，如果原数值有小数，那么转换为整数类型的时候就会报错。

### 用SQL函数对数据做处理

显示英雄以及他的物攻成长

```

SQL：SELECT name, ROUND(attack_growth,1) FROM heros
```

显示英雄最大生命值的最大值

```

SQL：SELECT MAX(hp_max) FROM heros
```

最大生命值最大的是哪个英雄，以及对应的数值

```

SQL：SELECT name, hp_max FROM heros WHERE hp_max = (SELECT MAX(hp_max) FROM heros)
```

显示英雄的名字，以及他们的名字字数

```

SQL：SELECT CHAR_LENGTH(name), name FROM heros
```

提取英雄上线日期（对应字段 birthdate）的年份，只显示有上线日期的英雄即可（有些英雄没有上线日期的数据，不需要显示）

```

SQL： SELECT name, EXTRACT(YEAR FROM birthdate) AS birthdate FROM heros WHERE birthdate is NOT NULL

------ 或者
SQL: SELECT name, YEAR(birthdate) AS birthdate FROM heros WHERE birthdate is NOT NULL
```

找出在 2016 年 10 月 1 日之后上线的所有英雄

```

SQL： SELECT * FROM heros WHERE DATE(birthdate)>'2016-10-01'

错误演示：
SELECT * FROM heros WHERE birthdate>'2016-10-01'

很多时候无法确认DATE类型，使用DATE(birthdate)来进行比较是更安全的！
```

查询在 2016 年 10 月 1 日之后上线英雄的平均最大生命值、平均最大法力和最高物攻最大值

```

SQL： SELECT AVG(hp_max), AVG(mp_max), MAX(attack_max) FROM heros WHERE DATE(birthdate)>'2016-10-01'
```

## 为什么使用 SQL 函数会带来问题

存在版本前后兼容问题，内置函数需要小心。

## 关于大小写的规范

Linux 是严格规范大小写的，Windows 无限定。建议使用统一的字段命名规则。

建议使用的命名规则：

-   关键字和函数名称全部大写；
-   数据库名、表名、字段名称全部小写；
-   SQL 语句必须以分号结尾。

## 总结

-   不同的 DBMS，内置函数的操作方式，输出内容是不尽相同的。

## 疑问

什么是SQL函数？



内置的SQL函数都包括哪些？



如何使用SQL函数对一个数据表进行操作，比如针对一个王者荣耀的英雄数据库，我们可以使用这些函数完成哪些操作？



什么情况下使用SQL函数？为什么使用SQL函数有时候会带来问题？