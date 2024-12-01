### Kubernetes with minikube

Деплой приложения:

1. `kubectl create -f kube/postgres.yaml`
2. `kubectl get pods`
![pg_up2.png](img%2Fpg_up2.png)
3. `kubectl create -f kube/app-deployment.yaml `
![app_up.png](img%2Fapp_up.png)
4. `kubectl expose deployment app-metrics --type=LoadBalancer --port=80`
5. `minikube service app-metrics --url`
![url.png](img%2Furl.png)
6. `kubectl get service`![service_up.png](img%2Fservice_up.png)
7. Приложение поднято
![app_use.png](img%2Fapp_use.png)