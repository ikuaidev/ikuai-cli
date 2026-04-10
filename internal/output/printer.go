package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Format controls the output mode.
type Format int

const (
	Table Format = iota
	JSON
	YAML
)

// FormatFromString parses a format name. Returns an error for unknown values.
func FormatFromString(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "table":
		return Table, nil
	case "json":
		return JSON, nil
	case "yaml":
		return YAML, nil
	default:
		return Table, fmt.Errorf("unknown format %q: must be one of: table, json, yaml", s)
	}
}

// String returns the format name.
func (f Format) String() string {
	switch f {
	case JSON:
		return "json"
	case YAML:
		return "yaml"
	default:
		return "table"
	}
}

// Printer routes output to the selected format renderer.
type Printer struct {
	stdout    io.Writer
	stderr    io.Writer
	format    Format
	HumanTime bool     // When true, converts "timestamp" columns to human-readable time across all formats.
	Columns   []string // When set, only these columns are shown in table format (in this order).
	Wide      bool     // When true, show all columns regardless of Columns setting.
	TermWidth int      // Terminal width in columns; 0 means no limit (no auto-fit).
}

// New creates a Printer with the given format.
func New(stdout, stderr io.Writer, format Format) *Printer {
	return &Printer{
		stdout: stdout,
		stderr: stderr,
		format: format,
	}
}

// Print renders raw JSON bytes according to the configured format.
func (p *Printer) Print(raw json.RawMessage) {
	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		// Not valid JSON — print as-is.
		_, _ = fmt.Fprintln(p.stdout, string(raw))
		return
	}
	p.PrintValue(v)
}

// PrintValue renders an already-decoded Go value according to the configured format.
// For table mode, Go typed maps (e.g. map[string]string) are normalized via JSON
// round-trip so renderTable receives map[string]interface{}.
func (p *Printer) PrintValue(v interface{}) {
	// --human-time: convert timestamp fields to human-readable time across all formats.
	if p.HumanTime {
		v = transformTimestamps(v)
	}

	switch p.format {
	case JSON:
		p.printJSON(v)
	case YAML:
		// Normalize integer-valued float64 to int64 to prevent scientific notation.
		p.printYAML(normalizeIntegers(v))
	default:
		// Normalize typed Go values to interface{} via JSON round-trip.
		b, err := json.Marshal(v)
		if err != nil {
			_, _ = fmt.Fprintln(p.stderr, "json error:", err)
			return
		}
		var normalized interface{}
		if err := json.Unmarshal(b, &normalized); err != nil {
			_, _ = fmt.Fprintln(p.stdout, string(b))
			return
		}
		p.printTable(normalized)
	}
}

// PrintPrettyJSON outputs indented JSON (used for --raw mode).
func (p *Printer) PrintPrettyJSON(raw json.RawMessage) {
	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		_, _ = fmt.Fprintln(p.stdout, string(raw))
		return
	}
	enc := json.NewEncoder(p.stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		_, _ = fmt.Fprintln(p.stderr, "json error:", err)
	}
}

func (p *Printer) printJSON(v interface{}) {
	enc := json.NewEncoder(p.stdout)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		_, _ = fmt.Fprintln(p.stderr, "json error:", err)
	}
}

func (p *Printer) printYAML(v interface{}) {
	enc := yaml.NewEncoder(p.stdout)
	enc.SetIndent(2)
	if err := enc.Encode(v); err != nil {
		_, _ = fmt.Fprintln(p.stderr, "yaml error:", err)
	}
	_ = enc.Close()
}

func (p *Printer) printTable(v interface{}) {
	var cols []string
	tw := 0
	if p.Wide {
		// --wide: show all columns, no auto-fit.
	} else {
		cols = p.Columns
		tw = p.TermWidth
	}
	renderTable(p.stdout, v, p.HumanTime, cols, tw)
}

// transformTimestamps recursively walks the data and replaces "timestamp"
// float64 values with human-readable time strings. Used when HumanTime is true
// for json/yaml output.
func transformTimestamps(v interface{}) interface{} {
	switch data := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(data))
		for k, val := range data {
			if k == "timestamp" {
				if f, ok := val.(float64); ok {
					out[k] = time.Unix(int64(f), 0).Local().Format("2006-01-02 15:04:05")
					continue
				}
			}
			out[k] = transformTimestamps(val)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(data))
		for i, item := range data {
			out[i] = transformTimestamps(item)
		}
		return out
	default:
		return v
	}
}

// normalizeIntegers recursively converts integer-valued float64 to int64
// to prevent YAML scientific notation (e.g. 1.7756e+09 → 1775600000).
func normalizeIntegers(v interface{}) interface{} {
	switch data := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(data))
		for k, val := range data {
			out[k] = normalizeIntegers(val)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(data))
		for i, item := range data {
			out[i] = normalizeIntegers(item)
		}
		return out
	case float64:
		if data == float64(int64(data)) {
			return int64(data)
		}
		return data
	default:
		return v
	}
}
