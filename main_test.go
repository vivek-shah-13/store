package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/vivek-shah-13/store/internal/migration"
	"gotest.tools/v3/assert"
)

var (
	dbName = "google"
)

func tearDownCustomers() error {
	db, err := connectDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM Customers")
	if err != nil {
		return err
	}
	_, err = db.Exec("ALTER TABLE Customers AUTO_INCREMENT=1")
	if err != nil {
		return err
	}
	return nil
}

func tearDownProducts() error {
	db, err := connectDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM Products")
	if err != nil {
		return err
	}
	_, err = db.Exec("ALTER TABLE Products AUTO_INCREMENT=1")
	if err != nil {
		return err
	}
	return nil
}

func tearDownOrders() error {
	db, err := connectDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM Orders")
	if err != nil {
		return err
	}
	_, err = db.Exec("ALTER TABLE Orders AUTO_INCREMENT=1")
	if err != nil {
		return err
	}
	err = tearDownCustomers()
	if err != nil {
		return err
	}
	err = tearDownProducts()
	if err != nil {
		return err
	}
	return nil
}

func createCustomersData(args [][]string) error {
	db, err := connectDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	app := &cli.App{
		Commands: []*cli.Command{
			newCreateCustomerCommand(&db),
		},
	}

	for _, val := range args {
		err = app.Run(val)
		if err != nil {
			return err
		}
	}
	return nil
}

func createCustomersDataV2(t *testing.T, args [][]string) {
	assert.NilError(t, createCustomersDataV3(args))
}

func createCustomersDataV3(args [][]string) error {
	db, err := connectDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, val := range args {
		_, err = db.Exec("INSERT INTO Customers (email, state) VALUES (?, ?)", val[0], val[1])
		if err != nil {
			return err
		}
	}
	return nil
}

func createProductsData(args [][]string) error {
	db, err := connectDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	app := &cli.App{
		Commands: []*cli.Command{
			newCreateProductCommand(&db),
		},
	}

	for _, val := range args {
		err = app.Run(val)
		if err != nil {
			return err
		}
	}
	return nil
}

func createProductsDataV2(t *testing.T, args [][]any) {
	assert.NilError(t, createProductsDataV3(args))
}

func createProductsDataV3(args [][]any) error {
	db, err := connectDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, val := range args {
		_, err = db.Exec("INSERT INTO Products (name, price, sku) VALUES (?, ?, ?)", val[0], val[1], val[2])
		if err != nil {
			return err
		}
	}
	return nil
}

func createOrdersData(customers [][]string, products [][]string, orders [][]string) error {
	err := createCustomersData(customers)
	if err != nil {
		return err
	}
	err = createProductsData(products)
	if err != nil {
		return err
	}

	db, err := connectDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	app := &cli.App{
		Commands: []*cli.Command{
			newCreateOrderCommand(&db),
		},
	}

	for _, val := range orders {
		err = app.Run(val)
		if err != nil {
			return err
		}
	}

	return nil
}

func createOrdersDataV2(t *testing.T, args [][]int) {
	assert.NilError(t, createOrdersDataV3(args))

}

func createOrdersDataV3(args [][]int) error {
	db, err := connectDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, val := range args {
		_, err = db.Exec("INSERT INTO Customers (email, state) VALUES ('vivek.shah@outreach.io', 'WA')")
		if err != nil {
			return err
		}
		_, err = db.Exec("INSERT INTO Products (name, price, sku) VALUES ('laptop', 25.5, 'abcde')")
		if err != nil {
			return err
		}
		_, err = db.Exec("INSERT INTO Orders (customer_id, product_id) VALUES (?, ?)", val[0], val[1])
		if err != nil {
			return err
		}

	}
	return nil
}

// ### We should have create commands for customers

//   1. Command Syntax: `store create-customer <email> <state>`

// This command should add a new customer to the customers table.
// The email and state arguments are required.

// Example:

//     > store create-customer zack.patrick@outreach.io WA
//     CUSTOMER_ID       EMAIL                       STATE
//     1                 zack.patrick@outreach.io    WA

// 1. When I run it with correct inputs, does the new customer get added to the database correctly?
// 2. When I run it with correct inputs, does the new customer get printed as I expected?
// 3. When I run it without a state that's a 2 letter code: it should return an error (TODO: error handling)
// 4. When I run it with a state that's not a valid abbreviation: it should return an error
// 5. When I run it without an email arg: it should return an error
// 6. When I run it without a state arg: it shoudl return an error

