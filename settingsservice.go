package main

import (
	"context"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"

	"ttpattendance/internal/config"
	"ttpattendance/internal/sheets"
)

// Settings is the setup state the Settings screen renders.
type Settings struct {
	// SpreadsheetID is the attendance spreadsheet this app writes to.
	SpreadsheetID string `json:"spreadsheetId"`

	// HasCredentials reports whether a Google OAuth credentials.json has been added.
	HasCredentials bool `json:"hasCredentials"`

	// IsAuthorized reports whether someone has signed in with a Google account.
	IsAuthorized bool `json:"isAuthorized"`

	// ConfigDir is where all of the above is stored, shown so it can be found on disk.
	ConfigDir string `json:"configDir"`
}

// Ready reports whether the app has everything it needs to record attendance.
func (s Settings) Ready() bool {
	return s.SpreadsheetID != "" && s.HasCredentials && s.IsAuthorized
}

// SettingsService configures the Google Sheets connection. The Flask version
// required hand-editing sheetInfo.json and dropping credentials.json next to the
// source; a packaged app has neither, so this is exposed in the UI instead.
type SettingsService struct{}

// Get returns the current setup state.
func (s *SettingsService) Get() (Settings, error) {
	info, err := config.ReadSheetInfo()
	if err != nil {
		return Settings{}, err
	}

	dir, err := config.Dir()
	if err != nil {
		return Settings{}, err
	}

	return Settings{
		SpreadsheetID:  info.SpreadsheetID,
		HasCredentials: config.HasCredentials(),
		IsAuthorized:   sheets.IsAuthorized(),
		ConfigDir:      dir,
	}, nil
}

// spreadsheetURLPattern pulls the id out of a full Google Sheets URL, so the id
// can be pasted straight from the browser's address bar.
var spreadsheetURLPattern = regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)

// extractSpreadsheetID normalises what someone pasted into a spreadsheet id.
func extractSpreadsheetID(value string) string {
	value = strings.TrimSpace(value)

	if match := spreadsheetURLPattern.FindStringSubmatch(value); match != nil {
		return match[1]
	}

	return value
}

// SetSpreadsheetID stores the attendance spreadsheet id. It accepts either a
// bare id or a full spreadsheet URL.
func (s *SettingsService) SetSpreadsheetID(value string) (Settings, error) {
	id := extractSpreadsheetID(value)
	if id == "" {
		return Settings{}, errors.New("please enter the spreadsheet id")
	}

	if err := config.WriteSpreadsheetID(id); err != nil {
		return Settings{}, err
	}

	return s.Get()
}

// ChooseCredentialsFile opens a native file picker for the credentials.json
// downloaded from the Google Cloud console and stores a copy.
//
// Choosing nothing is not an error: it returns the unchanged settings.
func (s *SettingsService) ChooseCredentialsFile() (Settings, error) {
	path, err := application.Get().Dialog.OpenFile().
		SetTitle("Choose credentials.json").
		AddFilter("Google OAuth credentials", "*.json").
		CanChooseFiles(true).
		PromptForSingleSelection()
	if err != nil {
		return Settings{}, err
	}
	if path == "" {
		return s.Get()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Settings{}, err
	}

	if err := config.WriteCredentials(data); err != nil {
		return Settings{}, err
	}

	return s.Get()
}

// Authorize signs in to the Google account that owns the attendance spreadsheet,
// via the browser. It blocks until the browser flow finishes.
func (s *SettingsService) Authorize(ctx context.Context) (Settings, error) {
	app := application.Get()

	if err := sheets.Authorize(ctx, app.Browser.OpenURL); err != nil {
		return Settings{}, err
	}

	return s.Get()
}

// SignOut forgets the saved Google authorisation.
func (s *SettingsService) SignOut() (Settings, error) {
	if err := config.DeleteToken(); err != nil {
		return Settings{}, err
	}

	return s.Get()
}

// OpenConfigDir reveals the configuration directory in the file manager.
func (s *SettingsService) OpenConfigDir() error {
	dir, err := config.Dir()
	if err != nil {
		return err
	}

	return application.Get().Browser.OpenFile(dir)
}

// ServiceName identifies this service in Wails' logs.
func (s *SettingsService) ServiceName() string {
	return "settings"
}
