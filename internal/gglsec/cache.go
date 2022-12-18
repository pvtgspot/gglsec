package gglsec

import (
	"github.com/xanzy/go-gitlab"
)

var (
	groupRequestsCache *GitlabGroupRequestsCache
)

type GitlabGroupRequestsCache struct {
	cache map[string]gitlab.Group
}

func GetGitlabGroupRequestsCache() *GitlabGroupRequestsCache {
	if groupRequestsCache == nil {
		groupRequestsCache = &GitlabGroupRequestsCache{
			cache: make(map[string]gitlab.Group),
		}
	}
	return groupRequestsCache
}

func (ggrc *GitlabGroupRequestsCache) Get(key string) (gitlab.Group, bool) {
	group, ok := ggrc.cache[key]
	return group, ok
}

func (ggrc *GitlabGroupRequestsCache) Set(key string, value gitlab.Group) {
	ggrc.cache[key] = value
}
