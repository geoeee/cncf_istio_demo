# clean:
# 	@echo "# iterate clean all demo projects"
# 	@make clean_bookinfo
# 	@make clean_echodemo
# 	@make clean_insoledemo

# clean_bookinfo:
# 	cd install/bookinfo && make clean

# clean_echodemo:
# 	cd install/echodemo && make clean

# clean_insoledemo:
# 	cd install/insoledemo && make clean

PWD = $(shell pwd)

istio/apply:
	echo "istio apply"
	kubectl apply -f $(PWD)/crds.yaml
	kubectl apply -f $(PWD)/authz-crds.yaml
	kubectl apply -f $(PWD)/istio.yaml


istio/clean:
	echo "istio clean"
	kubectl delete -f $(PWD)/istio.yaml
	kubectl delete -f $(PWD)/crds.yaml
	kubectl delete -f $(PWD)/authz-crds.yaml

istio/check:
	echo "istio check pods"
	kubectl get pods -n istio-system

kube_status:
	minikube status

kube_start:
	minikube start --logtostderr --v 0 --memory=6144 --bootstrapper=localkube
