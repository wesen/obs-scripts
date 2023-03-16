package main

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/andreykaipov/goobs"
)

type UI struct {
	window      fyne.Window
	statusLabel *widget.Label
	timerLabel  *widget.Label
	client      *goobs.Client
}

func main() {
	a := app.New()
	w := a.NewWindow("OBS Recorder")

	client, err := goobs.New("localhost:4455", goobs.WithPassword("qTh2WVY6SSC3JNmy"))
	if err != nil {
		w.SetContent(widget.NewLabel("Error connecting to OBS: " + err.Error()))
		w.ShowAndRun()
		return
	}

	ui := &UI{
		window:      w,
		statusLabel: widget.NewLabel(""),
		timerLabel:  widget.NewLabel(""),
		client:      client,
	}

	startButton := widget.NewButton("Start Recording", ui.startRecording)
	stopButton := widget.NewButton("Stop Recording", ui.stopRecording)

	buttons := container.NewHBox(startButton, stopButton)
	content := container.NewVBox(buttons, ui.statusLabel, ui.timerLabel)

	w.SetContent(content)
	w.Resize(fyne.NewSize(300, 200))
	w.ShowAndRun()
}

func (ui *UI) startRecording() {
	err := startRecording(ui.client)
	if err != nil {
		ui.statusLabel.SetText("Error starting recording: " + err.Error())
	} else {
		ui.statusLabel.SetText("Recording started.")
		go ui.runCountdownTimer(context.Background())
	}
}

func (ui *UI) stopRecording() {
	err := stopRecording(ui.client)
	if err != nil {
		ui.statusLabel.SetText("Error stopping recording: " + err.Error())
	} else {
		ui.statusLabel.SetText("Recording stopped.")
	}
}

func (ui *UI) runCountdownTimer(ctx context.Context) error {
	timerDuration := 5 * time.Second
	timer := time.NewTicker(1 * time.Second)
	defer timer.Stop()

	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case currentTime := <-timer.C:
			elapsed := currentTime.Sub(startTime)
			remaining := timerDuration - elapsed
			if remaining <= 0 {
				ui.timerLabel.SetText("")
				ui.showTimerExpiredDialog(ctx)
				return nil
			} else {
				ui.timerLabel.SetText(fmt.Sprintf("Time remaining: %v", remaining.Round(time.Second)))
			}
		}
	}
}

func (ui *UI) showTimerExpiredDialog(ctx context.Context) {
	var popUp dialog.Dialog
	continueButton := widget.NewButton("Continue", func() {
		ui.window.Canvas().SetOnTypedKey(nil)
		popUp.Hide()
		ui.runCountdownTimer(ctx)
	})

	stopButton := widget.NewButton("Stop", func() {
		ui.window.Canvas().SetOnTypedKey(nil)
		popUp.Hide()
		ui.stopRecording()
	})

	buttons := container.NewHBox(continueButton, stopButton)
	content := container.NewVBox(widget.NewLabel("Timer ran out"), buttons)

	popUp = dialog.NewCustom("Timer Expired", "Close", content, ui.window)
	// show popup
	ui.window.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		if event.Name == fyne.KeyEscape {
			popUp.Hide()
		}
	})
	popUp.Show()
}

func startRecording(client *goobs.Client) error {
	_, err := client.Record.StartRecord()
	return err
}

func stopRecording(client *goobs.Client) error {
	_, err := client.Record.StopRecord()
	return err
}
