# kubeconfig : Create and Update Kubernetes deployment
This sample project demonstrates how to deploy nginx server on Kubernetes cluster using a Go tool. 
The Go tool provides an option to update pod specifications by means of accepting CLI arguments

Also it provides standalone config files using which the nginx server can be deployed using kubectl


## Running this example
Make sure you have a Kubernetes cluster and kubectl is configured:

    kubectl get nodes

Compile this sample application on your workstation:

```
cd kubeproject
go build -o ./deploynginx main.go
```

Now, run this application on your workstation with kubeconfig file provided in this project:

```
./deploynginx
# or specify path for kubeconfig file with flag
./deploynginx -kubeconfig=./.kube/config

#To update nginx version or to scale the pods, use below options
./deploynginx -version=1.13 -scale=3

```
