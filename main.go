// main.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Define a struct for your data (e.g., a simple "Item" struct).
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Create an in-memory "database" (a slice) to store your items.
var items []Item

func main() {
	// Define HTTP routes and their corresponding handlers.
	http.HandleFunc("/items", GetItems)
	http.HandleFunc("/items/", GetItem)
	http.HandleFunc("/items/create", CreateItem)
	http.HandleFunc("/items/update/", UpdateItem)

	// Start the HTTP server.
	port := 8080
	fmt.Printf("Server listening on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// GetItems handles GET requests to retrieve all items.
func GetItems(w http.ResponseWriter, r *http.Request) {
	// Convert items slice to JSON and send it in the response.
	json.NewEncoder(w).Encode(items)
}

// GetItem handles GET requests to retrieve a specific item by ID.
func GetItem(w http.ResponseWriter, r *http.Request) {
	// Extract the item ID from the request URL.
	id := r.URL.Path[len("/items/"):]

	// Find the item with the given ID.
	var item Item
	for _, i := range items {
		if i.ID == id {
			item = i
			break
		}
	}

	if item.ID == "" {
		http.NotFound(w, r)
		return
	}

	// Convert the item to JSON and send it in the response.
	json.NewEncoder(w).Encode(item)
}

// CreateItem handles POST requests to create a new item.
func CreateItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a unique ID (you can use a library like uuid for this).
	newItem.ID = "unique_id_here"

	// Append the new item to the items slice.
	items = append(items, newItem)

	// Respond with the created item.
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

// UpdateItem handles PUT requests to update an existing item by ID.
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	// Extract the item ID from the request URL.
	id := r.URL.Path[len("/items/update/"):]

	// Find the item with the given ID.
	var updatedItem Item
	for i, item := range items {
		if item.ID == id {
			err := json.NewDecoder(r.Body).Decode(&updatedItem)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Update the item in the slice.
			items[i] = updatedItem

			// Respond with the updated item.
			json.NewEncoder(w).Encode(updatedItem)
			return
		}
	}

	// If the item is not found, return a 404 Not Found response.
	http.NotFound(w, r)
}
