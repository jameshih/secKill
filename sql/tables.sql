-- name: create-app-database
create database app;

-- name: use-app-database
use app;

-- name: create-product-table
CREATE TABLE product (
  id int AUTO_INCREMENT PRIMARY KEY,
  name varchar(1024) NOT NULL,
  total int(11) DEFAULT 0,
  status int(11) DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;
