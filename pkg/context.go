package pkg

const (
	PluginName       = "ks-devops"
	PluginActionName = "action.kubesphere.io/devops"
)

const (
	ActionSetName        = "action.kubesphere.io/setName"
	ActionCreatePipeline = "action.kubesphere.io/create-pipeline"
)

// PluginContext is the context of this plugin
type PluginContext struct {
	Name      string
	Namespace string
}
