package main

import (
	"database/sql"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

type Transaction struct {
	ID          int
	Date        string
	Category    string
	Amount      float64
	Description string
	Type        string
}

func main() {
	// Инициализация базы данных
	db, err := sql.Open("sqlite3", "./finance.db")
	if err != nil {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Ошибка",
			Content: "Не удалось подключиться к базе данных",
		})
		return
	}
	defer db.Close()
	createTable(db)

	// Инициализация приложения
	myApp := app.New()
	customTheme := newCustomTheme(true) // Начинаем с темной темы
	myApp.Settings().SetTheme(customTheme)
	myWindow := myApp.NewWindow("Finance Tracker")
	myWindow.Resize(fyne.NewSize(800, 600))

	// Заголовок
	title := widget.NewLabel("Учет доходов и расходов")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	// Переключатель темы
	themeSwitch := widget.NewCheck("Темная тема", func(checked bool) {
		customTheme := newCustomTheme(checked)
		myApp.Settings().SetTheme(customTheme)
	})
	themeSwitch.Checked = true // Начинаем с темной темы

	// Кнопки для главного окна
	addButton := widget.NewButtonWithIcon("Добавить транзакцию", theme.ContentAddIcon(), func() {
		addTransactionWindow(myApp, db).Show()
	})
	addButtonContainer := container.NewMax(addButton)
	addButtonContainer.Resize(fyne.NewSize(200, 60))
	addButtonAligned := container.NewHBox(addButtonContainer, widget.NewLabel("")) // Выравнивание влево

	viewButton := widget.NewButtonWithIcon("Просмотреть транзакции", theme.ViewFullScreenIcon(), func() {
		viewTransactionsWindow(myApp, db).Show()
	})
	viewButtonContainer := container.NewMax(viewButton)
	viewButtonContainer.Resize(fyne.NewSize(200, 60))
	viewButtonAligned := container.NewHBox(viewButtonContainer, widget.NewLabel(""))

	reportButton := widget.NewButtonWithIcon("Сгенерировать отчет", theme.DocumentIcon(), func() {
		reportWindow(myApp, db).Show()
	})
	reportButtonContainer := container.NewMax(reportButton)
	reportButtonContainer.Resize(fyne.NewSize(200, 60))
	reportButtonAligned := container.NewHBox(reportButtonContainer, widget.NewLabel(""))

	fullScreenButton := widget.NewButtonWithIcon("Полноэкранный режим", theme.ViewFullScreenIcon(), func() {
		myWindow.SetFullScreen(!myWindow.FullScreen())
	})
	fullScreenButtonContainer := container.NewMax(fullScreenButton)
	fullScreenButtonContainer.Resize(fyne.NewSize(200, 60))
	fullScreenButtonAligned := container.NewHBox(fullScreenButtonContainer, widget.NewLabel(""))

	exitButton := widget.NewButtonWithIcon("Выход", theme.LogoutIcon(), func() {
		myApp.Quit()
	})
	exitButtonContainer := container.NewMax(exitButton)
	exitButtonContainer.Resize(fyne.NewSize(200, 60))
	exitButtonAligned := container.NewHBox(exitButtonContainer, widget.NewLabel(""))

	// Добавляем новые кнопки
	statisticsButton := widget.NewButtonWithIcon("Статистика", theme.DocumentIcon(), func() {
		statisticsWindow(myApp, db).Show()
	})
	statisticsButtonContainer := container.NewMax(statisticsButton)
	statisticsButtonContainer.Resize(fyne.NewSize(200, 60))
	statisticsButtonAligned := container.NewHBox(statisticsButtonContainer, widget.NewLabel(""))

	budgetButton := widget.NewButtonWithIcon("Бюджет", theme.SettingsIcon(), func() {
		budgetWindow(myApp, db).Show()
	})
	budgetButtonContainer := container.NewMax(budgetButton)
	budgetButtonContainer.Resize(fyne.NewSize(200, 60))
	budgetButtonAligned := container.NewHBox(budgetButtonContainer, widget.NewLabel(""))

	exportButton := widget.NewButtonWithIcon("Экспорт данных", theme.DocumentSaveIcon(), func() {
		exportDataWindow(myApp, db).Show()
	})
	exportButtonContainer := container.NewMax(exportButton)
	exportButtonContainer.Resize(fyne.NewSize(200, 60))
	exportButtonAligned := container.NewHBox(exportButtonContainer, widget.NewLabel(""))

	// Обновляем вертикальное расположение кнопок
	buttons := container.NewVBox(
		container.NewHBox(widget.NewLabel(""), themeSwitch), // Добавляем переключатель темы
		addButtonAligned,
		viewButtonAligned,
		reportButtonAligned,
		statisticsButtonAligned,
		budgetButtonAligned,
		exportButtonAligned,
		fullScreenButtonAligned,
		exitButtonAligned,
	)

	// Загрузка изображения из той же директории
	imageData, err := os.ReadFile("jpeg.jpg")
	var imageContainer fyne.CanvasObject
	if err != nil {
		// Если изображение не удалось загрузить, используем заглушку
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Ошибка",
			Content: "Не удалось загрузить изображение: " + err.Error(),
		})
		placeholderImage := canvas.NewRectangle(&color.NRGBA{R: 150, G: 150, B: 150, A: 255})
		placeholderImage.CornerRadius = 20
		placeholderImage.SetMinSize(fyne.NewSize(300, 300))
		imageLabel := widget.NewLabel("Не удалось загрузить изображение")
		imageContainer = container.NewCenter(
			container.NewVBox(
				placeholderImage,
				imageLabel,
			),
		)
	} else {
		// Создание ресурса для изображения
		imageResource := fyne.NewStaticResource("jpeg.jpg", imageData)
		img := canvas.NewImageFromResource(imageResource)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(300, 300))
		// Примечание: canvas.Image не поддерживает CornerRadius
		// Для закругленных углов нужно предварительно обработать изображение
		imageContainer = container.NewCenter(img)
	}

	// Разделение окна: кнопки слева, изображение справа
	split := container.NewHSplit(buttons, imageContainer)
	split.SetOffset(0.5) // Делим окно пополам

	// Основной контейнер с заголовком и разделением
	content := container.NewVBox(
		container.NewCenter(title),
		split,
	)

	// Увеличиваем отступы по краям с помощью canvas.Rectangle в качестве спейсера
	spacer := canvas.NewRectangle(&color.NRGBA{R: 0, G: 0, B: 0, A: 0}) // Прозрачный спейсер
	spacer.SetMinSize(fyne.NewSize(0, 20))                              // 20 пикселей отступа
	customPaddedContent := container.NewBorder(
		spacer, spacer, spacer, spacer, // Отступы: сверху, снизу, слева, справа
		content,
	)

	myWindow.SetContent(customPaddedContent)
	myWindow.ShowAndRun()
}

