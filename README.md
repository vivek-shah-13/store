# store
## Get mysql setup
> brew install mysql
> mysql.server start

## Connect to mysql
> mysql -u root
mysql> create database store;
mysql> use store;

## Create user
> CREATE USER 'store' IDENTIFIED BY 'password123'