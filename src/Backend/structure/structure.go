package structure

var DepthColor = map[int]string{
	14: "#7BC9FF",
	13: "#90EE90",
	12: "#401F71",
	11: "#824D74",
	10: "#C65BCF",
	9:  "#F27BBD",
	8:  "#FDAF7B",
	7:  "#F9DBBB",
	6:  "#CECECE",
	5:  "#F7F9FB",
	4:  "#10439F",
	3:  "#FCE94F",
	2:  "#874CCC",
	1:  "#2ECC71",
	0:  "#39CCCC",
}

type Node struct {
	Id           string `json:"id"`
	TitleArticle string `json:"label"`
	UrlArticle   string `json:"title"`
	Shape        string `json:"shape"`
	Size         int    `json:"size"`
	Color        Color  `json:"color"`
	Font         Font   `json:"font"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Color struct {
	Border     string `json:"border"`
	Background string `json:"background"`
}

type Font struct {
	Color string `json:"color"`
	Size  int    `json:"size"`
}

type GraphView struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}
