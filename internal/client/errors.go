package client

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gravitino/terraform-provider-gravitino/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewResourceError(operation, resource string, err error) diag.Diagnostics {
	var diags diag.Diagnostics

	var errResp models.ErrorResponse
	if errors.As(err, &errResp) {
		diags = append(diags, diag.NewErrorDiagnostic(
			fmt.Sprintf("Failed %s %q", operation, resource),
			fmt.Sprintf("Server returned [%d %s]: %s", errResp.Code, errResp.Type, errResp.Message),
		))
		if len(errResp.Stack) > 0 {
			diags = append(diags, diag.NewErrorDiagnostic(
				"Server Stack Trace",
				strings.Join(errResp.Stack, "\n"),
			))
		}
		return diags
	}

	diags = append(diags, diag.NewErrorDiagnostic(
		fmt.Sprintf("Failed %s %q", operation, resource),
		err.Error(),
	))

	if hint := errorHint(err); hint != "" {
		diags = append(diags, diag.NewErrorDiagnostic("Hint", hint))
	}

	return diags
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Not Found")
}

func errorHint(err error) string {
	if err == nil {
		return ""
	}
	errStr := err.Error()

	if strings.Contains(errStr, "connection refused") {
		return "The Gravitino server is unreachable. Verify the URI and that the server is running."
	}
	if strings.Contains(errStr, "no such host") {
		return "Host not found. Verify the Gravitino server hostname is correct."
	}
	if strings.Contains(errStr, "certificate") || strings.Contains(errStr, "x509") {
		return "TLS certificate error. If using self-signed certificates, configure your HTTP client accordingly."
	}
	if strings.Contains(errStr, "401") || strings.Contains(errStr, "Unauthorized") {
		return "Authentication failed. Check your credentials (username/password or OAuth token)."
	}
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
		return "Request timed out. The Gravitino server may be overloaded or unreachable. Try increasing the timeout."
	}
	if strings.Contains(errStr, "invalid character") {
		return "Invalid response from server. The URI may be incorrect (e.g., pointing to a web page instead of the API)."
	}

	return ""
}
