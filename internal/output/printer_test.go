package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// --- Format parsing ---

func TestFormatFromString(t *testing.T) {
	tests := []struct {
		input string
		want  Format
		err   bool
	}{
		{"table", Table, false},
		{"json", JSON, false},
		{"yaml", YAML, false},
		{"TABLE", Table, false},
		{"JSON", JSON, false},
		{"invalid", Table, true},
		{"", Table, true},
	}
	for _, tt := range tests {
		got, err := FormatFromString(tt.input)
		if (err != nil) != tt.err {
			t.Errorf("FormatFromString(%q) error = %v, wantErr %v", tt.input, err, tt.err)
		}
		if got != tt.want {
			t.Errorf("FormatFromString(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

// --- JSON output ---

func TestPrintJSON(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, JSON)
	p.Print(json.RawMessage(`{"message":"ok","code":0}`))

	got := strings.TrimSpace(out.String())
	if got != `{"code":0,"message":"ok"}` {
		t.Fatalf("JSON Print() = %q, want compact JSON", got)
	}
}

func TestPrintValueJSON(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, JSON)
	p.PrintValue(map[string]interface{}{"message": "ok"})

	got := strings.TrimSpace(out.String())
	if got != `{"message":"ok"}` {
		t.Fatalf("JSON PrintValue() = %q", got)
	}
}

// --- YAML output ---

func TestPrintYAML(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, YAML)
	p.Print(json.RawMessage(`{"name":"test","value":42}`))

	got := out.String()
	if !strings.Contains(got, "name: test") || !strings.Contains(got, "value: 42") {
		t.Fatalf("YAML Print() = %q, want YAML with name and value", got)
	}
}

// --- Table output: array ---

func TestPrintTableArray(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`[{"ip":"1.1.1.1","name":"gw"},{"ip":"2.2.2.2","name":"dns"}]`))

	got := out.String()
	if !strings.Contains(got, "IP") || !strings.Contains(got, "NAME") {
		t.Fatalf("table array headers missing: %q", got)
	}
	if !strings.Contains(got, "1.1.1.1") || !strings.Contains(got, "dns") {
		t.Fatalf("table array data missing: %q", got)
	}
}

// --- Table output: single object ---

func TestPrintTableObject(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"host":"192.168.1.1","status":"ok"}`))

	got := out.String()
	if !strings.Contains(got, "HOST") || !strings.Contains(got, "192.168.1.1") {
		t.Fatalf("table object missing host: %q", got)
	}
	if !strings.Contains(got, "STATUS") || !strings.Contains(got, "ok") {
		t.Fatalf("table object missing status: %q", got)
	}
}

// --- Table output: null / empty ---

func TestPrintTableNull(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`null`))

	got := strings.TrimSpace(out.String())
	if got != "(no data)" {
		t.Fatalf("table null = %q, want (no data)", got)
	}
}

func TestPrintTableEmptyArray(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`[]`))

	got := strings.TrimSpace(out.String())
	if got != "(no data)" {
		t.Fatalf("table [] = %q, want (no data)", got)
	}
}

func TestPrintTableEmptyObject(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{}`))

	got := strings.TrimSpace(out.String())
	if got != "(no data)" {
		t.Fatalf("table {} = %q, want (no data)", got)
	}
}

// --- Table: nested values ---

func TestPrintTableNestedValue(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`[{"name":"test","tags":["a","b"]}]`))

	got := out.String()
	if !strings.Contains(got, `["a","b"]`) {
		t.Fatalf("nested array not JSON-serialized: %q", got)
	}
}

// --- Table: scalar ---

func TestPrintTableScalar(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`"hello"`))

	got := strings.TrimSpace(out.String())
	if got != "hello" {
		t.Fatalf("table scalar = %q, want hello", got)
	}
}

// --- Column name conversion ---

func TestToUpperSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"ip_addr", "IP_ADDR"},
		{"lanIp", "LAN_IP"},
		{"srcIP", "SRC_IP"},
		{"name", "NAME"},
		{"HTTPStatus", "HTTP_STATUS"},
		{"id", "ID"},
		{"pageSize", "PAGE_SIZE"},
	}
	for _, tt := range tests {
		got := toUpperSnakeCase(tt.input)
		if got != tt.want {
			t.Errorf("toUpperSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// --- Pretty JSON (for --raw mode) ---

func TestPrintPrettyJSON(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table) // format doesn't matter for PrintPrettyJSON
	p.PrintPrettyJSON(json.RawMessage(`{"code":0,"data":{"uptime":123}}`))

	got := out.String()
	if !strings.Contains(got, "  ") {
		t.Fatalf("PrintPrettyJSON not indented: %q", got)
	}
	if !strings.Contains(got, `"code": 0`) {
		t.Fatalf("PrintPrettyJSON missing code field: %q", got)
	}
}

// --- JSON null/empty via json/yaml ---

func TestPrintJSONNull(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, JSON)
	p.Print(json.RawMessage(`null`))

	got := strings.TrimSpace(out.String())
	if got != "null" {
		t.Fatalf("JSON null = %q, want null", got)
	}
}

func TestPrintInvalidJSONFallsThrough(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`not valid json`))

	got := strings.TrimSpace(out.String())
	if got != "not valid json" {
		t.Fatalf("Print(invalid) = %q, want raw passthrough", got)
	}
}

func TestPrintPrettyJSONInvalidFallsThrough(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.PrintPrettyJSON(json.RawMessage(`not valid json`))

	got := strings.TrimSpace(out.String())
	if got != "not valid json" {
		t.Fatalf("PrintPrettyJSON(invalid) = %q, want raw passthrough", got)
	}
}

func TestPrintValueTableWithTypedMap(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	// map[string]string requires JSON round-trip normalization.
	p.PrintValue(map[string]string{"host": "192.168.1.1", "status": "ok"})

	got := out.String()
	if !strings.Contains(got, "HOST") || !strings.Contains(got, "192.168.1.1") {
		t.Fatalf("table from typed map missing data: %q", got)
	}
}

func TestPrintTableHeterogeneousArray(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	// Second object has an extra key that should show as a column.
	p.Print(json.RawMessage(`[{"name":"wan1","type":"dhcp"},{"name":"wan2","type":"pppoe","user":"admin"}]`))

	got := out.String()
	if !strings.Contains(got, "USER") {
		t.Fatalf("table missing column from second object: %q", got)
	}
	if !strings.Contains(got, "admin") {
		t.Fatalf("table missing value from second object: %q", got)
	}
}

func TestPrintTableSanitizesControlChars(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`[{"name":"has\ttab","desc":"line1\nline2"}]`))

	got := out.String()
	// Tab in value should be replaced with space, not break column alignment.
	if strings.Contains(got, "has\ttab") {
		t.Fatalf("raw tab not sanitized in output: %q", got)
	}
	if !strings.Contains(got, "has tab") {
		t.Fatalf("tab should be replaced with space: %q", got)
	}
	// Newline in value should be escaped as literal \n.
	if strings.Count(got, "\n") > 3 {
		// Header + 1 data row + trailing = 3 newlines max; more means raw newline leaked
		t.Fatalf("raw newline not sanitized — too many line breaks: %q", got)
	}
}

// --- Auto-flatten: single key wrapping array ---

func TestPrintTableAutoFlattenSingleKeyArray(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"cpu": [{"cpu": 0.5}, {"cpu": 0.3}]}`))

	got := out.String()
	if !strings.Contains(got, "CPU") {
		t.Fatalf("auto-flatten missing header: %q", got)
	}
	if !strings.Contains(got, "0.5") || !strings.Contains(got, "0.3") {
		t.Fatalf("auto-flatten missing data rows: %q", got)
	}
}

func TestPrintTableAutoFlattenMultiColumn(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"disk_space_used": [{"partition":"/","used":"3G"},{"partition":"/boot","used":"200M"}]}`))

	got := out.String()
	if !strings.Contains(got, "PARTITION") || !strings.Contains(got, "USED") {
		t.Fatalf("auto-flatten missing column headers: %q", got)
	}
	if !strings.Contains(got, "/boot") || !strings.Contains(got, "200M") {
		t.Fatalf("auto-flatten missing data: %q", got)
	}
}

