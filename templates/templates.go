package templates

import (
	_ "embed"
)

//go:embed data_store.tmpl
var DataStore string

//go:embed directory.tmpl
var Directory string

//go:embed enum.tmpl
var Enum string

//go:embed gob.tmpl
var Gob string

//go:embed state.tmpl
var State string

//go:embed stats.tmpl
var Stats string

//go:embed store.tmpl
var Store string

//go:embed version.tmpl
var Version string
