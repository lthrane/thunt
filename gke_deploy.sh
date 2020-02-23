#!/bin/sh
gcloud container clusters create thunt-cluster --num-nodes=2
kubectl create deployment thunt --image=registry.hub.docker.com/lthrane/thunt
kubectl expose deployment thunt --type=LoadBalancer --port 8080 --target-port 8080
kubectl get service