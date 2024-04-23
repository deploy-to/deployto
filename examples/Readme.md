# helm версии 3, podman, minikube
#

Registry serrings https://minikube.sigs.k8s.io/docs/handbook/pushing/

minikube config set driver podman
minikube addons enable registry
minikube start --driver podman  --insecure-registry "192.168.0.0/0"


Иснтрукция по установке serviceb
запустить СI 
cd serviceb
make build
make push
cd helm
helm dependency build 
helm upgrade --install serviceb ./ -f values.yaml 


Иснтрукция по установке servicea
запустить СI 
cd servicea
make build
make push

helm upgrade --install servicea ./ -f values.yaml 

Итого получаем :

                    NGINX
                    |
                    |
                    ↓
 serviceb <---- servicea
  |    |            |
  |    |            |
  db   s3           db


проверить запрос к servicea через nginx:
minikube service servicea-ingress-nginx-controller 
curl {{URL из вывода предыдущей команды}}  -H "Host:servicea.deploy.to"