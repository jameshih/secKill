-- name: create-app-database
create database app;

-- name: use-app-database
use app;

-- name: create-product-table
CREATE TABLE product (
  id int AUTO_INCREMENT PRIMARY KEY,
  name varchar(1024) NOT NULL UNIQUE,
  total int(11) DEFAULT 0,
  status int(11) DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;

-- name: create-event-table
CREATE TABLE event (
  id int AUTO_INCREMENT PRIMARY KEY,
  name varchar(1024) NOT NULL UNIQUE,
  product_id int(11) DEFAULT 0,
  start_time int(11) DEFAULT 0,
  end_time int(11) DEFAULT 0,
  total int(11) DEFAULT 0,
  status int(11) DEFAULT 0,
  req_limit int DEFAULT 100,
  buy_limit int DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;

