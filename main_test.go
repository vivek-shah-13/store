package main

import (
	"testing"

	"github.com/urfave/cli/v2"
	"gotest.tools/v3/assert"
)

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

	db, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			newCreateCustomerCommand(db),
		},
	}

	err = app.Run([]string{"store", "create-customer", "WA"})
	assert.ErrorContains(t, err, "Must specify email and state")
}

func TestCreateCustomer_withAMissingStateArg_returnsAnError(t *testing.T) {
	db, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			newCreateCustomerCommand(db),
		},
	}
	err = app.Run([]string{"store", "create-customer", "vivek.shah@outreach.io"})
	assert.ErrorContains(t, err, "Must specify email and state")
}

func TestCreateCustomer_wthAInvalidStateAbbreviation_returnsAnError(t *testing.T) {
	db, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			newCreateCustomerCommand(db),
		},
	}
	err = app.Run([]string{"store", "create-customer", "vivek.shah@outreach.io", "PP"})
	assert.ErrorContains(t, err, "State must be a valid U.S. State or Territory")
}

func TestCreateCustomer_withANonTwoLetterStateCode_returnsAnError(t *testing.T) {
	db, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			newCreateCustomerCommand(db),
		},
	}
	err = app.Run([]string{"store", "create-customer", "vivek.shah@outreach.io", "PPP"})
	assert.ErrorContains(t, err, "State length must be 2")
}

func TestCreateCustomer_withValidInput_entersDatabaseCorrectly(t *testing.T) {
	db, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			newCreateCustomerCommand(db),
		},
	}
	_, err = db.Exec("DELETE FROM Customers")
	assert.NilError(t, err)
	_, err = db.Exec("ALTER TABLE Customers AUTO_INCREMENT=1")
	assert.NilError(t, err)
	err = app.Run([]string{"store", "create-customer", "vivek.s@outreach.io", "WA"})
	assert.NilError(t, err)

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
	db, err := connectDB()
	if err != nil {
		t.Fatal(err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			newCreateCustomerCommand(db),
		},
	}
	_, err = db.Exec("DELETE FROM Customers")
	assert.NilError(t, err)
	_, err = db.Exec("ALTER TABLE Customers AUTO_INCREMENT=1")
	assert.NilError(t, err)
	err = app.Run([]string{"store", "create-customer", "vivek.s@outreach.io", "WA"})
	assert.NilError(t, err)

	err = app.Run([]string{"store", "create-customer", "v.s@outreach.io", "CA"})
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
	db, _ := connectDB()

	app := &cli.App{
		Commands: []*cli.Command{
			newCreateCustomerCommand(db),
		},
	}
	db.Exec("DELETE FROM Customers")

	db.Exec("ALTER TABLE Customers AUTO_INCREMENT=1")

	app.Run([]string{"store", "create-customer", "vivek.s@outreach.io", "WA"})

	//Output:
	//ID  |Email                                              |State
	//1   |vivek.s@outreach.io                                |WA

}

func Example_CreateMultipleCustomer_hasCorrectPrintOutput() {
	db, _ := connectDB()

	app := &cli.App{
		Commands: []*cli.Command{
			newCreateCustomerCommand(db),
		},
	}
	db.Exec("DELETE FROM Customers")

	db.Exec("ALTER TABLE Customers AUTO_INCREMENT=1")

	app.Run([]string{"store", "create-customer", "vivek.s@outreach.io", "WA"})
	app.Run([]string{"store", "create-customer", "v.s@verylonglonglonglongemail.com", "MN"})
	//Output:
	//ID  |Email                                              |State
	//1   |vivek.s@outreach.io                                |WA
	//ID  |Email                                              |State
	//2   |v.s@verylonglonglonglongemail.com                  |MN
}

// func ExamplePrintCustomer() {
// 	c := Customer{
// 		ID:    1,
// 		State: "WA",
// 		Email: "zack.patrick@outreach.io",
// 	}

// 	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
// 	c.Print(w)
// 	w.Flush()
// 	// Output:
// 	// 1| zack.patrick@outreach.io|WA
// }
