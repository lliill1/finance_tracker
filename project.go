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
	myApp.Settings().SetTheme(theme.DarkTheme())
	myWindow := myApp.NewWindow("Finance Tracker")
	myWindow.Resize(fyne.NewSize(800, 600))

	// Заголовок
	title := widget.NewLabel("Учет доходов и расходов")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

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

	// Вертикальное расположение кнопок в один столбец
	buttons := container.NewVBox(
		addButtonAligned,
		viewButtonAligned,
		reportButtonAligned,
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
