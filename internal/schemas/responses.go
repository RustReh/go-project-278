package schemas

type PingResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
