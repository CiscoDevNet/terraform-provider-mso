package tfplugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// configTargetPath resolves a rooted or relative expression to one configured attribute.
func configTargetPath(ctx context.Context, config tfsdk.Config, source path.Path, target path.Expression, diagnosticTitle string, diagnostics *diag.Diagnostics) (path.Path, bool) {
	matches, matchDiagnostics := config.PathMatches(ctx, source.Expression().Merge(target))
	diagnostics.Append(matchDiagnostics...)
	if matchDiagnostics.HasError() {
		return path.Path{}, false
	}
	if len(matches) != 1 {
		diagnostics.AddError(diagnosticTitle, fmt.Sprintf("Expected exactly one match for %s from %s, got %d.", target, source, len(matches)))
		return path.Path{}, false
	}
	return matches[0], true
}
