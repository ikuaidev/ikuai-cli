package cliapp

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

func TestMergeDataWithFlags_EmptyDataNoFlags(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "", "")
	result, err := MergeDataWithFlags("{}", cmd, map[string]string{"name": "tagname"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %v", result)
	}
}

func TestMergeDataWithFlags_DataPassthrough(t *testing.T) {
	cmd := &cobra.Command{}
	result, err := MergeDataWithFlags(`{"foo":"bar","num":42}`, cmd, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result["foo"] != "bar" {
		t.Fatalf("foo = %v, want bar", result["foo"])
	}
}

func TestMergeDataWithFlags_FlagOverridesData(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "", "")
	_ = cmd.Flags().Set("name", "override")
	result, err := MergeDataWithFlags(`{"tagname":"original","extra":"keep"}`, cmd, map[string]string{"name": "tagname"})
	if err != nil {
		t.Fatal(err)
	}
	if result["tagname"] != "override" {
		t.Fatalf("tagname = %v, want override", result["tagname"])
	}
	if result["extra"] != "keep" {
		t.Fatalf("extra = %v, want keep", result["extra"])
	}
}

func TestMergeDataWithFlags_UnchangedFlagNotIncluded(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "", "")
	// Do NOT call Set — flag is not Changed
	result, err := MergeDataWithFlags("{}", cmd, map[string]string{"name": "tagname"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result["tagname"]; ok {
		t.Fatal("unchanged flag should not appear in result")
	}
}

func TestMergeDataWithFlags_InvalidJSON(t *testing.T) {
	cmd := &cobra.Command{}
	_, err := MergeDataWithFlags("not json", cmd, nil)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestMergeDataWithFlags_NumericPreservation(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("port", "", "")
	_ = cmd.Flags().Set("port", "8022")
	result, err := MergeDataWithFlags("{}", cmd, map[string]string{"port": "sshd_port"})
	if err != nil {
		t.Fatal(err)
	}
	// Should be int64, not string
	val := result["sshd_port"]
	if _, ok := val.(int64); !ok {
		t.Fatalf("sshd_port type = %T, want int64, value = %v", val, val)
	}
	if val.(int64) != 8022 {
		t.Fatalf("sshd_port = %v, want 8022", val)
	}
}

func TestMergeDataWithFlags_StringNotCoerced(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "", "")
	_ = cmd.Flags().Set("name", "myhost")
	result, err := MergeDataWithFlags("{}", cmd, map[string]string{"name": "hostname"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result["hostname"].(string); !ok {
		t.Fatalf("hostname type = %T, want string", result["hostname"])
	}
}

func TestMergeDataWithFlags_NumericStringNotCoerced(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "", "")
	_ = cmd.Flags().Set("name", "12345")
	result, err := MergeDataWithFlags("{}", cmd, map[string]string{"name": "tagname"})
	if err != nil {
		t.Fatal(err)
	}
	// "tagname" is NOT in integerAPIFields, so "12345" must stay as string.
	if _, ok := result["tagname"].(string); !ok {
		t.Fatalf("tagname type = %T, want string (value %v)", result["tagname"], result["tagname"])
	}
	if result["tagname"] != "12345" {
		t.Fatalf("tagname = %v, want 12345", result["tagname"])
	}
}

func TestMergeDataWithFlags_EnabledYes(t *testing.T) {
	cmd := &cobra.Command{}
	AddEnabledFlag(cmd)
	_ = cmd.Flags().Set("enabled", "yes")
	result, err := MergeDataWithFlags("{}", cmd, map[string]string{"enabled": "enabled"})
	if err != nil {
		t.Fatal(err)
	}
	if result["enabled"] != "yes" {
		t.Fatalf("enabled = %v, want yes", result["enabled"])
	}
}

func TestAddListFlags(t *testing.T) {
	cmd := &cobra.Command{}
	AddListFlags(cmd)
	for _, name := range []string{"page", "page-size", "filter", "order", "order-by"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("missing flag %q", name)
		}
	}
	// -p short flag for page
	if cmd.Flags().ShorthandLookup("p") == nil {
		t.Error("missing -p short flag for --page")
	}
}

func TestListParams(t *testing.T) {
	p := ListParams(2, 50, "enabled==yes", "desc", "name")
	if p["page"] != "2" || p["page_size"] != "50" {
		t.Fatalf("page/page_size wrong: %v", p)
	}
	if p["filter"] != "enabled==yes" {
		t.Fatalf("filter = %v", p["filter"])
	}
	if p["order"] != "desc" || p["order_by"] != "name" {
		t.Fatalf("order/order_by wrong: %v", p)
	}
}

func TestListParamsOmitsEmpty(t *testing.T) {
	p := ListParams(1, 20, "", "", "")
	if _, ok := p["filter"]; ok {
		t.Fatal("empty filter should be omitted")
	}
	if _, ok := p["order"]; ok {
		t.Fatal("empty order should be omitted")
	}
}

func TestParseJSON(t *testing.T) {
	v, err := ParseJSON(`{"key":"val"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := v.(map[string]interface{})
	if m["key"] != "val" {
		t.Fatalf("key = %v", m["key"])
	}
}

func TestParseJSON_Empty(t *testing.T) {
	v, err := ParseJSON("{}")
	if err != nil {
		t.Fatal(err)
	}
	m := v.(map[string]interface{})
	if len(m) != 0 {
		t.Fatalf("expected empty map, got %v", m)
	}
}

func TestParseJSON_Invalid(t *testing.T) {
	_, err := ParseJSON("bad")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRequireFlags_AllChanged(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("interface", "", "")
	_ = cmd.Flags().Set("name", "foo")
	_ = cmd.Flags().Set("interface", "eth0")
	if err := RequireFlags(cmd, "name", "interface"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRequireFlags_OneMissing(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "", "")
	err := RequireFlags(cmd, "name")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if valErr.Message != "missing required flag: --name" {
		t.Fatalf("unexpected message: %q", valErr.Message)
	}
}

func TestRequireFlags_MultipleMissing(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("interface", "", "")
	err := RequireFlags(cmd, "name", "interface")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if valErr.Message != "missing required flags: --name, --interface" {
		t.Fatalf("unexpected message: %q", valErr.Message)
	}
}

func TestRequireFlags_EmptyList(t *testing.T) {
	cmd := &cobra.Command{}
	if err := RequireFlags(cmd); err != nil {
		t.Fatalf("expected nil for empty flag list, got %v", err)
	}
}

func TestRequireFlags_PartiallySet(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("interface", "", "")
	_ = cmd.Flags().Set("name", "foo")
	// interface not set
	err := RequireFlags(cmd, "name", "interface")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if valErr.Message != "missing required flag: --interface" {
		t.Fatalf("unexpected message: %q", valErr.Message)
	}
}

// Verify dry-run is tested via the API client test package, not here.
// This file tests the cliapp helper functions only.
