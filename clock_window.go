package main

import (
	"fmt"
	"image/color"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/container"

)

type ClockWindow struct {
	window fyne.Window
	timerText *canvas.Text
	subjectText *canvas.Text
	manuallyClosed bool
}

func NewClockWindow(app fyne.App) *ClockWindow {
	c := &ClockWindow{app.NewWindow("Clock"), nil, nil, false}

	c.Init()

	return c
}


func (c *ClockWindow) Show() {
	c.manuallyClosed = false
	c.window.Show()
}

func (c *ClockWindow) Init() {
	c.window.SetCloseIntercept(func() {
		c.manuallyClosed = true 
		c.window.Hide()
	})
	
	c.timerText = canvas.NewText(fmtDuration(0), color.Black)
	c.timerText.Alignment = fyne.TextAlignCenter
	c.timerText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	c.timerText.TextSize = 20
	
	c.subjectText = canvas.NewText("", color.Black)
	c.subjectText.Alignment = fyne.TextAlignCenter
	c.subjectText.TextStyle = fyne.TextStyle{Monospace: true}
	
	content := container.New(layout.NewVBoxLayout(), c.timerText, widget.NewSeparator(), c.subjectText)
	c.window.SetContent(content)
}

func (c *ClockWindow) Update(s *SubjectRunner) {
	duration := s.Subject.Duration - s.RunTime
	c.timerText.Text = fmtDuration(duration)

	if duration > 0 {
		c.timerText.Color = color.Black
	} else {
		c.timerText.Color = color.RGBA{255, 0, 0, 255}
	}
	c.subjectText.Text = fmt.Sprintf("Current: %s", s.Subject.Title)
	c.timerText.Refresh()
	c.subjectText.Refresh()
	if (!c.manuallyClosed){
		c.window.Show()}
}