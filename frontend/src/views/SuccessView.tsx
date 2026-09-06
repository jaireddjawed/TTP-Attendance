import { useEffect, useState } from "react";

import type { Student } from "../../bindings/ttpattendance/internal/sheets";

/** WELCOME_DURATION_MS is how long the welcome sits before the kiosk resets. */
const WELCOME_DURATION_MS = 3000;

interface Props {
  student: Student;
  onDone: () => void;
}

/**
 * displayFirstName fixes up a first name for the greeting only — what was
 * swiped is still what gets written to the sheet.
 *
 * Marisol's card truncates her first name to six characters, so the welcome
 * screen would otherwise greet her as MARISO.
 */
function displayFirstName(firstName: string): string {
  return firstName === "MARISO" ? "MARISOL" : firstName;
}

export default function SuccessView({ student, onDone }: Props) {
  const [signedInAt] = useState(() => new Date());

  useEffect(() => {
    const timer = window.setTimeout(onDone, WELCOME_DURATION_MS);
    return () => window.clearTimeout(timer);
  }, [onDone]);

  return (
    <main className="container container-centered">
      <div className="welcome-box">
        <img src="/img/ucr-logo.png" alt="UCR Logo" className="ucr-logo" />
        <div>
          <h1 className="welcome-title">
            Welcome {displayFirstName(student.firstName)} {student.lastName}!
          </h1>
          <h3 className="welcome-time">{signedInAt.toLocaleString()}</h3>
        </div>
      </div>
    </main>
  );
}
