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



istio/apply:
	echo "istio apply"
	kubectl apply -f ${ISTIO_DEMO_INSTALL}/crds.yaml
	kubectl apply -f ${ISTIO_DEMO_INSTALL}/authz-crds.yaml
	kubectl apply -f ${ISTIO_DEMO_INSTALL}/istio.yaml


istio/clean:
	echo "istio clean"
	kubectl delete -f ${ISTIO_DEMO_INSTALL}/istio.yaml
	kubectl delete -f ${ISTIO_DEMO_INSTALL}/crds.yaml
	kubectl delete -f ${ISTIO_DEMO_INSTALL}/authz-crds.yaml

istio/check:
	echo "istio check pods"
	kubectl get pods -n istio-system

kube_status:
	minikube status

kube_start:
	minikube start --logtostderr --v 0 --memory=16144 --bootstrapper=localkube
