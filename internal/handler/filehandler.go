package handler

import (
	"crypto/x509"
	"embed"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

//go:embed templates/*.html
var templateFS embed.FS

const maxSize = 10 << 20

type FileHandler struct {
	maxSize   int64
	templates *template.Template
	debug     bool
}

func (h *FileHandler) loadTemplates() error {
	if h.debug {
		var err error
		h.templates, err = template.ParseGlob("internal/handler/templates/*.html")
		return err
	}

	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return err
	}

	h.templates = tmpl
	return nil
}

func NewFileHandler(debug bool) (*FileHandler, error) {

	h := &FileHandler{
		// service: service,
		debug: debug,
	}

	if err := h.loadTemplates(); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *FileHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, h.maxSize)

	if err := r.ParseMultipartForm(h.maxSize); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	cert, err := parseCertificate(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid certificate: %v", err), http.StatusBadRequest)
		return
	}

	response := map[string]any{
		"subject": cert.Subject.String(),
		"issuer":  cert.Issuer.String(),
		// "serial_number":   cert.SerialNumber.String(),
		"not_before":  cert.NotBefore.Format("2006-01-02"),
		"not_after":   cert.NotAfter.Format("2006-01-02"),
		"common_name": cert.Subject.CommonName,
		// "organization":    cert.Subject.Organization,
		// "email_addresses": cert.EmailAddresses,
		// "dns_names":       cert.DNSNames,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func parseCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block != nil {
		return x509.ParseCertificate(block.Bytes)
	}

	return x509.ParseCertificate(data)
}
