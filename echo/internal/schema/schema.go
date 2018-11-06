package schema

type Meta struct {
	Version string `json:"Version"`
	PodName string `json:"PodName"`
	Headers map[string][]string `json:"Headers"`
}
type Product struct {
	Meta Meta
}
type Company struct {
	CompanyName string `json:"company_name"`
	Meta        Meta
	Products    []*Product
}
type echoResponse struct {
	Companyies []*Company `json:"elements"`
}

