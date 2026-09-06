import { useCallback, useEffect, useRef, useState } from "react";

import { AttendanceService, SignInOutcome } from "../../bindings/ttpattendance";
import type { Student } from "../../bindings/ttpattendance/internal/sheets";
import { errorMessage } from "../errors";
import { parseIdCard } from "../idCard";

/**
 * INACTIVITY_LIMIT_MS is how long the kiosk waits before clearing the form and
 * putting focus back on the hidden card-reader field. Without this, someone who
 * taps a text box and walks away leaves the reader unfocused, and the next
 * student's swipe goes nowhere.
 */
const INACTIVITY_LIMIT_MS = 15_000;

interface Props {
  /** ready is false until the Google Sheets connection is configured. */
  ready: boolean;
  onNeedsMajor: (student: Student) => void;
  onSignedIn: (student: Student) => void;
  onOpenSettings: () => void;
}

interface FieldErrors {
  firstName?: string;
  lastName?: string;
  studentId?: string;
}

/** validate mirrors the checks the Go service runs, so mistakes are caught locally. */
function validate(student: Student): FieldErrors {
  const errors: FieldErrors = {};

  if (!student.firstName.trim()) {
    errors.firstName = "Please enter your first name.";
  }
  if (!student.lastName.trim()) {
    errors.lastName = "Please enter your last name.";
  }

  const studentId = student.studentId.trim();
  if (!studentId) {
    errors.studentId = "Please enter your student id.";
  } else if (!/^\d+$/.test(studentId)) {
    errors.studentId = "Your student id number can only contain digits.";
  } else if (studentId.length !== 9) {
    errors.studentId = "Your student id number must be 9 digits long.";
  }

  return errors;
}

export default function SignInView({ ready, onNeedsMajor, onSignedIn, onOpenSettings }: Props) {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [studentId, setStudentId] = useState("");
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({});
  const [submitError, setSubmitError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const cardReaderRef = useRef<HTMLInputElement>(null);

  // Held in a ref rather than state: it changes on every mouse move, and
  // re-rendering the form for that would fight with typing.
  const lastActivityRef = useRef(Date.now());

  const clearForm = useCallback(() => {
    setFirstName("");
    setLastName("");
    setStudentId("");
    setFieldErrors({});
    setSubmitError("");
  }, []);

  const submit = useCallback(
    async (student: Student) => {
      setSubmitting(true);
      setSubmitError("");

      try {
        const outcome = await AttendanceService.SignIn(student);

        if (outcome === SignInOutcome.OutcomeNeedsMajor) {
          onNeedsMajor(student);
        } else {
          onSignedIn(student);
        }

        clearForm();
      } catch (error) {
        setSubmitError(errorMessage(error, "There was an error signing you in. Please try again."));
      } finally {
        setSubmitting(false);
      }
    },
    [clearForm, onNeedsMajor, onSignedIn],
  );

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();

    const student: Student = {
      firstName: firstName.trim(),
      lastName: lastName.trim(),
      studentId: studentId.trim(),
    };

    const errors = validate(student);
    setFieldErrors(errors);

    if (Object.keys(errors).length === 0) {
      void submit(student);
    }
  };

  /**
   * The reader types the whole track and the field is off-screen, so rather than
   * waiting for a trailing Enter that not every model sends, each keystroke is
   * offered to the parser and the swipe is submitted as soon as one parses. The
   * student id is the last thing on the track, so a partial swipe cannot match.
   */
  const handleSwipe = (event: React.ChangeEvent<HTMLInputElement>) => {
    const card = parseIdCard(event.target.value);
    if (!card) {
      return;
    }

    event.target.value = "";
    void submit(card);
  };

  // Focus the reader on arrival, and again whenever the kiosk has been idle.
  useEffect(() => {
    cardReaderRef.current?.focus();

    const noteActivity = () => {
      lastActivityRef.current = Date.now();
    };

    window.addEventListener("keypress", noteActivity);
    window.addEventListener("mousemove", noteActivity);
    window.addEventListener("click", noteActivity);

    const timer = window.setInterval(() => {
      if (Date.now() - lastActivityRef.current < INACTIVITY_LIMIT_MS) {
        return;
      }

      lastActivityRef.current = Date.now();
      clearForm();
      cardReaderRef.current?.focus();
    }, 1000);

    return () => {
      window.removeEventListener("keypress", noteActivity);
      window.removeEventListener("mousemove", noteActivity);
      window.removeEventListener("click", noteActivity);
      window.clearInterval(timer);
    };
  }, [clearForm]);

  return (
    <main className="container">
      <div className="logo">
        <img src="/img/ttp-logo.jpg" alt="Transfer Student Center Logo" />
      </div>

      {!ready && (
        <div className="notice" role="status">
          This app is not connected to a Google Sheet yet.{" "}
          <button type="button" className="link" onClick={onOpenSettings}>
            Open Settings
          </button>
        </div>
      )}

      <form className="form" autoComplete="off" onSubmit={handleSubmit} noValidate>
        <div className="form-group">
          <label htmlFor="first-name">First Name</label>
          <input
            id="first-name"
            className="form-control"
            type="text"
            value={firstName}
            onChange={(event) => setFirstName(event.target.value)}
          />
          {fieldErrors.firstName && <p className="error">{fieldErrors.firstName}</p>}
        </div>

        <div className="form-group">
          <label htmlFor="last-name">Last Name</label>
          <input
            id="last-name"
            className="form-control"
            type="text"
            value={lastName}
            onChange={(event) => setLastName(event.target.value)}
          />
          {fieldErrors.lastName && <p className="error">{fieldErrors.lastName}</p>}
        </div>

        <div className="form-group">
          <label htmlFor="student-id">Student ID</label>
          <input
            id="student-id"
            className="form-control"
            type="text"
            inputMode="numeric"
            value={studentId}
            onChange={(event) => setStudentId(event.target.value)}
          />
          {fieldErrors.studentId && <p className="error">{fieldErrors.studentId}</p>}
        </div>

        {submitError && <p className="error error-banner">{submitError}</p>}

        <div className="form-group">
          <button className="btn btn-primary" type="submit" disabled={submitting}>
            {submitting ? "Signing in…" : "Sign In"}
          </button>
        </div>
      </form>

      {/* The card reader types into this field; it is present but never seen. */}
      <input
        ref={cardReaderRef}
        className="card-reader"
        type="text"
        aria-hidden="true"
        tabIndex={-1}
        autoComplete="off"
        onChange={handleSwipe}
      />

      <button type="button" className="settings-button" onClick={onOpenSettings} title="Settings">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
          <circle cx="12" cy="12" r="3" />
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.6a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
        </svg>
        <span className="sr-only">Settings</span>
      </button>
    </main>
  );
}
