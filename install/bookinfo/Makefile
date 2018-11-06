apply:
	istioctl kube-inject -f bookinfo.yaml > istio-bookinfo.yaml
	kubectl apply -f istio-bookinfo.yaml
	kubectl apply -f bookinfo-gateway.yaml
	kubectl apply -f destination-rules-all.yaml
	kubectl apply -f virtual-service-all-v1.yaml

clean:
	@echo "# clean bookinfo demo project"
	kubectl delete -f istio-bookinfo.yaml
	kubectl delete -f bookinfo-gateway.yaml
	kubectl delete -f destination-rules-all.yaml
	kubectl delete -f virtual-service-all-v1.yaml