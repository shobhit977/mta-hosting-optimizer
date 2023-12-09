package models

type IpConfig struct {
	Ip       string `json:"ip"`
	Hostname string `json:"hostname"`
	Active   bool   `json:"active"`
}

type HostnameResponse struct {
	Hostnames []string `json:"hostnames"`
}
