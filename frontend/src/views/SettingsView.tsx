import { useState } from "react";

import { SettingsService } from "../../bindings/ttpattendance";
import type { Settings } from "../../bindings/ttpattendance";
import { errorMessage } from "../errors";

interface Props {
  settings: Settings;
  onChange: (settings: Settings) => void;
  onClose: () => void;
}

/**
 * SettingsView connects the app to a Google Sheet. The Flask version needed
 * credentials.json dropped next to the source and a spreadsheet id hand-edited
 * into sheetInfo.json; a packaged app has neither, so it is done here instead.
 */
export default function SettingsView({ settings, onChange, onClose }: Props) {
  const [spreadsheetId, setSpreadsheetId] = useState(settings.spreadsheetId);
  const [busy, setBusy] = useState("");
  const [error, setError] = useState("");
  const [saved, setSaved] = useState(false);

  const ready = Boolean(settings.spreadsheetId && settings.hasCredentials && settings.isAuthorized);

  /** run performs one settings action, showing progress and any failure. */
  const run = async (label: string, action: () => Promise<Settings>) => {
    setBusy(label);
    setError("");
    setSaved(false);

    try {
      onChange(await action());
      return true;
    } catch (caught) {
      setError(errorMessage(caught, "Something went wrong. Please try again."));
      return false;
    } finally {
      setBusy("");
    }
  };

  const saveSpreadsheetId = async () => {
    if (await run("spreadsheet", () => SettingsService.SetSpreadsheetID(spreadsheetId))) {
      setSaved(true);
    }
  };

  return (
    <main className="container container-scroll">
      <div className="panel">
        <h1 className="panel-title">Settings</h1>
        <p className="panel-intro">
          TTP Attendance writes each sign-in to a Google Sheet. Complete the three steps below to
          connect it.
        </p>

        <section className="setting">
          <h2>
            <StepMark done={settings.hasCredentials} /> Google API credentials
          </h2>
          <p>
            Create an OAuth client of type <em>Desktop app</em> in the Google Cloud console with the
            Google Sheets API enabled, download its JSON file, and choose it here.
          </p>
          <div className="setting-row">
            <button
              type="button"
              className="btn"
              onClick={() => void run("credentials", SettingsService.ChooseCredentialsFile)}
              disabled={busy !== ""}
            >
              {settings.hasCredentials ? "Replace credentials.json" : "Choose credentials.json"}
            </button>
            <span className="status">
              {settings.hasCredentials ? "Credentials added." : "No credentials yet."}
            </span>
          </div>
        </section>

        <section className="setting">
          <h2>
            <StepMark done={Boolean(settings.spreadsheetId)} /> Attendance spreadsheet
          </h2>
          <p>Paste the spreadsheet's id, or its full URL — the id will be pulled out of it.</p>
          <div className="setting-row">
            <input
              className="form-control"
              type="text"
              value={spreadsheetId}
              placeholder="Spreadsheet id or URL"
              onChange={(event) => {
                setSpreadsheetId(event.target.value);
                setSaved(false);
              }}
            />
            <button
              type="button"
              className="btn"
              onClick={() => void saveSpreadsheetId()}
              disabled={busy !== "" || spreadsheetId.trim() === ""}
            >
              Save
            </button>
          </div>
          {saved && <p className="status">Saved.</p>}
        </section>

        <section className="setting">
          <h2>
            <StepMark done={settings.isAuthorized} /> Google account
          </h2>
          <p>
            Sign in as the account that owns the spreadsheet. This opens your browser, and only asks
            for access to Google Sheets.
          </p>
          <div className="setting-row">
            <button
              type="button"
              className="btn"
              onClick={() => void run("authorize", SettingsService.Authorize)}
              disabled={busy !== "" || !settings.hasCredentials}
            >
              {busy === "authorize" ? "Waiting for your browser…" : "Sign in with Google"}
            </button>
            {settings.isAuthorized && (
              <button
                type="button"
                className="link"
                onClick={() => void run("signout", SettingsService.SignOut)}
                disabled={busy !== ""}
              >
                Sign out
              </button>
            )}
            <span className="status">{settings.isAuthorized ? "Signed in." : "Not signed in."}</span>
          </div>
        </section>

        {error && <p className="error error-banner">{error}</p>}

        <footer className="panel-footer">
          <button type="button" className="btn btn-primary" onClick={onClose} disabled={!ready}>
            {ready ? "Start taking attendance" : "Finish the steps above"}
          </button>
          <button
            type="button"
            className="link"
            onClick={() => void SettingsService.OpenConfigDir().catch(console.error)}
          >
            Show config folder
          </button>
        </footer>

        <p className="panel-path">
          Settings are stored in <code>{settings.configDir}</code>
        </p>
      </div>
    </main>
  );
}

/** StepMark shows whether a setup step is done. */
function StepMark({ done }: { done: boolean }) {
  return (
    <span className={done ? "step-mark step-mark-done" : "step-mark"} aria-hidden="true">
      {done ? "✓" : "•"}
    </span>
  );
}
