# K8S

Kubernetes manifests for the Continuous Deployment pipeline are in this [repository](https://github.com/bruce-mig/argocd-go-micro-config)

- Running Postgres on the host machine:

    To spin up database on local machine, run the following command

    ```bash
    docker-compose -f postgres.yml up -d
    ```

    To connect to database from  minikube, change the host  to `host.minikube.internal` in the deployment manifest file.
