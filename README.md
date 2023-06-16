# store
## Get mysql setup
> brew install mysql
> mysql.server start

## Connect to mysql
> mysql -u root
mysql> create database store;
mysql> use store;

## Create user
> CREATE USER 'admin'@'%' IDENTIFIED BY 'password123';
> GRANT ALL PRIVILEGES ON *.* TO 'admin'@'%';