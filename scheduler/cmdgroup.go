package scheduler

import (
	"strings"
)

type CmdGroup struct {
	isHandleNode bool
	ignoreCase   bool
	fullMatch    bool
	Keywords     []string
	BaseHandlers []HandleFunc
	SelfHandlers []HandleFunc
	subCmdGroups []*CmdGroup
	scheduler    *Scheduler
}

func (group *CmdGroup) Use(middlewares ...HandleFunc) {
	group.BaseHandlers = combineHandlers(group.BaseHandlers, middlewares...)
}

func (group *CmdGroup) Bind(keyword string, handlers ...HandleFunc) *CmdGroup {
	leafCmd := group.Group(keyword)
	leafCmd.isHandleNode = true
	leafCmd.SelfHandlers = handlers
	return leafCmd
}

func (group *CmdGroup) EqualsTo(keyword string, handlers ...HandleFunc) *CmdGroup {
	leafCmd := group.Bind(keyword, handlers...)
	leafCmd.fullMatch = true
	return leafCmd
}

func (group *CmdGroup) Group(keyword string) *CmdGroup {
	cmdGroup := &CmdGroup{
		isHandleNode: false,
		ignoreCase:   false,
		fullMatch:    false,
		Keywords:     []string{keyword},
		BaseHandlers: combineHandlers(group.BaseHandlers),
		subCmdGroups: []*CmdGroup{},
		scheduler:    group.scheduler,
	}
	group.subCmdGroups = append(group.subCmdGroups, cmdGroup)
	return cmdGroup
}

func (group *CmdGroup) Alias(alias ...string) *CmdGroup {
	group.Keywords = append(group.Keywords, alias...)
	return group
}

func (group *CmdGroup) IgnoreCase() *CmdGroup {
	group.ignoreCase = true
	return group
}

func (group *CmdGroup) dealIgnoreCase(s string) (string, []string) {
	if group.ignoreCase {
		dealKeywords := make([]string, len(group.Keywords))
		for i, _ := range group.Keywords {
			dealKeywords[i] = strings.ToLower(group.Keywords[i])
		}
		return strings.ToLower(s), dealKeywords
	} else {
		return s, group.Keywords
	}
}

func (group *CmdGroup) SearchHandlerChain(message string) (string, []HandleFunc, string, bool) {
	dealMessage, dealKeywords := group.dealIgnoreCase(message)
	if group.fullMatch {
		if prefix, has := whatPrefixEqualsTo(dealMessage, dealKeywords...); has {
			var keyword string
			var handlers []HandleFunc
			var pretreatedMessage string
			var inSubGroup bool
			for _, subGroup := range group.subCmdGroups {
				keyword, handlers, pretreatedMessage, inSubGroup = subGroup.SearchHandlerChain(prefix)
				if inSubGroup {
					return keyword, handlers, pretreatedMessage, true
				}
			}
			if group.isHandleNode {
				return prefix, combineHandlers(group.BaseHandlers, group.SelfHandlers...), strings.TrimSpace(message[len(prefix):]), true
			}
		}
		return "", nil, "", false
	} else if prefix, has := whatPrefixIn(dealMessage, dealKeywords...); has {
		var keyword string
		var handlers []HandleFunc
		var pretreatedMessage string
		var inSubGroup bool
		for _, subGroup := range group.subCmdGroups {
			keyword, handlers, pretreatedMessage, inSubGroup = subGroup.SearchHandlerChain(strings.TrimSpace(message[len(prefix):]))
			if inSubGroup {
				return keyword, handlers, pretreatedMessage, true
			}
		}
		if group.isHandleNode {
			return prefix, combineHandlers(group.BaseHandlers, group.SelfHandlers...), strings.TrimSpace(message[len(prefix):]), true
		} else {
			return "", nil, "", false
		}
	}
	return "", nil, "", false
}
