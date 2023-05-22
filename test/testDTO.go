package main

type TestDTO struct {
	ID         string  `json:"id"`
	DataString string  `json:"dataString"`
	DataInt    int     `json:"dataInt"`
	DataBool   bool    `json:"dataBool"`
	DataFloat  float64 `json:"dataFloat"`
}