func TestCreateCustomer_withAMissingEmailArg_returnsAnError(t *testing.T) {
	err := createCustomersData([][]string{{"store", "create-customer", "WA"}})
	assert.ErrorContains(t, err, "Must specify email and state")
}

func TestCreateCustomer_withAMissingStateArg_returnsAnError(t *testing.T) {

	err := createCustomersData([][]string{{"store", "create-customer", "vivek.shah@outreach.io"}})
	assert.ErrorContains(t, err, "Must specify email and state")
}

func TestCreateCustomer_wthAInvalidStateAbbreviation_returnsAnError(t *testing.T) {

	err := createCustomersData([][]string{{"store", "create-customer", "vivek.shah@outreach.io", "PP"}})
	assert.ErrorContains(t, err, "State must be a valid U.S. State or Territory")
}

func TestCreateCustomer_withANonTwoLetterStateCode_returnsAnError(t *testing.T) {

	err := createCustomersData([][]string{{"store", "create-customer", "vivek.shah@outreach.io", "PPP"}})
	assert.ErrorContains(t, err, "State length must be 2")
}

func TestCreateCustomer_withValidInput_entersDatabaseCorrectly(t *testing.T) {

	err := tearDownCustomers()
	assert.NilError(t, err)
	err = createCustomersData([][]string{{"store", "create-customer", "vivek.s@outreach.io", "WA"}})
	assert.NilError(t, err)
	db, err := connectDB(dbName)
	assert.NilError(t, err)
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM Customers ORDER BY ID DESC LIMIT 1")

	assert.NilError(t, rows.Err())
	var id int
	var email string
	var state string
	rows.Scan(&id, &email, &state)
	assert.Equal(t, email, "vivek.s@outreach.io")
	assert.Equal(t, state, "WA")
	assert.Equal(t, id, 1)

}

func TestCreateMultipleCustomer_withValidInput_entersDatabaseCorrectly(t *testing.T) {
	db, err := connectDB(dbName)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = tearDownCustomers()
	assert.NilError(t, err)
	err = createCustomersData([][]string{{"store", "create-customer", "vivek.s@outreach.io", "WA"}, {"store", "create-customer", "v.s@outreach.io", "CA"}})
	assert.NilError(t, err)

	rows := db.QueryRow("SELECT * FROM Customers ORDER BY ID DESC LIMIT 1")

	assert.NilError(t, rows.Err())
	var id int
	var email string
	var state string
	rows.Scan(&id, &email, &state)
	assert.Equal(t, email, "v.s@outreach.io")
	assert.Equal(t, state, "CA")
	assert.Equal(t, id, 2)

}

func Example_CreateCustomer_hasCorrectPrintOutput() {
	printCustomer(os.Stdout, &Customer{1, "vivek.s@outreach.io", "WA"})

	//Output:
	//ID  |Email                                              |State |
	//1   |vivek.s@outreach.io                                |WA    |

}

func Example_CreateMultipleCustomer_hasCorrectPrintOutput() {
	printCustomer(os.Stdout, &Customer{1, "vivek.s@outreach.io", "WA"})
	printCustomer(os.Stdout, &Customer{2, "v.s@verylonglonglonglongemail.com", "MN"})
	//Output:
	//ID  |Email                                              |State |
	//1   |vivek.s@outreach.io                                |WA    |
	//ID  |Email                                              |State |
	//2   |v.s@verylonglonglonglongemail.com                  |MN    |
}

// Command Syntax: `store create-product [--sku=<sku>] <name> <price>`

// This command should add a new product to the products table.
// The name and price arguments are required.
// The sku flag is optional.

// Examples:

//     > store create-product Bananas 3.99
//     PRODUCT_ID        NAME       PRICE       SKU
//     1                 Bananas     3.99

//     > store create-product --sku aa1 Bananas 3.99
//     PRODUCT_ID        NAME       PRICE      SKU
//     1                 Bananas     3.99      aa1

// 1. insert valid product with name price and sku
// 2. insert valid product with name price and no sku
// 3. invalid product (missing name or price)
// 4. insert valid product with name decimal price and sku/no sku (doesn't really matter as sku is not the important variable being changed)
// 5. correct print statement no sku
// 6. correct print statement w/ sku
func TestCreateProduct_noName(t *testing.T) {

	err := createProductsData([][]string{{"store", "create-product", "3.99"}})
	assert.ErrorContains(t, err, "Must specify name and price")

}

