inspect:
	kubectl get virtualservice
	kubectl get gateway
	kubectl get destinationrule
	kubectl get svc
	kubectl get pods

insole/apply:
	@echo "applying insole demo"
	kubectl apply -f namespace-demo.yaml
	istioctl kube-inject -f deploy.yaml > deploy-istio.yaml
	kubectl apply -f deploy-istio.yaml
	kubectl apply -f gateway.yaml
	kubectl apply -f virtual-services.yaml
	kubectl apply -f destination-rules.yaml

insole/clean:
	@echo "cleaning insole demo"
	kubectl delete -f deploy-istio.yaml
	kubectl delete -f gateway.yaml
	kubectl delete -f virtual-services.yaml
	kubectl delete -f destination-rules.yaml
	kubectl delete -f namespace-demo.yaml

weight_split/apply:
	@echo "==== apply weight split for product service"
	kubectl apply -f virtual-service-weight-split.yaml

weight_split/clean:
	@echo "==== clean back random split for product service"
	kubectl apply -f virtual-services.yaml

header_split/apply:
	@echo "==== apply header split for product service"
	kubectl apply -f virtual-service-header-split.yaml

header_split/clean:
	@echo "==== clean back random split for product service"
	kubectl apply -f virtual-services.yaml

policy_jwt/apply:
	@echo "==== apply jwt token auth" 
	kubectl apply -f policy-jwt.yaml

policy_jwt/clean:
	@echo "==== clean jwt token auth" 
	kubectl delete -f policy-jwt.yaml

white_list/apply:
	@echo "==== apply white list config"
	kubectl apply -f white-list.yaml

white_list/clean:
	@echo "==== clean white list config"
	kubectl delete -f white-list.yaml

authz/apply:
	@echo "==== apply authz config"
	kubectl apply -f authz-config.yaml

authz/clean:
	@echo "==== clean authz config"
	kubectl delete -f authz-config.yaml