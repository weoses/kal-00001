# Cluster prerequisites

These are one-time, cluster-wide installs the `memelo` app chart depends on
(ingress + TLS issuance). They are not part of the per-env `helm upgrade`
for the app and don't need to be re-run per deploy.

Context: no cloud-controller-manager pod is visible in `kube-system`, but
Vultr's VKE clusters do provision `Service type=LoadBalancer` (confirmed by
test-deploying one — it got a real external IP in ~2 minutes; the CCM runs
on Vultr's managed side). There's no way to pin that LB to a specific
pre-existing reserved IP, so **DNS must point at whatever IP gets assigned
below**, not at any IP chosen in advance.

## 1. ingress-nginx

```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
  -n ingress-nginx --create-namespace \
  -f ingress-nginx-values.yaml

# wait for the LoadBalancer IP, then point the app domains at it:
kubectl get svc -n ingress-nginx ingress-nginx-controller -w
```

## 2. cert-manager

```bash
helm repo add jetstack https://charts.jetstack.io
helm repo update

helm upgrade --install cert-manager jetstack/cert-manager \
  -n cert-manager --create-namespace \
  --set crds.enabled=true

kubectl apply -f cluster-issuer.yaml   # fill in the ACME email first
kubectl get clusterissuer letsencrypt-prod
```

Once both are up, the `memelo` app chart's Ingress objects (annotated with
`cert-manager.io/cluster-issuer: letsencrypt-prod` and
`ingressClassName: nginx`) will get certs issued automatically on install.
