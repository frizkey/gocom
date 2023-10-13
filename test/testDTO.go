package main

type TestDTO struct {
	ID          string  `json:"id"`
	DataString  string  `json:"dataString"`
	DataString2 string  `json:"dataString2"`
	DataInt     int     `json:"dataInt"`
	DataBool    bool    `json:"dataBool"`
	DataFloat   float64 `json:"dataFloat"`
}
