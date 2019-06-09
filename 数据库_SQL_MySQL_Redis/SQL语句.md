## SQL 语句

SQL语句：数据库核心语句

SQL语句对表的增删改查：CURD：C CREATE(添加)：U UPDATE(更新)：R READ(查询)：D DELETE(删除)

在 SQL 语句中，表名和字段名加上 \` \` ，加快SQL语句的执行速度

查询：SELECT 字段，字段 FROM 表名

​	select `name` from 'tablename'

​	select * from `tablename`

​	select * from `tablename` where `id`=2

​	select * from `tablename` where `id`<>2

​	select * from `tablename` where `id`!=2

​	select * from `tablename` where `id`<2

​	select * from `tablename` where `id`>2

​	select * from `tablename` where `id` in (1,3,4)

​	select * from `tablename` where `id` not in (1,3,4)

​	select * from `tablename` where `id` between 2 and 4

​	select * from `tablename` where `name`='admin' and `pwd`='passwd'

​	select * from `tablename` where `name`='admin' or `name`='users'

​	select * from `tablename` order by `time` desc

​	select * from `tablename` where `id` in (1,3,4) order by `time` desc

​	select * from `tablename` limit 0,3

​	select * from `tablename` where `id` in (1,3,4) order by `id` desc limit 0,2



​	select count(\*) from `tablename`

​	select count(\*) as count from `tablename`

​	select max(\*) from `tablename`

​	select min(\*) from `tablename`

​	select avg(\*) from `tablename`

​	select sum(\*) from `tablename`



添加：insert into 表名 (字段名1,字段名2,字段名3) value(值1,值2,值3)
	insert into `admin` (`name`, `pwd`) value ('root', 'root123')
	
	
修改：update 表名 set 字段1=值1, 字段2=值2, 字段3=值3
	update `admin` set `name`='user', `pwd`='user123' where `id`=2
	
	
删除：delete from 表名 where 条件
	delete from `tablename` where `id`=2