func TestCreateProduct_noPrice(t *testing.T) {

	err := createProductsData([][]string{{"store", "create-product", "banana"}})
	assert.ErrorContains(t, err, "Must specify name and price")
}

func TestCreateProduct_validInputNoSku(t *testing.T) {
	db, err := connectDB(dbName)
	assert.NilError(t, err)
	defer db.Close()

	err = tearDownProducts()
	assert.NilError(t, err)
	err = createProductsData([][]string{{"store", "create-product", "banana", "5"}})
	assert.NilError(t, err)

	rows := db.QueryRow("SELECT * FROM Products ORDER BY ID DESC LIMIT 1")

	assert.NilError(t, rows.Err())
	var id int
	var name string
	var price float64
	var sku string
	rows.Scan(&id, &name, &price, &sku)
	assert.Equal(t, name, "banana")
	assert.Equal(t, price, float64(5))
	assert.Equal(t, id, 1)
	assert.Equal(t, sku, "")
}

func TestCreateProduct_validInputWithSku(t *testing.T) {
	db, err := connectDB(dbName)
	assert.NilError(t, err)
	defer db.Close()

	err = tearDownProducts()
	assert.NilError(t, err)
	err = createProductsData([][]string{{"store", "create-product", "--sku=abcde", "banana", "5"}})
	assert.NilError(t, err)

	rows := db.QueryRow("SELECT * FROM Products ORDER BY ID DESC LIMIT 1")

	assert.NilError(t, rows.Err())
	var id int
	var name string
	var price float64
	var sku string
	rows.Scan(&id, &name, &price, &sku)
	assert.Equal(t, name, "banana")
	assert.Equal(t, price, float64(5))
	assert.Equal(t, id, 1)
	assert.Equal(t, sku, "abcde")
}

func TestCreateProduct_decimalPrice(t *testing.T) {
	db, err := connectDB(dbName)
	assert.NilError(t, err)
	defer db.Close()

	err = tearDownProducts()
	assert.NilError(t, err)
	err = createProductsData([][]string{{"store", "create-product", "--sku=abcde", "banana", "5.50"}})
	assert.NilError(t, err)

	rows := db.QueryRow("SELECT * FROM Products ORDER BY ID DESC LIMIT 1")

	assert.NilError(t, rows.Err())
	var id int
	var name string
	var price float64
	var sku string
	rows.Scan(&id, &name, &price, &sku)
	assert.Equal(t, name, "banana")
	assert.Equal(t, price, float64(5.50))
	assert.Equal(t, id, 1)
	assert.Equal(t, sku, "abcde")
}

func Example_CreateProduct_CorrectOutputWithSku() {
	printProduct(os.Stdout, &Product{1, "laptop", 0, sql.NullString{String: "abcde", Valid: true}})

	//Output:
	//ID  |Name                      |Price         |Sku                       |
	//1   |laptop                    |0.00          |abcde                     |

}

func ExampleCreateProductCorrectOutputWithoutSku() {
	printProduct(os.Stdout, &Product{1, "laptop", 0, sql.NullString{String: "", Valid: false}})

	//Output:
	//ID  |Name                      |Price         |Sku                       |
	//1   |laptop                    |0.00          |                          |
}

// 3. Command Syntax: `store create-order <customer_id> <product_id>`

// This command should add a new order to the orders table.
// The customer_id and product_id arguments are required.

// Example:

//	> store create-order 1 3
//	ORDER_ID      CUSTOMER_ID         PRODUCT_ID
//	1             1                   3
//
// 1. creating new order with all inputs valid
// 2. creating new order missing customer id
// 3. creating new order missing product id
// 4. checking correct print status
// 5. creating new order with foreign key violation in customer-id
// 6. creating new order with foreign key violation in product-id
// 7. creating new order with customer_id not an integer
// 8. creating new order with product_id not an integer
func TestCreateNewOrderMissingProductId(t *testing.T) {

	err := createOrdersData([][]string{}, [][]string{}, [][]string{{"store", "create-order", "1"}})
	assert.ErrorContains(t, err, "Must specify product_id and customer_id")
}

