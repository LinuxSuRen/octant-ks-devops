package pipeline

import (
	"fmt"
	"github.com/linuxsuren/octant-ks-devops/pkg"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// OverviewHandler is the handler of Pipeline overview
func (h *PipelineHandler) OverviewHandler(request service.Request) (response component.ContentResponse, e error) {
	com1 := h.gen("tab 1", "tab1", request)

	var title = component.TitleFromString("sss" + request.Path())

	contentResponse := component.NewContentResponse(title)
	contentResponse.Add(com1)
	return *contentResponse, nil
}

func (h *PipelineHandler) gen(name, accessor string, request service.Request) component.Component {
	table := component.NewTable("pipeline", "xx", []component.TableCol{{
		Name: "Name",
	}, {
		Name: "CreateTime",
	}, {
		Name: "Detail",
	}})

	if nss, err := request.DashboardClient().List(request.Context(), store.Key{
		Namespace:  h.Context.Namespace,
		Kind:       "pipeline",
		APIVersion: "devops.kubesphere.io/v1alpha3",
	}); err == nil {
		for i := range nss.Items {
			item := nss.Items[i]

			row := component.TableRow{}
			row["Name"] = component.NewText(item.GetName())
			row["CreateTime"] = component.NewTimestamp(item.GetCreationTimestamp().Time)
			row["Detail"] = component.NewLink("View", "View",
				fmt.Sprintf("/%s/namespace/%s/pipeline/%s", pkg.PluginName, h.Context.Namespace, item.GetName()))
			table.Add(row)
		}
	}

	card := component.NewCard(component.TitleFromString("table"))
	card.SetBody(table)

	// fetch pipeline objects
	cardList := component.NewCardList(name)
	cardList.AddCard(*card)
	cardList.SetAccessor(accessor)
	return cardList
}
