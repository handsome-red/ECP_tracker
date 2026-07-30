package main

import (
	"database/sql"
	"ecp_app/internal/handler"
	// "ecp_app/internal/repository/sqlite"
	// "ecp_app/internal/service"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := "./app.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err) // Показывает причину
	}
	defer db.Close()

	// repo, err := sqlite.NewRepository(db)
	if err != nil {
		fmt.Println("repository error")
		return
	}

	// serviceEcp := service.NewEcpService(repo)
	// hand := handler.NewEcpHandler(serviceEcp)
	hand, _ := handler.NewFileHandler(false)

	r := chi.NewRouter()
	// r.Get("/", hand.StartPage)

	// r.Get("/download", hand.RegisterECP)  // Прикрепление ЭЦП и указание отдела
	r.Post("/download", hand.HandleUpload) // Загрузка ЭЦП

	// r.Get("/all", hand.AllECP)      // Увидеть все ЭЦП с датами, ФИО, отделами
	// r.Get("/expired", hand.Expired) // Увидеть все ЭЦП истекающие в следующем месяце
	// r.Get("/expired", hand.Files)   // Скачать Exel таблицу и Word заявление

	addr := ":8181"

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	fmt.Println("test")
	server.ListenAndServe()
}
