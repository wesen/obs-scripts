package main

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/andreykaipov/goobs"
)

func RunCountdownTimer(ctx context.Context, w fyne.Window, client *goobs.Client) error {
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		dialog.ShowInformation("Timer", "Timer ran out", w)
		return nil
	}
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

	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel

	statusLabel := widget.NewLabel("")

	startButton := widget.NewButton("Start Recording", func() {
		err := startRecording(client)
		if err != nil {
			statusLabel.SetText("Error starting recording: " + err.Error())
		} else {
			statusLabel.SetText("Recording started.")
			go func() {
				err := RunCountdownTimer(ctx, w, client)
				if err != nil {
					statusLabel.SetText("Error running countdown timer: " + err.Error())
				}
			}()
		}
	})

	stopButton := widget.NewButton("Stop Recording", func() {
		err := stopRecording(client)
		if err != nil {
			statusLabel.SetText("Error stopping recording: " + err.Error())
		} else {
			statusLabel.SetText("Recording stopped.")
		}
	})

	buttons := container.NewHBox(startButton, stopButton)
	content := container.NewVBox(buttons, statusLabel)

	w.SetContent(content)
	w.Resize(fyne.NewSize(300, 200))
	w.ShowAndRun()
}

func startRecording(client *goobs.Client) error {
	_, err := client.Record.StartRecord()
	return err
}

func stopRecording(client *goobs.Client) error {
	_, err := client.Record.StopRecord()
	return err
}
