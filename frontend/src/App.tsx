import { useCallback, useEffect, useState } from "react";

import { SettingsService } from "../bindings/ttpattendance";
import type { Settings } from "../bindings/ttpattendance";
import type { Student } from "../bindings/ttpattendance/internal/sheets";

import MajorView from "./views/MajorView";
import SettingsView from "./views/SettingsView";
import SignInView from "./views/SignInView";
import SuccessView from "./views/SuccessView";

/** Screen is the view currently on the kiosk, and what it needs to render. */
type Screen =
  | { name: "signIn" }
  | { name: "major"; student: Student }
  | { name: "success"; student: Student }
  | { name: "settings" };

/** isReady reports whether the Google Sheets connection is fully set up. */
function isReady(settings: Settings | null): boolean {
  return Boolean(settings?.spreadsheetId && settings.hasCredentials && settings.isAuthorized);
}

export default function App() {
  const [screen, setScreen] = useState<Screen>({ name: "signIn" });
  const [settings, setSettings] = useState<Settings | null>(null);

  // Load the setup state once at launch. Until it arrives nothing is rendered,
  // which avoids flashing the "not set up" notice on a configured machine.
  useEffect(() => {
    SettingsService.Get()
      .then((loaded) => {
        setSettings(loaded);
        if (!isReady(loaded)) {
          setScreen({ name: "settings" });
        }
      })
      .catch(console.error);
  }, []);

  const returnToSignIn = useCallback(() => setScreen({ name: "signIn" }), []);

  if (settings === null) {
    return <main className="container" />;
  }

  return (
    <>
      {screen.name === "signIn" && (
        <SignInView
          ready={isReady(settings)}
          onNeedsMajor={(student) => setScreen({ name: "major", student })}
          onSignedIn={(student) => setScreen({ name: "success", student })}
          onOpenSettings={() => setScreen({ name: "settings" })}
        />
      )}

      {screen.name === "major" && (
        <MajorView
          student={screen.student}
          onSignedIn={(student) => setScreen({ name: "success", student })}
          onCancel={returnToSignIn}
        />
      )}

      {screen.name === "success" && <SuccessView student={screen.student} onDone={returnToSignIn} />}

      {screen.name === "settings" && (
        <SettingsView settings={settings} onChange={setSettings} onClose={returnToSignIn} />
      )}

      <p className="made-with">Made with ❤️ by Jaired Jawed.</p>
    </>
  );
}
