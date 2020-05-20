module github.com/hth0919/migrationcontroller

go 1.13

require (
	github.com/hth0919/checkpointproto v0.0.2
	github.com/hth0919/migrationclient v0.0.8 // indirect
	github.com/operator-framework/operator-sdk v0.17.0
	github.com/spf13/pflag v1.0.5
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/api => k8s.io/api v0.15.12
	k8s.io/apimachinery => k8s.io/apimachinery v0.15.12
	k8s.io/client-go => k8s.io/client-go v0.15.12 // Required by prometheus-operator
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.3.0
)
