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




## Requirements
Create CLI application with following commands:
 a) We should have create commands for customers, products, and orders:
    
  1. Command Syntax: store create-customer <email> <state> 

    This command should add a new customer to the customers table.
    The email and state arguments are required. 
    
  Example:
    > store create-customer zack.patrick@outreach.io WA
    CUSTOMER_ID       EMAIL                       STATE
    1                 zack.patrick@outreach.io    WA  

  
  2. Command Syntax: store create-product [--sku=<sku>] <name> <price> 

    This command should add a new product to the products table.
    The name and price arguments are required. 
    The sku flag is optional. 

  Example:
    > store create-product Bananas 3.99 
    PRODUCT_ID        NAME       PRICE       SKU
    1                 Bananas     3.99      

    > store create-product --sku aa1 Bananas 3.99 
    PRODUCT_ID        NAME       PRICE       SKU
    1  

  3. Command Syntax: store create-order <customer_id> <product_id>

    This command should add a new order to the orders table.
    The customer_id and product_id arguments are required. 

  Example:
    > store create-order 1 3
    ORDER_ID      CUSTOMER_ID         PRODUCT_ID
    1             1                   3

  b) We should have show commands for customers, products, and orders:
    
    1. Command syntax: store show-orders [--name=<name>]

    This command should display all of the orders in the orders table. 
    This command should have an optional --name flag which allows users to filter by name. 
    This parameter should support wildcard matching (e.g. 'foo%')
  
    Example:
      > store show-products
      PRODUCT_ID        NAME       PRICE       SKU
      1                 Bananas     3.99       aa1
      2                 Milk        4.50       
      3                 Cookies     1.99       ab2

    Example (using name flag):
      > store show-products --name "M%"
      PRODUCT_ID        NAME       PRICE       SKU
      2                 Milk        4.50       

    2. Command Syntax: store show-customers [--email=<email>] [--state=<state>]
    
    Show customers should display all of the customers in the customers table. 
    This command should have an optional --email flag which allows users to filter by email. 
    This command should have an optional --state flag which allows users to filter by state.
    These flags should support wildcard matching (e.g. 'foo%')
    
    Example:
      > store show-users
        CUSTOMER_ID       EMAIL                       STATE
        1                 zack.patrick@outreach.io    WA  
        2                 kevin.kerr@outreach.io      WA 
        3                 foo@bar.com                 CA  

    Example (using name flag):
      > store show-users --email "%@outreach.io"
        CUSTOMER_ID       EMAIL                       STATE
        1                 zack.patrick@outreach.io    WA  
        2                 kevin.kerr@outreach.io      WA  

    Example (using state flag):
      > store show-users --state WA
        CUSTOMER_ID       EMAIL                       STATE
        1                 zack.patrick@outreach.io    WA  
        2                 kevin.kerr@outreach.io      WA  

    3. Command Syntax: store show-orders [--customer-id=<customer_id>] [--product-id=<product_id>]
    
    Show orders should displays all of the orders in the orders table. 
    This command should have an optional --product-id flag which allows users to filter by products. 
    This should have an optional --customer-id flag which allows users to filter by customers. 

    Example:
      > store show-orders
      ORDER_ID      CUSTOMER_ID         PRODUCT_ID
      1             1                   2
      2             1                   3
      3             1                   1
      4             3                   1

    Example (with customer-id flag):
      > store show-orders --customer-id 1
      ORDER_ID      CUSTOMER_ID         PRODUCT_ID
      1             1                   2
      2             1                   3
      3             1                   1

    Example (with product-id flag):
      > store show-orders --product-id 3
      ORDER_ID      CUSTOMER_ID         PRODUCT_ID
      2             1                   3
