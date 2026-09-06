package sheets

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	sheetsapi "google.golang.org/api/sheets/v4"

	"ttpattendance/internal/config"
)

// ErrNotAuthorized is returned when nobody has signed in with a Google account yet,
// or the saved token can no longer be refreshed.
var ErrNotAuthorized = errors.New("not signed in to Google: authorise the app in Settings")

// oauthConfig builds the OAuth client configuration from the stored credentials.json.
func oauthConfig() (*oauth2.Config, error) {
	creds, err := config.ReadCredentials()
	if err != nil {
		return nil, err
	}

	conf, err := google.ConfigFromJSON(creds, sheetsapi.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("could not read credentials.json: %w", err)
	}

	return conf, nil
}

// savedToken returns the token persisted by a previous authorisation.
func savedToken() (*oauth2.Token, error) {
	data, err := config.ReadToken()
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrNotAuthorized
	}

	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, ErrNotAuthorized
	}
	if token.RefreshToken == "" && !token.Valid() {
		return nil, ErrNotAuthorized
	}

	return &token, nil
}

func saveToken(token *oauth2.Token) error {
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}

	return config.WriteToken(append(data, '\n'))
}

// persistingTokenSource writes the token back to disk whenever the underlying
// source refreshes it, so the app keeps working across restarts without asking
// anyone to sign in again.
type persistingTokenSource struct {
	source  oauth2.TokenSource
	current *oauth2.Token
}

func (p *persistingTokenSource) Token() (*oauth2.Token, error) {
	token, err := p.source.Token()
	if err != nil {
		// A revoked or expired refresh token cannot be recovered from here;
		// surface it as "not authorised" so the UI points at Settings.
		return nil, fmt.Errorf("%w (%v)", ErrNotAuthorized, err)
	}

	if p.current == nil || token.AccessToken != p.current.AccessToken {
		p.current = token
		if err := saveToken(token); err != nil {
			return nil, err
		}
	}

	return token, nil
}

// httpClient returns an HTTP client that authenticates as the signed-in account.
func httpClient(ctx context.Context) (*http.Client, error) {
	conf, err := oauthConfig()
	if err != nil {
		return nil, err
	}

	token, err := savedToken()
	if err != nil {
		return nil, err
	}

	source := &persistingTokenSource{source: conf.TokenSource(ctx, token), current: token}

	return oauth2.NewClient(ctx, oauth2.ReuseTokenSource(nil, source)), nil
}

// IsAuthorized reports whether a usable token is stored.
func IsAuthorized() bool {
	_, err := savedToken()
	return err == nil
}

// Authorize runs the installed-application OAuth flow: it starts a loopback
// listener, hands openURL the consent page to open in the default browser and
// waits for Google to redirect back with an authorisation code.
//
// This is the Go counterpart of the Python flow.run_local_server(port=0) call.
func Authorize(ctx context.Context, openURL func(string) error) error {
	conf, err := oauthConfig()
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("could not start the local sign-in listener: %w", err)
	}
	defer listener.Close()

	// Desktop OAuth clients may redirect to any loopback port, so point Google
	// at the one we just claimed.
	conf.RedirectURL = fmt.Sprintf("http://127.0.0.1:%d/", listener.Addr().(*net.TCPAddr).Port)

	state, err := randomState()
	if err != nil {
		return err
	}

	type result struct {
		code string
		err  error
	}
	results := make(chan result, 1)

	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if errParam := query.Get("error"); errParam != "" {
			writeBrowserPage(w, "Sign-in cancelled", "You can close this tab and try again from the app.")
			results <- result{err: fmt.Errorf("google returned an error: %s", errParam)}
			return
		}

		if query.Get("state") != state {
			writeBrowserPage(w, "Sign-in failed", "The sign-in response did not match this request. Please close this tab and try again.")
			results <- result{err: errors.New("the sign-in response did not match this request")}
			return
		}

		code := query.Get("code")
		if code == "" {
			return // Not the redirect we are waiting for (a favicon request, say).
		}

		writeBrowserPage(w, "You're signed in", "You can close this tab and return to TTP Attendance.")
		results <- result{code: code}
	})}
	defer server.Close()

	go server.Serve(listener)

	if err := openURL(conf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)); err != nil {
		return fmt.Errorf("could not open the browser to sign in: %w", err)
	}

	// Give whoever is at the machine a few minutes to work through Google's consent screens.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	select {
	case res := <-results:
		if res.err != nil {
			return res.err
		}

		token, err := conf.Exchange(ctx, res.code)
		if err != nil {
			return fmt.Errorf("could not complete sign-in: %w", err)
		}

		return saveToken(token)

	case <-ctx.Done():
		return errors.New("timed out waiting for Google sign-in")
	}
}

func randomState() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("could not start sign-in: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// writeBrowserPage renders the page Google's redirect lands on.
func writeBrowserPage(w http.ResponseWriter, title, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!doctype html>
<html lang="en">
<head><meta charset="utf-8"><title>%[1]s</title></head>
<body style="margin:0;display:grid;place-items:center;height:100vh;background:#00274d;color:#fff;font-family:system-ui,sans-serif;text-align:center">
  <div><h1 style="color:#ffb81c">%[1]s</h1><p>%[2]s</p></div>
</body>
</html>`, title, message)
}
