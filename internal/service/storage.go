package service

type Word struct {
	Title       string `json:"title"`
	Translation string `json:"translation"`
}

type Report struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type noteReq struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type noteResp struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type qotdResp struct {
	Quote struct {
		Body string `json:"body"`
	} `json:"quote"`
}
