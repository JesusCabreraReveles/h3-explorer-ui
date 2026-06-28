// Package openapi embeds the API contract so the server can serve it and the
// binary stays self-contained (no runtime file dependency).
package openapi

import _ "embed"

// Spec is the raw OpenAPI 3.1 document for the H3 Explorer API.
//
//go:embed openapi.yaml
var Spec []byte