func createTable(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT,
			category TEXT,
			amount REAL,
			description TEXT,
			type TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS budget_limits (
			category TEXT PRIMARY KEY,
			limit_amount REAL
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Ошибка",
				Content: "Не удалось создать таблицу: " + err.Error(),
			})
		}
	}
}

func addTransactionWindow(a fyne.App, db *sql.DB) fyne.Window {
	window := a.NewWindow("Добавить транзакцию")
	window.Resize(fyne.NewSize(600, 400))

	typeSelect := widget.NewSelect([]string{"Доход", "Расход"}, nil)
	categoryEntry := widget.NewEntry()
	categoryEntry.SetPlaceHolder("Категория")
	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("Сумма")
	descriptionEntry := widget.NewEntry()
	descriptionEntry.SetPlaceHolder("Описание")
	dateEntry := widget.NewEntry()
	dateEntry.SetPlaceHolder("Дата (YYYY-MM-DD)")

	saveButton := widget.NewButton("Сохранить", func() {
		amount, err := strconv.ParseFloat(amountEntry.Text, 64)
		if err != nil || amount <= 0 {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Ошибка",
				Content: "Неверная сумма",
			})
			return
		}
		if dateEntry.Text == "" {
			dateEntry.Text = time.Now().Format("2006-01-02")
		}
		query := `INSERT INTO transactions (date, category, amount, description, type) VALUES (?, ?, ?, ?, ?)`
		_, err = db.Exec(query, dateEntry.Text, categoryEntry.Text, amount, descriptionEntry.Text, typeSelect.Selected)
		if err != nil {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Ошибка",
				Content: "Не удалось сохранить транзакцию",
			})
			return
		}
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Успех",
			Content: "Транзакция сохранена",
		})
		window.Close()
	})

	content := container.NewVBox(
		typeSelect,
		categoryEntry,
		amountEntry,
		descriptionEntry,
		dateEntry,
		saveButton,
	)
	window.SetContent(content)
	return window
}