func TestPrintTableAutoFlattenEmptyArray(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"items": []}`))

	got := strings.TrimSpace(out.String())
	if got != "(no data)" {
		t.Fatalf("auto-flatten empty array = %q, want (no data)", got)
	}
}

func TestPrintTableAutoFlattenNullValue(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"items": null}`))

	got := out.String()
	// null value is not []interface{}, so falls through to vertical KV table.
	if !strings.Contains(got, "ITEMS") {
		t.Fatalf("null value should render as KV table: %q", got)
	}
}

func TestPrintTableMultiKeyNotFlattened(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"host":"1.1.1.1","status":"ok"}`))

	got := out.String()
	// Multi-key object stays as vertical KV table.
	if !strings.Contains(got, "HOST") || !strings.Contains(got, "STATUS") {
		t.Fatalf("multi-key should stay as KV table: %q", got)
	}
}

func TestPrintTableSingleKeyNonArrayNotFlattened(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"message": "ok"}`))

	got := out.String()
	// Single key but value is a string, not an array — vertical KV table.
	if !strings.Contains(got, "MESSAGE") || !strings.Contains(got, "ok") {
		t.Fatalf("single key non-array should stay as KV table: %q", got)
	}
}

// --- Auto-flatten: array key + scalar keys (list response pattern) ---

func TestPrintTableAutoFlattenArrayPlusScalar(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"data": [{"id":1,"name":"foo"},{"id":2,"name":"bar"}], "total": 2}`))

	got := out.String()
	if !strings.Contains(got, "ID") || !strings.Contains(got, "NAME") {
		t.Fatalf("array+scalar flatten missing table headers: %q", got)
	}
	if !strings.Contains(got, "foo") || !strings.Contains(got, "bar") {
		t.Fatalf("array+scalar flatten missing data rows: %q", got)
	}
	if !strings.Contains(got, "Total: 2") {
		t.Fatalf("array+scalar flatten missing footer: %q", got)
	}
}

func TestPrintTableAutoFlattenArrayPlusMultipleScalars(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"data": [{"id":1}], "total": 5, "page": 1}`))

	got := out.String()
	if !strings.Contains(got, "ID") {
		t.Fatalf("missing table header: %q", got)
	}
	if !strings.Contains(got, "Page: 1") || !strings.Contains(got, "Total: 5") {
		t.Fatalf("missing footer lines: %q", got)
	}
}

func TestPrintTableAutoFlattenArrayPlusScalarEmpty(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"data": [], "total": 0}`))

	got := strings.TrimSpace(out.String())
	if got != "(no data)" {
		t.Fatalf("empty array+scalar = %q, want (no data)", got)
	}
}

func TestPrintTableAutoFlattenSingleKeyStillWorks(t *testing.T) {
	// Regression: original single-key auto-flatten must still work.
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"cpu": [{"val": 0.5}]}`))

	got := out.String()
	if !strings.Contains(got, "VAL") || !strings.Contains(got, "0.5") {
		t.Fatalf("single-key regression: %q", got)
	}
	// No footer line should appear.
	if strings.Contains(got, ":") && strings.Contains(got, "Cpu") {
		t.Fatalf("single-key should not have footer: %q", got)
	}
}

