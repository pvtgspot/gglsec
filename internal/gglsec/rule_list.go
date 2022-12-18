package gglsec

type RuleList struct {
	rules []Rule
}

func NewRuleList(rules ...Rule) *RuleList {
	rl := &RuleList{
		rules: make([]Rule, 0),
	}
	for _, rule := range rules {
		rl.Append(rule)
	}
	return rl
}

func (rl *RuleList) Get() []Rule {
	return rl.rules
}

func (rl *RuleList) Append(rules ...Rule) {
	rl.rules = append(rl.rules, rules...)
}
