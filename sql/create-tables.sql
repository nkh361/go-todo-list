/*
$ mysql -u root
$ use ticketing (command to show tables is 'show tables;')
mysql> source sql/create-tables.sql
*/


DROP TABLE IF EXISTS tickets;
CREATE TABLE tickets (
  id          INT AUTO_INCREMENT NOT NULL,
  username    VARCHAR(128) NOT NULL,
  title       VARCHAR(128) NOT NULL,
  priority    INT,
  PRIMARY KEY (`id`)
);

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id INT AUTO_INCREMENT NOT NULL,
    username VARCHAR(128) NOT NULL,
    password VARCHAR(128) NOT NULL,
    PRIMARY KEY (`id`)
);

ALTER TABLE users ADD CONSTRAINT uq_username UNIQUE (username);
CREATE UNIQUE INDEX ui_username ON users (username);