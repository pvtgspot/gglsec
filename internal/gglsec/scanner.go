package gglsec

type Scanner interface {
	Run() []*RuleResult
}

type GitlabConfigScanner struct {
	rl *RuleList
}

func NewGitlabConfigScanner(ruleList *RuleList) *GitlabConfigScanner {
	return &GitlabConfigScanner{
		rl: ruleList,
	}
}

func (gcs *GitlabConfigScanner) Run() *RuleResults {
	results := NewRuleResults()
	for _, rule := range gcs.rl.Get() {
		results.Append(rule.Apply())
	}
	return results
}
