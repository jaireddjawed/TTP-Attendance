// Package sheets records attendance in the Google Sheet backing the TTP Center.
//
// The spreadsheet holds one "Students" tab acting as a directory of everyone who
// has ever signed in, plus one tab per month ("Jan 2026") holding that month's
// sign-ins. Month tabs are created on demand.
package sheets

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/option"
	sheetsapi "google.golang.org/api/sheets/v4"

	"ttpattendance/internal/config"
)

// directorySheet is the tab holding every student who has signed in before.
const directorySheet = "Students"

// monthSheetLayout is the Go equivalent of Python's "%b %Y" — e.g. "Jan 2026".
const monthSheetLayout = "Jan 2006"

const (
	signInDateLayout = "01/02/2006"
	signInTimeLayout = "03:04 PM"
)

// columns in the directory sheet.
const (
	colFirstName = 0
	colLastName  = 1
	colStudentID = 2
	colMajor     = 3
)

// Student is one person signing in.
type Student struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	StudentID string `json:"studentId"`
	Major     string `json:"major,omitempty"`
}

// Client talks to one attendance spreadsheet as the signed-in Google account.
type Client struct {
	service       *sheetsapi.Service
	spreadsheetID string
}

// NewClient builds a client from the stored credentials, token and spreadsheet id.
// It fails with ErrNotAuthorized, config.ErrNoCredentials or config.ErrNoSpreadsheetID
// when the app has not been set up yet, which the UI turns into a prompt to open Settings.
func NewClient(ctx context.Context) (*Client, error) {
	spreadsheetID, err := config.SpreadsheetID()
	if err != nil {
		return nil, err
	}

	client, err := httpClient(ctx)
	if err != nil {
		return nil, err
	}

	service, err := sheetsapi.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("could not reach Google Sheets: %w", err)
	}

	return &Client{service: service, spreadsheetID: spreadsheetID}, nil
}

// IsInDirectory reports whether a student with this id has signed in before.
func (c *Client) IsInDirectory(ctx context.Context, studentID string) (bool, error) {
	rows, err := c.directory(ctx)
	if err != nil {
		return false, err
	}

	for _, row := range rows {
		// Rows are ragged — Sheets omits trailing empty cells — and the first row
		// is the header, so both need guarding before reading the id column.
		if len(row) <= colStudentID {
			continue
		}

		if id, ok := row[colStudentID].(string); ok && strings.EqualFold(strings.TrimSpace(id), studentID) {
			return true, nil
		}
	}

	return false, nil
}

// directory returns every row of the student directory.
func (c *Client) directory(ctx context.Context) ([][]any, error) {
	values, err := c.service.Spreadsheets.Values.
		Get(c.spreadsheetID, rangeRef(directorySheet, "A1:E2000")).
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("could not read the student directory: %w", wrapAPIError(err))
	}

	return values.Values, nil
}

// AddToDirectory records a student and their major so they are recognised next time.
func (c *Client) AddToDirectory(ctx context.Context, student Student) error {
	row := []any{student.FirstName, student.LastName, student.StudentID, student.Major}

	_, err := c.service.Spreadsheets.Values.
		Append(c.spreadsheetID, rangeRef(directorySheet, "A1:E2000"), &sheetsapi.ValueRange{Values: [][]any{row}}).
		ValueInputOption("USER_ENTERED").
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("could not add you to the student directory: %w", wrapAPIError(err))
	}

	return nil
}

// RecordSignIn appends a timestamped sign-in to the current month's sheet,
// creating that sheet first if this is the month's first sign-in.
func (c *Client) RecordSignIn(ctx context.Context, student Student, at time.Time) error {
	sheetName := at.Format(monthSheetLayout)

	if err := c.ensureMonthSheet(ctx, sheetName); err != nil {
		return err
	}

	row := []any{
		student.FirstName,
		student.LastName,
		student.StudentID,
		at.Format(signInDateLayout),
		at.Format(signInTimeLayout),
	}

	_, err := c.service.Spreadsheets.Values.
		Append(c.spreadsheetID, rangeRef(sheetName, "A1:E1000"), &sheetsapi.ValueRange{Values: [][]any{row}}).
		ValueInputOption("USER_ENTERED").
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("could not record your sign-in: %w", wrapAPIError(err))
	}

	return nil
}

// ensureMonthSheet creates the given month's sheet, with headers, if it is missing.
func (c *Client) ensureMonthSheet(ctx context.Context, sheetName string) error {
	spreadsheet, err := c.service.Spreadsheets.Get(c.spreadsheetID).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("could not open the attendance spreadsheet: %w", wrapAPIError(err))
	}

	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties != nil && sheet.Properties.Title == sheetName {
			return nil
		}
	}

	_, err = c.service.Spreadsheets.BatchUpdate(c.spreadsheetID, &sheetsapi.BatchUpdateSpreadsheetRequest{
		Requests: []*sheetsapi.Request{{
			AddSheet: &sheetsapi.AddSheetRequest{
				Properties: &sheetsapi.SheetProperties{
					Title: sheetName,
					GridProperties: &sheetsapi.GridProperties{
						RowCount:    1000,
						ColumnCount: 6,
					},
				},
			},
		}},
	}).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("could not create the %s sheet: %w", sheetName, wrapAPIError(err))
	}

	headers := []any{"First Name", "Last Name", "Student ID", "Sign In Date", "Sign In Time"}

	_, err = c.service.Spreadsheets.Values.
		Update(c.spreadsheetID, rangeRef(sheetName, "A1:E1"), &sheetsapi.ValueRange{Values: [][]any{headers}}).
		ValueInputOption("USER_ENTERED").
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("could not add headers to the %s sheet: %w", sheetName, wrapAPIError(err))
	}

	return nil
}

// rangeRef builds an A1 reference. Month sheet names contain a space, so the
// name is always quoted; any literal quote in it is escaped by doubling.
func rangeRef(sheetName, cells string) string {
	return fmt.Sprintf("'%s'!%s", strings.ReplaceAll(sheetName, "'", "''"), cells)
}
