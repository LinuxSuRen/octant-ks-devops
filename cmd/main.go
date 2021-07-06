package main

import (
	"fmt"
	"github.com/linuxsuren/octant-ks-devops/pkg"
	"github.com/linuxsuren/octant-ks-devops/pkg/config"
	"github.com/linuxsuren/octant-ks-devops/pkg/jenkins"
	"github.com/linuxsuren/octant-ks-devops/pkg/pipeline"
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"github.com/vmware-tanzu/octant/pkg/view/flexlayout"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"
)

type Handlers struct {
	Context *pkg.PluginContext
}

func (h *Handlers) actions(request *service.ActionRequest) error {
	switch request.ActionName {
	case action.RequestSetNamespace:
		h.Context.Namespace, _ = request.Payload.String("namespace")
	case pkg.ActionSetName:
		h.Context.Namespace, _ = request.Payload.String("name")
	}
	return nil
}

func (h *Handlers) InitRoutes(router *service.Router) {
	pipelineHandler := pipeline.PipelineHandler{Context: h.Context}
	router.HandleFunc("/overview", pipelineHandler.OverviewHandler)
	router.HandleFunc(fmt.Sprintf("/namespace/*/pipeline/*", ), pipelineHandler.DetailHandler)

	configHandler := config.ConfigHandler{
		Context: h.Context,
	}
	router.HandleFunc("/config", configHandler.Dashboard)

	jenkinsHandler := jenkins.JenkinsHandler{
		Context: h.Context,
	}
	router.HandleFunc("/jenkins", jenkinsHandler.Dashboard)
}

func main() {
	devopsGVK := schema.GroupVersionKind{
		Group:   "devops.kubesphere.io",
		Version: "v1alpha3",
		Kind:    "pipeline",
	}
	cmGVK := schema.GroupVersionKind{
		Version: "v1",
		Kind:    "ConfigMap",
	}

	capabilities := &plugin.Capabilities{
		SupportsPrinterConfig: []schema.GroupVersionKind{devopsGVK, cmGVK},
		SupportsPrinterStatus: []schema.GroupVersionKind{devopsGVK, cmGVK},
		SupportsTab:           []schema.GroupVersionKind{devopsGVK, cmGVK},
		ActionNames:           []string{pkg.PluginActionName, action.RequestSetNamespace},
		IsModule:              true,
	}

	handlers := &Handlers{
		Context: &pkg.PluginContext{},
	}
	options := []service.PluginOption{
		service.WithPrinter(handlePrint),
		service.WithTabPrinter(handleTab),
		service.WithNavigation(handleNavigation, handlers.InitRoutes),
		service.WithActionHandler(handlers.actions),
	}

	if p, err := service.Register(pkg.PluginName, "Get more from http://github.com/linuxsuren/octant-ks-devops",
		capabilities, options...); err != nil {
		log.Fatal(err)
	} else {
		p.Serve()
	}
}

func handleNavigation(request *service.NavigationRequest) (x navigation.Navigation, e error) {
	return navigation.Navigation{
		Title:    "ks-devops",
		Path:     pkg.PluginName + "/overview",
		IconName: "cloud",
		Children: []navigation.Navigation{{
			Title: "Config",
			Path:  pkg.PluginName + "/config",
		}, {
			Title: "Pipeline",
			Path:  pkg.PluginName + "/overview",
		}, {
			Title:    "S2I",
			Path:     pkg.PluginName + "/s2i",
			IconName: "s2i",
		}, {
			Title:    "Jenkins",
			Path:     pkg.PluginName + "/jenkins",
			IconName: "jenkins",
		}},
	}, nil
}

func handlePrint(request *service.PrintRequest) (response plugin.PrintResponse, e error) {
	if request.Object == nil {
		return plugin.PrintResponse{}, errors.Errorf("object is nil")
	}

	card := component.NewCard(component.TitleFromString("hello"))
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
