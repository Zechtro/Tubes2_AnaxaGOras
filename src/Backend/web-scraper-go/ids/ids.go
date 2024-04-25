package ids

import (
	"fmt"
	"strings"
	"sync" // Import package sync untuk mengatur goroutine
	"time"

	"github.com/gocolly/colly"
	"github.com/PuerkitoBio/goquery"
)

var (
	childNparent = make(map[string][]string)
	depthNode    = make(map[string]int)
	baseLink 	 = "https://en.wikipedia.org/wiki/"
	limiter		 = make(chan int, 15)
	pageScraped  = 0
	alrFound     = false
	// unvisitedPath []string
	mutex        sync.Mutex // Mutex untuk mengatur akses ke variabel bersama
)

func IDS(inputTitle string, searchTitle string, iteration int, wg *sync.WaitGroup) {
	defer wg.Done() // Menandai bahwa goroutine telah selesai
	pageToScrape := baseLink + inputTitle

	c := colly.NewCollector()

	// c.OnRequest(func(r *colly.Request) { 
	// 	fmt.Println("Visiting: ", r.URL) 
	// }) 
	// c.OnError(func(_ *colly.Response, err error) { 
	// 	fmt.Println("XXXXXXXX Something went wrong: ", err) 
	// }) 
	// c.OnResponse(func(r *colly.Response) { 
	// 	fmt.Println(iteration, "Page visited: ", r.Request.URL) 
	// }) 

	c.OnHTML("#bodyContent", func(e *colly.HTMLElement) {
		e.DOM.Find("a").Each(func(_ int, s *goquery.Selection) {
			if attr, ok := s.Attr("href"); ok {
				if isWiki(attr) {
					var foundTitle string = getArticleTitle(attr)
		
					mutex.Lock() // Mengunci akses ke variabel bersama
					pageScraped = pageScraped + 1
					val, exists := depthNode[foundTitle]
					newVal := depthNode[inputTitle] + 1
		
					if (!exists && val == newVal) {
						childNparent[foundTitle] = append(childNparent[foundTitle], inputTitle)
					} else if (!exists || val > newVal) {
						depthNode[foundTitle] = newVal
						childNparent[foundTitle] = []string{inputTitle}
					}
		
					if (foundTitle == searchTitle) {
						alrFound = true
						fmt.Println(inputTitle)
						fmt.Println(foundTitle)
						fmt.Println(iteration)
					} else if (!alrFound && iteration != 1 && !(!exists || val > newVal)) {
						wg.Add(1) // Menambahkan goroutine baru ke wait group
						
						go IDS(foundTitle, searchTitle, iteration-1, wg)
					}
					mutex.Unlock() // Membuka kunci akses ke variabel bersama
				}
			}
		})
		
		
	})
	limiter <- 1
	c.Visit(pageToScrape)
	<-limiter
}

func MainIDS(inputTitle string, searchTitle string) {
	childNparent[inputTitle] = []string{inputTitle}
	depthNode[inputTitle] = 1
	iteration := 1

	start := time.Now()
	var wg sync.WaitGroup

	for !alrFound {
		wg.Add(1) // Menambahkan goroutine pertama ke wait group
		go IDS(inputTitle, searchTitle, iteration, &wg)
		wg.Wait() // Menunggu sampai semua goroutine selesai
		fmt.Println("TESTING")
		iteration += 1
	}

	// Setelah semua goroutine selesai, lanjutkan dengan menampilkan hasil
	// for !alrFound {
	// 	for _, input := range unvisitedPath {
	// 		wg.Add(1) // Menambahkan goroutine baru ke wait group
	// 		go IDS(input, searchTitle, iteration, &wg)
	// 	}
	// 	wg.Wait() // Menunggu sampai semua goroutine selesai
	// 	unvisitedPath = []string{}
	// }

	end := time.Now()
	durasi := end.Sub(start)
	fmt.Println("Waktu eksekusi:", durasi)

	var a = childNparent[searchTitle]
	fmt.Print(searchTitle, ", ")
	for a[0] != inputTitle {
		fmt.Println(len(a))
		fmt.Print(a[0], ", ")
		a = childNparent[a[0]]
	}
	fmt.Print(a[0])
	fmt.Println("\nPage Scraped: ", pageScraped)

}

// HELPER FUNCTIONS
func isWiki(link string) bool {
	if len(link) <= 6 {
		return false
	} else if link[:6] == "/wiki/" {
		if strings.ContainsRune(link[6:], ':') {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

func getArticleTitle(link string) string {
	return link[6:]
}
