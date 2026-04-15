package tools

import "charm.land/fantasy"

func GetTools() []fantasy.AgentTool {
	tools := []fantasy.AgentTool{}

	tools = append(
		tools,
		NewLsTool(),
		NewGlobTool(),
	)

	return tools
}
