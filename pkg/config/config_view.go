package config

import (
	"github.com/linuxsuren/octant-ks-devops/pkg"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// ConfigHandler is the handler of pipeline
type ConfigHandler struct {
	Context *pkg.PluginContext
}

func (h *ConfigHandler) Dashboard(request service.Request) (response component.ContentResponse, e error) {
	response = *component.NewContentResponse(component.TitleFromString("Config"))

	var data map[string]map[string]string
	data = getConfigMap(request)

	response.Add(component.NewSummary("Summary", component.SummarySection{
		Header:  "JWTSecret",
		Content: component.NewText(data["authentication"]["jwtSecret"]),
	}, component.SummarySection{
		Header:  "Jenkins Token",
		Content: component.NewText(data["devops"]["password"]),
	}))
	return
}

func getConfigMap(request service.Request) (result map[string]map[string]string) {
	result = map[string]map[string]string{}
	result["devops"] = make(map[string]string, 0)
	result["authentication"] = make(map[string]string, 0)
	if data, err := request.DashboardClient().Get(request.Context(), store.Key{
		Namespace:     "kubesphere-devops-system",
		APIVersion:    "v1",
		Kind:          "ConfigMap",
		Name:          "devops-config",
	}); err == nil {
		configMap := &v1.ConfigMap{}
		if jsonData, err := data.MarshalJSON(); err == nil {
			if err = yaml.Unmarshal(jsonData, configMap); err == nil {
				_ = yaml.Unmarshal([]byte(configMap.Data["kubesphere.yaml"]), &result)
			} else {
				result["devops"]["a"] = err.Error()
			}
		} else {
			result["devops"]["b"] = err.Error()
		}
	} else {
		result["devops"]["c"] = err.Error()
	}
	return
}