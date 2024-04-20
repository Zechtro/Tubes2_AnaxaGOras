package bfs

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// KEYWORD: <span class="mw-page-title-main">

// Initiating maps

// checkedNode["a"] = false jika tidak ada key "a"=
var checkedNode = make(map[string]bool)

var totalCheckedArticleTitle int = 0 // jml artikel yang diperika
var totalScrapedArticle int = 0      // jml artikel yang dilalui
var totalTryToScrapeArticle int = 0

var depthOfNode = make(map[string]int)

// var urlSolutionGraph = make(map[string][]string)

var urlToTitle = make(map[string]string)

var SolutionGraph = make(map[string][]string)

var totalPath int = 1

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
	var nextBreathList = []string{}
	var isFound bool = false
	var n int = len(start_page)
	limiter := make(chan int, 150)
	if currentDepth == 1 && n == 1 {
		wg.Add(1)
		limiter <- 1
		go func() {
			defer wg.Done()
			root = start_page[0]
			depthOfNode[root] = 0
			checkedNode[root] = true
			c1 := colly.NewCollector()
			extensions.RandomUserAgent(c1)

			c1.OnError(func(_ *colly.Response, err error) {
				fmt.Println("Invalid target")
			})

			c1.OnRequest(func(r *colly.Request) {
				// fmt.Println("Visiting: ", r.URL)
				target = r.URL.String()[24:]
				fmt.Println(r.URL.String())
				fmt.Println(target)
			})

			c1.Visit(baseLink + "/wiki/" + target_page)
			<-limiter
		}()
	}
	for i := 0; i < len(start_page); i++ {
		wg.Add(1)
		limiter <- 1
		go func(i int) {
			defer wg.Done()
			c := colly.NewCollector(
			// colly.Async(true), // Enable asynchronous requests
			)
			// extensions.RandomUserAgent(c)

			var currentPage string = start_page[i]

			// c.OnRequest(func(r *colly.Request) {
			// 	fmt.Println("Visiting: ", r.URL)
			// })

			// c.OnError(func(_ *colly.Response, err error) {
			// 	fmt.Println("Something went wrong: ", err)
			// })

			// c.OnResponse(func(r *colly.Response) {
			// 	// fmt.Println("Page visited: ", r.Request.URL)
			// })

			c.OnHTML("a", func(e *colly.HTMLElement) {
				// cek apakah link wikipedia
				mu.Lock()
				if isWiki(e.Attr("href")) {
					// title := getArticleTitle(e.Attr("href"))
					page := e.Attr("href")
					// var isChildOfOtherParent bool = isIn(currentPage, parentOf[page])
					// cek jika sudah pernah dicek dan berada pada depth yang sama tetapi beda parent
					if checkedNode[page] {
						if depthOfNode[page] == currentDepth && !child_parent_bool[page][currentPage] {
							child_parent_bool[page][currentPage] = true
							if page == target {
								// masukin ke solusi
								// insertSolution1(page, currentPage)
								insertSolution(page, currentPage)
								isFound = true

								fmt.Println(currentPage, target)
							}
						}
					} else if !checkedNode[page] {
						totalCheckedArticleTitle += 1

						checkedNode[page] = true
						depthOfNode[page] = currentDepth
						child_parent_bool[page] = make(map[string]bool)
						child_parent_bool[page][currentPage] = true

						nextBreathList = append(nextBreathList, page)
						if page == target {
							// masukin ke solusi
							// insertSolution1(page, currentPage)
							insertSolution(page, currentPage)
							isFound = true

							// fmt.Println(currentPage, target)
						}
					}
					// fmt.Println(e.Attr("href"))
					// fmt.Println(getArticleTitle(e.Attr("href")))
				}
				mu.Unlock()
			})

			c.OnScraped(func(r *colly.Response) {
				// fmt.Println(r.Request.URL, " scraped!")
				totalScrapedArticle += 1
			})

			// currentDepth += 1	masih gayakin kapan harus ini

			c.Visit(baseLink + currentPage)
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
	wg.Wait()
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
		return
	} else {
		currentDepth += 1
		BFS(nextBreathList, target)
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

// func getArticleTitle(link string) string {
// 	return link[6:]
// }
