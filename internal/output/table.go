package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
	"unicode"
)

// renderTable auto-detects the shape of v and prints a human-readable table.
//   - []interface{} (array of objects) → horizontal table with headers
//   - map[string]interface{} (single object) → vertical KEY | VALUE
//     (special case: single key wrapping an array → auto-flatten to horizontal table)
//   - nil / empty → "(no data)"
//   - scalar → printed as-is
//
// columns: if non-nil, only these columns are shown (in this order).
// termWidth: if > 0, auto-fit columns to terminal width.
func renderTable(w io.Writer, v interface{}, humanTime bool, columns []string, termWidth int) {
	if v == nil {
		_, _ = fmt.Fprintln(w, "(no data)")
		return
	}

	switch data := v.(type) {
	case []interface{}:
		if len(data) == 0 {
			_, _ = fmt.Fprintln(w, "(no data)")
			return
		}
		renderArrayTable(w, data, humanTime, columns, termWidth)
	case map[string]interface{}:
		if len(data) == 0 {
			_, _ = fmt.Fprintln(w, "(no data)")
			return
		}
		renderObjectTable(w, data, humanTime, columns, termWidth)
	default:
		// Scalar value — just print it.
		_, _ = fmt.Fprintln(w, formatCell(v))
	}
}

// renderArrayTable prints an array of objects as a horizontal table.
// columns: if non-nil, only these columns are shown (in this order).
// termWidth: if > 0, auto-fit columns to terminal width.
func renderArrayTable(w io.Writer, items []interface{}, humanTime bool, columns []string, termWidth int) {
	var keys []string
	if len(columns) > 0 {
		keys = columns
	} else {
		keys = collectKeys(items)
	}
	if len(keys) == 0 {
		_, _ = fmt.Fprintln(w, "(no data)")
		return
	}

	// Auto-fit: trim trailing columns that exceed terminal width.
	var hidden int
	if termWidth > 0 && len(keys) > 1 {
		keys, hidden = autoFitColumns(keys, items, humanTime, termWidth)
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	// Header row.
	headers := make([]string, len(keys))
	for i, k := range keys {
		headers[i] = toUpperSnakeCase(k)
	}
	_, _ = fmt.Fprintln(tw, strings.Join(headers, "\t"))

	// Data rows — handle multi-line cells by splitting into sub-rows.
	for _, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			_, _ = fmt.Fprintln(tw, formatCell(item))
			continue
		}
		cells := make([]string, len(keys))
		maxLines := 1
		for i, k := range keys {
			if humanTime && k == "timestamp" {
				cells[i] = formatTimestamp(obj[k])
			} else {
				cells[i] = formatCell(obj[k])
			}
			if n := strings.Count(cells[i], "\n") + 1; n > maxLines {
				maxLines = n
			}
		}
		if maxLines == 1 {
			_, _ = fmt.Fprintln(tw, strings.Join(cells, "\t"))
		} else {
			// Split multi-line cells into sub-rows.
			splitCells := make([][]string, len(cells))
			for i, c := range cells {
				splitCells[i] = strings.Split(c, "\n")
			}
			for row := 0; row < maxLines; row++ {
				parts := make([]string, len(keys))
				for i := range keys {
					if row < len(splitCells[i]) {
						parts[i] = splitCells[i][row]
					}
				}
				_, _ = fmt.Fprintln(tw, strings.Join(parts, "\t"))
			}
			// Blank separator line after multi-line logical row.
			_, _ = fmt.Fprintln(tw)
		}
	}
	_ = tw.Flush()

	if hidden > 0 {
		_, _ = fmt.Fprintf(w, "(%d columns hidden, use --wide to see all)\n", hidden)
	}
}

