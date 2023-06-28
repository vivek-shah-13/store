package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v2"
	"github.com/vivek-shah-13/store/internal/migration"
)

type MigrationRunner struct {
	path string
}

func NewMigrationRunner(path string) *MigrationRunner {
	return &MigrationRunner{path: path}
}

func (m *MigrationRunner) Run(ctx context.Context, conn *sql.Conn, lastRanId int) (int, error) {
	files, err := m.loadMigrationFiles()
	if err != nil {
		return lastRanId, err
	}

	sort.Slice(files, func(i, j int) bool {
		id1, err := extractMigrationID(files[i])
		if err != nil {
			log.Fatal(err)
		}
		id2, err := extractMigrationID(files[j])
		if err != nil {
			log.Fatal(err)
		}

		return id1 < id2
	})

	fileSlice := files[lastRanId+1:]
	updatedID := lastRanId

	for _, file := range fileSlice {
		readFile, err := os.Open(file)
		if err != nil {
			return updatedID, err
		}
		defer readFile.Close()
		fileScanner := bufio.NewScanner(readFile)

		fileScanner.Split(bufio.ScanLines)
		for fileScanner.Scan() {

			if strings.Trim(fileScanner.Text(), "") == "" {
				continue
			}
			if _, err := conn.ExecContext(ctx, fileScanner.Text()); err != nil {
				return updatedID, fmt.Errorf("failed to execute migration %s: %w", file, err)
			}
		}
		log.Printf("Executed migration: %s\n", file)
		updatedID, err = extractMigrationID(file)
		return updatedID, err
	}

	return updatedID, nil
}

// RunAll runs migrations for every org in orgs and then returns the updated migration state.
func (m *MigrationRunner) RunAll(ctx context.Context, state *migration.MigrationState) (*migration.MigrationState, error) {
	for _, org := range state.Orgs {
		// TODO: after each successful migration run, update the org.LastRanMigrationID
		dbName := "store_" + org.Name
		db, err := connectDB(dbName)
		if err != nil {
			return nil, err
		}
		defer db.Close()
		conn, err := db.Conn(ctx)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		newID, err := m.Run(ctx, conn, org.LastRanMigrationID)
		if err != nil {
			return nil, err
		}
		org.LastRanMigrationID = newID

	}

	return state, nil
}

func extractMigrationID(file string) (int, error) {
	regex := regexp.MustCompile(`_(\d+)\.sql$`)
	matches := regex.FindStringSubmatch(file)

	if len(matches) != 2 {
		return -1, errors.New("Didn't find match")
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, errors.New("Couldn't convert to integer")
	}

	return id, nil
}

func (m *MigrationRunner) loadMigrationFiles() ([]string, error) {
	files, err := ioutil.ReadDir(m.path)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		migrationFiles = append(migrationFiles, filepath.Join(m.path, file.Name()))
	}

	return migrationFiles, nil
}

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
}

func NewCustomer(rows *sql.Rows) (*Customer, error) {
	var c Customer
	if err := rows.Scan(&c.ID, &c.Email, &c.State); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Customer) Print(w *tabwriter.Writer) {
	fmt.Fprintf(w, "%-*v\t%-*s\t%-*s\t\n", 3, c.ID, 50, c.Email, 2, c.State)
}

func NewProduct(rows *sql.Rows) (*Product, error) {
	var p Product
	if err := rows.Scan(&p.ID, &p.name, &p.price, &p.sku); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Product) Print(w *tabwriter.Writer) {
	fmt.Fprintf(w, "%-*v\t%-*s\t%-*.2f\t%-*s\t\n", 3, p.ID, 15, p.name, 13, p.price, 25, p.sku.String)
}

func NewOrder(rows *sql.Rows) (*Order, error) {
	var o Order
	if err := rows.Scan(&o.ID, &o.created_at, &o.cID, &o.pID); err != nil {
		return nil, err
	}
	return &o, nil
}

func (o *Order) Print(w *tabwriter.Writer) {
	fmt.Fprintf(w, "%-*v\t%-*v\t%-*v\t\n", 3, o.ID, 12, o.pID, 13, o.cID)
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

func connectDB(name string) (*sql.DB, error) {
	cfg := mysql.Config{
		User:   "admin",
		Passwd: "password123",
		Addr:   "localhost",
		DBName: "store_" + name,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func printCustomer(w io.Writer, customers ...*Customer) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	customersPrintHelper(customers, tw)
	tw.Flush()
}

func printOrder(w io.Writer, orders ...*Order) {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.Debug)

	orderPrintHelper(orders, tw)
	tw.Flush()
}

func printProduct(w io.Writer, products ...*Product) {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.Debug)

	productPrintHelper(products, tw)
	tw.Flush()
}

func newCreateCustomerCommand(db *sql.DB) *cli.Command {
	return &cli.Command{
		Name:      "create-customer",
		Usage:     "Creates a new customer to go in the customers database, must specify email and state(2 letter code)",
		ArgsUsage: "EMAIL STATE",
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
			printCustomer(os.Stdout, &Customer{int(customerID), email, state})

			return nil
		},
	}
}

