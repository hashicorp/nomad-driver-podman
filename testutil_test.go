// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/go-viper/mapstructure/v2"
	"github.com/hashicorp/go-msgpack/v2/codec"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	hcl2 "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	hjson "github.com/hashicorp/hcl/v2/json"
	"github.com/hashicorp/nomad/helper/pluginutils/hclspecutils"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/hashicorp/nomad/plugins/drivers"
	"github.com/hashicorp/nomad/plugins/shared/hclspec"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

// hclParser is a test helper for parsing driver TaskConfig from HCL.
// This replaces the import of github.com/hashicorp/nomad/helper/pluginutils/hclutils
// to eliminate the transitive Docker dependency.
type hclParser struct {
	spec *hclspec.Spec
	vars map[string]cty.Value
}

// newConfigParser returns a helper for parsing driver TaskConfig.
// Parser is an immutable object that can be used in multiple tests.
func newConfigParser(spec *hclspec.Spec) *hclParser {
	return &hclParser{
		spec: spec,
	}
}

// ParseHCL parses the HCL config string and decodes it into the out parameter.
// out parameter should be a pointer to a driver-specific TaskConfig.
func (p *hclParser) ParseHCL(t *testing.T, configStr string, out interface{}) {
	t.Helper()
	config := hclConfigToInterface(t, configStr)
	p.parse(t, config, out)
}

func (p *hclParser) parse(t *testing.T, config, out interface{}) {
	t.Helper()

	decSpec, diags := hclspecutils.Convert(p.spec)
	if diags.HasErrors() {
		t.Fatalf("failed to convert spec: %v", diags.Error())
	}

	ctyValue, diag, errs := parseHclInterface(config, decSpec, p.vars)
	if len(errs) > 0 {
		t.Error("unexpected errors parsing config")
		for _, err := range errs {
			t.Errorf("  * %v", err)
		}
		t.FailNow()
	}
	if diag.HasErrors() {
		t.Fatalf("unexpected diagnostics: %v", diag.Error())
	}

	// encode
	dtc := &drivers.TaskConfig{}
	if err := dtc.EncodeDriverConfig(ctyValue); err != nil {
		t.Fatalf("failed to encode driver config: %v", err)
	}

	// decode
	if err := dtc.DecodeDriverConfig(out); err != nil {
		t.Fatalf("failed to decode driver config: %v", err)
	}
}

// hclConfigToInterface parses an HCL config string into a Go interface.
func hclConfigToInterface(t *testing.T, config string) interface{} {
	t.Helper()

	root, err := hcl.Parse(config)
	if err != nil {
		t.Fatalf("failed to hcl parse the config: %v", err)
	}

	list, ok := root.Node.(*ast.ObjectList)
	if !ok {
		t.Fatalf("root should be an object")
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, list.Items[0]); err != nil {
		t.Fatalf("failed to decode object: %v", err)
	}

	var m2 map[string]interface{}
	if err := mapstructure.WeakDecode(m, &m2); err != nil {
		t.Fatalf("failed to weak decode object: %v", err)
	}

	return m2["config"]
}

// parseHclInterface converts an interface value representing an HCL2 body
// and returns the interpolated cty.Value.
func parseHclInterface(val interface{}, spec hcldec.Spec, vars map[string]cty.Value) (cty.Value, hcl2.Diagnostics, []error) {
	evalCtx := &hcl2.EvalContext{
		Variables: vars,
		Functions: getStdlibFuncs(),
	}

	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, structs.JsonHandle)
	err := enc.Encode(val)
	if err != nil {
		errorMessage := fmt.Sprintf("Label encoding failed: %v", err)
		return cty.NilVal,
			hcl2.Diagnostics([]*hcl2.Diagnostic{{
				Severity: hcl2.DiagError,
				Summary:  "Failed to encode label value",
				Detail:   errorMessage,
			}}),
			[]error{errors.New(errorMessage)}
	}

	hclFile, diag := hjson.Parse(buf.Bytes(), "")
	if diag.HasErrors() {
		return cty.NilVal, diag, formattedDiagnosticErrors(diag)
	}

	value, decDiag := hcldec.Decode(hclFile.Body, spec, evalCtx)
	diag = diag.Extend(decDiag)
	if diag.HasErrors() {
		return cty.NilVal, diag, formattedDiagnosticErrors(diag)
	}

	return value, diag, nil
}

func formattedDiagnosticErrors(diag hcl2.Diagnostics) []error {
	var errs []error
	for _, d := range diag {
		if d.Severity == hcl2.DiagError {
			errs = append(errs, fmt.Errorf("%s: %s", d.Summary, d.Detail))
		}
	}
	return errs
}

func getStdlibFuncs() map[string]function.Function {
	return map[string]function.Function{
		"abs":        stdlib.AbsoluteFunc,
		"coalesce":   stdlib.CoalesceFunc,
		"concat":     stdlib.ConcatFunc,
		"hasindex":   stdlib.HasIndexFunc,
		"int":        stdlib.IntFunc,
		"jsondecode": stdlib.JSONDecodeFunc,
		"jsonencode": stdlib.JSONEncodeFunc,
		"length":     stdlib.LengthFunc,
		"lower":      stdlib.LowerFunc,
		"max":        stdlib.MaxFunc,
		"min":        stdlib.MinFunc,
		"reverse":    stdlib.ReverseFunc,
		"strlen":     stdlib.StrlenFunc,
		"substr":     stdlib.SubstrFunc,
		"upper":      stdlib.UpperFunc,
	}
}
