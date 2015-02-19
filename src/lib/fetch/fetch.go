package fetch

import (
	"frontend/state"
)

type Fetcher struct {
	SS *state.SharedState
}

func NewFetcher(ss *state.SharedState) *Fetcher {
	return &Fetcher{ss}
}
