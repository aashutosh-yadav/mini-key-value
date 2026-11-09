// a simple kv store to simulate the working of a kv store
// this is just a implementation of a simplest form of kv store without any dis sys kept in mind
// here we have used mutex locks but there is no concurrency

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type KVStore struct {
	data map[string]string
	mu   sync.RWMutex
	file string
}

// Create a new key-value store
func NewKVStore(filename string) *KVStore {
	store := &KVStore{
		data: make(map[string]string),
		file: filename,
	}

	store.load()

	return store
}

// PUT: Insert or update a key
func (s *KVStore) Put(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value

	return s.save()
}

// GET: Retrieve a value
func (s *KVStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	return value, exists
}

// DELETE: Remove a key
func (s *KVStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)

	return s.save()
}

// EXISTS: Check whether a key exists
func (s *KVStore) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.data[key]
	return exists
}

// Save data to a JSON file
func (s *KVStore) save() error {
	fileData, err := json.Marshal(s.data)
	if err != nil {
		return err
	}

	return os.WriteFile(s.file, fileData, 0644)
}

// Load data from a JSON file
func (s *KVStore) load() {
	fileData, err := os.ReadFile(s.file)

	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error loading database:", err)
		}
		return
	}

	err = json.Unmarshal(fileData, &s.data)

	if err != nil {
		fmt.Println("Error parsing database:", err)
	}
}

func main() {
	store := NewKVStore("database.json")

	// PUT
	err := store.Put("name", "Alice")
	if err != nil {
		fmt.Println(err)
		return
	}

	store.Put("city", "Pune")

	// GET
	value, exists := store.Get("name")

	if exists {
		fmt.Println("Name:", value)
	}

	// EXISTS
	fmt.Println("Name exists:", store.Exists("name"))

	// DELETE
	store.Delete("city")

	fmt.Println("City exists:", store.Exists("city"))
}
