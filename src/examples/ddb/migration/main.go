package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	apexlog "github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"golang-migrator-example/src/examples/ddb"
	"golang-migrator-example/src/migrator"
)

func main() {
	var set string
	flag.StringVar(&set, "set", "", "migration set name")

	var table string
	flag.StringVar(&table, "table", "", "migration history table name")

	flag.Parse()

	if set == "" {
		log.Fatal("empty set name")
	}
	if table == "" {
		log.Fatal("empty table name")
	}

	// make sure apex logger marshals and outputs JSON
	apexlog.SetHandler(json.New(os.Stderr))

	ctx := context.Background()

	conf, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db := dynamodb.NewFromConfig(conf)

	o := migrator.DefaultDDBProviderOptions{
		MigrationSet: set,
		Conf:         conf,
		DB:           db,
	}

	var defs []migrator.Definition
	switch set {
	case "example":
		defs, err = ddb.DefsExample(ctx, o.Conf, o.DB), nil
	default:
		defs, err = nil, fmt.Errorf("unknown migration set: %s", o.MigrationSet)
	}
	if err != nil {
		log.Fatal(err)
	}

	m := migrator.New(db, table)
	summary, err := m.Run(ctx, set, defs)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(summary)
}
