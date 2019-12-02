# 12丨视图在SQL中的作用是什么，它是怎样工作的？

[TOC]

## 疑问

什么是视图？如何创建、更新和删除视图？

何使用视图来简化我们的 SQL 操作？

视图和临时表的区别是什么，它们各自有什么优缺点？

## 如何创建，更新和删除视图

它相当于是一张表或多张表的数据结果集。视图的这一特点，可以帮我们简化复杂的 SQL 查询，不需要考虑视图中包含的基础查询的细节

### 创建视图：CREATE VIEW

语法

```

CREATE VIEW view_name AS
SELECT column1, column2
FROM table
WHERE condition
```

案例

```

CREATE VIEW player_above_avg_height AS
SELECT player_id, height
FROM player
WHERE height > (SELECT AVG(height) from player)


SELECT * FROM player_above_avg_height
```

### 嵌套视图

当我们创建好一张视图之后，还可以在它的基础上继续创建视图

```

CREATE VIEW player_above_above_avg_height AS
SELECT player_id, height
FROM player
WHERE height > (SELECT AVG(height) from player_above_avg_height)
```

### 修改视图：ALTER VIEW

语法

```

ALTER VIEW view_name AS
SELECT column1, column2
FROM table
WHERE condition
```

案例

```

ALTER VIEW player_above_avg_height AS
SELECT player_id, player_name, height
FROM player
WHERE height > (SELECT AVG(height) from player)

SELECT * FROM player_above_avg_height
```

### 删除视图：DROP VIEW

语法

```

DROP VIEW view_name
```

案例

```

DROP VIEW player_above_avg_height
```

SQLite 不支持视图的修改，仅支持只读视图，也就是说你只能使用 CREATE VIEW 和 DROP VIEW，如果想要修改视图，就需要先 DROP 然后再 CREATE

## 如何使用视图简化 SQL 操作

对 SELECT 语句进行了封装，方便我们重用它们

### 利用视图完成复杂的连接

球员以及对应身高等级的查询

```

CREATE VIEW player_height_grades AS
SELECT p.player_name, p.height, h.height_level
FROM player as p JOIN height_grades as h
ON height BETWEEN h.height_lowest AND h.height_highest
```

### 利用视图对数据进行格式化

输出球员姓名和对应的球队，对应格式为 player_name(team_name)

```

CREATE VIEW player_team AS 
SELECT CONCAT(player_name, '(' , team.team_name , ')') AS player_team FROM player JOIN team WHERE player.team_id = team.team_id


SELECT * FROM player_team
```

### 使用视图与计算字段

正确地使用视图可以帮我们简化复杂的数据处理。

统计每位球员在每场比赛中的二分球、三分球和罚球的得分

```

CREATE VIEW game_player_score AS
SELECT game_id, player_id, (shoot_hits-shoot_3_hits)*2 AS shoot_2_points, shoot_3_hits*3 AS shoot_3_points, shoot_p_hits AS shoot_p_points, score  FROM player_score;


SELECT * FROM game_player_score
```

有一点需要注意，视图是虚拟表，它只是封装了底层的数据表查询接口，因此有些 RDBMS 不支持对视图创建索引（有些 RDBMS 则支持，比如新版本的 SQL Server）

## 总结

使用视图有很多好处，比如安全、简单清晰：

-   安全性
    -   虚拟表是基于底层数据表的，我们在使用视图时，一般不会轻易通过视图对底层数据进行修改，即使是使用单表的视图，也会受到限制，这也在一定程度上保证了数据表的数据安全性
    -   针对不同用户开放不同的数据查询权限
-   简单清晰
    -   视图是对 SQL 查询的封装，它可以将原本复杂的 SQL 查询简化，在编写好查询之后，我们就可以直接重用它而不必要知道基本的查询细节。
    -   好比我们在进行模块化编程一样，不仅结构清晰，还提升了代码的复用率。

