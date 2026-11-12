package ui

import (
	"ecp_app/internal/cert"
	"ecp_app/internal/db"
	"ecp_app/internal/model"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func newAddTab(w fyne.Window, store *db.Database) fyne.CanvasObject {

	// --- Общие элементы ---
	districtSelect := widget.NewSelect(nil, nil)
	districtSelect.PlaceHolder = "Выберите отдел"

	saveBtn := widget.NewButton("Сохранить", nil)
	saveBtn.Disable()

	// Загружаем список отделов из БД
	if districts, err := store.LoadDistricts(); err == nil {
		districtSelect.Options = districts
	}

	// ============================================================
	// РЕЖИМ 1: Из файла
	// ============================================================
	dropHint := widget.NewLabel("Перетащите файл сертификата в это окно\nили нажмите кнопку ниже")
	dropHint.Alignment = fyne.TextAlignCenter

	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	// Поля формы
	fioLabel := widget.NewLabel("ФИО: —")
	dateLabel := widget.NewLabel("Действителен до: —")

	// Состояние формы
	var parsed model.Certificate
	var hasFile bool
	var updateSaveState func()

	// Обработка файла (drag-and-drop или кнопка)
	handleFile := func(filePath string) {
		fioLabel.SetText("ФИО: читаю...")
		dateLabel.SetText("Действителен до: читаю...")
		progress.Show()

		go func() {
			cert, err := cert.Parse(filePath)
			fyne.Do(func() {
				progress.Hide()
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				parsed = cert
				hasFile = true
				fio := cert.FullName()
				fioLabel.SetText("ФИО: " + fio)
				dateLabel.SetText("Действителен до: " + cert.Expired())
				updateSaveState()
			})
		}()
	}

	w.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		if len(uris) == 0 {
			return
		}
		handleFile(uris[0].Path())
	})

	chooseBtn := widget.NewButton("Выбрать файл...", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()
			handleFile(reader.URI().Path())
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".cer", ".crt", ".pem", ".der"}))
		fd.Show()
	})

	fileBox := container.NewVBox(
		dropHint,
		chooseBtn,
		widget.NewSeparator(),
		progress,
		fioLabel,
		dateLabel,
	)

	// ============================================================
	// РЕЖИМ 2: Вручную
	// ============================================================
	surnameEntry := widget.NewEntry()
	surnameEntry.SetPlaceHolder("Фамилия")

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Имя")

	patronymicEntry := widget.NewEntry()
	patronymicEntry.SetPlaceHolder("Отчество")

	// Дата: используем Entry с валидацией формата YYYY-MM-DD.
	// (В Fyne нет встроенного DatePicker, поэтому проще всего так.)
	dateEntry := widget.NewEntry()
	dateEntry.SetPlaceHolder("ГГГГ-ММ-ДД")

	manualBox := container.NewVBox(
		widget.NewLabel("Введите данные сотрудника:"),
		surnameEntry,
		nameEntry,
		patronymicEntry,
		dateEntry,
	)

	// ============================================================
	// Переключатель режимов
	// ============================================================
	fileModeActive := true // по умолчанию — режим из файла

	// Пересборка состояния кнопки «Сохранить»:
	// в файловом режиме нужен файл + отдел;
	// в ручном — все текстовые поля непустые + валидная дата + отдел.
	updateSaveState = func() {
		if districtSelect.Selected == "" {
			saveBtn.Disable()
			return
		}
		if fileModeActive {
			if hasFile {
				saveBtn.Enable()
			} else {
				saveBtn.Disable()
			}
			return
		}
		// Ручной режим
		if surnameEntry.Text == "" || nameEntry.Text == "" || dateEntry.Text == "" {
			saveBtn.Disable()
			return
		}
		if _, err := time.Parse("2006-01-02", dateEntry.Text); err != nil {
			saveBtn.Disable()
			return
		}
		saveBtn.Enable()
	}

	districtSelect.OnChanged = func(string) { updateSaveState() }

	// Переключение видимых элементов
	modeSelect := widget.NewRadioGroup([]string{"Из файла", "Вручную"}, func(selected string) {
		fileModeActive = selected == "Из файла"
		if fileModeActive {
			fileBox.Show()
			manualBox.Hide()
		} else {
			fileBox.Hide()
			manualBox.Show()
		}
		updateSaveState()
	})
	modeSelect.Horizontal = true
	modeSelect.SetSelected("Из файла")

	// ============================================================
	// Сохранение (общая логика для обоих режимов)
	// ============================================================
	saveBtn.OnTapped = func() {
		var cer model.Certificate

		if fileModeActive {
			cer = model.Certificate{
				SecondName: parsed.SecondName,
				FirstName:  parsed.FirstName,
				Patronymic: parsed.Patronymic,
				District:   districtSelect.Selected,
				Expiration: parsed.Expiration,
			}
		} else {
			exp, err := time.Parse("2006-01-02", dateEntry.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("неверный формат даты: %w", err), w)
				return
			}
			cer = model.Certificate{
				SecondName: surnameEntry.Text,
				FirstName:  nameEntry.Text,
				Patronymic: patronymicEntry.Text,
				District:   districtSelect.Selected,
				Expiration: exp,
			}
		}

		if err := store.SaveCert(cer); err != nil {
			dialog.ShowError(err, w)
			return
		}
		dialog.ShowInformation("Готово", "Сертификат сохранён", w)

		// Сброс формы
		hasFile = false
		parsed = model.Certificate{}
		surnameEntry.SetText("")
		nameEntry.SetText("")
		patronymicEntry.SetText("")
		dateEntry.SetText("")
		districtSelect.ClearSelected()
		fioLabel.SetText("ФИО: —")
		dateLabel.SetText("Действителен до: —")
		updateSaveState()
	}

	// Сборка вкладки
	// ============================================================
	return container.NewVBox(
		widget.NewSeparator(),
		modeSelect,
		widget.NewSeparator(),
		fileBox,
		manualBox,
		widget.NewSeparator(),
		districtSelect,
		saveBtn,
	)
}
