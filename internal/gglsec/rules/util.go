package rules

import (
	"github.com/pvtgspot/gglsec/internal/gglsec"
	"github.com/xanzy/go-gitlab"
)

func getGroup(gid string, client *gitlab.Client, cache *gglsec.GitlabGroupRequestsCache) (gitlab.Group, error) {
	group, ok := cache.Get(gid)
	if !ok {
		group, _, err := client.Groups.GetGroup(gid, &gitlab.GetGroupOptions{})
		if err != nil {
			return gitlab.Group{}, err
		}
		cache.Set(gid, *group)
		return *group, nil
	}
	return group, nil
}
