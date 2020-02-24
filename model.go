package main

type AdConfiguration struct {
	ID        uint   `json:"id"`
	PartnerId uint   `json:"partnerid"`
	AdCode    string `json:"adcode"`
	Component string `json:"component"`
	Platform  string `json:"platform"`
	Location  string `json:"location"`
}
