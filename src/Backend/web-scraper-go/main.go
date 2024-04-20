package main

import (
	// "fmt"
	"encoding/json"
	"fmt"
	"net/http"

	// b "lala/bfs"
	"time"
	b "web-scraper/bfs"
)

func main() {
	http.HandleFunc("/userInfo", userInfoHandler)

	var wiki string = "/wiki/"
	var start_title = "Neuroscience"
	// var start_title = "Joko_Widodo"
	// var start_title = "Indonesia"
	var page string = wiki + start_title
	var start_page = []string{page}
	startTime := time.Now()
	// depth 1 - Neuroscience, Jokowi
	// b.BFS(start_page, "Anatomy")
	// depth 2 - Neuroscience
	// b.BFS(start_page, "Antenna_(biology)")
	// depth 2 - Jokowi
	// b.BFS(start_page, "Adolf_Hitler")
	// depth 3 - Neuroscience
	b.BFS(start_page, "Springtail")
	// Failed(?)
	// b.BFS(start_page, "Collophore")
	// b.BFS(start_page, "Special:BookSources/978-0521570480")
	endTime := time.Now()
	fmt.Println("Executed in ", endTime.Sub(startTime).Seconds()*1000, " ms")
	j, _ := json.Marshal(b.SolutionGraph)
	fmt.Println(string(j))
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	// myMap := map[string]interface{}{
	// 	"name": "Alice",
	// 	"age":  30,
	// 	"city": "New York",
	// }

	jsonBytes, err := json.Marshal(b.SolutionGraph)
	if err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow requests from any origin (replace with your frontend origin if needed)
	w.Write(jsonBytes)
}

// func getMessage(w http.ResponseWriter, r *http.Request) {
// 	message := Message{Text: "Hello from Go backend!"}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow requests from any origin (replace with your frontend origin if needed)
// 	json.NewEncoder(w).Encode(message)
// }
