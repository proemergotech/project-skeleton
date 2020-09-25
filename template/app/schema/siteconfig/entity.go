package siteconfig

type SiteConfig struct {
	SiteGroups    map[string]*SiteGroup    `json:"site_groups"`
	Sites         map[string]*Site         `json:"sites"`
	PaymentGroups map[string]*PaymentGroup `json:"payment_groups"`
}

type SiteGroup struct {
	Code string `json:"code"`
}

type Site struct {
	Name          string `json:"name"`
	SiteGroupCode string `json:"site_group_code"`
	PaymentGroup  string `json:"payment_group"`
}

type PaymentGroup struct {
	Packages      []Package `json:"packages"`
	SiteGroupCode string    `json:"site_group_code"`
}

type Package struct {
	ID string `json:"id"`
}
