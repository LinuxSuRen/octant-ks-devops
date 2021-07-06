module github.com/linuxsuren/octant-ks-devops

go 1.15

require (
	github.com/elazarl/goproxy/ext v0.0.0-20210110162100-a92cc753f88e // indirect
	github.com/pkg/errors v0.9.1
	github.com/vmware-tanzu/octant v0.19.0
	k8s.io/api v0.19.3
	k8s.io/apimachinery v0.20.2
//kubesphere.io/devops v0.0.0-20190413051334-03b00cc2226d // indirect
)

//replace kubesphere.io/devops => gitee.com/linuxsuren/ks-devops v0.0.0-20210702035442-1a510b13c417
