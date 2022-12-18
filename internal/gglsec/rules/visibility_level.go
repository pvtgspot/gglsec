package rules

import (
	"fmt"

	"github.com/pvtgspot/gglsec/internal/gglsec"
	"github.com/xanzy/go-gitlab"
)

type VisibilityLevelRule struct {
	*gglsec.RuleMixin
}

func NewVisibilityLevelRule(gid string, client *gitlab.Client) *VisibilityLevelRule {
	return &VisibilityLevelRule{
		&gglsec.RuleMixin{
			EntityId: gid,
			Meta: &gglsec.RuleMeta{
				Name:        "GL1002",
				Description: "Group visibility should be \"Private\"",
				Severity:    gglsec.SEVERITY_WARNING,
			},
			GitlabClient: client,
		},
	}
}

func (vlr *VisibilityLevelRule) Apply() *gglsec.RuleResult {
	const (
		resultNoMessage     string = "No message"
		resultWrongLevel    string = "Group visibility is \"%s\", but it should be \"%s\""
		resultRequiredLevel string = "Group visibility is \"%s\""
	)

	result := &gglsec.RuleResult{
		Meta:    vlr.Meta,
		Status:  false,
		Message: resultNoMessage,
	}

	cache := gglsec.GetGitlabGroupRequestsCache()

	group, err := getGroup(vlr.EntityId, vlr.GitlabClient, cache)
	if err != nil {
		result.Message = err.Error()
		return result
	}

	visibilityLevel := group.Visibility
	if visibilityLevel != gitlab.PrivateVisibility {
		result.Message = fmt.Sprintf(resultWrongLevel, visibilityLevel, gitlab.PrivateVisibility)
		return result
	}

	result.Status = true
	result.Message = fmt.Sprintf(resultRequiredLevel, gitlab.PrivateVisibility)

	return result
}
