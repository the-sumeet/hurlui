package main

type HurlReport []HurlSession

// One session per file
type HurlSession struct {
	Cookies  []HurlCookie `json:"cookies"`
	Entries  []HurlEntry  `json:"entries"`
	Filename string       `json:"filename"`
	Success  bool         `json:"success"`
	Time     int          `json:"time"`
}

// Hurl file is a list of entries
// https://hurl.dev/docs/entry.html
type HurlEntry struct {
	Asserts  []interface{} `json:"asserts"` // Can be various assert types
	Calls    []HurlCall    `json:"calls"`
	Captures []interface{} `json:"captures"` // Can be various capture types
	CurlCmd  string        `json:"curl_cmd"`
	Index    int           `json:"index"`
	Line     int           `json:"line"`
	Time     int           `json:"time"`
}

type HurlCall struct {
	Request  HurlRequest  `json:"request"`
	Response HurlResponse `json:"response"`
	Timings  HurlTimings  `json:"timings"`
}

type HurlRequest struct {
	Cookies     []HurlCookie     `json:"cookies"`
	Headers     []HurlHeader     `json:"headers"`
	Method      string           `json:"method"`
	QueryString []HurlQueryParam `json:"query_string"`
	URL         string           `json:"url"`
}

type HurlResponse struct {
	BodyPath    string           `json:"body"`
	Body        string           `json:"bodyContent"`
	Certificate *HurlCertificate `json:"certificate,omitempty"`
	Cookies     []HurlCookie     `json:"cookies"`
	Headers     []HurlHeader     `json:"headers"`
	HTTPVersion string           `json:"http_version"`
	Status      int              `json:"status"`
}

type HurlTimings struct {
	AppConnect    int    `json:"app_connect"`
	BeginCall     string `json:"begin_call"`
	Connect       int    `json:"connect"`
	EndCall       string `json:"end_call"`
	NameLookup    int    `json:"name_lookup"`
	PreTransfer   int    `json:"pre_transfer"`
	StartTransfer int    `json:"start_transfer"`
	Total         int    `json:"total"`
}

type HurlHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HurlCookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain,omitempty"`
	Path     string `json:"path,omitempty"`
	Expires  string `json:"expires,omitempty"`
	MaxAge   int    `json:"max_age,omitempty"`
	HttpOnly bool   `json:"http_only,omitempty"`
	Secure   bool   `json:"secure,omitempty"`
	SameSite string `json:"same_site,omitempty"`
}

type HurlQueryParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HurlCertificate struct {
	ExpireDate   string `json:"expire_date"`
	Issuer       string `json:"issuer"`
	SerialNumber string `json:"serial_number"`
	StartDate    string `json:"start_date"`
	Subject      string `json:"subject"`
}
