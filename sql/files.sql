CREATE TABLE IF NOT EXISTS files(
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    size INT NOT NULL,
    upLoadTime INT NOT NULL,
    fileType VARCHAR(128) NOT NULL,
    owner int NOT NULL,
    invisible int NOT NULL,
    FOREIGN KEY(fileType) REFERENCES fileType(typeName),
    FOREIGN KEY(owner) REFERENCES users(id)
);