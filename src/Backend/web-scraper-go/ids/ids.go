package ids

import (
	"fmt"
	"strings"
	"sync"
	"time"
	. "web-scraper/structure"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var (
	// urlToTitle   = make(map[string]string)
	childNparent = make(map[string][]string)
	depthNode    = make(map[string]int)
	baseLink     = "https://en.wikipedia.org"
	limiter      = make(chan int, 150)
	alrFound     = false
	targetTitle  string
	rootTitle    string
	mutex        sync.Mutex

	GraphSolusi             = GraphView{Nodes: []Node{}, Edges: []Edge{}}
	PageScraped             = 0
	ResultDepth             int
	Status                  string
	Err_msg                 string
	urlToTitle              = make(map[string]string)
	solutionParentChildBool = make(map[string]map[string]bool)
	insertedNodeToJSON      = make(map[string]bool)
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
					var foundTitle string = attr

					mutex.Lock() // Mengunci akses ke variabel bersama
					PageScraped = PageScraped + 1
					val, exists := depthNode[foundTitle]
					newVal := depthNode[inputTitle] + 1

					if !exists && val == newVal {
						childNparent[foundTitle] = append(childNparent[foundTitle], inputTitle)
					} else if !exists || val > newVal {
						depthNode[foundTitle] = newVal
						childNparent[foundTitle] = []string{inputTitle}
					}

					if foundTitle == searchTitle {
						alrFound = true
						fmt.Println(inputTitle)
						fmt.Println(foundTitle)
						fmt.Println(iteration)
					} else if !alrFound && iteration != 1 && !(!exists || val > newVal) {
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
	targetTitle = searchTitle
	rootTitle = inputTitle
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

	end := time.Now()
	durasi := end.Sub(start)
	fmt.Println("Waktu eksekusi:", durasi.Milliseconds())

	var a = childNparent[searchTitle]
	fmt.Print(searchTitle, ", ")
	for a[0] != inputTitle {
		fmt.Println(len(a))
		fmt.Print(a[0], ", ")
		a = childNparent[a[0]]
	}
	fmt.Print(a[0])
	fmt.Println("\nPage Scraped: ", PageScraped)

	for _, parentTemp := range childNparent[searchTitle] {
		insertToSolution(searchTitle, parentTemp)
	}
	Status = "OK"
	Err_msg = ""
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

func insertToSolution(child string, parent string) {
	// Mendapatkan Judul Artikel Child
	_, cExist := urlToTitle[child]
	if !cExist {
		cc := colly.NewCollector(
			colly.Async(true),
		)

		cc.OnHTML("#firstHeading", func(e *colly.HTMLElement) {
			if e.ChildText(".mw-page-title-main") != "" {
				urlToTitle[child] = e.ChildText(".mw-page-title-main")
			} else if e.Text != "" {
				urlToTitle[child] = e.Text
			}
		})

		cc.Visit(baseLink + child)
		cc.Wait()
	}

	// Mendapatkan Judul Artikel Parent
	_, pExist := urlToTitle[parent]
	if !pExist {
		cp := colly.NewCollector(
			colly.Async(true),
		)

		cp.OnHTML("#firstHeading", func(e *colly.HTMLElement) {
			if e.ChildText(".mw-page-title-main") != "" {
				urlToTitle[parent] = e.ChildText(".mw-page-title-main")
			} else if e.Text != "" {
				urlToTitle[parent] = e.Text
			}
		})

		cp.Visit(baseLink + parent)
		cp.Wait()
	}

	if urlToTitle[child] == urlToTitle[parent] {
		return
	}
	// Cek apakah child-parent sudah pernah dimasukkan ke solusi
	_, existChild := solutionParentChildBool[urlToTitle[parent]]
	if !existChild {
		solutionParentChildBool[urlToTitle[parent]] = make(map[string]bool)
	}
	if !solutionParentChildBool[urlToTitle[parent]][urlToTitle[child]] {
		solutionParentChildBool[urlToTitle[parent]][urlToTitle[child]] = true
		// Masukkan ke struktur JSON
		insertToJSON(child, parent)
	} else {
		return
	}

	// cek parentnya dari parent
	var n int = len(childNparent[parent])
	if n == 0 || urlToTitle[parent] == rootTitle {
		return
	} else {
		for _, parentTemp := range childNparent[parent] {
			insertToSolution(parent, parentTemp)
		}
	}
}

func insertToJSON(child string, parent string) {
	if urlToTitle[child] == rootTitle {
		depthNode[child] = 0
	} else if urlToTitle[parent] == rootTitle {
		depthNode[parent] = 0
	}
	_, existChildNode := insertedNodeToJSON[urlToTitle[child]]
	if !existChildNode {
		insertedNodeToJSON[urlToTitle[child]] = true
		var font_size int
		var node_size int
		if urlToTitle[child] == targetTitle || urlToTitle[child] == rootTitle {
			font_size = 15
			node_size = 15
		} else {
			font_size = 10
			node_size = 10
		}
		GraphSolusi.Nodes = append(GraphSolusi.Nodes, Node{
			Id:           urlToTitle[child],
			TitleArticle: urlToTitle[child],
			UrlArticle:   baseLink + child,
			Shape:        "star",
			Size:         node_size,
			Color: Color{
				Border:     DepthColor[depthNode[child]],
				Background: DepthColor[depthNode[child]],
			},
			Font: Font{
				Color: DepthColor[depthNode[child]],
				Size:  font_size,
			},
		})
	}
	_, existParentNode := insertedNodeToJSON[urlToTitle[parent]]
	if !existParentNode {
		insertedNodeToJSON[urlToTitle[parent]] = true
		var font_size int
		var node_size int
		if urlToTitle[parent] == targetTitle || urlToTitle[parent] == rootTitle {
			font_size = 15
			node_size = 15
		} else {
			font_size = 10
			node_size = 10
		}
		GraphSolusi.Nodes = append(GraphSolusi.Nodes, Node{
			Id:           urlToTitle[parent],
			TitleArticle: urlToTitle[parent],
			UrlArticle:   baseLink + parent,
			Shape:        "star",
			Size:         node_size,
			Color: Color{
				Border:     DepthColor[depthNode[parent]],
				Background: DepthColor[depthNode[parent]],
			},
			Font: Font{
				Color: DepthColor[depthNode[parent]],
				Size:  font_size,
			},
		})
	}
	GraphSolusi.Edges = append(GraphSolusi.Edges, Edge{
		From: urlToTitle[parent],
		To:   urlToTitle[child],
	})
}

func ResetData() {
	childNparent = make(map[string][]string)
	depthNode = make(map[string]int)
	baseLink = "https://en.wikipedia.org"
	limiter = make(chan int, 150)
	alrFound = false
	GraphSolusi = GraphView{Nodes: []Node{}, Edges: []Edge{}}
	PageScraped = 0
	urlToTitle = make(map[string]string)
	solutionParentChildBool = make(map[string]map[string]bool)
	insertedNodeToJSON = make(map[string]bool)
}
