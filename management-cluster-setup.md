# Steps to setup a TKG management cluster

These steps walk through the required pre-reqs to use the cpservice-annotator in order to use different avi setting for control plane load balancers in tkg 1.4.x


1. deploy the `aviinfrasettings` CRD into the mgmt cluster

```bash
kubectl apply -f ./aviinfrasettings/aviinfracrd.yml
```

2. patch the ako controller data values with an overlay that will inject the correct permissions to use `aviinfrasettings`

```bash
kubectl patch secret load-balancer-and-ingress-service-data-values -n tkg-system -p '{"data": {"overlays.yaml": "'$(cat ./aviinfrasettings/avi-cr-overlay.yml | base64)'"}}'
```

3. deploy an instance of aviinfrasettings for your avi network you want to use and customize it with any settings you want. modify the file refernced below for your needs.

```bash
kubectl apply -f ./aviinfrasettings/infrasetting.example.yml
```

4. add the TKG overlays so that you can specify the aviinfrasetting name when creating a cluster

 ```bash
 # copy the default values file
cp ./tkg-overlays/aviinfra_default_values.yml ~/.config/tanzu/tkg/providers/infrastructure-vsphere/ytt/.

# copy the overlay
cp ./tkg-overlays/aviinfra.yml ~/.config/tanzu/tkg/providers/infrastructure-vsphere/ytt/.

 ```

5. add the `AVI_CP_SETTING_NAME: <yourinfrasettingname>` to the cluster cofnig file 

6. the cluster is now setup to deploy the webhook