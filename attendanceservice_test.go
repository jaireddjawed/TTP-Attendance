package main

import (
	"testing"

	"ttpattendance/internal/sheets"
)

func TestValidate(t *testing.T) {
	valid := sheets.Student{FirstName: "Jane", LastName: "Doe", StudentID: "861234567"}

	tests := []struct {
		name         string
		student      sheets.Student
		requireMajor bool
		wantErr      string
	}{
		{
			name:    "accepts a complete sign-in",
			student: valid,
		},
		{
			name:    "trims surrounding whitespace",
			student: sheets.Student{FirstName: "  Jane ", LastName: " Doe", StudentID: " 861234567 "},
		},
		{
			name:    "rejects a missing first name",
			student: sheets.Student{FirstName: "   ", LastName: "Doe", StudentID: "861234567"},
			wantErr: "please enter your first name",
		},
		{
			name:    "rejects a missing last name",
			student: sheets.Student{FirstName: "Jane", StudentID: "861234567"},
			wantErr: "please enter your last name",
		},
		{
			name:    "rejects a short student id",
			student: sheets.Student{FirstName: "Jane", LastName: "Doe", StudentID: "8612345"},
			wantErr: "your student id number must be 9 digits long",
		},
		{
			name:    "rejects a non-numeric student id",
			student: sheets.Student{FirstName: "Jane", LastName: "Doe", StudentID: "86123456x"},
			wantErr: "your student id number must be 9 digits long",
		},
		{
			name:         "rejects a missing major when one is required",
			student:      valid,
			requireMajor: true,
			wantErr:      "please select your major",
		},
		{
			name:         "accepts a major when one is required",
			student:      sheets.Student{FirstName: "Jane", LastName: "Doe", StudentID: "861234567", Major: "robotics"},
			requireMajor: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := validate(test.student, test.requireMajor)

			if test.wantErr != "" {
				if err == nil {
					t.Fatalf("validate() succeeded, want error %q", test.wantErr)
				}
				if err.Error() != test.wantErr {
					t.Fatalf("validate() error = %q, want %q", err, test.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("validate() error = %v, want nil", err)
			}
			if got.FirstName != "Jane" || got.LastName != "Doe" || got.StudentID != "861234567" {
				t.Errorf("validate() = %+v, want the fields trimmed to Jane/Doe/861234567", got)
			}
		})
	}
}

func TestExtractSpreadsheetID(t *testing.T) {
	const id = "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "keeps a bare id", value: id, want: id},
		{name: "trims whitespace", value: "  " + id + "\n", want: id},
		{
			name:  "pulls the id out of a full URL",
			value: "https://docs.google.com/spreadsheets/d/" + id + "/edit#gid=0",
			want:  id,
		},
		{
			name:  "pulls the id out of a URL without a trailing path",
			value: "https://docs.google.com/spreadsheets/d/" + id,
			want:  id,
		},
		{name: "leaves an empty value empty", value: "   ", want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := extractSpreadsheetID(test.value); got != test.want {
				t.Errorf("extractSpreadsheetID(%q) = %q, want %q", test.value, got, test.want)
			}
		})
	}
}