func viewTransactionsWindow(a fyne.App, db *sql.DB) fyne.Window {
	window := a.NewWindow("Просмотр транзакций")
	window.Resize(fyne.NewSize(1000, 600))

	rows, err := db.Query("SELECT id, date, type, category, amount, description FROM transactions")
	if err != nil {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Ошибка",
			Content: "Не удалось загрузить транзакции",
		})
		return window
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.Date, &t.Type, &t.Category, &t.Amount, &t.Description); err != nil {
			continue
		}
		transactions = append(transactions, t)
	}

	list := widget.NewList(
		func() int { return len(transactions) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			t := transactions[i]
			o.(*widget.Label).SetText(fmt.Sprintf("%s | %s | %s | %.2f | %s", t.Date, t.Type, t.Category, t.Amount, t.Description))
		},
	)

	scroll := container.NewScroll(list)
	window.SetContent(scroll)
	return window
}

func reportWindow(a fyne.App, db *sql.DB) fyne.Window {
	window := a.NewWindow("Сгенерировать отчет")
	window.Resize(fyne.NewSize(600, 400))

	periodEntry := widget.NewEntry()
	periodEntry.SetPlaceHolder("Период (YYYY-MM)")
	generateButton := widget.NewButton("Сгенерировать", func() {
		period := periodEntry.Text
		query := `
		SELECT type, SUM(amount) 
		FROM transactions 
		WHERE strftime('%Y-%m', date) = ? 
		GROUP BY type`
		rows, err := db.Query(query, period)
		if err != nil {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Ошибка",
				Content: "Не удалось сгенерировать отчет",
			})
			return
		}
		defer rows.Close()

		var totalIncome, totalExpense float64
		for rows.Next() {
			var tType string
			var amount float64
			if err := rows.Scan(&tType, &amount); err != nil {
				continue
			}
			if tType == "Доход" {
				totalIncome = amount
			} else if tType == "Расход" {
				totalExpense = amount
			}
		}

		balance := totalIncome - totalExpense
		reportText := fmt.Sprintf("Отчет за %s:\nДоход: %.2f\nРасходы: %.2f\nБаланс: %.2f", period, totalIncome, totalExpense, balance)
		reportWindow := a.NewWindow("Отчет за " + period)
		reportWindow.Resize(fyne.NewSize(600, 400))
		reportLabel := widget.NewLabel(reportText)
		scroll := container.NewScroll(reportLabel)
		reportWindow.SetContent(scroll)
		reportWindow.Show()
	})

	content := container.NewVBox(periodEntry, generateButton)
	window.SetContent(content)
	return window
}

