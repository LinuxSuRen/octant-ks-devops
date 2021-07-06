package pipeline

import (
	"context"
	//"github.com/ghodss/yaml"
	"github.com/linuxsuren/octant-ks-devops/pkg/path"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	//"kubesphere.io/devops/pkg/api/devops/v1alpha3"
)

func (h *PipelineHandler) DetailHandler(request service.Request) (response component.ContentResponse, e error) {
	var title = component.TitleFromString("Details")
	contentResponse := component.NewContentResponse(title)
	contentResponse.Add(h.pipelineDetails(request)...)
	return *contentResponse, nil
}

func (h *PipelineHandler) pipelineDetails(request service.Request) (components []component.Component) {
	h.Context.Namespace, h.Context.Name = path.GetPipelineNamespaced(request.Path())

	components = make([]component.Component, 0)
	if pipeline, err := h.getPipeline(request.DashboardClient(), request.Context()); err == nil {
		components = append(components, createSummary(pipeline),
			createPipelineEditor(pipeline),
			createRawYAML(pipeline))
	} else {
		components = append(components, component.NewText("Cannot found Pipeline" + err.Error()))
	}
	return
}

func (h *PipelineHandler) getPipeline(client service.Dashboard, ctx context.Context) (pipeline *Pipeline, err error) {
	var data *unstructured.Unstructured
	if data, err = client.Get(ctx, store.Key{
		Namespace:  h.Context.Namespace,
		Name:       h.Context.Name,
		Kind:       "pipeline",
		APIVersion: "devops.kubesphere.io/v1alpha3",
	}); err == nil {
		var rawData []byte
		if rawData, err = data.MarshalJSON(); err == nil {
			pipeline = &Pipeline{}
			err = yaml.Unmarshal(rawData, pipeline)
		}
	}
	return
}

func createSummary(pipeline *Pipeline) (com component.Component) {
	summary := component.NewSummary("Summary", component.SummarySection{
		Header: "Type",
		Content: component.NewText(pipeline.Spec.Type),
	})
	if pipeline.Spec.Pipeline != nil {
		for _, param := range pipeline.Spec.Pipeline.Parameters {
			summary.AddSection(param.Name, component.NewText(param.Type))
		}
	}
	com = summary
	return
}

func createPipelineEditor(pipeline *Pipeline) (com component.Component) {
	editor := component.NewEditor(component.TitleFromString("Jenkinsfile"),
		pipeline.Spec.Pipeline.Jenkinsfile, true)
	editor.Config.Language = "groovy"
	editor.SetAccessor("groovy")
	com = editor
	return
}

func createRawYAML(pipeline *Pipeline) (com component.Component) {
	editor := component.NewEditor(component.TitleFromString("YAML"), "", true)
	pipeline.SetManagedFields(nil)
	_ = editor.SetValueFromObject(pipeline)
	com = editor
	return
}
