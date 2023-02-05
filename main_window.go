package main

import (
	"image/color"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/container"

)

type MainWindow struct {
	window fyne.Window
	runner *ScheduleRunner
	clockWindow *ClockWindow
}


func NewMainWindow(app fyne.App, runner *ScheduleRunner,clockWindow *ClockWindow) *MainWindow {
	m := &MainWindow{app.NewWindow("Schedule"), runner, clockWindow}

	m.window.SetCloseIntercept(func() {
		m.window.Hide()
	})
	m.UpdateAndShow()
	return m
}


func (m *MainWindow) UpdateAndShow() {
	content := m.BuildMainWindow()
	m.window.SetContent(content)
	m.window.Show()
}


func (m *MainWindow) BuildMainWindow() fyne.CanvasObject {

	var StatusColors = map[SubjectStatus]color.Color{
	  Ready: color.Black,
	  Running:  color.RGBA{0, 200, 0, 255},
	  Done: color.RGBA{0, 0, 255, 255},
	  Skipped: color.RGBA{82, 82, 82, 255},
	}

	var ButtonText = map[SubjectStatus]string{
	  Ready: "Start",
	  Running:  "Continue",
	  Done: "Start Again",
	  Skipped: "",
	}


	var subjects []fyne.CanvasObject
	subjects = append(subjects,
		widget.NewLabel(""),
		widget.NewLabelWithStyle("Title", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Author", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Duration", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)
	for i := 0; i < len(m.runner.Subjects); i++ {
		i := i
		subj := m.runner.Subjects[i]
		text := canvas.NewText(subj.Subject.Title, StatusColors[subj.Status])
		text.Alignment = fyne.TextAlignCenter
		text.TextStyle = fyne.TextStyle{Monospace: true}

		subjects = append(subjects,
			widget.NewButton(ButtonText[subj.Status], func() { 
				m.runner.StartSubject(i) 
				m.window.Hide()
				m.clockWindow.Show()
			}),
			text,
			widget.NewLabelWithStyle(subj.Subject.Author, fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}),
			widget.NewLabelWithStyle(fmtDuration(subj.Subject.Duration), fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}),
		)
	}
	content := container.New(layout.NewGridLayout(4), subjects...)

	return content
}