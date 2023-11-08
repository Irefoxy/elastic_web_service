package model

type Config struct {
	Secret         string `yaml:"Secret"`
	Address        string `yaml:"Address"`
	ElasticAddress string `yaml:"ElasticAddress"`
}

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
