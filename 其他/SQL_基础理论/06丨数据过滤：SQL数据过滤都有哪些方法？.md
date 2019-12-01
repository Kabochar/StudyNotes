# 06丨数据过滤：SQL数据过滤都有哪些方法？

[TOC]

## 比较运算符

![1575191109063](D:\Documents\笔记本\offer学习复习\其他\SQL_基础理论\1575191109063.png)

样例

```

SQL：SELECT name, hp_max FROM heros WHERE hp_max > 6000
```

在xx 之间

```

SQL：SELECT name, hp_max FROM heros WHERE hp_max BETWEEN 5399 AND 6811
```

## 逻辑运算符

![1575191174698](D:\Documents\笔记本\offer学习复习\其他\SQL_基础理论\1575191174698.png)

样例

假设想要筛选最大生命值大于6000，最大法力大于1700的英雄，然后按照最大生命值和最大法力值之和从高到低进行排序。

```

SQL：SELECT name, hp_max, mp_max FROM heros WHERE hp_max > 6000 AND mp_max > 1700 ORDER BY (hp_max+mp_max) DESC
```

当 WHERE 子句中同时出现 AND 和 OR 操作符的时候，你需要考虑到执行的先后顺序，也就是两个操作符执行的优先级。一般来说（）优先级最高，其次优先级是 AND，然后是 OR。

```

SQL：
SELECT name, role_main, role_assist, hp_max, mp_max, birthdate
FROM heros 
WHERE (role_main IN ('法师', '射手') OR role_assist IN ('法师', '射手')) 
AND DATE(birthdate) NOT BETWEEN '2016-01-01' AND '2017-01-01'
ORDER BY (hp_max + mp_max) DESC
```

### 使用通配符进行过滤

使用到 LIKE 操作符。如果我们想要匹配任意字符串出现的任意次数，需要使用（%）通配符。

```

SQL：SELECT name FROM heros WHERE name LIKE '% 太 %'
```

如果我们想要匹配单个字符，就需要使用下划线（_）通配符。（%）和（_）的区别在于，（%）代表零个或多个字符，而（_）只代表一个字符。

```

SQL：SELECT name FROM heros WHERE name LIKE '_% 太 %'
```

建议你尽量少用通配符，因为它需要消耗数据库更长的时间来进行匹配。即使你对LIKE检索的字段进行了索引，索引的价值也可能会失效。

## 总结

检索的代价也是很高的，通常都需要用到全表扫描，所以效率很低。只有当LIKE语句后面不用通配符，并且对字段进行索引的时候才不会对全表进行扫描。

保持高效率的一个很重要的原因，就是要避免全表扫描，所以我们会**考虑在WHERE及ORDER BY涉及到的列上增加索引。**



## 疑问解答

学会使用WHERE子句，如何使用比较运算符对字段的数值进行比较筛选？



如何使用逻辑运算符，进行多条件的过滤？



学会使用通配符对数据条件进行复杂过滤？



为什么要在 WHERE 和 ORDER BY 涉及的列上增加索引？

>   Index排序：索引可以保证数据的有序性，因此不需要再进行排序。
>
>   FileSort排序：一般在内存中进行排序，占用CPU较多。如果待排结果较大，会产生临时文件I/O到磁盘进行排序，效率较低。

1、SQL中，可以在WHERE子句和ORDER BY子句中使用索引，目的是在WHERE子句中避免全表扫描，ORDER BY子句避免使用FileSort排序。
当然，某些情况下全表扫描，或者FileSort排序不一定比索引慢。但总的来说，我们还是要避免，以提高查询效率。

一般情况下，优化器会帮我们进行更好的选择，当然我们也需要建立合理的索引。

2、尽量Using Index完成ORDER BY排序。

如果WHERE和ORDER BY相同列就使用单索引列；如果不同使用联合索引。

3、无法Using Index时，对FileSort方式进行调优。