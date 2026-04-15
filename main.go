/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"charm.land/fantasy"
	"github.com/TZGyn/kode/cmd"
	"github.com/TZGyn/kode/db"
	"github.com/TZGyn/kode/internal/agent/tools"
	"github.com/charmbracelet/fang"

	_ "modernc.org/sqlite"
)

func main() {

	tool := tools.NewViewTool()

	response, _ := tool.Run(context.Background(), fantasy.ToolCall{
		ID:    "test",
		Name:  "test",
		Input: "{\"file_path\":\"./internal/agent/tools/view.md\",\"offset\":0,\"limit\":20}"},
	)

	fmt.Println(response)

	return
	if err := db.InitDB(); err != nil {
		log.Panic(err)
	}

	if err := fang.Execute(context.Background(), cmd.RootCmd); err != nil {
		os.Exit(1)
	}
}
