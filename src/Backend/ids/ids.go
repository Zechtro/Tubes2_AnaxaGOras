package ids

import (
	"fmt"
	"strings"
	"sync"
	. "web-scraper/structure"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var (
	childParentBool = make(map[string]map[string]bool)
	depthNode       = make(map[string]int)
	scrapedNode     = make(map[string]bool)
	baseLink        = "https://en.wikipedia.org"
	limiter         = make(chan int, 150)
	alrFound        = false
	targetTitle     string
	rootTitle       string
	mutex           sync.Mutex

	target string
	root   string

	GraphSolusi             = GraphView{Nodes: []Node{}, Edges: []Edge{}}
	PageScraped             = 0
	ResultDepth             int
	Status                  string
	Err_msg                 string
	isInit                  bool = false
	urlToTitle                   = make(map[string]string)
	solutionParentChildBool      = make(map[string]map[string]bool)
	insertedNodeToJSON           = make(map[string]bool)
	depthTarget             int  = 999
	totalErrorScrape        int  = 0
)

func IDS(inputTitle string, target string, iteration int, wg *sync.WaitGroup) {
	defer wg.Done() // Menandai bahwa goroutine telah selesai
	pageToScrape := baseLink + inputTitle

	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		totalErrorScrape += 1
	})

	c.OnHTML("link", func(e *colly.HTMLElement) {
		if e.Attr("rel") == "canonical" {
			foundTitle := e.Attr("href")[24:]
			mutex.Lock()
			if foundTitle != inputTitle && foundTitle != root {
				// mutex.Lock()
				val := depthNode[foundTitle] // val bernilai nol jika foundTitle belum pernah discrape
				newVal := depthNode[inputTitle] + 1
				if val == 0 && foundTitle != root {
					PageScraped += 1
					depthNode[foundTitle] = newVal
				}
				_, existParentKey := childParentBool[foundTitle]
				if !existParentKey {
					childParentBool[foundTitle] = make(map[string]bool)
				}
				depthNode[foundTitle] = depthNode[inputTitle]

				for keyParent, _ := range childParentBool[inputTitle] {
					childParentBool[foundTitle][keyParent] = true
				}
				childParentBool[inputTitle] = childParentBool[foundTitle]
				if foundTitle == target {
					// masukin ke solusi
					fmt.Println(inputTitle, target)
					// insertToSolution(foundTitle, inputTitle)
					alrFound = true
					depthTarget = depthNode[inputTitle]
					fmt.Println("REDIRECT:", depthTarget)
					e.Request.Abort()
				} else {
					scrapedNode[foundTitle] = true
				}
				// mutex.Unlock()
			}
			mutex.Unlock()
		}
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		mutex.Lock() // Mengunci akses ke variabel bersama
		if e.Attr("class") != "mw-file-description" {
			if isWiki(e.Attr("href")) && e.Attr("href") != root && e.Attr("href") != "/wiki/Main_Page" {
				var foundTitle string = e.Attr("href")

				val := depthNode[foundTitle] // val bernilai nol jika foundTitle belum pernah discrape
				newVal := depthNode[inputTitle] + 1

				if val != 0 && val == newVal {
					if !childParentBool[foundTitle][inputTitle] {
						childParentBool[foundTitle][inputTitle] = true
					}
				} else if val == 0 {
					PageScraped = PageScraped + 1
					depthNode[foundTitle] = newVal
					childParentBool[foundTitle] = make(map[string]bool)
					childParentBool[foundTitle][inputTitle] = true
				}
				if foundTitle == target && newVal <= depthTarget {
					alrFound = true
					depthTarget = newVal
					fmt.Println(inputTitle)
					fmt.Println(foundTitle)
					fmt.Println(iteration)
				} else if iteration != 1 && !scrapedNode[foundTitle] {
					scrapedNode[foundTitle] = true
					wg.Add(1) // Menambahkan goroutine baru ke wait group
					go IDS(foundTitle, target, iteration-1, wg)
				}
			}
		}
		mutex.Unlock() // Membuka kunci akses ke variabel bersama
	})
	limiter <- 1
	c.Visit(pageToScrape)
	<-limiter
}

