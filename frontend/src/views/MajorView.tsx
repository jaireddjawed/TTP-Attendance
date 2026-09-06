import { useState } from "react";

import { AttendanceService } from "../../bindings/ttpattendance";
import type { Student } from "../../bindings/ttpattendance/internal/sheets";
import { errorMessage } from "../errors";
import { MAJORS, OTHER_MAJOR } from "../majors";

interface Props {
  /** student is the first-time student whose sign-in is waiting on a major. */
  student: Student;
  onSignedIn: (student: Student) => void;
  onCancel: () => void;
}

/**
 * MajorView asks a student the app has not seen before for their major. Their
 * sign-in is only recorded once they answer, which also adds them to the
 * directory so they go straight through next time.
 */
export default function MajorView({ student, onSignedIn, onCancel }: Props) {
  const [selected, setSelected] = useState("");
  const [otherMajor, setOtherMajor] = useState("");
  const [fieldError, setFieldError] = useState("");
  const [submitError, setSubmitError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (!selected) {
      setFieldError("Please select a major.");
      return;
    }

    const major = selected === OTHER_MAJOR ? otherMajor.trim() : selected;
    if (!major) {
      setFieldError("Please enter your major.");
      return;
    }

    setFieldError("");
    setSubmitError("");
    setSubmitting(true);

    try {
      const withMajor: Student = { ...student, major };
      await AttendanceService.SignInWithMajor(withMajor);
      onSignedIn(withMajor);
    } catch (error) {
      setSubmitError(errorMessage(error, "There was an error signing you in. Please try again."));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <main className="container">
      <div className="logo">
        <img src="/img/ttp-logo.jpg" alt="Transfer Student Center Logo" />
      </div>

      <p className="prompt">
        Welcome, {student.firstName}! Tell us your major and we'll remember it next time.
      </p>

      <form className="form" autoComplete="off" onSubmit={handleSubmit} noValidate>
        <div className="form-group">
          <label htmlFor="engineering-major">Major</label>
          <select
            id="engineering-major"
            className="form-control"
            value={selected}
            onChange={(event) => setSelected(event.target.value)}
          >
            <option value="" />
            {MAJORS.map((major) => (
              <option key={major.value} value={major.value}>
                {major.label}
              </option>
            ))}
          </select>
        </div>

        {selected === OTHER_MAJOR && (
          <div className="form-group">
            <input
              className="form-control"
              type="text"
              placeholder="Enter your major here"
              value={otherMajor}
              onChange={(event) => setOtherMajor(event.target.value)}
              autoFocus
            />
          </div>
        )}

        {fieldError && <p className="error">{fieldError}</p>}
        {submitError && <p className="error error-banner">{submitError}</p>}

        <div className="form-group">
          <button className="btn btn-primary" type="submit" disabled={submitting}>
            {submitting ? "Signing in…" : "Submit"}
          </button>
        </div>

        <button type="button" className="link" onClick={onCancel}>
          Cancel
        </button>
      </form>
    </main>
  );
}
