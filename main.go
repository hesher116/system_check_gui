package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func getSystemStatus() string {
	v, err := mem.VirtualMemory()
	if err != nil {
		return "Помилка читання памʼяті"
	}

	uptimeSec, err := host.Uptime()
	if err != nil {
		return "Помилка читання аптайму"
	}
	uptime := time.Duration(uptimeSec) * time.Second

	cpuPercents, err := cpu.Percent(0, false)
	if err != nil || len(cpuPercents) == 0 {
		return "Помилка CPU"
	}
	cpuLoad := cpuPercents[0]

	status := fmt.Sprintf(
		"Аптайм: %s\n💾 RAM: %.2f GB із %.2f GB\n🔥 CPU: %.1f%%\n",
		uptime.Round(time.Second),
		float64(v.Used)/1024/1024/1024,
		float64(v.Total)/1024/1024/1024,
		cpuLoad,
	)

	if v.Available < 512*1024*1024 {
		status += "Мало доступної памʼяті."
	} else if v.Available < 2*1024*1024*1024 {
		status += "Памʼяті небагато."
	} else if cpuLoad > 80 {
		status += "Високе навантаження CPU!"
	} else {
		status += "Все добре."
	}

	return status
}

func main() {
	a := app.New()
	w := a.NewWindow("System Check")

	label := widget.NewLabel(getSystemStatus())

	refreshBtn := widget.NewButton("🔁 Оновити", func() {
		label.SetText(getSystemStatus())
	})

	w.SetContent(container.NewVBox(
		label,
		refreshBtn,
	))

	w.Resize(fyne.NewSize(400, 200))
	w.ShowAndRun()
}
