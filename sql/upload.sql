CREATE TABLE IF NOT EXISTS upload(
    id int PRIMARY KEY AUTO_INCREMENT,
    fileId int NOT NULL,
    size int NOT NULL,
    time int NOT NULL,
    FOREIGN KEY(fileId) REFERENCES files(id)
);
size以4mb为1单位