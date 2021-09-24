# Control Plane Service Annotator

The cpservice-annotator is a k8s mutating webhook that will inject the right AKO annotations for specifying [aviinfrasettings](https://avinetworks.com/docs/ako/1.4/custom-resource-definitions/#avi-infra-setting-with-services-ingress-or-routes) into services created by the ako-operator-controller in TKG. This allows for having the control plane load balancer services created on different networks or service engines etc. 


## Prereqs

* TKG 1.4.x
* AVI enabled on the mgmt cluster
* AVI set up with the networks you want to use
* steps [here](./management-cluster-setup.md) followed to install the extra k8s objects needed


## Architecture


## Install

this will install the webhook, configure it with a self signed cert using cert manager, and  setup the rbac needed. 

***Be sure to have completed the prereqs above***

1. `make deploy-cert`

2. `make deploy`


## Usage


1. create a TKG cluster

```
tanzu create cluster -f yourclusterconfig.yml
```


you should see the cluster get inject with an annotation `cpservicemutate.field.vmware.com/aviinfrasetting: <yourinfrasettingname>` this will be picked up by the webhook and translated to an annotation on the service that configures avi


