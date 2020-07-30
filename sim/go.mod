module skenario

go 1.12

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.10.2 // indirect
	github.com/bvinc/go-sqlite-lite v0.6.1
	github.com/go-chi/chi v4.0.3+incompatible
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/google/go-containerregistry v0.0.0-20190222233527-d3e6a441f49f // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190212212710-3befbb6ad0cc // indirect
	github.com/hashicorp/go-plugin v1.0.1
	github.com/josephburnett/sk-plugin v0.0.0
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/nyarly/spies v0.0.0-20200413230442-112961b2b018 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/sclevine/agouti v3.0.0+incompatible
	github.com/sclevine/spec v1.4.0
	github.com/stretchr/testify v1.5.1
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/apimachinery v0.17.9
	k8s.io/client-go v10.0.0+incompatible
	k8s.io/klog v1.0.0
	sigs.k8s.io/yaml v1.2.0 // indirect
)

require (
	github.com/prometheus/client_golang v1.0.0
	k8s.io/api v0.17.9
	k8s.io/autoscaler/vertical-pod-autoscaler v0.0.0-20200723095539-0e5f460aac8c
	k8s.io/kubernetes v0.0.0
	k8s.io/metrics v0.17.9
)

// https://github.com/kubernetes/kubernetes/issues/79384#issuecomment-505627280
replace (
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	k8s.io/api => k8s.io/api v0.17.9
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.9
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.9
	k8s.io/apiserver => k8s.io/apiserver v0.17.9
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.17.9
	k8s.io/client-go => k8s.io/client-go v0.17.9
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.17.9
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.9
	k8s.io/code-generator => k8s.io/code-generator v0.17.9
	k8s.io/component-base => k8s.io/component-base v0.17.9
	k8s.io/cri-api => k8s.io/cri-api v0.17.9
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.17.9
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.17.9
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.17.9
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.17.9
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.17.9
	k8s.io/kubectl => k8s.io/kubectl v0.17.9
	k8s.io/kubelet => k8s.io/kubelet v0.17.9
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.17.9
	k8s.io/metrics => k8s.io/metrics v0.17.9
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.17.9
)

// Checkout https://github.com/josephburnett/kubernetes.git branch plugin
replace k8s.io/kubernetes => ../../kubernetes

replace k8s.io/autoscaler/vertical-pod-autoscaler => ../../autoscaler/vertical-pod-autoscaler

// TODO: replace this import with github.com/skenario/plugin
replace github.com/josephburnett/sk-plugin => ../plugin
