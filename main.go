package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v2"
)

var states = map[string]bool{
	"AL": true,
	"AK": true,
	"AZ": true,
	"AR": true,
	"AS": true,
	"CA": true,
	"CO": true,
	"CT": true,
	"DE": true,
	"DC": true,
	"FL": true,
	"GA": true,
	"GU": true,
	"HI": true,
	"ID": true,
	"IL": true,
	"IN": true,
	"IA": true,
	"KS": true,
	"KY": true,
	"LA": true,
	"ME": true,
	"MD": true,
	"MA": true,
	"MI": true,
	"MN": true,
	"MS": true,
	"MO": true,
	"MT": true,
	"NE": true,
	"NV": true,
	"NH": true,
	"NJ": true,
	"NM": true,
	"NY": true,
	"NC": true,
	"ND": true,
	"MP": true,
	"OH": true,
	"OK": true,
	"OR": true,
	"PA": true,
	"PR": true,
	"RI": true,
	"SC": true,
	"SD": true,
	"TN": true,
	"TX": true,
	"VT": true,
	"UT": true,
	"VA": true,
	"VI": true,
	"WA": true,
	"WV": true,
	"WI": true,
	"WY": true,
}

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

	// rows, err := db.QueryContext(ctx, "SELECT * FROM customers")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var id int
	// 	var email string
	// 	var state string

	// 	if err := rows.Scan(&id, &email, &state); err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	fmt.Println(id, email, state)
	// }

	// if err := rows.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	app := &cli.App{
		Name: "store",
		Commands: []*cli.Command{
			{
				Name:  "create-customer",
				Usage: "Creates a new customer to go in the customers database, must specify email and state(2 letter code)",
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() < 2 {
						return errors.New("Must specify email and state")
					}
					email := cCtx.Args().Get(0)
					state := cCtx.Args().Get(1)
					if len(state) != 2 {
						return errors.New("State length must be 2")
					}
					_, ok := states[strings.ToUpper(state)]
					if !ok {
						return errors.New("State must be a valid U.S. State or Territory")
					}

					insertStatement := "INSERT INTO Customers (email, state) VALUES (?, ?)"
					res, err := db.Exec(insertStatement, email, state)

					if err != nil {
						return err
					}
					customerID, err := res.LastInsertId()
					if err != nil {
						return err
					}
					w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
					fmt.Fprintln(w, "Customer_id\tEmail\tState\t")
					fmt.Fprintf(w, "%d\t%s\t%s\t\n", customerID, email, state)
					w.Flush()

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
					price, err := strconv.Atoi(cCtx.Args().Get(1))
					if err != nil {
						return err
					}
					sku := cCtx.String("sku")
					if sku != "" {
						insertStatement := "INSERT INTO Products (name, price, sku) VALUES (?, ?, ?)"
						res, err := db.Exec(insertStatement, name, price, sku)
						if err != nil {
							return err
						}

						productId, err := res.LastInsertId()
						if err != nil {
							return err
						}
						w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
						fmt.Fprintln(w, "Product_id\tName\tPrice\tSku\t")
						fmt.Fprintf(w, "%d\t%s\t%v\t%v\t\n", productId, name, price, sku)
						w.Flush()
					} else {
						insertStatement := "INSERT INTO Products (name, price) VALUES (?, ?)"
						res, err := db.Exec(insertStatement, name, price)
						if err != nil {
							return err
						}
						productId, err := res.LastInsertId()
						if err != nil {
							return err
						}
						w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
						fmt.Fprintln(w, "Product_id\tName\tPrice\t")
						fmt.Fprintf(w, "%d\t%s\t%v\t\n", productId, name, price)
						w.Flush()
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
					pID, err := strconv.Atoi(os.Args[2])
					if err != nil {
						return err
					}
					cID, err := strconv.Atoi(os.Args[3])
					if err != nil {
						return err
					}
					insertStatement := "INSERT INTO Orders (customer_id, product_id) VALUES (?, ?)"
					res, err := db.Exec(insertStatement, cID, pID)
					if err != nil {
						return err
					}
					orderID, err := res.LastInsertId()
					if err != nil {
						return err
					}

					w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
					fmt.Fprintln(w, "Order_id\tProduct_id\tCustomer_id\t")
					fmt.Fprintf(w, "%v\t%v\t%v\t\n", orderID, pID, cID)
					w.Flush()

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
						w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)

						rows, err := db.QueryContext(ctx, "SELECT * FROM PRODUCTS WHERE Name LIKE CONCAT('%', ?, '%')", name)
						if err != nil {
							return err
						}
						defer rows.Close()

						productPrintHelper(rows, w)
						w.Flush()
					} else {

						rows, err := db.QueryContext(ctx, "SELECT * FROM PRODUCTS")
						w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
						if err != nil {
							return err
						}
						defer rows.Close()
						productPrintHelper(rows, w)
						w.Flush()

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

						w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
						rows, err := db.QueryContext(ctx, "SELECT * FROM Customers WHERE Email LIKE CONCAT('%', ?, '%') AND State LIKE CONCAT('%', ?, '%')", email, state)
						if err != nil {
							return err
						}
						defer rows.Close()

						customersPrintHelper(rows, w)
						w.Flush()

					} else if cCtx.String("email") != "" {
						email := cCtx.String("email")
						w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
						rows, err := db.QueryContext(ctx, "SELECT * FROM Customers WHERE Email LIKE CONCAT('%', ?, '%')", email)
						if err != nil {
							return err
						}
						defer rows.Close()

						customersPrintHelper(rows, w)
						w.Flush()
					} else if cCtx.String("state") != "" {
						state := cCtx.String("state")
						w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
						rows, err := db.QueryContext(ctx, "SELECT * FROM Customers WHERE State LIKE CONCAT('%', ?, '%')", state)
						if err != nil {
							return err
						}
						defer rows.Close()

						customersPrintHelper(rows, w)
						w.Flush()
					} else {

						w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
						rows, err := db.QueryContext(ctx, "SELECT * FROM Customers")
						if err != nil {
							return err
						}
						defer rows.Close()

						customersPrintHelper(rows, w)
						w.Flush()
					}
					return nil
				},
			},
			{
				Name:  "show-orders",
				Usage: "displays all the orders within the orders database with an optional customerId and productId filter",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "customer-id",
						Usage: "the customer-id of the customer",
					},
					&cli.IntFlag{
						Name:  "product-id",
						Usage: "the product-id of the product",
					},
				}, Action: func(cCtx *cli.Context) error {
					if cCtx.Int("customer-id") != 0 && cCtx.Int("product-id") != 0 {
						cId := cCtx.Int("customer-id")
						pId := cCtx.Int("product-id")
						w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
						rows, err := db.QueryContext(ctx, "SELECT * FROM ORDERS WHERE customer_id=? AND product_id=?", cId, pId)
						if err != nil {
							return err
						}
						defer rows.Close()

						printHelper(rows, w)
						w.Flush()
					} else if cCtx.Int("customer-id") != 0 {
						cId := cCtx.Int("customer-id")

						w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
						rows, err := db.QueryContext(ctx, "SELECT * FROM ORDERS WHERE customer_id =?", cId)
						if err != nil {
							return err
						}
						defer rows.Close()

						printHelper(rows, w)
						w.Flush()
					} else if cCtx.Int("product-id") != 0 {

						pId := cCtx.Int("product-id")

						w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
						rows, err := db.QueryContext(ctx, "SELECT * FROM ORDERS WHERE product_id =?", pId)
						if err != nil {
							return err
						}
						defer rows.Close()

						printHelper(rows, w)
						w.Flush()
					} else {

						w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
						rows, err := db.QueryContext(ctx, "SELECT * FROM ORDERS ")
						if err != nil {
							return err
						}
						defer rows.Close()

						printHelper(rows, w)
						w.Flush()
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

func printHelper(rows *sql.Rows, w *tabwriter.Writer) error {
	fmt.Fprintln(w, "Order_ID\tCustomer_ID\tProduct_ID\t")
	for rows.Next() {
		var id int
		var date sql.NullString
		var cID int
		var pID int
		var sku sql.NullString

		if err := rows.Scan(&id, &date, &cID, &pID, &sku); err != nil {
			return err
		}

		fmt.Fprintf(w, "%d\t%d\t%d\t\n", id, cID, pID)

	}
	if rows.Err() != nil {
		return rows.Err()
	}
	return nil
}

func customersPrintHelper(rows *sql.Rows, w *tabwriter.Writer) error {
	fmt.Fprintln(w, "Customer_ID\tEmail\tState")
	for rows.Next() {
		var id int
		var email string
		var state string
		if err := rows.Scan(&id, &email, &state); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d\t%v\t%v\n", id, email, state)

	}
	if rows.Err() != nil {
		return rows.Err()
	}
	return nil
}

func productPrintHelper(rows *sql.Rows, w *tabwriter.Writer) error {
	fmt.Fprintln(w, "Product_id\tName\tPrice\tSku\t")
	for rows.Next() {
		var id int
		var name string
		var price float64
		var sku sql.NullString

		if err := rows.Scan(&id, &name, &price, &sku); err != nil {
			return err
		}
		if !sku.Valid {

			fmt.Fprintf(w, "%d\t%v\t%v\t\n", id, name, price)
		} else {
			fmt.Fprintf(w, "%d\t%v\t%v\t%v\n", id, name, price, sku.String)

		}
	}
	if rows.Err() != nil {
		return rows.Err()
	}
	return nil
}
