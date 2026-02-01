package assets

import _ "embed"

// SwaggerSpec contains the embedded swagger.json content
//go:embed swagger.json
var SwaggerSpec []byte
