package cliapp

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/ikuaidev/ikuai-cli/internal/api"
	"github.com/ikuaidev/ikuai-cli/internal/output"
	"github.com/ikuaidev/ikuai-cli/internal/session"
)

type Runtime struct {
	Session        *session.Session
	APIClient      *api.Client
	CredSource     string // "token", "env", or "none"
	Format         output.Format
	RawMode        bool
	HumanTime      bool
	DryRun         bool
	DefaultColumns []string // Per-command default columns for table output.
	UserColumns    []string // User-specified columns via --columns flag.
	WideMode       bool     // Show all columns (--wide).
	TermWidth      int      // Terminal width; 0 = unknown.
	Stdout         io.Writer
	Stderr         io.Writer
}

func New(stdout, stderr io.Writer) *Runtime {
	return &Runtime{
		Stdout: stdout,
		Stderr: stderr,
	}
}

func (r *Runtime) SyncSession() error {
	s, err := session.Load()
	if err != nil {
		return err
	}
	if s == nil {
		s = &session.Session{}
	}

	// Determine credential source: session file > env vars.
	if s.Token != "" {
		r.CredSource = "token"
	} else if envURL, envToken := envCred(); envURL != "" && envToken != "" {
		// Env fallback: populate in-memory only — never call save().
		s.BaseURL = envURL
		s.Token = envToken
		r.CredSource = "env"
	} else {
		r.CredSource = "none"
	}

	r.Session = s
	if s.BaseURL != "" {
		r.APIClient = api.New(s.BaseURL, s.Token)
		r.APIClient.RawMode = r.RawMode
		r.APIClient.DryRun = r.DryRun
	} else {
		r.APIClient = nil
	}
	return nil
}

// envCred reads IKUAI_CLI_BASE_URL and IKUAI_CLI_TOKEN from the environment.
// Both must be non-empty after trimming for the env path to activate.
func envCred() (baseURL, token string) {
	baseURL = strings.TrimSpace(os.Getenv("IKUAI_CLI_BASE_URL"))
	token = strings.TrimSpace(os.Getenv("IKUAI_CLI_TOKEN"))
	baseURL = strings.TrimRight(baseURL, "/")
	return
}

func (r *Runtime) NewClient(baseURL, token string) *api.Client {
	return api.New(baseURL, token)
}

func (r *Runtime) LoadSession() (*session.Session, error) {
	return session.Load()
}

func (r *Runtime) RequireURL() error {
	if r.Session == nil || r.Session.BaseURL == "" {
		return &AuthError{Message: "no URL configured. Run: ikuai-cli auth set-url https://192.168.1.1"}
	}
	if r.APIClient == nil {
		r.APIClient = api.New(r.Session.BaseURL, r.Session.Token)
		r.APIClient.RawMode = r.RawMode
		r.APIClient.DryRun = r.DryRun
	}
	return nil
}

func (r *Runtime) RequireAuth() error {
	if err := r.RequireURL(); err != nil {
		return err
	}
	if r.Session.Token == "" {
		return &AuthError{Message: "not authenticated. Run: ikuai-cli auth set-token <TOKEN>"}
	}
	return nil
}

// newPrinter creates a Printer with the current Runtime config.
func (r *Runtime) newPrinter() *output.Printer {
	p := output.New(r.Stdout, r.Stderr, r.Format)
	p.HumanTime = r.HumanTime

	// Column config: --columns > DefaultColumns; --wide overrides all.
	if len(r.UserColumns) > 0 {
		p.Columns = r.UserColumns
	} else if len(r.DefaultColumns) > 0 {
		p.Columns = r.DefaultColumns
	}
	p.Wide = r.WideMode
	p.TermWidth = r.TermWidth
	return p
}

// PrintJSON renders a Go value (map, struct) using the configured format.
func (r *Runtime) PrintJSON(v interface{}) {
	r.newPrinter().PrintValue(v)
}

// PrintRaw renders raw API JSON bytes using the configured format.
// In --raw mode, outputs pretty-printed full envelope.
func (r *Runtime) PrintRaw(raw json.RawMessage) {
	p := r.newPrinter()
	if r.RawMode {
		p.PrintPrettyJSON(raw)
		return
	}
	p.Print(raw)
}
