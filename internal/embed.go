package internal

import _ "embed"

//go:embed templates/typeharness.ts
var TypeHarness string

//go:embed templates/runner.mjs
var RunnerMJS string
