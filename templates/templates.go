package templates

import (
	_ "embed"
)

//go:embed data_store.go.tmpl
var DataStoreGo string

//go:embed directory.go.tmpl
var DirectoryGo string

//go:embed enum.go.tmpl
var EnumGo string

//go:embed enum.ts.tmpl
var EnumTs string

//go:embed gob.go.tmpl
var GobGo string

//go:embed state.go.tmpl
var StateGo string

//go:embed stats.go.tmpl
var StatsGo string

//go:embed store.go.tmpl
var StoreGo string
