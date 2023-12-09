package models

type IpConfig struct {
	Ip       string `json:"ip"`
	Hostname string `json:"hostname"`
	Active   bool   `json:"active"`
}

type ServerResponse struct {
	Hostnames []string `json:"hostnames"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"statusCode"`
}
