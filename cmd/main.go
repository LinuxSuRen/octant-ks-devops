package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"github.com/vmware-tanzu/octant/pkg/view/flexlayout"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"
)

const pluginName = "ks-devops"
const pluginActionName = "action.kubesphere.io/devops"

type pluginContext struct {
	Namespace string
}

type Handlers struct {
	Context *pluginContext
}

func (h *Handlers) actions(request *service.ActionRequest) error {
	switch request.ActionName {
	case action.RequestSetNamespace:
		h.Context.Namespace, _ = request.Payload.String("namespace")
		log.Println("===", h.Context.Namespace)
	}
	return nil
}

func (h *Handlers) InitRoutes(router *service.Router) {
	gen := func(name, accessor string, request service.Request) component.Component {
		table := component.NewTable("pipeline", "xx", []component.TableCol{{
			Name: "name",
		}, {
			Name: "status",
		}})

		if nss, err := request.DashboardClient().List(request.Context(), store.Key{
			Namespace:  h.Context.Namespace,
			Kind:       "pipeline",
			APIVersion: "devops.kubesphere.io/v1alpha3",
		}); err == nil {
			for i := range nss.Items {
				ns := nss.Items[i].GetName()
				row := component.TableRow{}
				row["name"] = component.NewText(ns)
				row["status"] = component.NewText(h.Context.Namespace)
				table.Add(row)
			}
		}

		card_1 := component.NewCard(component.TitleFromString("table"))
		card_1.SetBody(table)

		// fetch pipeline objects

		cardList := component.NewCardList(name)
		cardList.AddCard(*card_1)
		cardList.SetAccessor(accessor)
		return cardList
	}

	router.HandleFunc("*", func(request service.Request) (response component.ContentResponse, e error) {
		com1 := gen("tab 1", "tab1", request)

		var title = component.TitleFromString("sss")

		contentResponse := component.NewContentResponse(title)
		contentResponse.Add(com1)
		return *contentResponse, nil
	})
}

func main() {
	devopsGVK := schema.GroupVersionKind{
		Group:   "devops.kubesphere.io",
		Version: "v1alpha3",
		Kind:    "pipeline",
	}

	capabilities := &plugin.Capabilities{
		SupportsPrinterConfig: []schema.GroupVersionKind{devopsGVK},
		SupportsPrinterStatus: []schema.GroupVersionKind{devopsGVK},
		SupportsTab:           []schema.GroupVersionKind{devopsGVK},
		ActionNames:           []string{pluginActionName, action.RequestSetNamespace},
		IsModule:              true,
	}

	handlers := &Handlers{
		Context: &pluginContext{},
	}
	options := []service.PluginOption{
		service.WithPrinter(handlePrint),
		service.WithTabPrinter(handleTab),
		service.WithNavigation(handleNavigation, handlers.InitRoutes),
		service.WithActionHandler(handlers.actions),
	}

	if p, err := service.Register(pluginName, "ks devops demo plugin",
		capabilities, options...); err != nil {
		log.Fatal(err)
	} else {
		p.Serve()
	}
}

func handleNavigation(request *service.NavigationRequest) (x navigation.Navigation, e error) {
	return navigation.Navigation{
		Title:    "ks-devops",
		Path:     request.GeneratePath(),
		IconName: "cloud",
	}, nil
}

func handlePrint(request *service.PrintRequest) (response plugin.PrintResponse, e error) {
	if request.Object == nil {
		return plugin.PrintResponse{}, errors.Errorf("object is nil")
	}

	card := component.NewCard(component.TitleFromString("hello"))

	key, err := store.KeyFromObject(request.Object)
	fmt.Println("key", key)
	if err != nil {
		return plugin.PrintResponse{}, err
	}

	u, err := request.DashboardClient.Get(request.Context(), key)
	if err != nil {
		return plugin.PrintResponse{}, err
	}
	fmt.Println("===", u.GetName())

	return plugin.PrintResponse{
		Config: []component.SummarySection{
			{Header: "from-plugin", Content: component.NewText("hello ks-devops")},
		},
		Status: []component.SummarySection{
			{Header: "from-plugin", Content: component.NewText("hello ks-devops")},
		},
		Items: []component.FlexLayoutItem{
			{Width: component.WidthHalf, View: card},
		},
	}, nil
}

func handleTab(request *service.PrintRequest) (response plugin.TabResponse, err error) {
	if request.Object == nil {
		response = plugin.TabResponse{}
		err = errors.New("object is nil")
		return
	}

	layout := flexlayout.New()
	section := layout.AddSection()
	content := component.NewMarkdownText("hello md")

	_ = section.Add(content, component.WidthFull)

	tab := component.NewTabWithContents(*layout.ToComponent("KS DevOps"))
	response = plugin.TabResponse{Tab: tab}
	return
}
