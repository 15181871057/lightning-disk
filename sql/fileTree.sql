CREATE TABLE IF NOT EXISTS fileTree(
    father INT NOT NULL,
    child INT PRIMARY KEY NOT NULL,
    FOREIGN KEY(father) REFERENCES files(id),
    FOREIGN KEY(child) REFERENCES files(id)
);


// 遍历树
WITH RECURSIVE cte_select(id) AS (
select father from fileTree WHERE fileTree.father=3
UNION ALL
select f.child from fileTree AS f JOIN cte_select AS c ON f.father=c.id)
select * from cte_select;

// 删除文件
WITH RECURSIVE cte_select(id) AS (
select father from fileTree WHERE fileTree.father=2
UNION ALL
select f.child from fileTree AS f JOIN cte_select AS c ON f.father=c.id)
DELETE FROM files WHERE id in (SELECT * FROM cte_select);

// 删除树
WITH RECURSIVE cte_select(id) AS (
select father from fileTree WHERE fileTree.father=2
UNION ALL
select f.child from fileTree AS f JOIN cte_select AS c ON f.father=c.id)
DELETE FROM fileTree WHERE father in (SELECT * FROM cte_select);

// 多表删除
WITH RECURSIVE cte_select(id) AS (
select father from fileTree WHERE fileTree.father=3
UNION ALL
select f.child from fileTree AS f JOIN cte_select AS c ON f.father=c.id)
DELETE fileTree, Files FROM fileTree JOIN  Files WHERE father in (SELECT * FROM cte_select);

select * from fileTree WHERE father=?