func newCreateProductCommand(db *sql.DB) *cli.Command {
	return &cli.Command{
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
			price, err := strconv.ParseFloat(cCtx.Args().Get(1), 32)
			price = (price * 100) / 100
			if err != nil {
				return err
			}
			sku := cCtx.String("sku")

			p := Product{
				name:  name,
				price: price,
				sku:   sql.NullString{String: sku, Valid: true},
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
			var insertStatement string
			if sku != "" {
				insertStatement = "INSERT INTO Products (name, price, sku) VALUES (?, ?, ?)"

			} else {

				insertStatement = "INSERT INTO Products (name, price) VALUES (?, ?)"

			}

			return productsInsertHelper(db, p, insertStatement, w)

		},
	}
}

func newCreateOrderCommand(db *sql.DB) *cli.Command {
	return &cli.Command{
		Name:      "create-order",
		Usage:     "Creates a new order to go in the order database, must specify customer_id and product_id",
		ArgsUsage: "PRODUCT_ID CUSTOMER_ID",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() < 2 {
				return errors.New("Must specify product_id and customer_id")
			}
			pID, err := strconv.Atoi(cCtx.Args().Get(0))
			if err != nil {
				return errors.New("Must be valid integer")
			}
			cID, err := strconv.Atoi(cCtx.Args().Get(1))
			if err != nil {
				return errors.New("Must be valid integer")
			}
			insertStatement := "INSERT INTO Orders (customer_id, product_id) VALUES (?, ?)"
			res, err := db.Exec(insertStatement, cID, pID)
			if err != nil {
				row := db.QueryRow("SELECT * FROM Customers ORDER BY ID LIMIT 1")
				var id int
				var email string
				var state string

				err := row.Scan(&id, &email, &state)
				if err != nil {
					return errors.New("Customer ID does not exist")
				}
				log.Print(id)
				if id < cID {
					return errors.New("Customer ID does not exist")

				}

				return errors.New("Product ID does not exist")

			}
			orderID, err := res.LastInsertId()
			if err != nil {
				return err
			}
			printOrder(os.Stdout, &Order{int(orderID), sql.NullString{String: "", Valid: false}, cID, pID})

			return nil
		},
	}
}

func newShowCustomerCommand(db *sql.DB, ctx context.Context) *cli.Command {
	return &cli.Command{
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
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
			statement := "SELECT * FROM Customers WHERE Email LIKE CONCAT('%', ?, '%') AND State LIKE CONCAT('%', ?, '%')"
			email := cCtx.String("email")
			state := cCtx.String("state")
			err := customerHelper(w, db, statement, email, state, ctx)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func newShowProductCommand(db *sql.DB, ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "show-products",
		Usage: "Shows the products from the products database, optional flag name to filter by name",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "the name of the product",
			},
		},
		Action: func(cCtx *cli.Context) error {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
			statement := "SELECT * FROM PRODUCTS WHERE Name LIKE CONCAT('%', ?, '%')"
			name := cCtx.String("name")
			productHelper(w, db, statement, name, ctx)
			return nil
		},
	}
}

func newShowOrderCommand(db *sql.DB, ctx context.Context) *cli.Command {
	return &cli.Command{
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
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
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
			err := ordersHelper(w, db, statement, args, ctx)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
func runMigrations(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "run-migrations",
		Usage: "runs migrations for all orgs",
		Action: func(cCtx *cli.Context) error {
			path := "migrations"
			runner := NewMigrationRunner(path)

			state, err := migration.LoadMigrationState(ctx, migration.DefaultMigrationStatePath)
			if err != nil {
				return err
			}

			state, err = runner.RunAll(ctx, state)
			if err != nil {
				return err
			}

			return  migration.SaveMigrationState(ctx, state, migration.DefaultMigrationStatePath)
		},
	}

}

const lastRanMigrationID = 1

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var db *sql.DB

	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	// TODO(zpatrick): get all orgs
	// TODO(zpatrick): migrate all orgs on run-migrations call
	// orgs, err := getAllOrgs(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var org string = "default"

	app := &cli.App{
		Name: "store",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "org",
				Value:       "default",
				Usage:       "org to connect to",
				Destination: &org,
			},
		},
		Before: func(cCtx *cli.Context) error {
			var err error
			db, err = connectDB(org)

			return err
		},
		Commands: []*cli.Command{
			runMigrations(ctx),
			newCreateCustomerCommand(db),
			newCreateProductCommand(db),
			newCreateOrderCommand(db),
			newShowCustomerCommand(db, ctx),
			newShowProductCommand(db, ctx),
			newShowOrderCommand(db, ctx),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func orderPrintHelper(orders []*Order, w *tabwriter.Writer) error {
	fmt.Fprintf(w, "%-*s\t%-*s\t%-*s\t\n", 10, "OrderID", 12, "ProductID", 13, "CustomerID")
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
	fmt.Fprintf(w, "%-*s\t%-*s\t%-*s\t%-*s\t\n", 3, "ID", 25, "Name", 13, "Price", 25, "Sku")

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
	fmt.Fprintf(w, "%-*s\t%-*s\t%-*s\t\n", 3, "ID", 50, "Email", 2, "State")
	for _, c := range customers {
		c.Print(w)
	}
	return nil
}

func productsInsertHelper(db *sql.DB, p Product, statement string, w *tabwriter.Writer) error {
	var res sql.Result
	var err error
	if p.sku.String != "" {
		res, err = db.Exec(statement, p.name, p.price, p.sku)
	} else {
		res, err = db.Exec(statement, p.name, p.price)
	}
	if err != nil {
		return err
	}

	productId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	printProduct(os.Stdout, &Product{int(productId), p.name, p.price, p.sku})

	return nil

}
