/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"log"
	"os"

	"github.com/TZGyn/kode/cmd"
	"github.com/TZGyn/kode/db"
	"github.com/charmbracelet/fang"

	_ "modernc.org/sqlite"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Panic(err)
	}

	if err := fang.Execute(context.Background(), cmd.RootCmd); err != nil {
		os.Exit(1)
	}
}
