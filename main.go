// TTP Attendance records student attendance at the BCOE Transfer Student Center
// at UC Riverside. Students sign in by filling in a form or by swiping their
// student id card, and each sign-in is appended to a Google Sheet.
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "TTP Attendance",
		Description: "Student attendance for the BCOE Transfer Student Center at UC Riverside",
		Services: []application.Service{
			application.NewService(&AttendanceService{}),
			application.NewService(&SettingsService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// The app runs unattended on the front desk, so it opens full screen with a
	// dark background matching the sign-in page to avoid a white flash on launch.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "TTP Attendance",
		Width:            1280,
		Height:           800,
		StartState:       application.WindowStateFullscreen,
		BackgroundColour: application.NewRGB(0, 39, 77),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
