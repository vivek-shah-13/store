package main

import (
	"os"
	"testing"
	"text/tabwriter"

	"github.com/urfave/cli/v2"
	"gotest.tools/v3/assert"
)

func TestExample(t *testing.T) {
	t.Log("Test Example")
}

func TestExample2(t *testing.T) {
	t.Log("Test Example 2")
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
// 6. When I run it would a state arg: it shoudl return an error

func TestCreateCustomer_withAMissingEmailArg_returnsAnError(t *testing.T) {
	// TODO: populate cCtx and db

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

	// example
	assert.Equal(t, 1, 1)
}

func ExamplePrintCustomer() {
	c := Customer{
		ID:    1,
		State: "WA",
		Email: "zack.patrick@outreach.io",
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	c.Print(w)
	w.Flush()
	// Output:
	// 1| zack.patrick@outreach.io|WA
}
