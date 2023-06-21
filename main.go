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

type Customer struct {
	ID    int
	Email string
	State string
}

type Product struct {
	ID    int
	name  string
	price float64
	sku   sql.NullString
}

type Order struct {
	ID         int
	created_at sql.NullString
	cID        int
	pID        int
	sku        sql.NullString
}

func NewCustomer(rows *sql.Rows) (*Customer, error) {
	var c Customer
	if err := rows.Scan(&c.ID, &c.Email, &c.State); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Customer) Print(w *tabwriter.Writer) {
	fmt.Fprintf(w, "%d\t%v\t%v\n", c.ID, c.Email, c.State)
}

func NewProduct(rows *sql.Rows) (*Product, error) {
	var p Product
	if err := rows.Scan(&p.ID, &p.name, &p.price, &p.sku); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Product) Print(w *tabwriter.Writer) {
	if p.sku.Valid {
		fmt.Fprintf(w, "%d\t%v\t%v\t%v\t\n", p.ID, p.name, p.price, p.sku.String)
	} else {
		fmt.Fprintf(w, "%d\t%v\t%v\t%v\t\n", p.ID, p.name, p.price, "")
	}
}

func NewOrder(rows *sql.Rows) (*Order, error) {
	var o Order
	if err := rows.Scan(&o.ID, &o.created_at, &o.cID, &o.pID, &o.sku); err != nil {
		return nil, err
	}
	return &o, nil
}

func (o *Order) Print(w *tabwriter.Writer) {
	fmt.Fprintf(w, "%d\t%d\t%d\t\n", o.ID, o.cID, o.pID)
}

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
					args := []any{}

					w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
					var insertStatement string
					if sku != "" {
						insertStatement = "INSERT INTO Products (name, price, sku) VALUES (?, ?, ?)"

					} else {

						insertStatement = "INSERT INTO Products (name, price) VALUES (?, ?)"

					}
					args = append(args, name, price, sku)

					err = productsInsertHelper(db, args, insertStatement, w)
					if err != nil {
						return err
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
						Usage: "the name of the product",
					},
				},
				Action: func(cCtx *cli.Context) error {
					w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
					statement := "SELECT * FROM PRODUCTS WHERE Name LIKE CONCAT('%', ?, '%')"
					name := cCtx.String("name")
					productHelper(w, db, statement, name, ctx)
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
					w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
					statement := "SELECT * FROM Customers WHERE Email LIKE CONCAT('%', ?, '%') AND State LIKE CONCAT('%', ?, '%')"
					email := cCtx.String("email")
					state := cCtx.String("state")
					err := customerHelper(w, db, statement, email, state, ctx)
					if err != nil {
						return err
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
					w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
					args := []any{}
					cID := cCtx.Int("customer-id")
					pID := cCtx.Int("product-id")
					var statement string
					if cCtx.Int("customer-id") != 0 && cCtx.Int("product-id") != 0 {

						statement = "SELECT * FROM ORDERS WHERE customer_id=? AND product_id=?"
						args = append(args, cID, pID)

					} else if cCtx.Int("customer-id") != 0 {

						statement = "SELECT * FROM ORDERS WHERE customer_id =?"
						args = append(args, cID)

					} else if cCtx.Int("product-id") != 0 {

						statement = "SELECT * FROM ORDERS WHERE product_id =?"
						args = append(args, pID)

					} else {

						statement = "SELECT * FROM ORDERS "

					}
					err = ordersHelper(w, db, statement, args, ctx)
					if err != nil {
						return err
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

func orderPrintHelper(orders []*Order, w *tabwriter.Writer) error {
	fmt.Fprintln(w, "Order_ID\tCustomer_ID\tProduct_ID\t")
	for _, o := range orders {
		o.Print(w)
	}
	return nil
}

func ordersHelper(w *tabwriter.Writer, db *sql.DB, statement string, args []any, ctx context.Context) error {
	var rows *sql.Rows
	var err error
	if len(args) == 1 {
		rows, err = db.QueryContext(ctx, statement, args[0])
	} else if len(args) == 2 {
		rows, err = db.QueryContext(ctx, statement, args[0], args[1])
	} else {
		rows, err = db.QueryContext(ctx, statement)
	}
	if err != nil {
		return err
	}
	defer rows.Close()
	orders := []*Order{}
	for rows.Next() {
		o, err := NewOrder(rows)
		if err != nil {
			return err
		}
		orders = append(orders, o)
	}
	if rows.Err() != nil {
		return err
	}
	err = orderPrintHelper(orders, w)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

func productPrintHelper(products []*Product, w *tabwriter.Writer) error {
	fmt.Fprintln(w, "Product_id\tName\tPrice\tSku\t")

	for _, p := range products {
		p.Print(w)
	}
	return nil
}

func productHelper(w *tabwriter.Writer, db *sql.DB, statement string, name string, ctx context.Context) error {
	rows, err := db.QueryContext(ctx, statement, name)
	if err != nil {
		return err
	}
	defer rows.Close()

	products := []*Product{}
	for rows.Next() {
		p, err := NewProduct(rows)
		if err != nil {
			return err
		}
		products = append(products, p)
	}
	if rows.Err() != nil {
		return err
	}

	err = productPrintHelper(products, w)
	if err != nil {
		return err
	}
	w.Flush()
	return nil

}

func customerHelper(w *tabwriter.Writer, db *sql.DB, statement string, email string, state string, ctx context.Context) error {
	rows, err := db.QueryContext(ctx, statement, email, state)
	if err != nil {
		return err
	}
	defer rows.Close()

	customers := []*Customer{}
	for rows.Next() {
		c, err := NewCustomer(rows)
		if err != nil {
			return err
		}
		customers = append(customers, c)
	}
	if rows.Err() != nil {
		return err
	}

	err = customersPrintHelper(customers, w)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func customersPrintHelper(customers []*Customer, w *tabwriter.Writer) error {
	fmt.Fprintln(w, "Customer_id\tEmail\tState")
	for _, c := range customers {
		c.Print(w)
	}
	return nil
}

func productsInsertHelper(db *sql.DB, args []any, statement string, w *tabwriter.Writer) error {
	var res sql.Result
	var err error
	if args[2] != "" {
		res, err = db.Exec(statement, args[0], args[1], args[2])
	} else {
		res, err = db.Exec(statement, args[0], args[1])
	}
	if err != nil {
		return err
	}

	productId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Product_id\tName\tPrice\tSku\t")

	fmt.Fprintf(w, "%d\t%s\t%v\t%v\t\n", productId, args[0], args[1], args[2])
	w.Flush()
	return nil
}
