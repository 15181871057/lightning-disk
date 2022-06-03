CREATE TABLE IF NOT EXISTS users(
    id INT AUTO_INCREMENT,
    name VARCHAR(32) NOT NULL UNIQUE,
    passwd VARCHAR(32) NOT NULL,
    registerTime int NOT NULL,
    rootFileId int NOT NULL,
    PRIMARY KEY(id, name)
);