package main

import (
	"ecp_app/internal/handler"
	"ecp_app/internal/repository/mem_repository"
	"ecp_app/internal/service"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	repo, err := memrepository.NewMemRepository()
	if err != nil {
		fmt.Println("repository error")
		return
	}

	serviceEcp := service.NewEcpService(repo)
	hand := handler.NewEcpHandler(serviceEcp)

	r := chi.NewRouter()
	r.Get("/", hand.StartPage)

	r.Get("/download", hand.RegisterECP)  // Прикрепление ЭЦП и указание отдела
	r.Post("/download", hand.DownloadECP) // Загрузка ЭЦП

	r.Get("/all", hand.AllECP)      // Увидеть все ЭЦП с датами, ФИО, отделами
	r.Get("/expired", hand.Expired) // Увидеть все ЭЦП истекающие в следующем месяце
	r.Get("/expired", hand.Files)   // Скачать Exel таблицу и Word заявление

	addr := ":8181"

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	server.ListenAndServe()
}