func TestPrintTableTwoArrayKeysRenderedAsSections(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"snap_lan": [{"id":1,"name":"lan1"}], "snap_wan": [{"id":2,"name":"wan1"}]}`))

	got := out.String()
	// Multiple array keys render as named sections.
	if !strings.Contains(got, "== SNAP_LAN ==") || !strings.Contains(got, "== SNAP_WAN ==") {
		t.Fatalf("missing section headers: %q", got)
	}
	if !strings.Contains(got, "lan1") || !strings.Contains(got, "wan1") {
		t.Fatalf("missing section data: %q", got)
	}
}

func TestPrintTableSingleKeyWrappingMap(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"sysinfo": {"hostname":"router","uptime":12345}}`))

	got := out.String()
	// Single key wrapping a map should unwrap to show inner fields.
	if !strings.Contains(got, "HOSTNAME") || !strings.Contains(got, "router") {
		t.Fatalf("inner map not unwrapped: %q", got)
	}
	if !strings.Contains(got, "UPTIME") || !strings.Contains(got, "12345") {
		t.Fatalf("inner fields missing: %q", got)
	}
	// The outer key "sysinfo" should NOT appear as a row header.
	if strings.Contains(got, "SYSINFO") {
		t.Fatalf("outer single key should be unwrapped, not shown: %q", got)
	}
}

func TestPrintTableSingleKeyWrappingEmptyMap(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`{"data": {}}`))

	got := strings.TrimSpace(out.String())
	if got != "(no data)" {
		t.Fatalf("empty inner map = %q, want (no data)", got)
	}
}

// --- Rule 5: One map-of-maps + scalar siblings ---

func TestPrintTableMapOfMapsWithScalarFooter(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	// Pattern from `monitor flow-shunting`: scalar + single map-of-maps.
	p.Print(json.RawMessage(`{"clean_time": 0, "data": {"today": {"conn_cnt": 0}, "week": {"conn_cnt": 1550}, "yesterday": {"conn_cnt": 1550}}}`))

	got := out.String()
	if !strings.Contains(got, "NAME") || !strings.Contains(got, "CONN_CNT") {
		t.Fatalf("missing table headers: %q", got)
	}
	if !strings.Contains(got, "today") || !strings.Contains(got, "week") || !strings.Contains(got, "yesterday") {
		t.Fatalf("missing row names: %q", got)
	}
	if !strings.Contains(got, "1550") {
		t.Fatalf("missing row data: %q", got)
	}
	if !strings.Contains(got, "Clean_time: 0") {
		t.Fatalf("missing scalar footer: %q", got)
	}
}

// --- Rule 6: Map of maps (all values are non-empty maps) ---

func TestPrintTableMapOfMaps(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	// Pattern from `monitor wireless-stats`: multiple map keys, all maps.
	p.Print(json.RawMessage(`{"ap_status": {"ap_count": 0, "ap_online": 0}, "clt_status": {"clt_count": 0, "clt_active": 0}}`))

	got := out.String()
	if !strings.Contains(got, "NAME") {
		t.Fatalf("missing NAME column: %q", got)
	}
	if !strings.Contains(got, "ap_status") || !strings.Contains(got, "clt_status") {
		t.Fatalf("missing row names: %q", got)
	}
	if !strings.Contains(got, "AP_COUNT") || !strings.Contains(got, "CLT_COUNT") {
		t.Fatalf("missing inner headers: %q", got)
	}
}

func TestPrintTableMapOfMapsUsesConfiguredColumns(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Columns = []string{"name", "interface", "link"}
	p.Print(json.RawMessage(`{"ether_info":{"eth0":{"interface":"lan1","link":1,"driver":"e1000e"},"eth1":{"interface":"wan1","link":0,"driver":"igb"}}}`))

	got := out.String()
	if !strings.Contains(got, "NAME") || !strings.Contains(got, "INTERFACE") || !strings.Contains(got, "LINK") {
		t.Fatalf("missing configured headers: %q", got)
	}
	if strings.Contains(got, "DRIVER") || strings.Contains(got, "e1000e") {
		t.Fatalf("map-of-maps should respect configured columns: %q", got)
	}
	if !strings.Contains(got, "eth0") || !strings.Contains(got, "lan1") {
		t.Fatalf("missing configured row values: %q", got)
	}
}

