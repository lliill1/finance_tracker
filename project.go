package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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
	myWindow := myApp.NewWindow("Finance Tracker")

	// Кнопки для главного окна
	addButton := widget.NewButton("Добавить транзакцию", func() {
		addTransactionWindow(myApp, db).Show()
	})
	viewButton := widget.NewButton("Просмотреть транзакции", func() {
		viewTransactionsWindow(myApp, db).Show()
	})
	reportButton := widget.NewButton("Сгенерировать отчет", func() {
		reportWindow(myApp, db).Show()
	})

	// Макет главного окна
	content := container.NewVBox(addButton, viewButton, reportButton)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		category TEXT,
		amount REAL,
		description TEXT,
		type TEXT
	)`
	_, err := db.Exec(query)
	if err != nil {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Ошибка",
			Content: "Не удалось создать таблицу",
		})
	}
}

func addTransactionWindow(a fyne.App, db *sql.DB) fyne.Window {
	window := a.NewWindow("Добавить транзакцию")

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

	window.SetContent(list)
	return window
}

func reportWindow(a fyne.App, db *sql.DB) fyne.Window {
	window := a.NewWindow("Сгенерировать отчет")

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
		reportLabel := widget.NewLabel(reportText)
		reportWindow.SetContent(reportLabel)
		reportWindow.Show()
	})

	content := container.NewVBox(periodEntry, generateButton)
	window.SetContent(content)
	return window
}
