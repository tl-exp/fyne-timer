package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"gopkg.in/yaml.v3"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/systray"
)


type Subject struct {
    Title string `yaml:"title"`
    Duration time.Duration `yaml:"duration"`
    Author string `yaml:"author"`
}

type Schedule struct {
    Subjects []Subject `yaml:"subjects"`
    Version string `yaml:"version"`
}


func ReadConfig(filename string) (*Schedule, error) {
    buf, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    schedule := &Schedule{}
    err = yaml.Unmarshal(buf, schedule)
    if err != nil {
        return nil, fmt.Errorf("in file %q: %w", filename, err)
    }
    if schedule.Version != "1"  {
    	return nil, fmt.Errorf("in file %q: Unsupported version %s", filename, schedule.Version)
    }

    return schedule, err
}


func fmtDuration(d time.Duration) string {
	sign := ""
	if d < 0 {
		sign = "-"
		d = -d
	}
    d = d.Round(time.Second)
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    d -= m * time.Minute
    s := d / time.Second
    return fmt.Sprintf("%s%02d:%02d:%02d", sign, h, m, s)
}


const AppName = "Meeting Timer"

func main() {
	schedule, err := ReadConfig("./schedule.yaml")

    if err != nil {
        log.Fatal(err)
    }

	a := app.New()
	trayUpdater := Tray{}

	runner := NewScheduleRunner(schedule)

	
	clockWindow := NewClockWindow(a)
	mainWindow := NewMainWindow(a, runner, clockWindow)

	runner.SetCallbackOnTick(func (s *SubjectRunner) {
		trayUpdater.Update(s)
		clockWindow.Update(s)
	})
	runner.SetCallbackOnSubjectFinish(func (s *SubjectRunner) {mainWindow.UpdateAndShow()})

	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu(AppName,
			fyne.NewMenuItem("Show", func () {mainWindow.UpdateAndShow()}))
		desk.SetSystemTrayMenu(m)
	}

	runner.Start()
	defer runner.Stop()

	a.Run()
}

// just to not update too often
type Tray struct {
	LastIcon string
	LastTitle string
}


func (t *Tray) Update(s *SubjectRunner) {
	duration := s.Subject.Duration - s.RunTime
	rest := fmtDuration(duration)
	var resource *fyne.StaticResource

	if duration > 0 {
		resource = resourceGreenIco
	} else {
		resource = resourceRedIco
	}

	title := fmt.Sprintf("%s: %s", s.Subject.Title, rest)
	if t.LastTitle != title {
		systray.SetTitle(title)
		systray.SetTooltip(title)
		t.LastTitle = title
	}

	if t.LastIcon != resource.StaticName {
		t.LastIcon = resource.StaticName
		systray.SetIcon(resource.StaticContent)
	}
}