// renderObjectTable prints a single object as a vertical KEY | VALUE table.
// Auto-flatten rules (checked in code order):
//  1. Single key wrapping a map → recursively render inner map.
//  2. Single key wrapping an array of objects → render array as horizontal table.
//  3. One object-array key + scalar keys (no maps) → render array as table + scalar footer.
//  4. Multiple object-array keys (no maps, no scalar-arrays) → render each array as a named section.
//  5. One map-of-maps key + scalar keys (no arrays) → render inner as NAME table + scalar footer.
//  6. Map of maps (all values are non-empty maps) → render as horizontal table
//     with a NAME column holding the outer key.
//  7. Otherwise → vertical KEY | VALUE table.
func renderObjectTable(w io.Writer, obj map[string]interface{}, humanTime bool, columns []string, termWidth int) {
	// Rule 1: Single key wrapping a map → unwrap (recursive).
	if len(obj) == 1 {
		for _, v := range obj {
			if inner, ok := v.(map[string]interface{}); ok {
				if len(inner) == 0 {
					_, _ = fmt.Fprintln(w, "(no data)")
					return
				}
				renderObjectTable(w, inner, humanTime, columns, termWidth)
				return
			}
		}
	}

	// Classify keys: object-array vs scalar-array vs map vs scalar.
	objectArrayKeys := make([]string, 0)
	scalarArrayKeys := make([]string, 0)
	scalarKeys := make([]string, 0)
	mapKeys := make([]string, 0)
	for k, v := range obj {
		switch val := v.(type) {
		case []interface{}:
			if arrayHasObjects(val) {
				objectArrayKeys = append(objectArrayKeys, k)
			} else {
				scalarArrayKeys = append(scalarArrayKeys, k)
			}
		case map[string]interface{}:
			mapKeys = append(mapKeys, k)
		default:
			scalarKeys = append(scalarKeys, k)
		}
	}

	// Rules 2-4 require no map keys and no scalar-array keys.
	// Scalar arrays can't render as horizontal tables; maps would be lost.
	canTriggerArrayRules := len(mapKeys) == 0 && len(scalarArrayKeys) == 0

	if canTriggerArrayRules {
		switch len(objectArrayKeys) {
		case 1:
			// Rule 2 or 3: one array key → table + optional footer.
			k := objectArrayKeys[0]
			arr := obj[k].([]interface{})
			if len(arr) == 0 {
				_, _ = fmt.Fprintln(w, "(no data)")
				return
			}
			renderArrayTable(w, arr, humanTime, columns, termWidth)
			renderFooter(w, obj, scalarKeys)
			return
		default:
			if len(objectArrayKeys) >= 2 {
				// Rule 4: multiple object-array keys → named sections.
				sort.Strings(objectArrayKeys)
				for i, k := range objectArrayKeys {
					if i > 0 {
						_, _ = fmt.Fprintln(w)
					}
					_, _ = fmt.Fprintf(w, "== %s ==\n", toUpperSnakeCase(k))
					arr := obj[k].([]interface{})
					if len(arr) == 0 {
						_, _ = fmt.Fprintln(w, "(no data)")
					} else {
						renderArrayTable(w, arr, humanTime, columns, termWidth)
					}
				}
				renderFooter(w, obj, scalarKeys)
				return
			}
		}
	}

	// Rule 5: One map-of-maps key + scalar siblings (no arrays) → unwrap inner
	// map as NAME-column table and print scalars as footer.
	if len(mapKeys) == 1 && len(objectArrayKeys) == 0 && len(scalarArrayKeys) == 0 {
		k := mapKeys[0]
		inner := obj[k].(map[string]interface{})
		if isMapOfMaps(inner) {
			renderMapOfMapsTable(w, inner, humanTime, columns, termWidth)
			renderFooter(w, obj, scalarKeys)
			return
		}
	}

	// Rule 6: Map of maps → horizontal table with NAME column.
	if len(mapKeys) >= 2 && len(scalarKeys) == 0 && len(objectArrayKeys) == 0 && len(scalarArrayKeys) == 0 {
		if isMapOfMaps(obj) {
			renderMapOfMapsTable(w, obj, humanTime, columns, termWidth)
			return
		}
	}

	// Rule 7: vertical KEY | VALUE.
	renderObjectTableInner(w, obj, humanTime)
}

// isMapOfMaps reports whether every value in m is a non-empty map.
func isMapOfMaps(m map[string]interface{}) bool {
	if len(m) == 0 {
		return false
	}
	for _, v := range m {
		inner, ok := v.(map[string]interface{})
		if !ok || len(inner) == 0 {
			return false
		}
	}
	return true
}

