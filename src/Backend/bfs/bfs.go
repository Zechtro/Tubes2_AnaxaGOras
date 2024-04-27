package bfs

import (
	"fmt"
	"strings"
	"sync"
	. "web-scraper/structure"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var insertedNodeToJSON = make(map[string]bool)

var GraphSolusi = GraphView{Nodes: []Node{}, Edges: []Edge{}}

// checkedNode["a"] = false jika tidak ada key "a"=
var checkedNode = make(map[string]bool)

var TotalCheckedArticleTitle int = 0 // jml link artikel yang dilakukan perbandingan
var totalScrapedArticle int = 0      // jml artikel yang discrape
var totalTryToScrapeArticle int = 0
var totalErrorScrape int = 0

var depthOfNode = make(map[string]int)
var depthTarget int = 999

var urlToTitle = make(map[string]string)

var child_parent_bool = make(map[string]map[string]bool)
var solutionParentChildBool = make(map[string]map[string]bool)

var baseLink string = "https://en.wikipedia.org"

var root string

var target string

var currentDepth int = 1

var rootTitle string
var targetTitle string

// Untuk keperluan response
var Status string
var Err_msg string
var ResultDepth int

func insertToSolution(child string, parent string) {
	// Mendapatkan Judul Artikel Child
	_, cExist := urlToTitle[child]
	if !cExist {
		cc := colly.NewCollector(
			colly.Async(true),
		)

		cc.OnError(func(_ *colly.Response, err error) {
			urlToTitle[child] = child[6:]
		})

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

		cp.OnError(func(_ *colly.Response, err error) {
			urlToTitle[parent] = parent[6:]
		})

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
	_, existKey := child_parent_bool[child]
	if !existKey || urlToTitle[parent] == rootTitle {
		return
	} else {
		if child_parent_bool[parent][root] {
			insertToSolution(parent, root)
		} else {
			for key, _ := range child_parent_bool[parent] {
				if depthOfNode[parent]-1 == depthOfNode[key] {
					insertToSolution(parent, key)
				}
			}
		}
	}
}

func insertToJSON(child string, parent string) {
	if urlToTitle[child] == rootTitle {
		depthOfNode[child] = 0
	} else if urlToTitle[parent] == rootTitle {
		depthOfNode[parent] = 0
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
	fmt.Println("EDGES:", parent, child)
	GraphSolusi.Edges = append(GraphSolusi.Edges, Edge{
		From: urlToTitle[parent],
		To:   urlToTitle[child],
	})
}

func GetSolutionAndConvertToJSON() {
	for parent, _ := range child_parent_bool[target] {
		fmt.Println("Solusi", parent, depthOfNode[parent], depthTarget-1)
		if depthOfNode[parent] == depthTarget-1 {
			insertToSolution(target, parent)
		}
	}
}

func BFS(start_page []string, target_page string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var invalidStart bool = false
	var invalidTarget bool = false
	var nextBreadthList = []string{}
	var isFound bool = false
	var n int = len(start_page)
	limiter := make(chan int, 150)
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
	if !invalidStart && !invalidTarget {
		var wg1 sync.WaitGroup
		for i := 0; i < len(start_page); i++ {
			wg1.Add(1)
			limiter <- 1
			go func(i int) {
				defer wg1.Done()
				c := colly.NewCollector(
					colly.Async(true),
				)
				extensions.RandomUserAgent(c)

				var currentPage string = start_page[i]

				c.OnError(func(r *colly.Response, err error) {
					totalErrorScrape += 1
				})

				c.OnHTML("link", func(e *colly.HTMLElement) {
					if e.Attr("rel") == "canonical" {
						page := e.Attr("href")[24:]
						if page != currentPage && page != root {
							mu.Lock()
							if !checkedNode[page] && page != root {
								TotalCheckedArticleTitle += 1
								checkedNode[page] = true
							}
							_, existParentKey := child_parent_bool[page]
							if !existParentKey {
								child_parent_bool[page] = make(map[string]bool)
							}
							depthOfNode[page] = depthOfNode[currentPage]

							for keyParent, _ := range child_parent_bool[currentPage] {
								child_parent_bool[page][keyParent] = true
							}
							child_parent_bool[currentPage] = child_parent_bool[page]
							if page == target {
								isFound = true
								depthTarget = depthOfNode[currentPage]
								fmt.Println("REDIRECT:", depthTarget)
								e.Request.Abort()
							}
							mu.Unlock()
						}
					}
				})

				c.OnHTML("a", func(e *colly.HTMLElement) {
					mu.Lock()
					if e.Attr("class") != "mw-file-description" {
						// cek apakah link wikipedia
						if isWiki(e.Attr("href")) && e.Attr("href") != root && e.Attr("href") != "/wiki/Main_Page" {
							page := e.Attr("href")
							// cek jika sudah pernah dicek dan berada pada depth yang sama tetapi beda parent
							if checkedNode[page] && page != root {
								// _, existParentTemp := child_parent_bool[page]
								// if !existParentTemp {
								// 	child_parent_bool[page] = make(map[string]bool)
								// }
								if depthOfNode[page] == currentDepth && !child_parent_bool[page][currentPage] {
									child_parent_bool[page][currentPage] = true
									if page == target && currentDepth <= depthTarget {
										isFound = true
										depthTarget = currentDepth

										fmt.Println(currentPage, target)
									}
								}
							} else if !checkedNode[page] && page != root {
								TotalCheckedArticleTitle += 1

								checkedNode[page] = true
								depthOfNode[page] = currentDepth
								child_parent_bool[page] = make(map[string]bool)
								child_parent_bool[page][currentPage] = true

								if page == target && currentDepth <= depthTarget {
									// masukin ke solusi
									depthTarget = currentDepth
									fmt.Println(currentPage, target)
									isFound = true
								} else {
									nextBreadthList = append(nextBreadthList, page)
								}
							}
						}
						// }
					}
					mu.Unlock()
				})

				c.OnScraped(func(r *colly.Response) {
					totalScrapedArticle += 1
				})

				c.Visit(baseLink + currentPage)
				c.Wait()
				totalTryToScrapeArticle += 1
				<-limiter
			}(i)
			// single answer
			// if isFound {
			// 	fmt.Println("Depth: ", currentDepth)
			// 	fmt.Println("Total Checked Article: ", TotalCheckedArticleTitle)
			// 	fmt.Println("Total Scraped Article: ", totalScrapedArticle)
			// 	fmt.Println(urlSolutionGraph)
			// 	return
			// }
		}
		wg1.Wait()
		close(limiter)

		// multiple answer
		if isFound {
			fmt.Println("Depth: ", currentDepth)
			fmt.Println("Total Checked Article: ", TotalCheckedArticleTitle)
			fmt.Println("Total Scraped Article: ", totalScrapedArticle)
			fmt.Println("Total Visited Article: ", totalTryToScrapeArticle)
			fmt.Println("Total Error Scrape: ", totalErrorScrape)
			Status = "OK"
			Err_msg = ""
			ResultDepth = depthOfNode[target]
			return
		} else {
			currentDepth += 1
			BFS(nextBreadthList, target)
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

func ResetData() {
	insertedNodeToJSON = make(map[string]bool)
	GraphSolusi = GraphView{Nodes: []Node{}, Edges: []Edge{}}
	checkedNode = make(map[string]bool)
	TotalCheckedArticleTitle = 0
	totalScrapedArticle = 0
	totalTryToScrapeArticle = 0
	totalErrorScrape = 0
	depthOfNode = make(map[string]int)

	urlToTitle = make(map[string]string)

	child_parent_bool = make(map[string]map[string]bool)
	solutionParentChildBool = make(map[string]map[string]bool)
	currentDepth = 1
	depthTarget = 999
}
