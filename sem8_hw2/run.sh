#!/bin/bash

istioctl install --set profile=logs -y
kubectl label namespace default istio-injection=enabled

kubectl apply -f k8s/deployment/deployment.yaml 
kubectl apply -f k8s/service/service.yaml
kubectl apply -f k8s/gateway/gateway.yaml
kubectl apply -f k8s/gateway/virtual-service.yaml
kubectl apply -f k8s/gateway/destination-rule.yaml

kubectl wait --for=condition=ready pod -l app=logs-app --timeout=60s
