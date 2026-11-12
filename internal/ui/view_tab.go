package ui

import (
	"ecp_app/internal/db"
	"ecp_app/internal/model"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func certTable(w fyne.Window, loader func() ([]model.Certificate, error)) fyne.CanvasObject {
	var rows []model.Certificate

	table := widget.NewTable(
		func() (int, int) { return len(rows) + 1, 5 },
		func() fyne.CanvasObject {
			l := widget.NewLabel("")
			l.Wrapping = fyne.TextWrapWord
			return l
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			l := obj.(*widget.Label)
			if id.Row == 0 {
				headers := []string{"Фамилия", "Имя", "Отчество", "Отдел", "Действителен до"}
				l.SetText(headers[id.Col])
				l.TextStyle = fyne.TextStyle{Bold: true}
				return
			}
			l.TextStyle = fyne.TextStyle{}
			row := rows[id.Row-1]
			switch id.Col {
			case 0:
				l.SetText(row.SecondName)
			case 1:
				l.SetText(row.FirstName)
			case 2:
				l.SetText(row.Patronymic)
			case 3:
				l.SetText(row.District)
			case 4:
				l.SetText(row.Expired())
			}
		},
	)
	table.SetColumnWidth(0, 150)
	table.SetColumnWidth(1, 120)
	table.SetColumnWidth(2, 150)
	table.SetColumnWidth(3, 500)
	table.SetColumnWidth(4, 170)

	reload := func() {
		loaded, err := loader()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		rows = loaded
		table.Refresh()
	}
	reload()

	refreshBtn := widget.NewButton("Обновить", reload)
	return container.NewBorder(refreshBtn, nil, nil, nil, table)
}

func newViewTab(w fyne.Window, store *db.Database) fyne.CanvasObject {
	return certTable(w, store.LoadCerts)
}

func newExpiringTab(w fyne.Window, store *db.Database) fyne.CanvasObject {
	title := widget.NewLabel("")
	next := time.Now().AddDate(0, 1, 0)
	title.SetText(fmt.Sprintf("Сертификаты, истекающие в %s:", monthNameRU(next.Month())+" "+fmt.Sprint(next.Year())))

	table := certTable(w, store.LoadExpiringNextMonth)
	return container.NewBorder(title, nil, nil, nil, table)
}

func monthNameRU(m time.Month) string {
	months := []string{
		"январе", "феврале", "марте", "апреле", "мае", "июне",
		"июле", "августе", "сентябре", "октябре", "ноябре", "декабре",
	}
	return months[m-1]
}
