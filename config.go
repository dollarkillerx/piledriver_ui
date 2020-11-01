package main

type config struct {
	Address    string `json:"address"`
	UserID     string `json:"user_id"`
	Password   string `json:"password"`
	Socks5Addr string `json:"socks_5_addr"`
}
