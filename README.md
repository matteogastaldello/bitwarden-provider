# Bitwarden Provider
Bitwarden Provider that handles bitwarden secrets (till now😀)

## Notes
Le API di Bitwareden non sono pubbliche in versione free. Bitwarden mette comunque a disposizione tramite il client "cli" di fare il redirect di simulare le chiamate ad una api creando in locale un server che poi redirige le chiamate al server pubblico di bitwarden. 
Per questo motivo il funzionamento è leggermente diverso rispetto ad una api pubblica. Di seguito la spiegazione:
- il vault viene prima sbloccato tramite una chiamata POST con payload `{"password": "my-password"}` a `localhost:PORT/unlock`. Questa chiamata di fatto sblocca il Vault per tutte le successive chiamate (a cui non viene più richiesta alcuna credenziale). Ho implementato la funzione Unlock in `internal/clients/unlocker` per portare a termine questa operazione.
- per simulare il funzionamento dell'operator è quindi necessario installare Bitwarden CLI e lanciare il comando `bw serve` (https://bitwarden.com/help/cli/#serve) e impostare lo spec della risorsa ConnectorConfig->apiUrl (vedi `samples/connectorConfig.yaml`)
- La password di bitwarden deve essere impostata come k8s secret tramite `kubectl create secret generic bitwarden-password --from-literal=password=my-password`

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Start your cluster: (KIND or MINIKUBE)

```sh
minikube start
```

2. Install CRDs:

```sh
make install
```
3. Edit Samples in `samples/` folder according to your bw serve config


4. Install Instances of Custom Resources:

```sh
kubectl apply -f samples/
```

5. Install Instances of Custom Resources:

```sh
make install
```

6. Start your controller:

```sh
make run
```


### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

