# 04丨使用 DDL 创建数据库&数据表时需要注意什么？

[TOC]

## 疑问解答

了解DDL的基础语法，它如何定义数据库和数据表？

使用DDL定义数据表时，都有哪些约束性？

使用DDL 设计数据库时，都有哪些重要原则？

## DDL的基础语法及设计工具

DDL的英文全称是Data Definition Language，中文是数据定义语言。它定义了数据库的结构和数据表的结构。
在DDL中，我们常用的功能是增删改，分别对应的命令是CREATE、DROP和ALTER。需要注意的是，在执行DDL的时候，不需要COMMIT，就可以完成执行任务。

1，对数据库进行定义

```
CREATE DATABASE nba; // 创建一个名为 nba 的数据库
DROP DATABASE nba; // 删除一个名为 nba 的数据库
```

2，对数据表进行定义

```
CREATE TABLE [](字段名 数据类型，......)
```

### 创建表结构

```
CREATE TABLE player  (
  player_id int(11) NOT NULL AUTO_INCREMENT,
  player_name varchar(255) NOT NULL
);


DROP TABLE IF EXISTS `player`;
CREATE TABLE `player`  (
  `player_id` int(11) NOT NULL AUTO_INCREMENT,
  `team_id` int(11) NOT NULL,
  `player_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `height` float(3, 2) NULL DEFAULT 0.00,
  PRIMARY KEY (`player_id`) USING BTREE,
  UNIQUE INDEX `player_name`(`player_name`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;
```

### 修改表结构

添加字段

```
ALTER TABLE player ADD (age int(11));
```

修改字段名

```
ALTER TABLE player RENAME COLUMN age to player_age
```

修改字段的数据类型

```
ALTER TABLE player MODIFY (player_age float(3,1));
```

删除字段

```
ALTER TABLE player DROP COLUMN player_age;
```

### 数据表的常见约束

#### 主键约束

主键起的作用是唯一标识一条记录，不能重复，不能为空，即 UNIQUE + NOT NULL。一个数据表的主键只能有一个。主键可以是一个字段，也可以由多个字段复合组成

#### 外键约束

外键确保了表与表之间引用的完整性。一个表中的外键对应另一张表的主键

#### 唯一性约束

唯一性约束表明了字段在表中的数值是唯一的，即使我们已经有了主键，还可以对其他字段进行唯一性约束

#### NOT NULL 约束

NOT NULL 约束。对字段定义了NOT NULL，即表明该字段不应为空，必须有取值

#### DEFUALT

DEFAULT，表明了字段的默认值。如果在插入数据的时候，这个字段没有取值，就设置为默认值

#### CHECK 约束

CHECK约束，用来检查特定字段取值范围的有效性，CHECK约束的结果不能为FALSE

## 设计数据表的原则

遵循三少一多原则

-   数据表的个数越少越好
    -   RDBMS的核心在于对实体和联系的定义，也就是E-R图（Entity Relationship Diagram），数据表越少，证明实体和联系设计得越简洁，既方便理解又方便操作。
-   数据表中的字段个数越少越好
    -   字段个数越多，数据冗余的可能性越大。
-   数据表中联合主键的字段个数越少越好
    -   联合主键中的字段越多，占用的索引空间越大，不仅会加大理解难度，还会增加运行时间和索引空间，因此联合主键的字段个数越少越好。
-   使用主键和外键越多越好（有待考究）
    -   数据库的设计实际上就是定义各种表，以及各种字段之间的关系。这些关系越多，证明这些实体之间的冗余度越低，利用度越高。
    -   解释为什么越多越好：
        -   越多越好的使用场景：外键本身是为了实现强一致性，所以如果需要**正确性>性能**的话，还是建议使用外键，它可以让我们在数据库的层面保证数据的完整性和一致性。
        -   对于学习SQL初期为了建立一个强一致性可靠性高的数据库而选择使用主键和外键约束。

