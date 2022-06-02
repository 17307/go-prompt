package prompt

import "strings"

type CommandAndFlags struct {
	allSuggest  map[string]Suggest
	value       string
	flags       []string
	commands    []string
	globalFlags []string
	parser      func(document Document) []Suggest
	flagGroup   [][]string
}

func NewCommandAndFlags(value string, sugs []Suggest, AllOptionsType *map[string]*CommandAndFlags) *CommandAndFlags {
	// value 为当前command的名字
	// sugs 为当前command下全部的命令集合
	res := &CommandAndFlags{
		value:      value,
		commands:   make([]string, 0),
		flags:      make([]string, 0),
		parser:     nil,
		allSuggest: make(map[string]Suggest),
	}

	// 解析传入的sugs
	for _, s := range sugs {
		// 判断是否存在颜色，不存在则设置默认值
		// 判断如果以 - 开头，则是flag，否是则commands
		if strings.HasPrefix(s.Text, "-") {
			res.flags = append(res.flags, s.Text)
		} else {
			res.commands = append(res.commands, s.Text)
		}
		res.allSuggest[s.Text] = s
	}
	// 全局注册当前command
	(*AllOptionsType)[value] = res
	return res
}

func (c *CommandAndFlags) SetParser(f func(d Document) []Suggest) {
	c.parser = f
}
func (c *CommandAndFlags) SetGlobalFlags(args ...string) {
	c.globalFlags = args
}
func (c *CommandAndFlags) SetFlagGroup(args [][]string) {
	c.flagGroup = args
}
func (c *CommandAndFlags) GetSuggest(args ...string) []Suggest {
	if len(args) == 0 {
		// 为空则返回全部
		args = c.flags
	}
	res := make([]Suggest, 0)
	for _, a := range args {
		if r, ok := c.allSuggest[a]; ok {
			res = append(res, r)
		}
	}
	return res
}
func (c *CommandAndFlags) GetBestSuggest(args ...string) []string {
	for _, i := range args {
		for _, groups := range c.flagGroup {
			for _, g := range groups {
				if strings.HasSuffix(i, g) {
					return groups
				}
			}
		}
	}
	// 没找到最佳的，返回全部flagGroup
	var res = make([]string, 0)
	for _, groups := range c.flagGroup {
		res = append(res, groups...)
	}
	return res
}

func ExcludeFlags(args []string, flags []string) []string {
	// 过滤已经存在的flag
	res := make([]string, 0)
	for _, i := range flags {
		f := true
		for _, j := range args {
			if i == j {
				f = false
			}
		}
		if f {
			res = append(res, i)
		}
	}
	return res
}
