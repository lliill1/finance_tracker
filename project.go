package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Инициализация приложения
	myApp := app.New()
	myWindow := myApp.NewWindow("Учет доходов и расходов")

	// Кнопки основного окна
	addButton := widget.NewButton("Добавить транзакцию", func() {
		addTransactionWindow(myApp).Show()
	})
	viewButton := widget.NewButton("Просмотреть транзакции", func() {
		viewTransactionsWindow(myApp).Show()
	})
	reportButton := widget.NewButton("Генерировать отчет", func() {
		reportWindow(myApp).Show()
	})

	// Компоновка основного окна
	content := container.NewVBox(addButton, viewButton, reportButton)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

// Окно для добавления транзакции
func addTransactionWindow(a fyne.App) fyne.Window {
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
		// Здесь будет логика сохранения в базу данных
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

// Окно для просмотра транзакций
func viewTransactionsWindow(a fyne.App) fyne.Window {
	window := a.NewWindow("Просмотреть транзакции")

	// Пример данных
	transactions := []string{
		"2023-10-01 - Доход - Зарплата - 50000 - Зарплата за месяц",
		"2023-10-02 - Расход - Продукты - 2000 - Покупка продуктов",
	}

	list := widget.NewList(
		func() int { return len(transactions) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(transactions[i])
		},
	)

	window.SetContent(list)
	return window
}

// Окно для генерации отчета
func reportWindow(a fyne.App) fyne.Window {
	window := a.NewWindow("Генерировать отчет")

	periodEntry := widget.NewEntry()
	periodEntry.SetPlaceHolder("Период (YYYY-MM)")
	generateButton := widget.NewButton("Сгенерировать", func() {
		period := periodEntry.Text
		reportWindow := a.NewWindow("Отчет за " + period)
		reportLabel := widget.NewLabel("Пример отчета:\nДоходы: 50000\nРасходы: 2000\nБаланс: 48000")
		reportWindow.SetContent(reportLabel)
		reportWindow.Show()
	})

	content := container.NewVBox(periodEntry, generateButton)
	window.SetContent(content)
	return window
}