func TestPrintTableMapOfMapsOnlyWhenAllInnerNonEmpty(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	// One inner map is empty — should fall through to vertical KV instead.
	p.Print(json.RawMessage(`{"a": {"x": 1}, "b": {}}`))

	got := out.String()
	// Should render as KV table, not map-of-maps.
	if strings.Contains(got, "NAME") {
		t.Fatalf("should not render as map-of-maps when one inner is empty: %q", got)
	}
}

// --- Human time: timestamp conversion ---

func TestPrintTableHumanTimeEnabled(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.HumanTime = true
	p.Print(json.RawMessage(`[{"name":"disk","timestamp":1775626740}]`))

	got := out.String()
	want := time.Unix(1775626740, 0).Local().Format("2006-01-02 15:04:05")
	if !strings.Contains(got, want) {
		t.Fatalf("HumanTime enabled: got %q, want timestamp converted to %q", got, want)
	}
}

func TestPrintTableHumanTimeDisabled(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	// HumanTime defaults to false.
	p.Print(json.RawMessage(`[{"name":"disk","timestamp":1775626740}]`))

	got := out.String()
	if !strings.Contains(got, "1775626740") {
		t.Fatalf("HumanTime disabled: timestamp should be raw number: %q", got)
	}
}

func TestPrintTableHumanTimeNonNumeric(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.HumanTime = true
	p.Print(json.RawMessage(`[{"timestamp":"not-a-number"}]`))

	got := out.String()
	if !strings.Contains(got, "not-a-number") {
		t.Fatalf("HumanTime non-numeric: should pass through unchanged: %q", got)
	}
}

func TestPrintTableHumanTimeEpochZero(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.HumanTime = true
	p.Print(json.RawMessage(`[{"timestamp":0}]`))

	got := out.String()
	want := time.Unix(0, 0).Local().Format("2006-01-02 15:04:05")
	if !strings.Contains(got, want) {
		t.Fatalf("HumanTime epoch 0: got %q, want %q", got, want)
	}
}

func TestPrintHumanTimeAppliesToJSON(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, JSON)
	p.HumanTime = true
	p.Print(json.RawMessage(`[{"timestamp":1775626740}]`))

	got := out.String()
	want := time.Unix(1775626740, 0).Local().Format("2006-01-02 15:04:05")
	if !strings.Contains(got, want) {
		t.Fatalf("JSON with HumanTime: got %q, want timestamp converted to %q", got, want)
	}
}

func TestPrintHumanTimeAppliesToYAML(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, YAML)
	p.HumanTime = true
	p.Print(json.RawMessage(`[{"timestamp":1775626740}]`))

	got := out.String()
	want := time.Unix(1775626740, 0).Local().Format("2006-01-02 15:04:05")
	if !strings.Contains(got, want) {
		t.Fatalf("YAML with HumanTime: got %q, want timestamp converted to %q", got, want)
	}
}

func TestPrintYAMLIntegersNotScientific(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, YAML)
	p.Print(json.RawMessage(`[{"timestamp":1775626740}]`))

	got := out.String()
	if strings.Contains(got, "e+") || strings.Contains(got, "E+") {
		t.Fatalf("YAML should not use scientific notation: %q", got)
	}
	if !strings.Contains(got, "1775626740") {
		t.Fatalf("YAML should contain integer value: %q", got)
	}
}

func TestPrintYAMLNull(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, YAML)
	p.Print(json.RawMessage(`null`))

	got := strings.TrimSpace(out.String())
	if got != "null" {
		t.Fatalf("YAML null = %q, want null", got)
	}
}

// --- Column filtering ---

func TestColumnFilterShowsOnlySpecified(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Columns = []string{"id", "name"}
	p.Print(json.RawMessage(`[{"id":1,"name":"foo","extra":"bar","z":"hidden"}]`))

	got := out.String()
	if !strings.Contains(got, "ID") || !strings.Contains(got, "NAME") {
		t.Fatalf("filtered table missing expected columns: %q", got)
	}
	if strings.Contains(got, "EXTRA") || strings.Contains(got, "hidden") {
		t.Fatalf("filtered table should not contain extra columns: %q", got)
	}
}

