/*
Copyright 2023 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

// https://github.com/aquasecurity/trivy-operator/blob/84df941b628441c285c08850bf73fd0e5fd3aa05/pkg/apis/aquasecurity/v1alpha1/common_types.go

package v1alpha1

import (
	"fmt"
	"strings"
)

const (
	TTLReportAnnotation = "trivy-operator.aquasecurity.github.io/report-ttl"
)

// Severity level of a vulnerability or a configuration audit check.
// +enum
type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityMedium   Severity = "MEDIUM"
	SeverityLow      Severity = "LOW"

	SeverityUnknown Severity = "UNKNOWN"
)

// StringToSeverity returns the enum constant of Severity with the specified
// name. The name must match exactly an identifier used to declare an enum
// constant. (Extraneous whitespace characters are not permitted.)
func StringToSeverity(name string) (Severity, error) {
	s := strings.ToUpper(name)
	switch s {
	case "CRITICAL", "HIGH", "MEDIUM", "LOW", "NONE", "UNKNOWN":
		return Severity(s), nil
	default:
		return "", fmt.Errorf("unrecognized name literal: %s", name)
	}
}

const ScannerNameTrivy = "Trivy"

// Scanner is the spec for a scanner generating a security assessment report.
type Scanner struct {
	// Name the name of the scanner.
	Name string `json:"name"`

	// Vendor the name of the vendor providing the scanner.
	Vendor string `json:"vendor"`

	// Version the version of the scanner.
	Version string `json:"version"`
}
