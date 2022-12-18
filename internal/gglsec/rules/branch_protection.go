package rules

import (
	"fmt"

	"github.com/pvtgspot/gglsec/internal/gglsec"
	"github.com/xanzy/go-gitlab"
)

const recomendedBranchProtection int = 2

type GroupBranchProtectionRule struct {
	*gglsec.RuleMixin
}

func NewGroupBranchProtectionRule(gid string, client *gitlab.Client) *GroupBranchProtectionRule {
	return &GroupBranchProtectionRule{
		&gglsec.RuleMixin{
			EntityId: gid,
			Meta: &gglsec.RuleMeta{
				Name:        "GL1001",
				Description: "Protection of the main branch in the group settings should be set to \"Fully protected\"",
				Severity:    gglsec.SEVERITY_WARNING,
			},
			GitlabClient: client,
		},
	}
}

func (bpc *GroupBranchProtectionRule) Apply() *gglsec.RuleResult {
	const (
		resultNoMessage          string = "No message"
		resultWrongProtection    string = "Default branch protection is set to \"%s\""
		resultRequiredProtection string = "Default branch protection is set to \"%s\" not \"%s\""
	)

	result := &gglsec.RuleResult{
		Meta:    bpc.Meta,
		Status:  false,
		Message: resultNoMessage,
	}

	cache := gglsec.GetGitlabGroupRequestsCache()

	group, err := getGroup(bpc.EntityId, bpc.GitlabClient, cache)
	if err != nil {
		result.Message = err.Error()
		return result
	}

	defaultBranchProtection := group.DefaultBranchProtection
	if defaultBranchProtection == recomendedBranchProtection {
		result.Status = true
		result.Message = fmt.Sprintf(
			resultWrongProtection,
			mapBranchProtectionToString(defaultBranchProtection),
		)
		return result
	}

	result.Message = fmt.Sprintf(
		resultRequiredProtection,
		mapBranchProtectionToString(defaultBranchProtection),
		mapBranchProtectionToString(recomendedBranchProtection),
	)

	return result
}

func mapBranchProtectionToString(mode int) string {
	var res string
	switch mode {
	case 0:
		res = "No protection"
	case 1:
		res = "Partial protection"
	case 2:
		res = "Full protection"
	case 3:
		res = "Protected against pushes"
	}
	return res
}