func statisticsWindow(a fyne.App, db *sql.DB) fyne.Window {
	window := a.NewWindow("Статистика")
	window.Resize(fyne.NewSize(1000, 800))

	// Элементы управления
	periodSelect := widget.NewSelect([]string{"Все время", "По годам", "По месяцам", "Выбрать период"}, nil)
	periodSelect.SetSelected("Все время")
	
	yearSelect := widget.NewSelect([]string{}, nil)
	monthSelect := widget.NewSelect([]string{"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", 
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"}, nil)
	
	startDateEntry := widget.NewEntry()
	startDateEntry.SetPlaceHolder("Начальная дата (YYYY-MM-DD)")
	endDateEntry := widget.NewEntry()
	endDateEntry.SetPlaceHolder("Конечная дата (YYYY-MM-DD)")
	
	refreshButton := widget.NewButton("Обновить", nil)
	
	// Контейнер для динамического изменения элементов управления
	filterContainer := container.NewVBox()
	
	// Обработчик изменения периода
	periodSelect.OnChanged = func(selected string) {
		filterContainer.Objects = nil
		
		switch selected {
		case "По годам":
			// Заполняем годами из базы данных
			rows, err := db.Query("SELECT DISTINCT strftime('%Y', date) FROM transactions ORDER BY date DESC")
			if err == nil {
				defer rows.Close()
				var years []string
				for rows.Next() {
					var year string
					if err := rows.Scan(&year); err == nil {
						years = append(years, year)
					}
				}
				yearSelect.Options = years
				if len(years) > 0 {
					yearSelect.SetSelected(years[0])
				}
				filterContainer.Add(yearSelect)
			}
			
		case "По месяцам":
			// Получаем текущий год
			currentYear := time.Now().Format("2006")
			yearSelect.SetSelected(currentYear)
			filterContainer.Add(container.NewHBox(
				widget.NewLabel("Год:"),
				yearSelect,
				widget.NewLabel("Месяц:"),
				monthSelect,
			))
			
		case "Выбрать период":
			filterContainer.Add(container.NewVBox(
				startDateEntry,
				endDateEntry,
			))
		}
		
		filterContainer.Refresh()
	}
	
	// Основной контейнер статистики
	statsContainer := container.NewVBox()
	
	// Функция обновления статистики
	updateStats := func() {
		var query string
		var args []interface{}
		
		switch periodSelect.Selected {
		case "Все время":
			query = `
				SELECT 
					type, 
					category, 
					SUM(amount) as total
				FROM transactions
				GROUP BY type, category
				ORDER BY type, total DESC
			`
			
		case "По годам":
			if yearSelect.Selected == "" {
				return
			}
			query = `
				SELECT 
					type, 
					category, 
					SUM(amount) as total
				FROM transactions
				WHERE strftime('%Y', date) = ?
				GROUP BY type, category
				ORDER BY type, total DESC
			`
			args = append(args, yearSelect.Selected)
			
		case "По месяцам":
			if yearSelect.Selected == "" || monthSelect.Selected == "" {
				return
			}
			monthNum := fmt.Sprintf("%02d", monthSelect.SelectedIndex()+1)
			query = `
				SELECT 
					type, 
					category, 
					SUM(amount) as total
				FROM transactions
				WHERE strftime('%Y', date) = ? 
				AND strftime('%m', date) = ?
				GROUP BY type, category
				ORDER BY type, total DESC
			`
			args = append(args, yearSelect.Selected, monthNum)
			
		case "Выбрать период":
			if startDateEntry.Text == "" || endDateEntry.Text == "" {
				return
			}
			query = `
				SELECT 
					type, 
					category, 
					SUM(amount) as total
				FROM transactions
				WHERE date BETWEEN ? AND ?
				GROUP BY type, category
				ORDER BY type, total DESC
			`
			args = append(args, startDateEntry.Text, endDateEntry.Text)
		}
		
		// Выполняем запрос
		rows, err := db.Query(query, args...)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		defer rows.Close()
		
		var stats []struct {
			Type     string
			Category string
			Total    float64
		}
		
		var totalIncome, totalExpense float64
		
		for rows.Next() {
			var stat struct {
				Type     string
				Category string
				Total    float64
			}
			if err := rows.Scan(&stat.Type, &stat.Category, &stat.Total); err != nil {
				continue
			}
			stats = append(stats, stat)
			
			if stat.Type == "Доход" {
				totalIncome += stat.Total
			} else {
				totalExpense += stat.Total
			}
		}
		
		// Очищаем контейнер
		statsContainer.Objects = nil
		
		// Добавляем общую информацию
		statsContainer.Add(widget.NewLabel(fmt.Sprintf("Общий доход: %.2f ₽", totalIncome)))
		statsContainer.Add(widget.NewLabel(fmt.Sprintf("Общий расход: %.2f ₽", totalExpense)))
		statsContainer.Add(widget.NewLabel(fmt.Sprintf("Баланс: %.2f ₽", totalIncome-totalExpense)))
		statsContainer.Add(widget.NewSeparator())
		
		// Создаем таблицу статистики
		table := widget.NewTable(
			func() (int, int) { return len(stats) + 1, 3 },
			func() fyne.CanvasObject {
				return widget.NewLabel("Шаблон")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				label := o.(*widget.Label)
				if i.Row == 0 {
					// Заголовки
					switch i.Col {
					case 0:
						label.SetText("Тип")
					case 1:
						label.SetText("Категория")
					case 2:
						label.SetText("Сумма")
					}
				} else {
					// Данные
					stat := stats[i.Row-1]
					switch i.Col {
					case 0:
						label.SetText(stat.Type)
					case 1:
						label.SetText(stat.Category)
					case 2:
						label.SetText(fmt.Sprintf("%.2f ₽", stat.Total))
					}
				}
			},
		)
		table.SetColumnWidth(0, 100)
		table.SetColumnWidth(1, 200)
		table.SetColumnWidth(2, 150)
		
		statsContainer.Add(table)
		statsContainer.Refresh()
	}
	
	refreshButton.OnTapped = updateStats
	
	// Первоначальное обновление
	periodSelect.OnChanged(periodSelect.Selected)
	
	// Основной макет
	content := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Период:"),
			periodSelect,
			refreshButton,
		),
		filterContainer,
		widget.NewSeparator(),
		statsContainer,
	)
	
	window.SetContent(content)
	
	// Автоматическое обновление при изменении фильтров
	yearSelect.OnChanged = func(string) { updateStats() }
	monthSelect.OnChanged = func(string) { updateStats() }
	startDateEntry.OnChanged = func(string) { updateStats() }
	endDateEntry.OnChanged = func(string) { updateStats() }
	
	return window
}

