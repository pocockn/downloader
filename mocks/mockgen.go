package mocks

import _ "github.com/golang/mock/mockgen/model" // Side effects required for mockgen dependencies

//go:generate mockgen --build_flags=--mod=mod -destination=../mocks/mock_store.go -package=mocks github.com/pocockn/downloader/store Store
