// Package cache temporarily stores items
package cache

import (
	"log"
	"sync"
	"time"
)

// Store is a cache store. It holds the items until their lifetime has expired
type Store struct {
	store     *sync.Map
	gcRound   time.Duration
	deleted   chan interface{}
	OnCollect func(interface{})
}

type Item struct {
	ExpiryTime time.Time
	Data       interface{}
}

func New(GCRound time.Duration, onClean func(interface{})) *Store {
	s := &Store{
		store:     &sync.Map{},
		gcRound:   GCRound,
		deleted:   make(chan interface{}, 10),
		OnCollect: onClean,
	}

	// Run cache collector
	go func() {
		for {
			s.OnCollect(<-s.deleted)
		}
	}()

	// Run GC
	go s.gc()

	return s
}

// Store stores an item with the assiociated key
func (s *Store) Store(key interface{}, value *Item) {
	s.store.Store(key, value)
}

// Retrieve searches for the store and returns the value stored.
func (s *Store) Retrieve(key interface{}) *Item {
	val, ok := s.store.Load(key)
	if !ok {
		return nil
	}

	// returns nil if the assertion fails
	// the underscore makes it so go doesn't panic
	item, _ := val.(*Item)
	return item
}

func (s *Store) gc() {
	ticker := time.NewTicker(s.gcRound)
	for range ticker.C {
		s.store.Range(func(key, value interface{}) bool {
			val, ok := value.(*Item)
			if !ok {
				log.Fatalf("something else than an item has been stored in the store: %#v, exiting", value)
			}

			if time.Now().After(val.ExpiryTime) {
				old, ok := s.store.LoadAndDelete(key)
				if ok {
					s.deleted <- old
				}
			}

			return true
		})
	}
}
