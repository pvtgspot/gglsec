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
			Meta: &gglsec.RuleMeta{
				Name:        "GL1002",
				Description: "Group visibility should be \"Private\"",
				Severity:    gglsec.SEVERITY_WARNING,
				EntityId:    gid,
				EntityType:  gglsec.ENTITY_TYPE_GROUP,
			},
			GitlabClient: client,
		},
	}
}

func (vlr *VisibilityLevelRule) Apply() *gglsec.RuleResult {
	const (
		resultWrongLevel    string = "Group visibility is \"%s\", but it should be \"%s\""
		resultRequiredLevel string = "Group visibility is \"%s\""
	)

	result := gglsec.NewRuleResult(vlr.Meta)
	cache := gglsec.GetGitlabGroupRequestsCache()

	group, err := getGroup(vlr.Meta.EntityId, vlr.GitlabClient, cache)
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