func budgetWindow(a fyne.App, db *sql.DB) fyne.Window {
	window := a.NewWindow("Управление бюджетом")
	window.Resize(fyne.NewSize(600, 400))

	categoryEntry := widget.NewEntry()
	categoryEntry.SetPlaceHolder("Категория")
	limitEntry := widget.NewEntry()
	limitEntry.SetPlaceHolder("Лимит бюджета")

	// Создаем таблицу для существующих лимитов
	createBudgetTable := func() *widget.Table {
		rows, err := db.Query("SELECT category, limit_amount FROM budget_limits")
		if err != nil {
			return nil
		}
		defer rows.Close()

		var limits []struct {
			Category string
			Limit    float64
		}

		for rows.Next() {
			var limit struct {
				Category string
				Limit    float64
			}
			if err := rows.Scan(&limit.Category, &limit.Limit); err != nil {
				continue
			}
			limits = append(limits, limit)
		}

		table := widget.NewTable(
			func() (int, int) { return len(limits) + 1, 2 },
			func() fyne.CanvasObject {
				return widget.NewLabel("Шаблон")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				label := o.(*widget.Label)
				if i.Row == 0 {
					switch i.Col {
					case 0:
						label.SetText("Категория")
					case 1:
						label.SetText("Лимит")
					}
				} else {
					limit := limits[i.Row-1]
					switch i.Col {
					case 0:
						label.SetText(limit.Category)
					case 1:
						label.SetText(fmt.Sprintf("%.2f ₽", limit.Limit))
					}
				}
			},
		)
		return table
	}

	// Создаем таблицу бюджетов
	createTable(db)
	table := createBudgetTable()

	var saveButton *widget.Button
	saveButton = widget.NewButton("Сохранить лимит", func() {
		limit, err := strconv.ParseFloat(limitEntry.Text, 64)
		if err != nil || limit <= 0 {
			dialog.ShowError(fmt.Errorf("неверная сумма"), window)
			return
		}

		_, err = db.Exec(`
			INSERT OR REPLACE INTO budget_limits (category, limit_amount)
			VALUES (?, ?)
		`, categoryEntry.Text, limit)

		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		// Обновляем таблицу
		window.SetContent(container.NewVBox(
			widget.NewLabel("Управление бюджетом"),
			categoryEntry,
			limitEntry,
			saveButton,
			createBudgetTable(),
		))
	})

	content := container.NewVBox(
		widget.NewLabel("Управление бюджетом"),
		categoryEntry,
		limitEntry,
		saveButton,
		table,
	)

	window.SetContent(content)
	return window
}

func exportDataWindow(a fyne.App, db *sql.DB) fyne.Window {
	window := a.NewWindow("Экспорт данных")
	window.Resize(fyne.NewSize(400, 200))

	formatSelect := widget.NewSelect([]string{"CSV", "JSON"}, nil)
	formatSelect.SetSelected("CSV")

	exportButton := widget.NewButton("Экспортировать", func() {
		// Здесь будет логика экспорта данных
		dialog.ShowInformation("Экспорт", "Функция экспорта будет реализована в следующем обновлении", window)
	})

	content := container.NewVBox(
		widget.NewLabel("Выберите формат экспорта:"),
		formatSelect,
		exportButton,
	)

	window.SetContent(content)
	return window
}
