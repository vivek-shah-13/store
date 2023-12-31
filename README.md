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


### Bonus
The migration runner should take a `lastRanID string` parameter, which will only execute the migrations starting from the given id. 

```go
func (m *MigrationRunner) Run(ctx context.Context, conn *sql.Conn, lastRanID int) error {
    // fill out as necessary
}
```

For example, if we have the following migrations available:
```
migrations/
    initial_0000.sql
    addProductSku_0001.sql
    example_0002.sql
```

We should expect the runner to perform the following operations for the given `lastRanIDs`:
```go
// An argument of -1 tells the runner we have not run any migrations yet. 
// This should run all of the migrations.
runner.Run(ctx, conn, -1) 

// An argument of 0 tells the runner we last ran migration initial_0000.sql.
// This should run the remaining migrations (addProductSku_0001.sql and example_0002.sql).
runner.Run(ctx, conn, 0)

// An argument of 2 tells the runner we last ran migration example_0002.sql.
// This should not run any migrations since example_0002.sql is the latest migration available.
runner.Run(ctx, conn, 2) 
```



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



# Refactor
Cleanup our application code so it is much easier to read, test, and extend. 
Our application's logic should be split up into separate packages, there should
be minimal copied code, and our application should be wired together in a cohesive manner. 

Consider adding the following:
* An abstraction over the data access layer
* Use a separate package to house mysql-related logic
* Use a separate package to house http-related logic
* Use separate package(s) to house domain models/entities: Customers, Orders, and Products


# HTTP Client
Create a new package named `client` which interfaces with your server's API. 
It should have the following methods:

```go
func (c *Client) createProduct(ctx context.Context, name string, price float64, sku string) (*Product, error)

func (c *Client) createCustomer(ctx context.Context, name, email string) (*Customer, error)

func (c *Client) createOrder(ctx context.Context, customerID, productID int) (*Order, error)

func (c *Client) showProducts(ctx context.Context, name string) ([]*Product, error)

func (c *Client) showCustomers(ctx context.Context, email, state string) ([]*Customer, error)

func (c *Client) showOrders(ctx context.Context, customerID, productID int) ([]*Order, error)
```

Each method should create a new http request given the parameters, send said request to your server,
and return the unmarshalled response. 

BONUS: Unmarshal errors returned by the server and convert them into specific error types (e.g. ErrCustomerDoesNotExist). 


Next, wire up each of these client methods to a new top-level `client` cli command. 
This `client` cli command will basically have the same methods at the main cli app:

```
store client [--host=localhost] create-product [--sku=SKU] NAME PRICE
store client [--host=localhost] create-customer EMAIL NAME
store client [--host=localhost] create-order CUSTOMER_ID PRODUCT_ID
store client [--host=localhost] show-products [--name=NAME]
store client [--host=localhost] show-customers [--email=EMAIL] [--STATE=state]
store client [--host=localhost] show-orders [--customer-id=CUSTOMER_ID] [--product-id=PRODUCT_ID]
```
