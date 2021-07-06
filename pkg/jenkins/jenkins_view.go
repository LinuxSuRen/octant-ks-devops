package jenkins

import (
	"github.com/linuxsuren/octant-ks-devops/pkg"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/yaml"
	"strconv"
)

// JenkinsHandler is the handler of pipeline
type JenkinsHandler struct {
	Context *pkg.PluginContext
}

func (h *JenkinsHandler) Dashboard(request service.Request) (response component.ContentResponse, e error) {
	response = *component.NewContentResponse(component.TitleFromString("Jenkins"))

	user, password := getUserAndPassword(request)
	response.Add(component.NewSummary("Summary", component.SummarySection{
		Header:  "User",
		Content: component.NewText(user),
	}, component.SummarySection{
		Header:  "Jenkins Token",
		Content: component.NewText(password),
	}))
	return
}

func getUserAndPassword(request service.Request) (user, password string) {
	if list, err := request.DashboardClient().List(request.Context(), store.Key{
		APIVersion: "v1",
		Kind:       "Secret",
		Selector:   &labels.Set{
			"devops.kubesphere.io/component=": "Jenkins",
		},
	}); err == nil && len(list.Items) > 0 {
		data := list.Items[0]
		secret := &v1.Secret{}
		if jsonData, err := data.MarshalJSON(); err == nil {
			if err = yaml.Unmarshal(jsonData, secret); err == nil {
				user = string(secret.Data["jenkins-admin-user"])
				password = string(secret.Data["jenkins-admin-password"])
			} else {
				user = err.Error()
			}
		} else {
			user = err.Error()
		}
	} else if err != nil {
		user = err.Error()
	} else {
		user = strconv.Itoa(len(list.Items))
	}
	return
}