// renderMapOfMapsTable renders a map-of-maps as a horizontal table keyed by
// a synthetic NAME column holding the outer key. Outer keys are sorted.
func renderMapOfMapsTable(w io.Writer, m map[string]interface{}, humanTime bool, columns []string, termWidth int) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	items := make([]interface{}, 0, len(keys))
	for _, k := range keys {
		inner := m[k].(map[string]interface{})
		row := make(map[string]interface{}, len(inner))
		row["name"] = k
		for ik, iv := range inner {
			if ik != "name" { // Don't overwrite existing "name" field.
				row[ik] = iv
			}
		}
		items = append(items, row)
	}
	renderArrayTable(w, items, humanTime, columns, termWidth)
}

// arrayHasObjects reports whether the first non-nil element of arr is a map.
// Empty arrays are treated as "object arrays" so they render as empty tables.
func arrayHasObjects(arr []interface{}) bool {
	if len(arr) == 0 {
		return true
	}
	for _, v := range arr {
		if v == nil {
			continue
		}
		_, ok := v.(map[string]interface{})
		return ok
	}
	return true // all-nil array → treat as empty
}

// renderFooter prints scalar keys as "Key: Value" lines below a table.
func renderFooter(w io.Writer, obj map[string]interface{}, scalarKeys []string) {
	if len(scalarKeys) == 0 {
		return
	}
	sort.Strings(scalarKeys)
	_, _ = fmt.Fprintln(w)
	for _, k := range scalarKeys {
		_, _ = fmt.Fprintf(w, "%s: %s\n", titleCase(k), formatCell(obj[k]))
	}
}

// renderObjectTableInner renders a map as a vertical KEY | VALUE table
// without applying any auto-flatten rules.
func renderObjectTableInner(w io.Writer, obj map[string]interface{}, humanTime bool) {
	keys := sortedKeys(obj)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, k := range keys {
		val := formatCell(obj[k])
		if humanTime && k == "timestamp" {
			val = formatTimestamp(obj[k])
		}
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", toUpperSnakeCase(k), val)
	}
	_ = tw.Flush()
}

