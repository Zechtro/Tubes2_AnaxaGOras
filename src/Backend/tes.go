package main 
 
import ( 
	"fmt" 
) 

func main() { 
	var a = make(map[string][]string)
	fmt.Println(a)
	fmt.Println(len(a["xoxo"]))
	a["xoxo"] = append(a["xoxo"],"LALALA")
	fmt.Println(len(a["xoxo"]))
	a["xoxo"] = append(a["xoxo"],"LALALA")
	fmt.Println(len(a["xoxo"]))

	var i int

}
