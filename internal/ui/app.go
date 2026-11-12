package ui

import (
	"ecp_app/internal/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Build(w fyne.Window, db *db.Database) {
	tabs := container.NewAppTabs(
		container.NewTabItem("Добавить", newAddTab(w, db)),
		container.NewTabItem("Все сотрудники", newViewTab(w, db)),
		container.NewTabItem("Истекают в след. месяце", newExpiringTab(w, db)),
	)
	w.SetContent(tabs)
}
