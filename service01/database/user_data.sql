CREATE DATABASE IF NOT EXISTS sd;
USE sd;

DROP TABLE IF EXISTS users;
CREATE TABLE users (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL
    ) engine = INNODB;

INSERT INTO
    users (name, email, password)
VALUES (
        'John Doe',
        'john@example.com',
        'password123'
    ), (
        'Jane Smith',
        'jane@example.com',
        'password456'
    ), (
        'Bob Johnson',
        'bob@example.com',
        'password789'
    ), (
        'Alice Brown',
        'alice@example.com',
        'passwordabc'
    ), (
        'Charlie Wilson',
        'charlie@example.com',
        'passworddef'
    );