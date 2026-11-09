// a simple kv store to simulate the working of a kv store
// this is just a implementation of a simplest form of kv store without any dis sys kept in mind
// here we have used mutex locks but there is no concurrency

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
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

// Save data to JSON file
func (s *KVStore) save() error {
	fileData, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.file, fileData, 0644)
}

// Load data from JSON file
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

// Generate a random 8-character name
func randomName() string {
	const letters = "abcdefghijklmnopqrstuvwxyz"

	name := make([]byte, 8)

	for i := range name {
		name[i] = letters[rand.Intn(len(letters))]
	}

	return string(name)
}

func main() {
	store := NewKVStore("database.json")

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1: Write 100 users
	go func() {
		defer wg.Done()

		for i := 1; i <= 100; i++ {
			user := fmt.Sprintf("user%d", i)
			name := randomName()
			age := 20 + rand.Intn(30)

			err := store.Put(user+":name", name)
			if err != nil {
				fmt.Println("Error writing name:", err)
			}

			err = store.Put(user+":age", fmt.Sprintf("%d", age))
			if err != nil {
				fmt.Println("Error writing age:", err)
			}
		}
	}()

	// Goroutine 2: Read user information
	go func() {
		defer wg.Done()

		for i := 1; i <= 100; i++ {
			user := fmt.Sprintf("user%d", i)

			name, nameExists := store.Get(user + ":name")
			age, ageExists := store.Get(user + ":age")

			if nameExists && ageExists {
				fmt.Printf(
					"%s -> Name: %s, Age: %s\n",
					user, name, age,
				)
			}
		}
	}()

	// Wait for both goroutines
	wg.Wait()

	fmt.Println("All goroutines finished!")

	// Check final values
	name, exists := store.Get("user1:name")
	if exists {
		fmt.Println("User1 name:", name)
	}

	age, exists := store.Get("user1:age")
	if exists {
		fmt.Println("User1 age:", age)
	}
}
