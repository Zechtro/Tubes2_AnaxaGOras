package ids 
 
import ( 
	"fmt" 
	"github.com/gocolly/colly"
	// "log"
) 
	
// Initiating maps	
	
//visitedNode["a"] = false jika tidak ada key "a"
var visitedNode = make(map[string]bool)

var nodeDepthVisitiedAt = make(map[string]int)

var solutionGraph = make(map[string][]string)

//wholeGraph["a"]["b"] = false jika "b" belum pernah menjadi child node dari "a"
var wholeGraph = make(map[string]map[string]bool)

var baseLink string = "https://en.wikipedia.org/wiki/"

func nodeVisited(title string, depth int){
	visitedNode[title] = true
	nodeDepthVisitiedAt[title] = depth
}

func BFS(title string) { 

	// pageToScrape := baseLink + "Google"

	c := colly.NewCollector()

	
	c.OnRequest(func(r *colly.Request) { 
		fmt.Println("Visiting: ", r.URL) 
		}) 
		
	c.OnError(func(_ *colly.Response, err error) { 
		fmt.Println("Something went wrong: ", err) 
	}) 
			
	c.OnResponse(func(r *colly.Response) { 
		fmt.Println("Page visited: ", r.Request.URL) 
	}) 
	
	c.OnHTML("a", func(e *colly.HTMLElement) { 
		// printing all URLs associated with the a links in the page
		if(isWiki(e.Attr("href"))){
			fmt.Println(e.Attr("href")) 
			fmt.Println(getArticleTitle(e.Attr("href")))
		}
	
	}) 
					
	c.OnScraped(func(r *colly.Response) { 
		fmt.Println(r.Request.URL, " scraped!") 
	})
					
	c.Visit("https://en.wikipedia.org/wiki/Neuroscience")
	fmt.Println(title)
}
				
				
// HELPER FUNCTIONS
func isWiki(link string) bool{
	if(len(link) <= 6){
		// fmt.Println("Length <= 6")
		return false
	}else if(link[:6] == "/wiki/"){
		// fmt.Println("Wiki link!")
		return true;
	}else{
		// fmt.Println("NOT Wiki link!")
		return false;
	}
}

func getArticleTitle(link string) string{
	return link[6:]
}


// SCRAPER FUNCTIONS

// // Base Web Link
// var baseLink string = "en.wikipedia.org/wiki/"

// // initializing the list of pages to scrape with an empty slice 
// var pagesToScrape []string 

// // the first pagination URL to scrape 
// pageToScrape := "https://scrapeme.live/shop/page/1/" 

// // initializing the list of pages discovered with a pageToScrape 
// pagesDiscovered := []string{ pageToScrape } 

// // current iteration 
// i := 1 
// // max pages to scrape 
// limit := 5 

// // initializing a Colly instance 
// c := colly.NewCollector() 

// c.OnRequest(func(r *colly.Request) { 
// 	fmt.Println("Visiting: ", r.URL) 
// }) 

// c.OnError(func(_ *colly.Response, err error) { 
// 	log.Println("Something went wrong: ", err) 
// }) 

// c.OnResponse(func(r *colly.Response) { 
// 	fmt.Println("Page visited: ", r.Request.URL) 
// }) 

// // iterating over the list of pagination links to implement the crawling logic 
// c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) { 
// 	// discovering a new page 
// 	newPaginationLink := e.Attr("href") 

// 	// if the page discovered is new 
// 	if !contains(pagesToScrape, newPaginationLink) { 
// 		// if the page discovered should be scraped 
// 		if !contains(pagesDiscovered, newPaginationLink) { 
// 			pagesToScrape = append(pagesToScrape, newPaginationLink) 
// 		} 
// 		pagesDiscovered = append(pagesDiscovered, newPaginationLink) 
// 	} 
// }) 

// c.OnHTML("li.product", func(e *colly.HTMLElement) { 
// 	// scraping logic... 
// }) 

// c.OnScraped(func(response *colly.Response) { 
// 	// until there is still a page to scrape 
// 	if len(pagesToScrape) != 0 && i < limit { 
// 		// getting the current page to scrape and removing it from the list 
// 		pageToScrape = pagesToScrape[0] 
// 		pagesToScrape = pagesToScrape[1:] 

// 		// incrementing the iteration counter 
// 		i++ 

// 		// visiting a new page 
// 		c.Visit(pageToScrape) 
// 	} 
// }) 

// // visiting the first page 
// c.Visit(pageToScrape) 

// // convert the data to CSV...