func TestColumnFilterPreservesOrder(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Columns = []string{"name", "id"}
	p.Print(json.RawMessage(`[{"id":1,"name":"foo"}]`))

	got := out.String()
	nameIdx := strings.Index(got, "NAME")
	idIdx := strings.Index(got, "ID")
	if nameIdx < 0 || idIdx < 0 {
		t.Fatalf("missing columns: %q", got)
	}
	if nameIdx > idIdx {
		t.Fatalf("column order should follow Columns slice, NAME before ID: %q", got)
	}
}

func TestColumnFilterWideOverride(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Columns = []string{"id"}
	p.Wide = true
	p.Print(json.RawMessage(`[{"id":1,"name":"foo","extra":"bar"}]`))

	got := out.String()
	if !strings.Contains(got, "EXTRA") || !strings.Contains(got, "NAME") {
		t.Fatalf("Wide=true should show all columns: %q", got)
	}
}

func TestColumnFilterNilShowsAll(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	// Columns is nil by default.
	p.Print(json.RawMessage(`[{"id":1,"name":"foo","extra":"bar"}]`))

	got := out.String()
	if !strings.Contains(got, "EXTRA") {
		t.Fatalf("nil Columns should show all: %q", got)
	}
}

func TestColumnFilterMissingKeyShowsEmptyCell(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Columns = []string{"id", "missing"}
	p.Print(json.RawMessage(`[{"id":1,"name":"foo"}]`))

	got := out.String()
	if !strings.Contains(got, "MISSING") {
		t.Fatalf("column header for missing key should appear: %q", got)
	}
	if strings.Contains(got, "NAME") {
		t.Fatalf("unrequested column should not appear: %q", got)
	}
}

func TestColumnFilterIgnoredForJSON(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, JSON)
	p.Columns = []string{"id"}
	p.Print(json.RawMessage(`[{"id":1,"name":"foo"}]`))

	got := out.String()
	if !strings.Contains(got, "name") {
		t.Fatalf("JSON should output all fields regardless of Columns: %q", got)
	}
}

func TestColumnFilterIgnoredForYAML(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, YAML)
	p.Columns = []string{"id"}
	p.Print(json.RawMessage(`[{"id":1,"name":"foo"}]`))

	got := out.String()
	if !strings.Contains(got, "name") {
		t.Fatalf("YAML should output all fields regardless of Columns: %q", got)
	}
}

func TestColumnFilterAutoFlattenWrapper(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Columns = []string{"id", "name"}
	// API response pattern: {data: [...], total: N}
	p.Print(json.RawMessage(`{"data":[{"id":1,"name":"foo","extra":"bar"}],"total":1}`))

	got := out.String()
	if !strings.Contains(got, "ID") && !strings.Contains(got, "NAME") {
		t.Fatalf("column filter should work through auto-flatten: %q", got)
	}
	if strings.Contains(got, "EXTRA") {
		t.Fatalf("filtered column should be hidden in auto-flatten: %q", got)
	}
}

// --- Auto-fit (terminal width) ---

func TestAutoFitTrimsColumns(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.TermWidth = 20 // Very narrow terminal
	p.Print(json.RawMessage(`[{"id":1,"name":"foo","extra":"bar","more":"baz"}]`))

	got := out.String()
	// Should still contain at least the first column.
	if !strings.Contains(got, "ID") {
		t.Fatalf("auto-fit should always show at least one column: %q", got)
	}
}

func TestAutoFitNoTrimWhenFits(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.TermWidth = 200 // Very wide terminal
	p.Print(json.RawMessage(`[{"id":1,"name":"foo"}]`))

	got := out.String()
	if strings.Contains(got, "columns hidden") {
		t.Fatalf("no columns should be hidden in wide terminal: %q", got)
	}
	if !strings.Contains(got, "ID") || !strings.Contains(got, "NAME") {
		t.Fatalf("all columns should show: %q", got)
	}
}

