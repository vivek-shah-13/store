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


# CLI Application
Create a cli application which interacts with the customer, product, and order tables in the database.

## Requirements

### We should have create commands for customers, products, and orders:

  1. Command Syntax: `store create-customer <email> <state>`


This command should add a new customer to the customers table.
The email and state arguments are required. 
    
Example:

    > store create-customer zack.patrick@outreach.io WA
    CUSTOMER_ID       EMAIL                       STATE
    1                 zack.patrick@outreach.io    WA  

  
2. Command Syntax: `store create-product [--sku=<sku>] <name> <price>`


This command should add a new product to the products table.
The name and price arguments are required. 
The sku flag is optional. 

Examples:

    > store create-product Bananas 3.99 
    PRODUCT_ID        NAME       PRICE       SKU
    1                 Bananas     3.99      

    > store create-product --sku aa1 Bananas 3.99 
    PRODUCT_ID        NAME       PRICE      SKU
    1                 Bananas     3.99      aa1

3. Command Syntax: `store create-order <customer_id> <product_id>`


This command should add a new order to the orders table.
The customer_id and product_id arguments are required. 

Example:

    > store create-order 1 3
    ORDER_ID      CUSTOMER_ID         PRODUCT_ID
    1             1                   3

 
 ### We should have show commands for customers, products, and orders:

1. Command syntax: `store show-orders [--name=<name>]`


This command should display all of the orders in the orders table. 
This command should have an optional --name flag which allows users to filter by name. 
This parameter should support wildcard matching (e.g. 'foo%')

Examples:

      > store show-products
      PRODUCT_ID        NAME       PRICE       SKU
      1                 Bananas     3.99       aa1
      2                 Milk        4.50       
      3                 Cookies     1.99       ab2

      > store show-products --name "M%"
      PRODUCT_ID        NAME       PRICE       SKU
      2                 Milk        4.50       
      5                 Markers     8.99       xi1          

2. Command Syntax: `store show-customers [--email=<email>] [--state=<state>]`   


Show customers should display all of the customers in the customers table. 
This command should have an optional --email flag which allows users to filter by email. 
This command should have an optional --state flag which allows users to filter by state.
These flags should support wildcard matching (e.g. 'foo%')

Examples:

    > store show-users
    CUSTOMER_ID       EMAIL                       STATE
    1                 zack.patrick@outreach.io    WA  
    2                 kevin.kerr@outreach.io      WA 
    3                 foo@bar.com                 CA  

    > store show-users --email "%@outreach.io"
    CUSTOMER_ID       EMAIL                       STATE
    1                 zack.patrick@outreach.io    WA  
    2                 kevin.kerr@outreach.io      WA  

    > store show-users --state WA
    CUSTOMER_ID       EMAIL                       STATE
    1                 zack.patrick@outreach.io    WA  
    2                 kevin.kerr@outreach.io      WA  

3. Command Syntax: `store show-orders [--customer-id=<customer_id>] [--product-id=<product_id>]`


Show orders should displays all of the orders in the orders table. 
This command should have an optional --product-id flag which allows users to filter by products. 
This should have an optional --customer-id flag which allows users to filter by customers. 

Examples:

    > store show-orders
    ORDER_ID      CUSTOMER_ID         PRODUCT_ID
    1             1                   2
    2             1                   3
    3             1                   1
    4             3                   1

    > store show-orders --customer-id 1
    ORDER_ID      CUSTOMER_ID         PRODUCT_ID
    1             1                   2
    2             1                   3
    3             1                   1

    > store show-orders --product-id 3
    ORDER_ID      CUSTOMER_ID         PRODUCT_ID
    2             1                   3


# Testing
Write unit and integration tests for each command in the cli application. 


# Migration Runner
We now want to manage our database migrations in a scalable way. 


## Requirements

Create a `/migrations` directory which will hold each of our migration files.
Each migration should live in a /migrations directory and have follow the pattern <name>_<id>.sql.
The first migration should be named 'initial_0000.sql' and setup our tables.
Add a second migration, `addProductSku_0001.sql`, which adds the SKU column to products table.

When completed, the directory structure should look like:
```
migrations/
    initial_0000.sql
    addProductSku_0001.sql
```

Create a `MigrationRunner` type with the following structure:
```go
type MigrationRunner struct {
    path string
}

func NewMigrationRunner(path string) *MigrationRunner {
    return &MigrationRunner{path: path}
}

func (m *MigrationRunner) Run(ctx context.Context, conn *sql.Conn) error {
    // fill out as necessary
}
```
The `MigrationRunner.Run` method should do the following:
1. Load each of the migration files in `m.path`
2. Sort the migration files based on the `id` in their filename (e.g. `initial_0000.sql` is `0`,`addProductSku_0001.sql` is `1`).
3. In order, load the sql command(s) from each migration file and execute the command(s) against the `conn` param.
			

# Team Exercise: Tenancy
We have 2 customers for our application: google and microsoft. 
We now need to make sure each customer has their own separate data.
	
## Requirements
* Each customer's data must reside in their own database.
* Migrations need to run for each database.
* We must a "default" database for developers to use. 
* The store's top-level command now takes a `--org <name>` argument. 
If no `--org` argument is set, the "default" org is used. 

## Working as a Team:
Together, the team should decide together how to implement this feature and how
to split up the work. 
One recommended way might be the following:
1. Decide on a common interface which can be used to split up the work.
2. Pair program a PR which adds any common interface/scaffolding as necessary.
3. Split up the work: migration runner, cli flag, database switching. 


# Web Application
We now would like to expose our application to the internet.
Create a HTTP interface for the application. 

## Requirements
Create the following HTTP endpoints:

    GET /customer[?org=<org>&state=<state>&email=<email>]
        This should return a list of customers satisfying the query parameters.
        If org is unset, the default org will be used. 

    GET /customer/:id[?org=<org>]
        This should return the customer with the given id.
        If org is unset, the default org will be used.
        
    GET /product[?org=<org>&name=<name>&sku=<sku>]
        This should return a list of products satisfying the query parameters.
        If org is unset, the default org will be used. 

    GET /product/:id[?org=<org>]
        This should return the product with the given id.
        If org is unset, the default org will be used.

    GET /order[?org=<org>&product_id=<product_id>&customer_id=<customer_id>]
        This should return a list of orders satisfying the query parameters.
        If org is unset, the default org will be used. 

    GET /order/:id[?org=<org>]
        This should return the order with the given id.
        If org is unset, the default org will be used.

* Each of the endpoints should have unit and integration tests.
* Use https://ngrok.com/ to expose your application to the public internet. 
Others should now be able to interact with your website. 