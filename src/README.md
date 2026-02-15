# E-Commerce API
## Description
This is a simple e-commerce API that allows you to manage products. You can create, read, update, and delete products. The API is built using Go and uses an in-memory store to manage the products.

## Instructions

### Prepare Credentials

Retrieve you temporary credentials from Learner's Lab.
Enter your configuration when prompted:
```
aws configure
```

And set your session token:
```
aws configure set aws_session_token <YOUR-TEMP-SESSION-TOKEBN>
```

### Apply Infrastructure
```
cd terraform
terraform init
terraform apply -auto-approve
```

### Send Requests
Make sure you are in terraform folder

Get public ip address
```
$taskArn = aws ecs list-tasks `
  --cluster $(terraform output -raw ecs_cluster_name) `
  --service-name $(terraform output -raw ecs_service_name) `
  --query 'taskArns[0]' `
  --output text

$networkInterfaceId = aws ecs describe-tasks `
  --cluster $(terraform output -raw ecs_cluster_name) `
  --tasks $taskArn `
  --query "tasks[0].attachments[0].details[?name=='networkInterfaceId'].value" `
  --output text

$publicIp = aws ec2 describe-network-interfaces `
  --network-interface-ids $networkInterfaceId `
  --query 'NetworkInterfaces[0].Association.PublicIp' `
  --output text

echo $publicIp
```

Send some requests using Postman or curl, using the public IP address you retrieved in the previous step:
```
curl http://<PUBLIC-IP-ADDRESS>:8080...
```
## Endpoints
- `GET /products/{id}`: Get a product by its ID.
- `POST /products/{id}/details`: Create or update a product's details by its ID.

## Testing
You can use the locustfile to perform load testing on the API. Make sure you have Locust installed, and then from the testing directory run:
```
docker-compose up
docker-compose up -scale worker={NUMBER_OF_WORKERS}
```

Then open your browser and navigate to `http://<PUBLIC-IP-ADDRESS>:8089` to access the Locust web interface. From there, you can start the load test by specifying the number of users and the spawn rate.

## Clean Up
```
terraform destroy -auto-approve
```
