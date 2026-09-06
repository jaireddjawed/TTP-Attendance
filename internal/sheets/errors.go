package sheets

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/api/googleapi"
)

// wrapAPIError turns the Google API's raw errors into something a student or a
// front-desk worker can act on. Anything unrecognised is passed through as-is.
func wrapAPIError(err error) error {
	var apiErr *googleapi.Error
	if !errors.As(err, &apiErr) {
		return err
	}

	switch apiErr.Code {
	case http.StatusUnauthorized:
		return fmt.Errorf("%w (Google rejected the saved token)", ErrNotAuthorized)

	case http.StatusForbidden:
		return errors.New("the signed-in Google account does not have access to this spreadsheet — share it with that account, or sign in as a different one in Settings")

	case http.StatusNotFound:
		return errors.New("no spreadsheet with the configured id was found — check the spreadsheet id in Settings")

	case http.StatusTooManyRequests:
		return errors.New("Google is rate limiting requests right now — wait a moment and try again")
	}

	if apiErr.Message != "" {
		return errors.New(apiErr.Message)
	}

	return err
}
