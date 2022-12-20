package rules

import (
	"github.com/pvtgspot/gglsec/internal/gglsec"
	"github.com/xanzy/go-gitlab"
)

type TwoFactorAuthRule struct {
	*gglsec.RuleMixin
}

func NewTwoFactorAuthRule(gid string, client *gitlab.Client) *TwoFactorAuthRule {
	return &TwoFactorAuthRule{
		&gglsec.RuleMixin{
			Meta: &gglsec.RuleMeta{
				Name:        "GL1003",
				Description: "Two factor authentication must be enabled",
				Severity:    gglsec.SEVERITY_WARNING,
				EntityId:    gid,
				EntityType:  gglsec.ENTITY_TYPE_GROUP,
			},
			GitlabClient: client,
		},
	}
}

func (tfa *TwoFactorAuthRule) Apply() *gglsec.RuleResult {
	const (
		resultTFADisabled string = "Two-factor authentication is disabled"
		resultTFAEnabled  string = "Two-factor authentication is enabled"
	)

	result := gglsec.NewRuleResult(tfa.Meta)
	cache := gglsec.GetGitlabGroupRequestsCache()

	group, err := getGroup(tfa.Meta.EntityId, tfa.GitlabClient, cache)
	if err != nil {
		result.Message = err.Error()
		return result
	}

	if group.RequireTwoFactorAuth {
		result.Status = true
		result.Message = resultTFAEnabled
		return result
	}

	result.Message = resultTFADisabled

	return result
}
