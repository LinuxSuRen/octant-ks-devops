package path

import (
	"regexp"
)

// GetPipelineNamespaced returns the namespace and name from path
func GetPipelineNamespaced(routePath string) (ns, name string) {
	r := regexp.MustCompile(`/namespace/(.*)/pipeline/(.*)`)
	if match := r.FindStringSubmatch(routePath); match != nil && len(match) == 3 {
		ns = match[1]
		name = match[2]
	}
	return
}
