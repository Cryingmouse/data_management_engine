package context

type HostContext struct {
	IP       string
	Username string
	Password string
}

type Pagination struct {
	Page     int
	PageSize int
}

type QueryFilter struct {
	// The attributes of DB table returned.
	Fields []string
	// The keyword for fuzzy search.
	Keyword map[string]string
	// Pagination
	Pagination *Pagination
	// The condition to filter the records by query.
	Conditions interface{}
}

type SystemInfo struct {
	ComputerName   string `json:"host_name"`
	Caption        string `json:"os_type"`
	OSArchitecture string `json:"os_arch"`
	Version        string `json:"version"`
	BuildNumber    string `json:"build_number"`
}

var SecurityKey = "MySecretForMagnascale!!!"