func MainIDS(inputTitle string, searchTitle string) {
	var invalidStart bool = false
	var invalidTarget bool = false
	if !isInit {
		isInit = true
		var wg sync.WaitGroup
		for i := 0; i < 2; i++ {
			wg.Add(1)
			limiter <- 1
			go func(i int) {
				defer wg.Done()
				c1 := colly.NewCollector(
					colly.Async(true),
				)
				extensions.RandomUserAgent(c1)

				c1.OnError(func(_ *colly.Response, err error) {
					fmt.Println("Invalid", err)
					if i == 0 {
						invalidStart = true
					} else {
						invalidTarget = true
					}
				})

				c1.OnHTML("#firstHeading", func(e *colly.HTMLElement) {
					if e.ChildText(".mw-page-title-main") != "" {
						if i == 0 {
							rootTitle = e.ChildText(".mw-page-title-main")
						} else {
							targetTitle = e.ChildText(".mw-page-title-main")
						}
					} else if e.Text != "" {
						if i == 0 {
							rootTitle = e.Text
						} else {
							targetTitle = e.Text
						}
					}
				})

				c1.OnHTML("link", func(e *colly.HTMLElement) {
					if e.Attr("rel") == "canonical" {
						if i != 0 {
							target = e.Attr("href")[24:]
							fmt.Println("Target", target)
						} else {
							root = e.Attr("href")[24:]
							depthNode[root] = 0
							fmt.Println("Root", root)
						}
					}
				})
				if i == 0 {
					c1.Visit(baseLink + inputTitle)
				} else {
					c1.Visit(baseLink + searchTitle)
				}
				c1.Wait()
				<-limiter
			}(i)
		}
		wg.Wait()
	}
	if !invalidStart && !invalidTarget {
		iteration := 1

		for !alrFound {
			var wg1 sync.WaitGroup
			scrapedNode = make(map[string]bool)
			fmt.Println("Iterasi ke-", iteration)
			wg1.Add(1) // Menambahkan goroutine pertama ke wait group
			limiter <- 1
			IDS(root, target, iteration, &wg1)
			wg1.Wait() // Menunggu sampai semua goroutine selesai
			iteration += 1
		}
		close(limiter)

		fmt.Println("\nPage Checked: ", PageScraped)
		fmt.Println("Total Error Scrape: ", totalErrorScrape)

		Status = "OK"
		Err_msg = ""
		ResultDepth = depthNode[target]
	} else {
		ResultDepth = 0
		if invalidStart && invalidTarget {
			Err_msg = "Start Page and Target Page Not Found"
		} else if invalidStart {
			Err_msg = "Start Page Not Found"
		} else {
			Err_msg = "Target Page Not Found"
		}
		Status = "ERROR"
		return
	}

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
	_, existKey := childParentBool[parent]
	if !existKey || urlToTitle[parent] == rootTitle {
		return
	} else {
		if childParentBool[parent][root] {
			insertToSolution(parent, root)
		} else {
			for key, _ := range childParentBool[parent] {
				if depthNode[parent]-1 == depthNode[key] {
					insertToSolution(parent, key)
				}
			}
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

func GetSolutionAndConvertToJSON() {
	for parent, _ := range childParentBool[target] {
		fmt.Println("Solusi", parent, depthNode[parent], depthTarget-1)
		if depthNode[parent] == depthTarget-1 {
			insertToSolution(target, parent)
		}
	}
}

func ResetData() {
	childParentBool = make(map[string]map[string]bool)
	depthNode = make(map[string]int)
	baseLink = "https://en.wikipedia.org"
	limiter = make(chan int, 150)
	alrFound = false
	GraphSolusi = GraphView{Nodes: []Node{}, Edges: []Edge{}}
	PageScraped = 0
	urlToTitle = make(map[string]string)
	solutionParentChildBool = make(map[string]map[string]bool)
	insertedNodeToJSON = make(map[string]bool)
	isInit = false
	depthTarget = 999
	scrapedNode = make(map[string]bool)
	totalErrorScrape = 0
}
