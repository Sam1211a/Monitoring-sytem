package model

import "sync"

var (
	LatestFrame = make(map[string][]byte)
	Mu          sync.RWMutex
)
