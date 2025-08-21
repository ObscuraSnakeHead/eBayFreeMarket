package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	_ "ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/util"
	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/webapp"
)

func main() {
	app := &cli.App{
		Name:  "tochka",
		Usage: "run tochka server",
		Action: func(*cli.Context) error {
			runServer()
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"srv"},
				Usage:   "run tochka server",
				Action: func(cCtx *cli.Context) error {
					runServer()
					return nil
				},
			},
			{
				Name:    "sync",
				Aliases: []string{"snc"},
				Usage:   "sync database models and views",
				Action: func(cCtx *cli.Context) error {
					syncModels()
					syncDatabaseViews()
					return nil
				},
			},
			{
				Name:    "update-deposits",
				Aliases: []string{"ud"},
				Usage:   "update deposit accounts",
				Action: func(cCtx *cli.Context) error {
					updateDeposits()
					return nil
				},
			},
			{
				Name:    "user",
				Aliases: []string{"ug"},
				Usage:   "`use $username $action $role` -- $action:(grant, revoke) $role:(seller, staff, admin) to $user",
				Action: func(cCtx *cli.Context) error {
					username, action, role := cCtx.Args().Get(0), cCtx.Args().Get(1), cCtx.Args().Get(2)
					manageRole(username, action, role)
					return nil
				},
			},
			{
				Name:    "index",
				Aliases: []string{"i"},
				Usage:   "index items for bleve search",
				Action: func(cCtx *cli.Context) error {
					indexItems()
					return nil
				},
			},
			{
				Name:    "import-metro",
				Aliases: []string{"im"},
				Usage:   "[deprecated] import metro stations",
				Action: func(cCtx *cli.Context) error {
					importMetroStations()
					return nil
				},
			},
			{
				Name:    "staff-stats",
				Aliases: []string{"ss"},
				Usage:   "run database report query for staff statistics",
				Action: func(cCtx *cli.Context) error {
					importMetroStations()
					return nil
				},
			},
			{
				Name:    "maintain-transactions",
				Aliases: []string{"mt"},
				Usage:   "maintain stuck transactions",
				Action: func(cCtx *cli.Context) error {
					maintainTransactions()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func manageRole(username, action, role string) {
	user, _ := webapp.FindUserByUsername(username)
	if user == nil {
		fmt.Println("No such user")
		return
	}

	if action == "grant" && role == "admin" {
		user.IsAdmin = !user.IsAdmin
	} else {
		fmt.Println("Wrong action")
		return
	}
	user.Save()
}

func indexItems() {
	println("[Index] Indexing items...")
	for _, item := range webapp.GetItemsForIndexing() {
		println("[Index] ", item.Name)
		err := item.Index()
		if err != nil {
			println("Error: ", err)
		}
	}
}

func syncModels() {
	webapp.SyncModels()
}

func syncDatabaseViews() {
	webapp.SyncDatabaseViews()
}

func importMetroStations() {
	webapp.ImportCityMetroStations(524901, "./dumps/moscow-metro.json")
	webapp.ImportCityMetroStations(498817, "./dumps/spb-metro.json")
}

func updateDeposits() {
	webapp.CommandUpdateDeposits()
}

func maintainTransactions() {
	webapp.TaskFreezeStuckCompletedTransactions()
}
