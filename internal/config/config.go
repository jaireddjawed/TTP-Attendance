// Package config locates and persists the files TTP Attendance needs to talk to
// Google Sheets: the OAuth client credentials, the saved authorisation token and
// the id of the attendance spreadsheet.
//
// A packaged .app has no writable working directory, so everything lives in the
// per-user application support directory instead. For the benefit of people
// upgrading from the Flask version, files sitting next to the executable are
// still read if the config directory has nothing yet.
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const appDirName = "TTP Attendance"

const (
	credentialsFileName = "credentials.json"
	tokenFileName       = "token.json"
	sheetInfoFileName   = "sheetInfo.json"
)

// ErrNoCredentials is returned when no OAuth client credentials have been added yet.
var ErrNoCredentials = errors.New("no credentials.json found: add your Google OAuth client credentials in Settings")

// ErrNoSpreadsheetID is returned when the attendance spreadsheet has not been chosen yet.
var ErrNoSpreadsheetID = errors.New("no spreadsheet id configured: add the attendance spreadsheet id in Settings")

// SheetInfo mirrors the sheetInfo.json file used by the previous version.
type SheetInfo struct {
	SpreadsheetID string   `json:"spreadsheetId"`
	Scopes        []string `json:"scopes,omitempty"`
}

// Dir returns the directory holding this application's configuration, creating
// it if necessary.
func Dir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(base, appDirName)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}

	return dir, nil
}

// CredentialsPath returns the path to the OAuth client credentials file.
func CredentialsPath() (string, error) { return resolve(credentialsFileName) }

// TokenPath returns the path to the saved authorisation token.
func TokenPath() (string, error) { return resolve(tokenFileName) }

// SheetInfoPath returns the path to the spreadsheet configuration file.
func SheetInfoPath() (string, error) { return resolve(sheetInfoFileName) }

// resolve returns the config directory path for name, falling back to a copy in
// the working directory when the config directory does not have one. This lets
// an existing checkout of the Flask version keep working without re-authorising.
func resolve(name string) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, name)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	if legacy, err := filepath.Abs(name); err == nil {
		if _, err := os.Stat(legacy); err == nil {
			return legacy, nil
		}
	}

	return path, nil
}

// ReadCredentials returns the raw OAuth client credentials.
func ReadCredentials() ([]byte, error) {
	path, err := CredentialsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNoCredentials
	}

	return data, err
}

// HasCredentials reports whether OAuth client credentials have been added.
func HasCredentials() bool {
	_, err := ReadCredentials()
	return err == nil
}

// WriteCredentials stores the given OAuth client credentials, replacing any
// existing ones. The contents are checked for the shape Google hands out so a
// mistakenly chosen file is rejected before it breaks the sign-in flow.
func WriteCredentials(data []byte) error {
	var probe struct {
		Installed json.RawMessage `json:"installed"`
		Web       json.RawMessage `json:"web"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return errors.New("that file is not valid JSON — download credentials.json from the Google Cloud console")
	}
	if probe.Installed == nil && probe.Web == nil {
		return errors.New("that file does not look like an OAuth client credentials file — it should contain an \"installed\" or \"web\" section")
	}

	dir, err := Dir()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, credentialsFileName), data, 0o600)
}

// ReadToken returns the saved authorisation token, or nil when there is none.
func ReadToken() ([]byte, error) {
	path, err := TokenPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}

	return data, err
}

// WriteToken saves the authorisation token for the next run.
func WriteToken(data []byte) error {
	dir, err := Dir()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, tokenFileName), data, 0o600)
}

// DeleteToken removes the saved authorisation token, forcing a fresh sign-in.
func DeleteToken() error {
	path, err := TokenPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

// ReadSheetInfo returns the spreadsheet configuration. A missing file yields a
// zero value rather than an error so Settings can render before it is set up.
func ReadSheetInfo() (SheetInfo, error) {
	var info SheetInfo

	path, err := SheetInfoPath()
	if err != nil {
		return info, err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return info, nil
	} else if err != nil {
		return info, err
	}

	if err := json.Unmarshal(data, &info); err != nil {
		return SheetInfo{}, err
	}

	return info, nil
}

// SpreadsheetID returns the configured attendance spreadsheet id.
func SpreadsheetID() (string, error) {
	info, err := ReadSheetInfo()
	if err != nil {
		return "", err
	}
	if info.SpreadsheetID == "" {
		return "", ErrNoSpreadsheetID
	}

	return info.SpreadsheetID, nil
}

// WriteSpreadsheetID stores the attendance spreadsheet id.
func WriteSpreadsheetID(id string) error {
	info, err := ReadSheetInfo()
	if err != nil {
		return err
	}
	info.SpreadsheetID = id

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}

	dir, err := Dir()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, sheetInfoFileName), append(data, '\n'), 0o600)
}
