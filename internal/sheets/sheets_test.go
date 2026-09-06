package sheets

import (
	"testing"
	"time"
)

func TestRangeRef(t *testing.T) {
	tests := []struct {
		name      string
		sheetName string
		cells     string
		want      string
	}{
		{
			name:      "quotes a month sheet name containing a space",
			sheetName: "Jan 2026",
			cells:     "A1:E1000",
			want:      "'Jan 2026'!A1:E1000",
		},
		{
			name:      "quotes the directory sheet",
			sheetName: directorySheet,
			cells:     "A1:E2000",
			want:      "'Students'!A1:E2000",
		},
		{
			name:      "escapes a quote in the sheet name by doubling it",
			sheetName: "Spring '26",
			cells:     "A1:E1",
			want:      "'Spring ''26'!A1:E1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := rangeRef(test.sheetName, test.cells); got != test.want {
				t.Errorf("rangeRef(%q, %q) = %q, want %q", test.sheetName, test.cells, got, test.want)
			}
		})
	}
}

// TestSignInLayouts pins the formats down: they have to keep matching the
// columns already written by the Python version, or a month's sheet ends up
// with two different date formats in it.
func TestSignInLayouts(t *testing.T) {
	at := time.Date(2026, time.January, 9, 15, 4, 5, 0, time.UTC)

	if got, want := at.Format(monthSheetLayout), "Jan 2026"; got != want {
		t.Errorf("month sheet name = %q, want %q", got, want)
	}
	if got, want := at.Format(signInDateLayout), "01/09/2026"; got != want {
		t.Errorf("sign-in date = %q, want %q", got, want)
	}
	if got, want := at.Format(signInTimeLayout), "03:04 PM"; got != want {
		t.Errorf("sign-in time = %q, want %q", got, want)
	}
}
