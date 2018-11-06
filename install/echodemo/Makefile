apply:
	istioctl kube-inject -f echo.yaml > istio-echo.yaml
	kubectl apply -f istio-echo.yaml
	kubectl apply -f echo-gateway.yaml
	kubectl apply -f echo-virtual-svc.yaml
	kubectl apply -f echo-destination-rule.yaml

clean:
	kubectl delete -f istio-echo.yaml
	kubectl delete -f echo-gateway.yaml
	kubectl delete -f echo-virtual-svc.yaml
	kubectl delete -f echo-destination-rule.yaml
	