func TestCreateNewOrderMissingCustomerId(t *testing.T) {
	err := createOrdersData([][]string{}, [][]string{}, [][]string{{"store", "create-order", "12"}})
	assert.ErrorContains(t, err, "Must specify product_id and customer_id")
}

func TestCreateNewOrderCustomerIdIsInt(t *testing.T) {

	err := createOrdersData([][]string{}, [][]string{}, [][]string{{"store", "create-order", "1", "4.5"}})
	assert.ErrorContains(t, err, "Must be valid integer")
}

func TestCreateNewOrderProductIdIsInt(t *testing.T) {
	err := createOrdersData([][]string{}, [][]string{}, [][]string{{"store", "create-order", "1.5", "4"}})
	assert.ErrorContains(t, err, "Must be valid integer")
}

func TestCreateNewOrderWithValidInputs(t *testing.T) {
	err := tearDownOrders()
	assert.NilError(t, err)

	db, err := connectDB(dbName)
	assert.NilError(t, err)
	defer db.Close()

	err = createOrdersData([][]string{{"store", "create-customer", "vivek.s@outreach.io", "WA"}}, [][]string{{"store", "create-product", "--sku=abcde", "banana", "5"}}, [][]string{{"store", "create-order", "1", "1"}})

	var orderID int
	var createdAt sql.NullString
	var customerID int
	var productID int

	rows := db.QueryRow("SELECT * FROM Orders ORDER BY ID DESC LIMIT 1")
	assert.NilError(t, rows.Err())

	err = rows.Scan(&orderID, &createdAt, &customerID, &productID)
	assert.NilError(t, err)
	assert.Equal(t, productID, 1)
	assert.Equal(t, customerID, 1)
	assert.Equal(t, orderID, 1)
}

func TestCreateNewOrder_ValidInputs_CustomerForeignKey_Error(t *testing.T) {
	err := tearDownOrders()
	assert.NilError(t, err)

	err = createOrdersData([][]string{}, [][]string{{"store", "create-product", "--sku=abcde", "banana", "5"}}, [][]string{{"store", "create-order", "1", "1"}})

	assert.ErrorContains(t, err, "Customer ID does not exist")
}

func TestCreateNewOrder_ValidInputs_ProductForeignKey_Error(t *testing.T) {
	err := tearDownOrders()
	assert.NilError(t, err)

	err = createOrdersData([][]string{{"store", "create-customer", "vivek.s@outreach.io", "WA"}}, [][]string{}, [][]string{{"store", "create-order", "1", "1"}})

	assert.ErrorContains(t, err, "Product ID does not exist")
}

func Example_CreateOrder_WithCorrectOutput() {
	printOrder(os.Stdout, &Order{1, sql.NullString{}, 1, 1})

	//Output:
	//OrderID    |ProductID    |CustomerID    |
	//1          |1            |1             |
}

// ### We should have show commands for customers, products, and orders:

// 1. Command syntax: `store show-orders [--name=<name>]`

// This command should display all of the orders in the orders table.
// This command should have an optional --name flag which allows users to filter by name.
// This parameter should support wildcard matching (e.g. 'foo%')

// Examples:

//       > store show-products
//       PRODUCT_ID        NAME       PRICE       SKU
//       1                 Bananas     3.99       aa1
//       2                 Milk        4.50
//       3                 Cookies     1.99       ab2

//       > store show-products --name "M%"
//       PRODUCT_ID        NAME       PRICE       SKU
//       2                 Milk        4.50
//       5                 Markers     8.99       xi1

// 1. test output of show-products no flag
// 2. test output of show-products with name flag

func Example_ShowProductsNoFlag() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err := tearDownProducts()
	if err != nil {
		log.Fatal(err)
	}
	db, err := connectDB(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createProductsDataV3([][]any{{"laptop", 25, "abcde"}, {"book", 12.5, "bcd"}})
	app := &cli.App{
		Commands: []*cli.Command{
			newShowProductCommand(&db, ctx),
		},
	}
	app.Run([]string{"store", "show-products"})
	//Output:
	//ID  |Name                      |Price         |Sku                       |
	//1   |laptop                    |25.00         |abcde                     |
	//2   |book                      |12.50         |bcd                       |
}

