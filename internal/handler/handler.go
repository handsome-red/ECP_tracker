package handler

import (
	"ecp_app/internal/service"
	"net/http"
)

type EcpHandler struct {
	serviceEcp service.EcpServiceInterface
}

func NewEcpHandler(serviceEcp service.EcpServiceInterface) *EcpHandler {
	return &EcpHandler{serviceEcp: serviceEcp}
}

func (h *EcpHandler) StartPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world\nAndrey\nLOL\n"))
}

func (h *EcpHandler) RegisterECP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Enter the department\n"))
} // Прикрепление ЭЦП и указание отдела

func (h *EcpHandler) DownloadECP(w http.ResponseWriter, r *http.Request) {

} // Загрузка ЭЦП

func (h *EcpHandler) AllECP(w http.ResponseWriter, r *http.Request) {

} // Увидеть все ЭЦП с датами, ФИО, отделами

func (h *EcpHandler) Expired(w http.ResponseWriter, r *http.Request) {

} // Увидеть все ЭЦП истекающие в следующем месяце

func (h *EcpHandler) Files(w http.ResponseWriter, r *http.Request) {

} // Скачать Exel таблицу и Word заявление
