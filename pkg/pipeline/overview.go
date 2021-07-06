package pipeline

import (
	"fmt"
	"github.com/linuxsuren/octant-ks-devops/pkg"
	"github.com/vmware-tanzu/octant/pkg/action"
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

	table.Config.ButtonGroup = component.NewButtonGroup()
	table.Config.ButtonGroup.AddButton(
		component.NewButton("ss", action.Payload{}))
	table.Config.ButtonGroup.AddButton(
		component.NewButton("Create a piepline", action.Payload{}, component.WithModal(h.pipelineCreateStepper())))
	//card := component.NewCard(component.TitleFromString("table"))
	//card.SetBody(table)

	// fetch pipeline objects
	//cardList := component.NewCardList(name)
	//cardList.AddCard(*card)
	//cardList.SetAccessor(accessor)

	flexLayout := component.NewFlexLayout("NewFlexLayout")
	flexLayout.AddSections(component.FlexLayoutSection{
		{Width: component.WidthFull, View: table},
	})
	return flexLayout
}

func (h *PipelineHandler) pipelineCreateStepper() (modal *component.Modal) {
	//typeForm := component.Form{
	//	Fields: []component.FormField{
	//		component.NewFormFieldCheckBox("Source", "source", []component.InputChoice{{
	//			Label:   "sss",
	//			Value:   "sss",
	//			Checked: true,
	//		}}),
	//	},
	//}

	//networkingForm := component.Form{
	//	Fields: []component.FormField{
	//		component.NewFormFieldRadio("IP Family", "ipFamily", []component.InputChoice{
	//			{
	//				"IPv4",
	//				"string(v1alpha4.IPv4Family)",
	//				false,
	//			},
	//			{
	//				"IPv6",
	//				"string(v1alpha4.IPv6Family)",
	//				false,
	//			},
	//		}),
	//		component.NewFormFieldText("API Server Address", "apiServerAddress", ""),
	//		component.NewFormFieldNumber("API Server Port", "apiServerPort", ""),
	//		component.NewFormFieldCheckBox("Disable Default CNI", "disableDefaultCNI", []component.InputChoice{
	//			{
	//				"Yes",
	//				"true",
	//				false,
	//			},
	//			{
	//				"No",
	//				"false",
	//				false,
	//			},
	//		}),
	//		component.NewFormFieldText("Pod Subnet", "podSubnet", ""),
	//		component.NewFormFieldText("Service Subnet", "serviceSubnet", ""),
	//		component.NewFormFieldRadio("", "", []component.InputChoice{
	//			{
	//				"iptables",
	//				"iptables",
	//				false,
	//			},
	//			{
	//				"ipvs",
	//				"ipvs",
	//				false,
	//			},
	//		}),
	//	},
	//}

	stepper := component.Stepper{
		Base: component.Base{
			Metadata: component.Metadata{
				Type:  "base",
				Title: component.TitleFromString("xxx"),
			},
		},
		Config: component.StepperConfig{
			Action: pkg.ActionCreatePipeline,
			Steps: []component.StepConfig{
				//{
				//	Name:        "type",
				//	//Form:        networkingForm,
				//	Title:       "type",
				//	Description: "choose pipeline type",
				//},
			},
		},
	}
	modal = component.NewModal(component.TitleFromString("Create Pipeline"))

	//if _, err := typeForm.MarshalJSON(); err !=nil {
	//	modal.SetBody(component.NewText(err.Error()))
	//} else if _, err = stepper.MarshalJSON(); err != nil {
	//	modal.SetBody(component.NewText(err.Error()))
	//} else {
	//	modal.SetBody(&stepper)
	//}
	modal.SetBody(&stepper)
	modal.SetSize(component.ModalSizeExtraLarge)
	return
}
