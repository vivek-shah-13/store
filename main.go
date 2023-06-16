package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v2"
	"log"
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
