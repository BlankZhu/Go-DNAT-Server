package entity

// Rule describe the NAT rule
type Rule struct {
	RuleID string `json:"ruleID"`
	CIDR   string `json:"CIDR"`        // IPv6 CIDR of incoming IP traffic
	Dest   string `json:"destination"` // destination IP of the DNAT
}
