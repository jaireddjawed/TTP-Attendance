package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"ttpattendance/internal/sheets"
)

// SignInOutcome tells the frontend which screen to show after a sign-in attempt.
type SignInOutcome string

const (
	// OutcomeSignedIn means the sign-in was recorded; show the welcome screen.
	OutcomeSignedIn SignInOutcome = "signed-in"

	// OutcomeNeedsMajor means this student is not in the directory yet and has to
	// pick a major before the sign-in can be recorded.
	OutcomeNeedsMajor SignInOutcome = "needs-major"
)

// AttendanceService records student attendance. It replaces the /submit-signin
// and /submit-major endpoints of the Flask version.
type AttendanceService struct{}

// SignIn records a sign-in for a student who has signed in before.
//
// A student the directory has never seen is not recorded: the frontend sends
// them to the major screen and calls SignInWithMajor instead.
func (a *AttendanceService) SignIn(ctx context.Context, student sheets.Student) (SignInOutcome, error) {
	student, err := validate(student, false)
	if err != nil {
		return "", err
	}

	client, err := sheets.NewClient(ctx)
	if err != nil {
		return "", err
	}

	known, err := client.IsInDirectory(ctx, student.StudentID)
	if err != nil {
		return "", err
	}
	if !known {
		return OutcomeNeedsMajor, nil
	}

	if err := client.RecordSignIn(ctx, student, time.Now()); err != nil {
		return "", err
	}

	return OutcomeSignedIn, nil
}

// SignInWithMajor adds a first-time student to the directory and records their
// sign-in in one step.
func (a *AttendanceService) SignInWithMajor(ctx context.Context, student sheets.Student) error {
	student, err := validate(student, true)
	if err != nil {
		return err
	}

	client, err := sheets.NewClient(ctx)
	if err != nil {
		return err
	}

	if err := client.AddToDirectory(ctx, student); err != nil {
		return err
	}

	return client.RecordSignIn(ctx, student, time.Now())
}

// validate trims the submitted fields and rejects anything the sheet should not
// receive. The frontend checks the same rules, but it is the only caller we
// control, so the service does not take them on trust.
func validate(student sheets.Student, requireMajor bool) (sheets.Student, error) {
	student.FirstName = strings.TrimSpace(student.FirstName)
	student.LastName = strings.TrimSpace(student.LastName)
	student.StudentID = strings.TrimSpace(student.StudentID)
	student.Major = strings.TrimSpace(student.Major)

	if student.FirstName == "" {
		return student, errors.New("please enter your first name")
	}
	if student.LastName == "" {
		return student, errors.New("please enter your last name")
	}
	if !isStudentID(student.StudentID) {
		return student, errors.New("your student id number must be 9 digits long")
	}
	if requireMajor && student.Major == "" {
		return student, errors.New("please select your major")
	}

	return student, nil
}

// isStudentID reports whether id is a 9-digit UCR student id.
func isStudentID(id string) bool {
	if len(id) != 9 {
		return false
	}

	for _, r := range id {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

// ServiceName identifies this service in Wails' logs.
func (a *AttendanceService) ServiceName() string {
	return "attendance"
}