func Example_ShowProductNameFlag() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err := tearDownProducts()
	if err != nil {
		log.Fatal(err)
	}
	db, err := connectDB(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createProductsDataV3([][]any{{"laptop", 25, "abcde"}, {"book", 12.5, "bcd"}})
	app := &cli.App{
		Commands: []*cli.Command{
			newShowProductCommand(&db, ctx),
		},
	}
	app.Run([]string{"store", "show-products", "--name=laptop"})
	//Output:
	//ID  |Name                      |Price         |Sku                       |
	//1   |laptop                    |25.00         |abcde                     |

}

// 2. Command Syntax: `store show-customers [--email=<email>] [--state=<state>]`

// Show customers should display all of the customers in the customers table.
// This command should have an optional --email flag which allows users to filter by email.
// This command should have an optional --state flag which allows users to filter by state.
// These flags should support wildcard matching (e.g. 'foo%')

// Examples:

//     > store show-users
//     CUSTOMER_ID       EMAIL                       STATE
//     1                 zack.patrick@outreach.io    WA
//     2                 kevin.kerr@outreach.io      WA
//     3                 foo@bar.com                 CA

//     > store show-users --email "%@outreach.io"
//     CUSTOMER_ID       EMAIL                       STATE
//     1                 zack.patrick@outreach.io    WA
//     2                 kevin.kerr@outreach.io      WA

//     > store show-users --state WA
//     CUSTOMER_ID       EMAIL                       STATE
//     1                 zack.patrick@outreach.io    WA
//     2                 kevin.kerr@outreach.io      WA
// 1. test output of show-users
// 2. test output of show-users with email param
// 3. test output of show-users with state param
// 4. test output of show-users with both email and state param

func Example_ShowCustomersNoFlag() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err := tearDownCustomers()
	if err != nil {
		log.Fatal(err)
	}
	db, err := connectDB(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createCustomersDataV3([][]string{{"vivek.shah@oureach.io", "WA"}, {"vivek.s@outlook.com", "MN"}})
	app := &cli.App{
		Commands: []*cli.Command{
			newShowCustomerCommand(&db, ctx),
		},
	}
	app.Run([]string{"store", "show-customers"})
	//Output:
	//ID  |Email                                              |State |
	//1   |vivek.shah@oureach.io                              |WA    |
	//2   |vivek.s@outlook.com                                |MN    |
}

func Example_ShowCustomersStateFlag() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err := tearDownCustomers()
	if err != nil {
		log.Fatal(err)
	}
	db, err := connectDB(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createCustomersDataV3([][]string{{"vivek.shah@oureach.io", "WA"}, {"vivek.s@outlook.com", "MN"}})
	app := &cli.App{
		Commands: []*cli.Command{
			newShowCustomerCommand(&db, ctx),
		},
	}
	app.Run([]string{"store", "show-customers", "--state=WA"})
	//Output:
	//ID  |Email                                              |State |
	//1   |vivek.shah@oureach.io                              |WA    |
}

func Example_ShowCustomersEmailFlag() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err := tearDownCustomers()
	if err != nil {
		log.Fatal(err)
	}
	db, err := connectDB(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createCustomersDataV3([][]string{{"vivek.shah@outreach.io", "WA"}, {"vivek.s@outlook.com", "MN"}})
	app := &cli.App{
		Commands: []*cli.Command{
			newShowCustomerCommand(&db, ctx),
		},
	}
	app.Run([]string{"store", "show-customers", "--email=outreach"})
	//Output:
	//ID  |Email                                              |State |
	//1   |vivek.shah@outreach.io                             |WA    |
}

func Example_ShowCustomers_EmailFlag_StateFlag() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err := tearDownCustomers()
	if err != nil {
		log.Fatal(err)
	}
	db, err := connectDB(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createCustomersDataV3([][]string{{"vivek.shah@outreach.io", "WA"}, {"vivek.s@outlook.com", "MN"}, {"v.s@outreach.io", "WA"}})
	app := &cli.App{
		Commands: []*cli.Command{
			newShowCustomerCommand(&db, ctx),
		},
	}
	app.Run([]string{"store", "show-customers", "-state=WA", "--email=vivek"})
	//Output:
	//ID  |Email                                              |State |
	//1   |vivek.shah@outreach.io                             |WA    |
}

/*
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

*/

// 1. show-orders with no flag
// 2. show-orders with customer-id flag
// 3. show-orders with product-id flag
// 4. show-orders with both flags

