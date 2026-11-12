package cert

import (
	"crypto/x509"
	"ecp_app/internal/model"
	"os"
	"strings"
)

func Parse(path string) (model.Certificate, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return model.Certificate{}, err
	}

	cert, err := x509.ParseCertificate(raw)
	if err != nil {
		return model.Certificate{}, err
	}

	parts := strings.Split(cert.Subject.CommonName, " ")

	return model.Certificate{
		SecondName: parts[0],
		FirstName:  parts[1],
		Patronymic: parts[2],
		District:   "",
		Expiration: cert.NotAfter,
	}, nil
}