// collectKeys returns the union of all keys across all maps in items, sorted.
func collectKeys(items []interface{}) []string {
	seen := make(map[string]bool)
	for _, item := range items {
		if obj, ok := item.(map[string]interface{}); ok {
			for k := range obj {
				seen[k] = true
			}
		}
	}
	if len(seen) == 0 {
		return nil
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// titleCase capitalizes the first letter of a string (e.g., "total" → "Total").
func titleCase(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// formatTimestamp converts a numeric value to a human-readable local time string.
// Non-numeric values are returned via formatCell unchanged.
func formatTimestamp(v interface{}) string {
	if v == nil {
		return ""
	}
	if f, ok := v.(float64); ok {
		return time.Unix(int64(f), 0).Local().Format("2006-01-02 15:04:05")
	}
	return formatCell(v)
}

// formatCell converts a value to a table cell string.
// Nested objects/arrays are JSON-serialized. Nil becomes empty string.
func formatCell(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		// Sanitize control characters that break tabwriter alignment.
		val = strings.ReplaceAll(val, "\t", " ")
		val = strings.ReplaceAll(val, "\n", "\\n")
		val = strings.ReplaceAll(val, "\r", "\\r")
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case map[string]interface{}:
		if s := flattenNestedCell(val); s != "" {
			return s
		}
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(b)
	case []interface{}:
		if s := flattenSimpleArray(val); s != "" {
			return s
		}
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(b)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// flattenNestedCell handles iKuai nested cell objects like src_addr/dst_addr/src_port/dst_port/time.
// Pattern: {"custom": ["val1", "val2"], "object": [{"gp_name": "grp", ...}]}
// Returns a multi-line string or "" if the map doesn't match this pattern.
func flattenNestedCell(m map[string]interface{}) string {
	custom, hasCustom := m["custom"]
	object, hasObject := m["object"]

	// Must have at least one of custom/object to match the pattern.
	if !hasCustom && !hasObject {
		return ""
	}

	var lines []string

	// Extract custom values (strings, scalars, or nested maps).
	if arr, ok := custom.([]interface{}); ok {
		for _, v := range arr {
			if obj, ok := v.(map[string]interface{}); ok {
				// Nested map: render each key: value on its own line.
				for k, val := range obj {
					s := fmt.Sprintf("%v", val)
					if s != "" {
						lines = append(lines, fmt.Sprintf("%s: %s", k, s))
					}
				}
			} else if s := fmt.Sprintf("%v", v); s != "" {
				lines = append(lines, s)
			}
		}
	}

	// Extract object references (array of objects with gp_name).
	if arr, ok := object.([]interface{}); ok {
		for _, v := range arr {
			if obj, ok := v.(map[string]interface{}); ok {
				if name, ok := obj["gp_name"]; ok {
					lines = append(lines, fmt.Sprintf("[%v]", name))
				}
			}
		}
	}

	if len(lines) == 0 {
		return "*"
	}
	return strings.Join(lines, "\n")
}

// flattenSimpleArray renders an array of single-key objects as a newline-separated list.
// e.g. [{"ip":"1.2.3.4"},{"ip":"5.6.7.8"}] → "1.2.3.4\n5.6.7.8"
// Returns "" if the array doesn't match this pattern.
func flattenSimpleArray(arr []interface{}) string {
	if len(arr) == 0 {
		return ""
	}
	var lines []string
	for _, item := range arr {
		obj, ok := item.(map[string]interface{})
		if !ok {
			return "" // not an object array
		}
		if len(obj) == 0 {
			continue
		}
		// Single-key objects: show "key: value".
		if len(obj) == 1 {
			for k, v := range obj {
				lines = append(lines, fmt.Sprintf("%s: %v", k, v))
			}
			continue
		}
		// Two-key objects with a "comment" key: show "value (comment)" or just "value".
		if len(obj) == 2 {
			if comment, hasComment := obj["comment"]; hasComment {
				for k, v := range obj {
					if k == "comment" {
						continue
					}
					c := fmt.Sprintf("%v", comment)
					if c != "" {
						lines = append(lines, fmt.Sprintf("%v (%s)", v, c))
					} else {
						lines = append(lines, fmt.Sprintf("%v", v))
					}
				}
				continue
			}
		}
		// Multi-key objects: render each key:value on its own line.
		for k, v := range obj {
			s := fmt.Sprintf("%v", v)
			if s != "" {
				lines = append(lines, fmt.Sprintf("%s: %s", k, s))
			}
		}
	}
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n")
}

// autoFitColumns trims trailing columns so the table fits within termWidth.
// Returns the subset of keys that fit and the number of hidden columns.
func autoFitColumns(keys []string, items []interface{}, humanTime bool, termWidth int) ([]string, int) {
	// Calculate max width per column (header vs cell values).
	widths := make([]int, len(keys))
	for i, k := range keys {
		widths[i] = len(toUpperSnakeCase(k))
	}
	for _, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		for i, k := range keys {
			var cell string
			if humanTime && k == "timestamp" {
				cell = formatTimestamp(obj[k])
			} else {
				cell = formatCell(obj[k])
			}
			// For multi-line cells, use the longest line's width (not total string length).
			for _, line := range strings.Split(cell, "\n") {
				if len(line) > widths[i] {
					widths[i] = len(line)
				}
			}
		}
	}

	// Progressively add columns: each needs its width + 2 (tabwriter padding).
	total := 0
	fit := 0
	for i, w := range widths {
		need := w
		if i > 0 {
			need += 2 // tabwriter inter-column padding
		}
		if total+need > termWidth {
			break
		}
		total += need
		fit++
	}
	if fit < 1 {
		fit = 1 // always show at least one column
	}
	if fit >= len(keys) {
		return keys, 0
	}
	return keys[:fit], len(keys) - fit
}

// toUpperSnakeCase converts a JSON key to UPPER_SNAKE_CASE.
//   - snake_case: "ip_addr" → "IP_ADDR"
//   - camelCase:  "lanIp"   → "LAN_IP"
//   - mixed:      "srcIP"   → "SRC_IP"
func toUpperSnakeCase(s string) string {
	var b strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) && i > 0 {
			prev := runes[i-1]
			// Insert underscore before uppercase if preceded by lowercase,
			// or if preceded by uppercase followed by lowercase (e.g., "IP" in "srcIP").
			if unicode.IsLower(prev) {
				b.WriteRune('_')
			} else if unicode.IsUpper(prev) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				b.WriteRune('_')
			}
		}
		b.WriteRune(unicode.ToUpper(r))
	}
	return b.String()
}
