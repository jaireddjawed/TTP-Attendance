/**
 * errorMessage turns a rejected service call into text worth showing someone.
 *
 * Wails rejects with an Error carrying the message the Go service returned, so
 * in practice this is the message written in attendanceservice.go or the sheets
 * package — already phrased for the person standing at the desk.
 */
export function errorMessage(error: unknown, fallback: string): string {
  if (error instanceof Error && error.message) {
    return error.message;
  }

  if (typeof error === "string" && error) {
    return error;
  }

  return fallback;
}
