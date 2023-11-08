package model

type Place struct {
	Name     string
	Address  string
	Phone    string
	Location geoData
}

type geoData struct {
	Lon float64
	Lat float64
}
