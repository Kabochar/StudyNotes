# 05丨检索数据：你还在SELECT * 么？

[TOC]

## SELECT 查询的基础语法

### 查询列

```

SQL：SELECT name FROM heros
```

### 起别名

```

SQL：SELECT name, hp_max, mp_max, attack_max, defense_max FROM heros;


SQL：SELECT * FROM heros
```

### 起别名

```

SQL：SELECT name AS n, hp_max AS hm, mp_max AS mm, attack_max AS am, defense_max AS dm FROM heros
```

### 查询常数

>   在SELECT查询结果中增加一列固定的常数列。
>
>   使用场景：想整合不同的数据源，用常数列作为这个表的标记。

```

SQL：SELECT '王者荣耀' as platform, name FROM heros
```

### 去除重复行

```

SQL：SELECT DISTINCT attack_range FROM heros
```

-   DISTINCT 需要放到所有列名的前面
-   DISTINCT 其实是对后面所有列名的组合进行去重

## 如何排序检索数据

使用 ORDER BY 语句

-   排序的列名
    -   如果存在多个，从第一个开始
-   排序的顺序
    -   DESC，降序
    -   ASC，升序（默认使用 ASC）
-   非选择列排序
    -   即使 SELET 语句没有，ORDER BY 同样可以使用
-   ORDER BY 的位置
    -   t通常 位于 SELECT 语句的最后一条子句

按照最大生命值排序

```

SQL：SELECT name, hp_max FROM heros ORDER BY hp_max DESC 
```

最大生命值和法力值排名

```

SQL：SELECT name, hp_max FROM heros ORDER BY mp_max, hp_max DESC  
```

### 约束返回结果数量

使用  LIMIT 关键词

```

SQL：SELECT name, hp_max FROM heros ORDER BY hp_max DESC LIMIT 5
```

## SELECT 的执行顺序

SELECT 查询时的两个顺序：

-   关键字的顺序是不能颠倒的

```

SELECT ... FROM ... WHERE ... GROUP BY ... HAVING ... ORDER BY ...
```

-   SELECT 语句的执行顺序（在 MySQL 和 Oracle 中，SELECT 执行顺序基本相同）

```

FROM > WHERE > GROUP BY > HAVING > SELECT的字段 > DISTINCT > ORDER BY > LIMIT
```

### 样例

```
SELECT DISTINCT player_id, player_name, count(*) as num #顺序5
FROM player JOIN team ON player.team_id = team.team_id #顺序1
WHERE height > 1.80 #顺序2
GROUP BY player.team_id #顺序3
HAVING num > 2 #顺序4
ORDER BY num DESC #顺序6
LIMIT 2 #顺序7
```

## 什么情况下用 SELECT*，如何提升 SELECT 查询效率？

练习时使用，生产开发不建议

## 疑问解答

SELECT 查询的基础语法？



如何排序检索数据？



什么情况下用SELECT*，如何提升SELECT查询效率？



SELECT COUNT(*) ＞ SELECT COUNT(1) ＞ SELECT COUNT(具体字段)解析？

翻阅专栏评论

