# 10丨常用的SQL标准有哪些，在SQL92中是如何使用连接的？

[TOC]

## 疑问

-   SQL 实际上存在不同的标准，不同标准下的连接定义也有不同。你首先需要了解常用的SQL标准有哪些？
-   了解了SQL的标准之后，我们从SQL92标准入门，来看下连接表的种类有哪些？
-   针对一个实际的数据库表，如果你想要做数据统计，需要学会使用跨表的连接进行操作？
-   表格中一共有 3 支球队，现在这 3 支球队需要进行比赛，请用一条 SQL 语句显示出所有可能的比赛组合？

```
区分主客队
SELECT CONCAT(kedui.team_name, ' VS ', zhudui.team_name) as '客队 VS 主队' FROM team as zhudui LEFT JOIN team as kedui on zhudui.team_id<>kedui.team_id;

不区分主客队
SELECT a.team_name as '队伍1' ,'VS' , b.team_name as '队伍2' FROM team as a ,team as b where a.team_id<b.team_id;
```



## 常用的SQL标准有哪些

应用最广泛：SQL92和SQL99，92，99代表标准指定时间

SQL92 形式更简单，但语句会偏长，SQL 99相对复杂，但可读性强

## 在SQL92中是如何使用连接的

### 笛卡尔积

笛卡尔乘积是一个数学运算。假设我有两个集合X和Y，那么X和Y的笛卡尔积就是X和Y的所有可能组合，也就是第一个对象来自于X，第二个对象来自于Y的所有可能。

```

SQL: SELECT * FROM player, team
```

笛卡尔积也称为交叉连接，英文是CROSS  JOIN，它的作用就是可以**把任意表进行连接**，即使这两张表不相关。

### 等值连接

player 表和team表都存在team_id 这一列使用等值链接

```

SQL: SELECT player_id, player.team_id, player_name, height, team_name FROM player, team WHERE player.team_id = team.team_id
```

需要注意的是，如果我们使用了表的别名，在查询字段中就只能使用别名进行代替，不能使用原有的表名！

### 非等值连接

如果连接多个表的条件是等号时，就是等值连接，其他的运算符连接就是非等值查询。

想要知道每个球员的身高的级别可以采用非等值连接查询

```

SQL：SELECT p.player_name, p.height, h.height_level
FROM player AS p, height_grades AS h
WHERE p.height BETWEEN h.height_lowest AND h.height_highest
```

### 外连接

除了查询满足条件的记录以外，外连接还可以查询某一方不满足条件的记录。两张表的外连接，会有一张是主表，另一张是从表。如果是多张表的外连接，那么第一张表是主表，即显示全部的行，而第剩下的表则显示对应连接的信息。

在SQL92中采用（+）代表从表所在的位置，而且在SQL92中，只有左外连接和右外连接，没有全外连接。

什么是左外连接，什么是右外连接呢？

左外连接，就是指左边的表是主表，需要显示左边表的全部行，而右侧的表是从表，（+）表示哪个是从表。

```

SQL：SELECT * FROM player, team where player.team_id = team.team_id(+)

SQL：SELECT * FROM player, team where player.team_id = team.team_id(+)
```

右外连接，指的就是右边的表是主表，需要显示右边表的全部行，而左侧的表是从表。

```

SQL：SELECT * FROM player, team where player.team_id(+) = team.team_id

SQL：SELECT * FROM player RIGHT JOIN team on player.team_id = team.team_id
```

需要注意的是，LEFT JOIN和RIGHT JOIN只存在于SQL99及以后的标准中，在SQL92中不存在，只能用（+）表示。

### 自连接

自连接可以对多个表进行操作，也可以对同一个表进行操作。也就是说查询条件使用了当前表的字段。

查看比布雷克·格里芬高的球员都有谁，以及他们的对应身高

```

SQL：SELECT b.player_name, b.height FROM player as a , player as b WHERE a.player_name = '布雷克 - 格里芬' and a.height < b.height
```

### 总结

SQL92和SQL99是经典的SQL标准，也分别叫做SQL-2和SQL-3标准。也正是在这两个标准发布之后，SQL影响力越来越大，甚至超越了数据库领域。