func Example_ShowOrdersNoFlag() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err := tearDownOrders()
	if err != nil {
		log.Fatal(err)
	}
	db, err := connectDB(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createOrdersDataV3([][]int{{1, 1}, {2, 2}})

	app := &cli.App{
		Commands: []*cli.Command{
			newShowOrderCommand(&db, ctx),
		},
	}
	app.Run([]string{"store", "show-orders"})
	//Output:
	//OrderID    |ProductID    |CustomerID    |
	//1          |1            |1             |
	//2          |2            |2             |
}

func Example_ShowOrdersCustomerIDFlag() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err := tearDownOrders()
	if err != nil {
		log.Fatal(err)
	}
	db, err := connectDB(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createOrdersDataV3([][]int{{1, 1}, {2, 2}})

	app := &cli.App{
		Commands: []*cli.Command{
			newShowOrderCommand(&db, ctx),
		},
	}
	app.Run([]string{"store", "show-orders", "--customer-id=2"})
	//Output:
	//OrderID    |ProductID    |CustomerID    |
	//2          |2            |2             |

}

func Example_ShowOrdersProductIDFlag() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err := tearDownOrders()
	if err != nil {
		log.Fatal(err)
	}
	db, err := connectDB(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createOrdersDataV3([][]int{{1, 1}, {2, 2}})

	app := &cli.App{
		Commands: []*cli.Command{
			newShowOrderCommand(&db, ctx),
		},
	}
	app.Run([]string{"store", "show-orders", "--product-id=1"})
	//Output:
	//OrderID    |ProductID    |CustomerID    |
	//1          |1            |1             |

}

func Example_ShowOrdersProductIDFlag_CustomerIDFlag() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err := tearDownOrders()
	if err != nil {
		log.Fatal(err)
	}
	db, err := connectDB(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createOrdersDataV3([][]int{{1, 1}, {2, 2}, {1, 2}})

	app := &cli.App{
		Commands: []*cli.Command{
			newShowOrderCommand(&db, ctx),
		},
	}
	app.Run([]string{"store", "show-orders", "-product-id=2", "--customer-id=1"})
	//Output:
	//OrderID    |ProductID    |CustomerID    |
	//3          |2            |1             |

}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/*
Testing main.go with multiple orgs

Things to test:
1. Running migrations - are tables created succesfully in all orgs, with correct columns
2. Connecting to a specific org
	a. Org that exists
	b. Org that doesn't exist
	c. Default org
3. All the previous tests on these new orgs, ensuring they are editing the correct databases and not any others

*/

func TestMigrations(t *testing.T) {
	orgs := []*migration.OrgMigrationState{{Name: "microsoft", LastRanMigrationID: -1}, {Name: "google", LastRanMigrationID: -1}, {Name: "default", LastRanMigrationID: -1}}
	state := &migration.MigrationState{Orgs: orgs}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	migration.SaveMigrationState(ctx, state, migration.DefaultMigrationStatePath)
	db, err := connectDB("microsoft")
	assert.NilError(t, err)
	defer db.Close()

	_, err = db.Exec("DROP TABLE IF EXISTS Orders, Products, Customers")
	assert.NilError(t, err)

	db, err = connectDB("google")
	assert.NilError(t, err)
	defer db.Close()

	_, err = db.Exec("DROP TABLE IF EXISTS Orders, Products, Customers")
	assert.NilError(t, err)

	app := &cli.App{
		Commands: []*cli.Command{
			runMigrations(ctx),
		},
	}
	app.Run([]string{"store", "run-migrations"})
	db, err = connectDB("microsoft")
	assert.NilError(t, err)
	defer db.Close()

	QueryRows(db, t)
	db, err = connectDB("google")
	assert.NilError(t, err)
	defer db.Close()

	QueryRows(db, t)

}

func QueryRows(db *sql.DB, t *testing.T) {
	_, err := db.Query("SELECT * FROM Orders")
	assert.NilError(t, err)
	_, err = db.Query("SELECT * FROM Customers")
	assert.NilError(t, err)
	_, err = db.Query("SELECT * FROM Products")
	assert.NilError(t, err)
	_, err = db.Query("SELECT sku FROM Products")
	assert.NilError(t, err)
}

func runMigrationsHelper(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	app := &cli.App{
		Commands: []*cli.Command{
			runMigrations(ctx),
		},
	}
	app.Run([]string{"store", "run-migrations"})
}

