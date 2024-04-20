package ids 
 
import ( 
	"fmt" 
	"github.com/gocolly/colly"
	"strings"
) 
	
// Initiating maps	
	
//visitedNode["a"] = false jika tidak ada key "a"
// var pathStack []string
// var unvisitedPath = make(map[string][]string)

// var nodeDepthVisitiedAt = make(map[string]int)

// var solutionGraph = make(map[string][]string)

//wholeGraph["a"]["b"] = false jika "b" belum pernah menjadi child node dari "a"
// var wholeGraph = make(map[string]map[string]bool)


// func nodeVisited(title string, depth int){
	// 	visitedNode[title] = true
	// 	nodeDepthVisitiedAt[title] = depth
	// }
var childNparent = make(map[string][]string)
var depthNode = make(map[string]int)
var unvisitedPath []string

var alrFound bool = false

func IDS(inputTitle string, searchTitle string, iteration int) { 
	var baseLink string = "https://en.wikipedia.org/wiki/"
	pageToScrape := baseLink + inputTitle

	c := colly.NewCollector()


	// c.OnRequest(func(r *colly.Request) { 
	// 	fmt.Println("Visiting: ", r.URL) 
	// 	}) 
		
	// c.OnError(func(_ *colly.Response, err error) { 
	// 	fmt.Println("Something went wrong: ", err) 
	// }) 
			
	// c.OnResponse(func(r *colly.Response) { 
	// 	if iteration == 2 {
	// 		fmt.Println(iteration, " Page visited: ", r.Request.URL) 
	// 	}
	// }) 

	c.OnHTML("a", func(e *colly.HTMLElement) {
		if(isWiki(e.Attr("href"))){
			// fmt.Println(getArticleTitle(e.Attr("href")))
			var foundTitle string = getArticleTitle(e.Attr("href"))
			
			val, exists := depthNode[foundTitle]
			newVal := depthNode[inputTitle] + 1

			if (val == newVal) {
				childNparent[foundTitle] = append(childNparent[foundTitle], inputTitle)
			} else if (!exists || val > newVal) {
				depthNode[foundTitle] = newVal
				childNparent[foundTitle] = []string{inputTitle}
			}

			if (foundTitle == searchTitle) {
				alrFound = true
			} else if (iteration > 1) && (!exists || val > newVal) {
				IDS(foundTitle, searchTitle, iteration-1)
			} else if (iteration == 1) {
				unvisitedPath = append(unvisitedPath, foundTitle)
			}
		}
	}) 
					
	// c.OnScraped(func(r *colly.Response) { 
	// 	fmt.Println(r.Request.URL, " scraped!") 
	// })
					
	c.Visit(pageToScrape)
}
				
func MainIDS(inputTitle string, searchTitle string, iteration int) {
	childNparent[inputTitle] = []string{inputTitle}
	depthNode[inputTitle] = 1
	IDS(inputTitle, searchTitle, iteration+1)

	for (!alrFound) {

		for _, input := range unvisitedPath {
			IDS(input, searchTitle, iteration)
		}

		unvisitedPath = []string{}
	}

	var a = childNparent[searchTitle]
	fmt.Print(searchTitle, ", ")
	for (a[0] != inputTitle) {
		fmt.Print(a[0], ", ")
		a = childNparent[a[0]]
	}
	fmt.Print(a[0])
}
				
// HELPER FUNCTIONS
func isWiki(link string) bool{
	if(len(link) <= 6){
		// fmt.Println("Length <= 6")
		return false
	}else if(link[:6] == "/wiki/"){
		if (strings.ContainsRune(link[6:], ':') || link[6:] == "Main_Page") {
			return false
		} else {
			return true;
		}
	}else{
		// fmt.Println("NOT Wiki link!")
		return false;
	}
}

func getArticleTitle(link string) string{
	return link[6:]
}