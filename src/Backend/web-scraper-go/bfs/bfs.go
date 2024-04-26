package bfs

import (
	"fmt"
	"sync"
	. "web-scraper/structure"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// KEYWORD : mw-file-description

var insertedNodeToJSON = make(map[string]bool)

var GraphSolusi = GraphView{Nodes: []Node{}, Edges: []Edge{}}

// checkedNode["a"] = false jika tidak ada key "a"=
var checkedNode = make(map[string]bool)
var scrapedNode = make(map[string]bool)

var totalCheckedArticleTitle int = 0 // jml artikel yang diperika
var totalScrapedArticle int = 0      // jml artikel yang dilalui
var totalTryToScrapeArticle int = 0

var depthOfNode = make(map[string]int)

// var urlSolutionGraph = make(map[string][]string)

var urlToTitle = make(map[string]string)

var SolutionGraph = make(map[string][]string)

var totalPath int = 0

// var queue = make(map[int]map[string]int)

// var parentOf = make(map[string][]string)

var child_parent_bool = make(map[string]map[string]bool)
var solutionParentChildBool = make(map[string]map[string]bool)

// Keperluan JSON
var titleID = make(map[string]int)

var baseLink string = "https://en.wikipedia.org"

var root string

var target string

var currentDepth int = 1

var thresholdRetry int = 2
var countRetry int = 0

var rootTitle string
var targetTitle string

// Untuk keperluan response
var Status string
var Err_msg string
var ResultDepth int

// func insertSolution1(child string, parent string) {
// 	_, existChild := solutionParentChildBool[parent]
// 	if !existChild {
// 		solutionParentChildBool[parent] = make(map[string]bool)
// 	}
// 	if !solutionParentChildBool[parent][child] {
// 		solutionParentChildBool[parent][child] = true
// 		urlSolutionGraph[parent] = append(urlSolutionGraph[parent], child)
// 	} else {
// 		return
// 	}
// 	// cek parentnya dari parent
// 	_, existKey := child_parent_bool[parent]
// 	if !existKey {
// 		// berarti udh root
// 		return
// 	} else {
// 		for key, _ := range child_parent_bool[parent] {
// 			insertSolution1(parent, key)
// 		}
// 	}
// }

func insertSolution(child string, parent string) {
	// fmt.Println("SOL", child, parent)

	_, cExist := urlToTitle[child]
	if !cExist {
		cc := colly.NewCollector()

		// cc.OnHTML("span", func(e *colly.HTMLElement) {
		// 	className := e.Attr("class")
		// 	if className == "mw-page-title-main" {
		// 		urlToTitle[child] = e.Text
		// 		fmt.Println("INSERT SOLUTION")
		// 		fmt.Println(child, e.Text)
		// 	}
		// })
		cc.OnHTML("#firstHeading", func(e *colly.HTMLElement) {
			if e.ChildText(".mw-page-title-main") != "" {
				urlToTitle[child] = e.ChildText(".mw-page-title-main")
			} else if e.Text != "" {
				// fmt.Println("SPECIAL CASE")
				// fmt.Println(child, e.Text)
				// fmt.Println("INSERT SOLUTION")
				urlToTitle[child] = e.Text
			}
		})

		cc.Visit(baseLink + child)
	}
	_, pExist := urlToTitle[parent]
	if !pExist {
		cp := colly.NewCollector()

		// cp.OnHTML("span", func(e *colly.HTMLElement) {
		// 	className := e.Attr("class")
		// 	if className == "mw-page-title-main" {
		// 		urlToTitle[parent] = e.ChildText("i")
		// 		fmt.Println("INSERT SOLUTION")
		// 		fmt.Println(parent, e.ChildText("i"))
		// 	}
		// })
		cp.OnHTML("#firstHeading", func(e *colly.HTMLElement) {
			if e.ChildText(".mw-page-title-main") != "" {
				urlToTitle[parent] = e.ChildText(".mw-page-title-main")
			} else if e.Text != "" {
				// fmt.Println("SPECIAL CASE")
				// fmt.Println(parent, e.Text)
				// fmt.Println("INSERT SOLUTION")
				urlToTitle[parent] = e.Text
			}
		})

		cp.Visit(baseLink + parent)
	}

	_, existChild := solutionParentChildBool[urlToTitle[parent]]
	if !existChild {
		solutionParentChildBool[urlToTitle[parent]] = make(map[string]bool)
	}
	if !solutionParentChildBool[urlToTitle[parent]][urlToTitle[child]] {
		if child == target {
			totalPath += 1
		}
		solutionParentChildBool[urlToTitle[parent]][urlToTitle[child]] = true
		SolutionGraph[urlToTitle[parent]] = append(SolutionGraph[urlToTitle[parent]], urlToTitle[child])
	} else {
		return
	}
	// cek parentnya dari parent
	_, existKey := child_parent_bool[parent]
	if !existKey {
		// berarti udh root
		return
	} else {
		for key, _ := range child_parent_bool[parent] {
			insertSolution(parent, key)
		}
	}

}
func insertToSolutionZ(child string, parent string) {
	// Mendapatkan Judul Artikel Child
	_, cExist := urlToTitle[child]
	if !cExist {
		cc := colly.NewCollector(
			colly.Async(true),
		)

		// cc.OnHTML("span", func(e *colly.HTMLElement) {
		// 	className := e.Attr("class")
		// 	if className == "mw-page-title-main" {
		// 		urlToTitle[child] = e.Text
		// 		fmt.Println("INSERT SOLUTION")
		// 		fmt.Println(child, e.Text)
		// 	}
		// })
		cc.OnHTML("#firstHeading", func(e *colly.HTMLElement) {
			if e.ChildText(".mw-page-title-main") != "" {
				urlToTitle[child] = e.ChildText(".mw-page-title-main")
			} else if e.Text != "" {
				// fmt.Println("SPECIAL CASE")
				// fmt.Println(child, e.Text)
				// fmt.Println("INSERT SOLUTION")
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

		// cp.OnHTML("span", func(e *colly.HTMLElement) {
		// 	className := e.Attr("class")
		// 	if className == "mw-page-title-main" {
		// 		urlToTitle[parent] = e.ChildText("i")
		// 		fmt.Println("INSERT SOLUTION")
		// 		fmt.Println(parent, e.ChildText("i"))
		// 	}
		// })
		cp.OnHTML("#firstHeading", func(e *colly.HTMLElement) {
			if e.ChildText(".mw-page-title-main") != "" {
				urlToTitle[parent] = e.ChildText(".mw-page-title-main")
			} else if e.Text != "" {
				// fmt.Println("SPECIAL CASE")
				// fmt.Println(parent, e.Text)
				// fmt.Println("INSERT SOLUTION")
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
		if child == target {
			totalPath += 1
		}
		solutionParentChildBool[urlToTitle[parent]][urlToTitle[child]] = true
		// Masukkan ke struktur JSON
		insertToJSON(child, parent)
	} else {
		return
	}

	// cek parentnya dari parent
	_, existKey := child_parent_bool[child]
	if !existKey || urlToTitle[parent] == rootTitle {
		return
	} else {
		for key, _ := range child_parent_bool[parent] {
			insertToSolutionZ(parent, key)
		}
	}
}

func insertToJSON(child string, parent string) {
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
				Border:     DepthColor[depthOfNode[child]],
				Background: DepthColor[depthOfNode[child]],
			},
			Font: Font{
				Color: DepthColor[depthOfNode[child]],
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
				Border:     DepthColor[depthOfNode[parent]],
				Background: DepthColor[depthOfNode[parent]],
			},
			Font: Font{
				Color: DepthColor[depthOfNode[parent]],
				Size:  font_size,
			},
		})
	}
	GraphSolusi.Edges = append(GraphSolusi.Edges, Edge{
		From: urlToTitle[parent],
		To:   urlToTitle[child],
	})
}

// func countTotalSolution() {
// 	var wg sync.WaitGroup
// 	for _, val := range SolutionGraph {
// 		wg.Add(1)
// 		go func(val []string) {
// 			defer wg.Done()
// 			// n := len(val)
// 			// fmt.Println(n)
// 			totalPath *= len(val)
// 		}(val)
// 	}
// 	wg.Wait()
// }

func BFS(start_page []string, target_page string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var invalidStart bool = false
	var invalidTarget bool = false
	var nextBreathList = []string{}
	var isFound bool = false
	var n int = len(start_page)
	limiter := make(chan int, 500)
	// Inisialisasi info root dan target
	if currentDepth == 1 && n == 1 {
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
						// fmt.Println("SPECIAL CASE")
						// fmt.Println(child, e.Text)
						// fmt.Println("INSERT SOLUTION")
						if i == 0 {
							rootTitle = e.Text
						} else {
							targetTitle = e.Text
						}
					}
				})

				c1.OnHTML("link", func(e *colly.HTMLElement) {
					// fmt.Println("LALA", e.Attr("title"))
					if e.Attr("rel") == "canonical" {
						// fmt.Println("++++++++++++++++++++++", e.Attr("href")[24:])
						if i != 0 {
							target = e.Attr("href")[24:]
							// checkedNode[target] = true
							fmt.Println("Target", target)
						} else {
							root = e.Attr("href")[24:]
							depthOfNode[root] = 0
							checkedNode[root] = true
							fmt.Println("Root", root)
						}
					}
				})
				if i == 0 {
					c1.Visit(baseLink + start_page[0])
				} else {
					c1.Visit(baseLink + "/wiki/" + target_page)
				}
				c1.Wait()
				<-limiter
			}(i)
		}
		wg.Wait()
	}
	// fmt.Println("SOMETHING")
	if !invalidStart && !invalidTarget {
		var wg1 sync.WaitGroup
		for i := 0; i < len(start_page); i++ {
			wg1.Add(1)
			limiter <- 1
			go func(i int) {
				defer wg1.Done()
				c := colly.NewCollector(
					colly.Async(true), // Enable asynchronous requests
				)
				// extensions.RandomUserAgent(c)

				var currentPage string = start_page[i]
				// fmt.Println(currentPage)

				// var isAbort bool = false
				// c.OnRequest(func(r *colly.Request) {
				// 	mu.Lock()
				// 	if scrapedNode[r.URL.String()] {
				// 		fmt.Println("ALR SCRAPED, ABORTING", currentPage)
				// 		isAbort = true
				// 		r.Abort()
				// 	} else {
				// 		// fmt.Println("SCRAPING...")
				// 		scrapedNode[r.URL.String()] = true
				// 	}
				// 	mu.Unlock()
				// })

				// c.OnHTML("#firstHeading", func(e *colly.HTMLElement) {
				// 	if e.ChildText(".mw-page-title-main") != "" {
				// 		mu.Lock()
				// 		if !scrapedNode[e.ChildText(".mw-page-title-main")] {
				// 			scrapedNode[e.ChildText(".mw-page-title-main")] = true
				// 			urlToTitle[currentPage] = e.ChildText(".mw-page-title-main")
				// 			mu.Unlock()
				// 		} else {
				// 			// fmt.Println("ALR SCRAPED, ABORTING", currentPage)
				// 			isAbort = true
				// 			e.Request.Abort()
				// 			mu.Unlock()
				// 		}

				// 	} else if e.Text != "" {
				// 		// fmt.Println("SPECIAL CASE")
				// 		// fmt.Println(parent, e.Text)
				// 		// fmt.Println("INSERT SOLUTION")
				// 		mu.Lock()
				// 		if !scrapedNode[e.Text] {
				// 			scrapedNode[e.Text] = true
				// 			urlToTitle[currentPage] = e.Text
				// 			mu.Unlock()
				// 		} else {
				// 			// fmt.Println("ALR SCRAPED, ABORTING", currentPage)
				// 			isAbort = true
				// 			e.Request.Abort()
				// 			mu.Unlock()
				// 		}
				// 	}
				// })

				// c.OnError(func(r *colly.Response, err error) {
				// 	mu.Lock()
				// 	if err.Error() != "Too Many Requests" {
				// 		if countRetry < thresholdRetry {
				// 			countRetry += 1
				// 			fmt.Println("RETRY", countRetry, err)
				// 			mu.Unlock()
				// 			r.Request.Retry()
				// 		} else {
				// 			countRetry = 0
				// 		}
				// 	}
				// 	if err == nil {
				// 		countRetry = 0
				// 	}
				// 	mu.Unlock()
				// })

				// c.OnResponse(func(r *colly.Response) {
				// 	// fmt.Println("Page visited: ", r.Request.URL)
				// })

				c.OnHTML("a", func(e *colly.HTMLElement) {
					mu.Lock()
					if e.Attr("class") != "mw-file-description" {
						// cek apakah link wikipedia
						if isWiki(e.Attr("href")) {
							// title := getArticleTitle(e.Attr("href"))
							page := e.Attr("href")

							// var isChildOfOtherParent bool = isIn(currentPage, parentOf[page])
							// cek jika sudah pernah dicek dan berada pada depth yang sama tetapi beda parent
							if checkedNode[page] && page != root {
								if depthOfNode[page] == currentDepth && !child_parent_bool[page][currentPage] {
									child_parent_bool[page][currentPage] = true
									if page == target {
										// masukin ke solusi
										// insertSolution1(page, currentPage)
										// insertSolution(page, currentPage)
										insertToSolutionZ(page, currentPage)
										isFound = true

										fmt.Println(currentPage, target)
									}
								}
							} else if !checkedNode[page] && page != root {
								totalCheckedArticleTitle += 1

								checkedNode[page] = true
								depthOfNode[page] = currentDepth
								child_parent_bool[page] = make(map[string]bool)
								child_parent_bool[page][currentPage] = true

								nextBreathList = append(nextBreathList, page)
								if page == target {
									// masukin ke solusi
									// insertSolution1(page, currentPage)
									// insertSolution(page, currentPage)
									insertToSolutionZ(page, currentPage)
									isFound = true

									// fmt.Println(currentPage, target)
								}
							}
							// fmt.Println(e.Attr("href"))
							// fmt.Println(getArticleTitle(e.Attr("href")))
						}
					}
					mu.Unlock()
				})

				c.OnScraped(func(r *colly.Response) {
					// fmt.Println(r.Request.URL, " scraped!")
					// mu.Lock()
					// if !isAbort {
					totalScrapedArticle += 1
					// } else {
					// 	isAbort = false
					// }
					// mu.Unlock()
				})

				// currentDepth += 1	masih gayakin kapan harus ini

				c.Visit(baseLink + currentPage)
				c.Wait()
				totalTryToScrapeArticle += 1
				// fmt.Println(currentPage)
				<-limiter
			}(i)
			// single answer
			// if isFound {
			// 	fmt.Println("Depth: ", currentDepth)
			// 	fmt.Println("Total Checked Article: ", totalCheckedArticleTitle)
			// 	fmt.Println("Total Scraped Article: ", totalScrapedArticle)
			// 	fmt.Println(urlSolutionGraph)
			// 	return
			// }
		}
		// limiter <- 1
		// go func() {
		wg1.Wait()
		// <-limiter
		close(limiter)
		// }()
		// single answer
		// currentDepth += 1
		// BFS(nextBreathList, target)

		// multiple answer
		if isFound {
			fmt.Println("Depth: ", currentDepth)
			fmt.Println("Total Checked Article: ", totalCheckedArticleTitle)
			fmt.Println("Total Scraped Article: ", totalScrapedArticle)
			fmt.Println("Total Visited Article: ", totalTryToScrapeArticle)
			// countTotalSolution()
			fmt.Println("Total Solution: ", totalPath)
			// fmt.Println(urlSolutionGraph)
			// for key, val := range SolutionGraph {
			// 	wg.Add(1)
			// 	go func(key string, val []string) {
			// 		defer wg.Done()
			// 		fmt.Println(key, val)
			// 	}(key, val)
			// }
			// wg.Wait()
			// j, err := json.Marshal(SolutionGraph)
			// fmt.Println(string(j), err)
			Status = "OK"
			Err_msg = ""
			ResultDepth = depthOfNode[target]
			return
		} else {
			currentDepth += 1
			BFS(nextBreathList, target)
		}
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
		// fmt.Println("Length <= 6")
		return false
	} else if link[:6] == "/wiki/" {
		// fmt.Println("Wiki link!")
		return true
	} else {
		// fmt.Println("NOT Wiki link!")
		return false
	}
}

func ResetData() {
	insertedNodeToJSON = make(map[string]bool)
	GraphSolusi = GraphView{Nodes: []Node{}, Edges: []Edge{}}
	checkedNode = make(map[string]bool)
	scrapedNode = make(map[string]bool)
	totalCheckedArticleTitle = 0
	totalScrapedArticle = 0
	totalTryToScrapeArticle = 0
	depthOfNode = make(map[string]int)

	urlToTitle = make(map[string]string)

	SolutionGraph = make(map[string][]string)

	totalPath = 0

	child_parent_bool = make(map[string]map[string]bool)
	solutionParentChildBool = make(map[string]map[string]bool)

	currentDepth = 1
}
