# city-population-app

## Testing the Golang App with `curl`

Assuming you have the Golang application running either locally or in Kubernetes (we'll cover Kubernetes deployment in the next steps), you can test the endpoints using `curl`.

a. **Health Endpoint**:
```
curl http://localhost:8080/health
```

b. **Insert or Update City Endpoint**:
```
curl -X POST -H "Content-Type: application/json" -d '{"city":"Jakarta","population": 8808000}' http://localhost:8080/city
```

c. **Retrieve Population of a City Endpoint**:
```
curl http://localhost:8080/city/Jakarta
```

## Building the Docker Image

Assuming you have Golang installed on your system and you are in the root directory of the Golang application:

a. **Build the Docker image**:
```
docker build -t your-dockerhub-username/city-population-app:latest .
```

b. **Push the Docker image to DockerHub** (optional, if you want to deploy it to Kubernetes from DockerHub):
```
docker push your-dockerhub-username/city-population-app:latest
```

Replace `your-dockerhub-username` with your actual DockerHub username.

## Deploying and Testing to Kubernetes with Docker Desktop

Assuming you have Docker Desktop installed and running on your system:

a. Start Docker Desktop and ensure that Kubernetes is enabled.

b. Open the terminal and navigate to the root directory of your Golang application.

c. **Build the Docker image** (if not done previously):
```
docker build -t your-dockerhub-username/city-population-app:latest .
```

d. **Deploy the application to Kubernetes using Helm** (make sure you have the Helm chart ready):
```
helm install city-population-app ./k8s/charts/city-population-app
```

e. **Check the status of the deployment**:
```
kubectl get pods
```

You should see the pods running. Make sure both the Golang application and Elasticsearch pods are in a "Running" state.

f. **Expose the service to access the Golang application**:
```
kubectl port-forward svc/city-population-app-service 8080:8080
```

The Golang application is now accessible at `localhost:8080`.

g. **Test the Golang application** using `curl` as described in the first step of testing.

h. **When you are done testing**, delete the Helm deployment:
```
helm delete city-population-app
```