func TestConnectingToDefault(t *testing.T) {
	var db *sql.DB
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "org",
				Usage: "org to connect to",
				Value: "default",
			},
		},
		Before: func(cCtx *cli.Context) error {
			org := cCtx.String("org")
			log.Println("connecting to org:", org)

			orgDB, err := connectDB(org)
			if err != nil {
				return err
			}

			db = orgDB
			return nil
		},
	}
	assert.NilError(t, app.Run([]string{"store"}))
	_ = db
}

func TestConnectingToMicrosoft(t *testing.T) {
	var db *sql.DB
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "org",
				Usage: "org to connect to",
				Value: "default",
			},
		},
		Before: func(cCtx *cli.Context) error {
			org := cCtx.String("org")
			log.Println("connecting to org:", org)

			orgDB, err := connectDB(org)
			if err != nil {
				return err
			}

			db = orgDB
			return nil
		},
	}
	assert.NilError(t, app.Run([]string{"store", "--org=microsoft"}))
	_ = db
}

func TestConnectingToGoogle(t *testing.T) {
	var db *sql.DB
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "org",
				Usage: "org to connect to",
				Value: "default",
			},
		},
		Before: func(cCtx *cli.Context) error {
			org := cCtx.String("org")
			log.Println("connecting to org:", org)

			orgDB, err := connectDB(org)
			if err != nil {
				return err
			}

			db = orgDB
			return nil
		},
	}
	assert.NilError(t, app.Run([]string{"store", "--org=google"}))
	_ = db
}

func TestConnectingToNonExistingDB(t *testing.T) {
	var db *sql.DB
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "org",
				Usage: "org to connect to",
				Value: "default",
			},
		},
		Before: func(cCtx *cli.Context) error {
			org := cCtx.String("org")
			log.Println("connecting to org:", org)

			orgDB, err := connectDB(org)
			if err != nil {
				return err
			}

			db = orgDB
			return nil
		},
	}
	assert.ErrorContains(t, app.Run([]string{"store", "--org=abc"}), "Error 1049 (42000): Unknown database 'store_abc'")
	_ = db
}

func Test_SuiteOfTests_Microsoft(t *testing.T) {
	runMigrationsHelper(t)
	dbName = "microsoft"
	Test_RunSuiteOfTests(t)
}

func Test_SuiteOfTests_Google(t *testing.T) {
	runMigrationsHelper(t)
	dbName = "google"
	Test_RunSuiteOfTests(t)
}

func Test_SuiteOfTests_Default(t *testing.T) {
	runMigrationsHelper(t)
	dbName = "default"
	Test_RunSuiteOfTests(t)
}

func Test_RunSuiteOfTests(t *testing.T) {
	TestCreateCustomer_withAMissingEmailArg_returnsAnError(t)
	TestCreateCustomer_withAMissingStateArg_returnsAnError(t)
	TestCreateCustomer_wthAInvalidStateAbbreviation_returnsAnError(t)
	TestCreateCustomer_withANonTwoLetterStateCode_returnsAnError(t)
	TestCreateCustomer_withValidInput_entersDatabaseCorrectly(t)
	TestCreateMultipleCustomer_withValidInput_entersDatabaseCorrectly(t)
	TestCreateProduct_noName(t)
	TestCreateProduct_noPrice(t)
	TestCreateProduct_validInputNoSku(t)
	TestCreateProduct_validInputWithSku(t)
	TestCreateProduct_decimalPrice(t)
	TestCreateNewOrderMissingProductId(t)
	TestCreateNewOrderMissingCustomerId(t)
	TestCreateNewOrderCustomerIdIsInt(t)
	TestCreateNewOrderProductIdIsInt(t)
	TestCreateNewOrderWithValidInputs(t)
	TestCreateNewOrder_ValidInputs_CustomerForeignKey_Error(t)
	TestCreateNewOrder_ValidInputs_ProductForeignKey_Error(t)
	Example_ShowProductNameFlag()
	Example_ShowProductsNoFlag()
	Example_ShowCustomersNoFlag()
	Example_ShowCustomersStateFlag()
	Example_ShowCustomersEmailFlag()
	Example_ShowCustomers_EmailFlag_StateFlag()
	Example_ShowOrdersNoFlag()
	Example_ShowOrdersCustomerIDFlag()
	Example_ShowOrdersProductIDFlag()
	Example_ShowOrdersProductIDFlag_CustomerIDFlag()
}
