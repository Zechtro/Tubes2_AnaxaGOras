package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"
	b "web-scraper/bfs"
	i "web-scraper/ids"
	s "web-scraper/structure"
)

type RequestInfo struct {
	Algorithm  string `json:"algorithm"`
	StartPage  string `json:"startPage"`
	TargetPage string `json:"targetPage"`
}

type ResponseInfo struct {
	Status         string      `json:"status"`
	Error_Message  string      `json:"err"`
	Graph          s.GraphView `json:"graph"`
	ResultDepth    int         `json:"depth"`
	ArticleChecked int         `json:"checked"`
	ExecutionTime  float64     `json:"time"`
}

// Untuk menerima request dari frontend
var reqInfo RequestInfo

// Untuk memberi response ke frontend
var respInfo ResponseInfo

var wiki string = "/wiki/"

func main() {
	http.HandleFunc("/api/process", request_response_Handler)
	fmt.Println("Server listening on port 8000")
	http.ListenAndServe(":8000", nil)
}

func request_response_Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&reqInfo)
		if err != nil {
			http.Error(w, "Error decoding input.", http.StatusBadRequest)
		} else {
			fmt.Println("Received input:")
			fmt.Println("Algorithm: ", reqInfo.Algorithm)
			fmt.Println("Start: ", reqInfo.StartPage)
			fmt.Println("Target: ", reqInfo.TargetPage)
			fmt.Println("Ready to GET")
			// ALGORITMA UTAMA
			if reqInfo.Algorithm == "bfs" {
				// ALGORITMA BFS
				startTime := time.Now()
				fmt.Println("Processing BFS...")
				b.BFS([]string{wiki + reqInfo.StartPage}, reqInfo.TargetPage)
				b.GetSolutionAndConvertToJSON()
				endTime := time.Now()
				respInfo.Status = b.Status
				respInfo.Error_Message = b.Err_msg
				respInfo.Graph = b.GraphSolusi
				respInfo.ResultDepth = b.ResultDepth
				respInfo.ArticleChecked = b.TotalCheckedArticleTitle
				respInfo.ExecutionTime = float64(endTime.Sub(startTime).Seconds() * 1000000 / 1000)
			} else if reqInfo.Algorithm == "ids" {
				// ALGORITMA IDS
				startTime := time.Now()
				fmt.Println("Processing IDS...")
				i.MainIDS(wiki+reqInfo.StartPage, wiki+reqInfo.TargetPage)
				i.GetSolutionAndConvertToJSON()
				endTime := time.Now()
				respInfo.Status = i.Status
				respInfo.Error_Message = i.Err_msg
				respInfo.Graph = i.GraphSolusi
				respInfo.ResultDepth = i.ResultDepth
				respInfo.ArticleChecked = i.PageScraped
				respInfo.ExecutionTime = float64(endTime.Sub(startTime).Seconds() * 1000000 / 1000)

			}
			// RESPONSE TO FRONTEND
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			err := json.NewEncoder(w).Encode(respInfo)
			b.ResetData()
			i.ResetData()
			fmt.Println("Reset Data")
			if err != nil {
				http.Error(w, "Error encoding data", http.StatusInternalServerError)
			}
		}
	} else if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println("NOT POST", r.Method)
	}
}