func TestAutoFitZeroTermWidthNoOp(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.TermWidth = 0 // No auto-fit
	p.Print(json.RawMessage(`[{"id":1,"name":"foo","extra":"bar"}]`))

	got := out.String()
	if strings.Contains(got, "columns hidden") {
		t.Fatalf("TermWidth=0 should not auto-fit: %q", got)
	}
}

func TestAutoFitDisabledWithWide(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.TermWidth = 10 // Tiny terminal
	p.Wide = true    // But --wide overrides
	p.Print(json.RawMessage(`[{"id":1,"name":"foo","extra":"bar"}]`))

	got := out.String()
	if strings.Contains(got, "columns hidden") {
		t.Fatalf("--wide should disable auto-fit: %q", got)
	}
}

// --- Nested cell flattening ---

func TestFlattenNestedCellCustomOnly(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Columns = []string{"id", "src_addr"}
	p.Print(json.RawMessage(`[{"id":1,"src_addr":{"custom":["192.168.1.1","10.0.0.1"],"object":[]}}]`))

	got := out.String()
	if !strings.Contains(got, "192.168.1.1") || !strings.Contains(got, "10.0.0.1") {
		t.Fatalf("nested custom values should be flattened: %q", got)
	}
	// Should NOT contain raw JSON braces.
	if strings.Contains(got, `{"custom"`) {
		t.Fatalf("nested cell should not show raw JSON: %q", got)
	}
}

func TestFlattenNestedCellCustomAndObject(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Columns = []string{"id", "src_addr"}
	p.Print(json.RawMessage(`[{"id":1,"src_addr":{"custom":["192.168.1.1"],"object":[{"gp_name":"mygroup","gid":"GPIP1"}]}}]`))

	got := out.String()
	if !strings.Contains(got, "192.168.1.1") {
		t.Fatalf("custom value missing: %q", got)
	}
	if !strings.Contains(got, "[mygroup]") {
		t.Fatalf("object reference should show as [gp_name]: %q", got)
	}
}

func TestFlattenNestedCellEmpty(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Columns = []string{"id", "src_addr"}
	p.Print(json.RawMessage(`[{"id":1,"src_addr":{"custom":[],"object":[]}}]`))

	got := out.String()
	if !strings.Contains(got, "*") {
		t.Fatalf("empty nested cell should show *: %q", got)
	}
}

func TestFlattenNestedCellNotMatching(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`[{"id":1,"config":{"key":"value"}}]`))

	got := out.String()
	// Non-matching map should fall through to JSON serialization.
	if !strings.Contains(got, `{"key":"value"}`) {
		t.Fatalf("non-matching map should render as JSON: %q", got)
	}
}

func TestMultiLineCellRendering(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Columns = []string{"id", "src_addr"}
	p.Print(json.RawMessage(`[{"id":1,"src_addr":{"custom":["192.168.1.1","10.0.0.1"],"object":[]}},{"id":2,"src_addr":{"custom":["172.16.0.1"],"object":[]}}]`))

	got := out.String()
	lines := strings.Split(got, "\n")
	// Row 1 has 2-line src_addr → should produce 2 sub-rows + blank separator.
	// Row 2 has 1-line src_addr → single row, no separator.
	foundSep := false
	for i, l := range lines {
		if l == "" && i > 0 && i < len(lines)-1 {
			foundSep = true
			break
		}
	}
	if !foundSep {
		t.Fatalf("multi-line row should have blank separator: %q", got)
	}
}

func TestSingleLineCellNoExtraBlanks(t *testing.T) {
	var out bytes.Buffer
	p := New(&out, &out, Table)
	p.Print(json.RawMessage(`[{"id":1,"name":"foo"},{"id":2,"name":"bar"}]`))

	got := out.String()
	// No blank lines between single-line rows (header + 2 data + trailing).
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	for i, l := range lines {
		if l == "" {
			t.Fatalf("unexpected blank line at position %d in single-line table: %q", i, got)
		}
	}
}
