package main

type InputData struct {
	DataArray []data `json:"DataArray"`
}

type data struct {
	ID    string `json:"ID"`
	Class string `json:"Class"`
	Value string `json:"Value"`
}
