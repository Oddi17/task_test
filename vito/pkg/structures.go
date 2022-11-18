package pkg

type Transaction struct {
	ID       int `json:"id_user"`
	ID_ORDER int `json:"id_order"`
}

type Balances struct {
	ID      int     `json:"id"`
	ACCOUNT float32 `json:"balance"`
}

type Reservation struct {
	ID         int `json:"id_user"`
	ID_SERVICE int `json:"id_service"`
}
