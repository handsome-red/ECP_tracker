package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Response struct {
	Items   []Certificate `json:"items"`
	HasMore bool          `json:"hasMore"`
}

type Certificate struct {
	// ID                                    string  `json:"id"`
	FullName string `json:"fullName"`
	// Snils                                 string  `json:"snils"`
	// OrganizationName                      string  `json:"organizationName"`
	// SubdivisionName                       *string `json:"subdivisionName"`
	// Position                              string  `json:"position"`
	// Type                                  int     `json:"type"`
	// SerialNumber                          string  `json:"serialNumber"`
	// CertificateRequestID                  string  `json:"certificateRequestId"`
	From string `json:"from"`
	To   string `json:"to"`
	// IsCertificateChangeStatusRequestExist bool    `json:"isCertificateChangeStatusRequestExist"`
}

func main() {

	dirName := "json_ecp"

	dirEntries, err := os.ReadDir(dirName)
	if err != nil {
		fmt.Println("Ошибка чтения директории:", err)
		return
	}

	var all []Certificate

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(dirName, entry.Name())

		data, err := os.ReadFile(fullPath)
		if err != nil {
			fmt.Printf("Не удалось прочитать %s: %v\n", fullPath, err)
			continue
		}

		var resp Response
		if err := json.Unmarshal(data, &resp); err != nil {
			fmt.Printf("Ошибка парсинга %s: %v\n", fullPath, err)
			continue
		}

		all = append(all, resp.Items...)
	}

	for _, value := range all {
		fmt.Println(value.FullName)
		dateByte := (value.To)[:10]
		fmt.Println(dateByte)
	}

	fmt.Println(len(all))
}
