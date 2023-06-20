package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v2"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	cfg := mysql.Config{
		User:   "admin",
		Passwd: "password123",
		Addr:   "localhost",
		DBName: "store",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("failed to open database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	rows, err := db.QueryContext(ctx, "SELECT * FROM customers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var email string
		var state string

		if err := rows.Scan(&id, &email, &state); err != nil {
			log.Fatal(err)
		}

		fmt.Println(id, email, state)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name: "store",
		Commands: []*cli.Command{
			{
				Name:  "example",
				Usage: "show a bunch of stuff",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "the name we want to display",
					},
				},
				Action: func(cCtx *cli.Context) error {
					name := "<none>"
					if cCtx.String("name") != "" {
						name = cCtx.String("name")
					}

					fmt.Println("the first argument you added: ", cCtx.Args().First())
					fmt.Println("the name flag you passed: ", name)
					return nil
				},
			},
			{
				Name:  "create-customer",
				Usage: "Creates a new customer to go in the customers database, must specify email and state(2 letter code)",
				Action: func(cCtx *cli.Context) error {
					if len(os.Args) != 4 {
						return errors.New("Must specify email and state")
					}
					email := os.Args[2]
					state := os.Args[3]
					insertStatement := "INSERT INTO Customers (email, state) VALUES (?, ?)"
					_, err := db.Exec(insertStatement, email, state)
					if err != nil {
						log.Fatal(err)
					}
					idQuery := "SELECT LAST_INSERT_ID()"
					var customerID int
					err = db.QueryRow(idQuery).Scan(&customerID)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("CUSTOMER_ID\tEMAIL\t\t\tSTATE")
					fmt.Printf("%d\t\t%s\t%s\n", customerID, email, state)

					return nil
				},
			},
			{
				Name:  "create-product",
				Usage: "Creates a new product to go in the products database, must specify name and price",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "sku",
						Usage: "the sku of the product",
					},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() < 2 {
						return errors.New("Must specify name and price")
					}

					name := cCtx.Args().Get(0)
					price, err1 := strconv.Atoi(cCtx.Args().Get(1))
					if err1 != nil {
						log.Fatal(err1)
					}
					sku := cCtx.String("sku")
					if sku != "" {
						insertStatement := "INSERT INTO Products (name, price, sku) VALUES (?, ?, ?)"
						_, err := db.Exec(insertStatement, name, price, sku)
						if err != nil {
							log.Fatal(err)
						}
						idQuery := "SELECT LAST_INSERT_ID()"
						var productId int
						err = db.QueryRow(idQuery).Scan(&productId)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println("Product_ID\tName\t\t\tPrice\t\t\tSku")
						fmt.Printf("%d\t\t%s\t%v\t%s\n", productId, name, price, sku)
					} else {
						insertStatement := "INSERT INTO Products (name, price) VALUES (?, ?)"
						_, err := db.Exec(insertStatement, name, price)
						if err != nil {
							log.Fatal(err)
						}
						idQuery := "SELECT LAST_INSERT_ID()"
						var productId int
						err = db.QueryRow(idQuery).Scan(&productId)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println("PRODUCT_ID\tNAME\t\t\tPRICE")
						fmt.Printf("%d\t\t%s\t%v\n", productId, name, price)
					}
					return nil

				},
			},
			{
				Name:  "create-order",
				Usage: "Creates a new order to go in the order database, must specify customer_id and product_id",
				Action: func(cCtx *cli.Context) error {
					if len(os.Args) != 4 {
						return errors.New("Must specify product_id and customer_id")
					}
					pId, error1 := strconv.Atoi(os.Args[2])
					cId, err1 := strconv.Atoi(os.Args[3])
					if error1 != nil {
						log.Fatal(error1)
					}
					if err1 != nil {
						log.Fatal(err1)
					}
					insertStatement := "INSERT INTO Orders (customer_id, product_id) VALUES (?, ?)"
					_, err := db.Exec(insertStatement, cId, pId)
					if err != nil {
						log.Fatal(err)
					}
					idQuery := "SELECT LAST_INSERT_ID()"
					var orderId int
					err = db.QueryRow(idQuery).Scan(&orderId)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("Order_id\tProduct_id\tCustomer_id")
					fmt.Printf("%d\t\t%v\t%v\n", orderId, pId, cId)

					return nil
				},
			},
			{
				Name:  "show-products",
				Usage: "Shows the proudcts from the products database, optional flag name to filter by name",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "the name of the order",
					},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("name") != "" {
						name := cCtx.String("name")

						fmt.Println("Product_ID\tName\t\t\tPrice\t\t\tSku")
						rows, err := db.QueryContext(ctx, "SELECT * FROM PRODUCTS WHERE Name LIKE CONCAT('%', ?, '%')", name)
						if err != nil {
							log.Fatal(err)
						}
						defer rows.Close()

						for rows.Next() {
							var id int
							var name string
							var price float64
							var sku sql.NullString

							if err := rows.Scan(&id, &name, &price, &sku); err != nil {
								log.Fatal(err)
							}
							if !sku.Valid {
								fmt.Printf("%d\t\t%v\t\t\t%v\n", id, name, price)
							} else {
								fmt.Printf("%d\t\t%v\t%v\t\t\t%v\n", id, name, price, sku.String)
							}
						}
					} else {
						fmt.Println("Product_ID\tName\t\t\tPrice\t\t\tSku")
						rows, err := db.QueryContext(ctx, "SELECT * FROM PRODUCTS")
						if err != nil {
							log.Fatal(err)
						}
						defer rows.Close()

						for rows.Next() {
							var id int
							var name string
							var price float64
							var sku sql.NullString

							if err := rows.Scan(&id, &name, &price, &sku); err != nil {
								log.Fatal(err)
							}
							if !sku.Valid {
								fmt.Printf("%d\t\t%v\t\t\t%v\n", id, name, price)
							} else {
								fmt.Printf("%d\t\t%v\t%v\t\t\t%vtn", id, name, price, sku.String)
							}
						}
					}
					return nil
				},
			},
			{
				Name:  "show-customers",
				Usage: "displays all the customers inside the customers database, optional flags to filter by email and state",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "email",
						Usage: "the email of the customer",
					},
					&cli.StringFlag{
						Name:  "state",
						Usage: "the state of the customer",
					},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("email") != "" && cCtx.String("state") != "" {
						email := cCtx.String("email")
						state := cCtx.String("state")

						fmt.Println("Customer_ID\tEmail\t\t\tState")
						rows, err := db.QueryContext(ctx, "SELECT * FROM Customers WHERE Email LIKE CONCAT('%', ?, '%') AND State LIKE CONCAT('%', ?, '%')", email, state)
						if err != nil {
							log.Fatal(err)
						}
						defer rows.Close()

						for rows.Next() {
							var id int
							var email string
							var state string
							if err := rows.Scan(&id, &email, &state); err != nil {
								log.Fatal(err)
							}
							fmt.Printf("%d\t\t%v\t%v\n", id, email, state)

						}
					} else if cCtx.String("email") != "" {
						email := cCtx.String("email")
						fmt.Println("Customer_ID\tEmail\t\t\tState")
						rows, err := db.QueryContext(ctx, "SELECT * FROM Customers WHERE Email LIKE CONCAT('%', ?, '%')", email)
						if err != nil {
							log.Fatal(err)
						}
						defer rows.Close()

						for rows.Next() {
							var id int
							var email string
							var state string
							if err := rows.Scan(&id, &email, &state); err != nil {
								log.Fatal(err)
							}
							fmt.Printf("%d\t\t%v\t%v\n", id, email, state)

						}
					} else if cCtx.String("state") != "" {
						state := cCtx.String("state")
						fmt.Println("Customer_ID\tEmail\t\t\tState")
						rows, err := db.QueryContext(ctx, "SELECT * FROM Customers WHERE State LIKE CONCAT('%', ?, '%')", state)
						if err != nil {
							log.Fatal(err)
						}
						defer rows.Close()

						for rows.Next() {
							var id int
							var email string
							var state string
							if err := rows.Scan(&id, &email, &state); err != nil {
								log.Fatal(err)
							}
							fmt.Printf("%d\t\t%v\t%v\n", id, email, state)

						}
					} else {

						fmt.Println("Customer_ID\tEmail\t\t\tState")
						rows, err := db.QueryContext(ctx, "SELECT * FROM Customers")
						if err != nil {
							log.Fatal(err)
						}
						defer rows.Close()

						for rows.Next() {
							var id int
							var email string
							var state string

							if err := rows.Scan(&id, &email, &state); err != nil {
								log.Fatal(err)
							}
							fmt.Printf("%d\t\t%v\t%v\n", id, email, state)

						}
					}
					return nil
				},
			},
			{
				Name:  "show-orders",
				Usage: "displays all the orders within the orders database with an optional customerId and productId filter",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "customer-id",
						Usage: "the customer-id of the customer",
					},
					&cli.StringFlag{
						Name:  "product-id",
						Usage: "the product-id of the product",
					},
				}, Action: func(cCtx *cli.Context) error {
					if cCtx.String("customer-id") != "" && cCtx.String("product-id") != "" {
						cId := cCtx.String("customer-id")
						pId := cCtx.String("product-id")

						fmt.Println("Order_ID\tCustomer_ID\t\t\tProduct_ID")
						rows, err := db.QueryContext(ctx, "SELECT * FROM ORDERS WHERE customer_id LIKE CONCAT('%', ?, '%') AND product_id LIKE CONCAT('%', ?, '%')", cId, pId)
						if err != nil {
							log.Fatal(err)
						}
						defer rows.Close()

						for rows.Next() {
							var id int
							var date sql.NullString
							var cId string
							var pId string
							var sku sql.NullString

							if err := rows.Scan(&id, &date, &cId, &pId, &sku); err != nil {
								log.Fatal(err)
							}

							fmt.Printf("%d\t\t%v\t\t\t%v\n", id, cId, pId)

						}
					} else if cCtx.String("customer-id") != "" {
						cId := cCtx.String("customer-id")

						fmt.Println("Order_ID\tCustomer_ID\t\t\tProduct_ID")
						rows, err := db.QueryContext(ctx, "SELECT * FROM ORDERS WHERE customer_id LIKE CONCAT('%', ?, '%')", cId)
						if err != nil {
							log.Fatal(err)
						}
						defer rows.Close()

						for rows.Next() {
							var id int
							var date sql.NullString
							var cId string
							var pId string
							var sku sql.NullString

							if err := rows.Scan(&id, &date, &cId, &pId, &sku); err != nil {
								log.Fatal(err)
							}

							fmt.Printf("%d\t\t%v\t\t\t%v\n", id, cId, pId)

						}
					} else if cCtx.String("product-id") != "" {

						pId := cCtx.String("product-id")

						fmt.Println("Order_ID\tCustomer_ID\t\t\tProduct_ID")
						rows, err := db.QueryContext(ctx, "SELECT * FROM ORDERS WHERE product_id LIKE CONCAT('%', ?, '%')", pId)
						if err != nil {
							log.Fatal(err)
						}
						defer rows.Close()

						for rows.Next() {
							var id int
							var date sql.NullString
							var cId string
							var pId string
							var sku sql.NullString

							if err := rows.Scan(&id, &date, &cId, &pId, &sku); err != nil {
								log.Fatal(err)
							}

							fmt.Printf("%d\t\t%v\t\t\t%v\n", id, cId, pId)

						}
					} else {

						fmt.Println("Order_ID\tCustomer_ID\t\t\tProduct_ID")
						rows, err := db.QueryContext(ctx, "SELECT * FROM ORDERS ")
						if err != nil {
							log.Fatal(err)
						}
						defer rows.Close()

						for rows.Next() {
							var id int
							var date sql.NullString
							var cId string
							var pId string
							var sku sql.NullString

							if err := rows.Scan(&id, &date, &cId, &pId, &sku); err != nil {
								log.Fatal(err)
							}

							fmt.Printf("%d\t\t%v\t\t\t%v\n", id, cId, pId)

						}
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
