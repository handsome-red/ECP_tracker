package main

import (
	"ecp_app/internal/db"
	"ecp_app/internal/ui"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	database, err := db.Init()
	if err != nil {
		fmt.Printf("Не удалось инициализировать базу данных: %v\n", err)
		os.Exit(1)
	}
	defer database.DB.Close()

	a := app.New()
	w := a.NewWindow("ЭЦП архив")
	ui.Build(w, database)
	w.Resize(fyne.NewSize(800, 500))
	w.ShowAndRun()
}
