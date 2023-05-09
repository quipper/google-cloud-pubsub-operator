PROJECT_NAME := google-cloud-pubsub-operator
CLUSTER_NAME := $(PROJECT_NAME)-e2e
KUBECONFIG := output/kubeconfig.yaml
export KUBECONFIG

all:

cluster: $(KUBECONFIG)
$(KUBECONFIG):
	kind create cluster --name $(CLUSTER_NAME)
	kubectl version

delete-cluster:
	kind delete cluster --name $(CLUSTER_NAME)
	-rm $(KUBECONFIG)

deploy:
	cd controller && kustomize edit set image controller="$(CONTROLLER_IMAGE)"
	kustomize build controller | kubectl apply -f -
	kubectl -n $(PROJECT_NAME)-system rollout status deployment $(PROJECT_NAME)-controller-manager
	kubectl -n $(PROJECT_NAME)-system get deployment
	kustomize build ../config/samples | kubectl apply -f -

test:
	echo TODO

logs-controller:
	-kubectl -n $(PROJECT_NAME)-system logs -l control-plane=controller-manager --all-containers --tail=-1
