package model

type Todo struct {
	Title string `json:"plain_text"`
}

type RichText struct {
	Todo    []Todo `json:"rich_text"`
	Checked bool   `json:"checked"`
}

type Block struct {
	Id       string   `json:"id"`
	RichText RichText `json:"to_do"`
}

type BlockResponse struct {
	Results []Block `json:"results"`
}
