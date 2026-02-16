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
- - Example: `curl http://<PUBLIC-IP-ADDRESS>:8080/products/1`
- - Responses:
  - `200 OK`: Returns the product details in JSON format.
    - "product_id": 1,
    - "sku": "sku1",
    - "manufacturer": "manufacturer1",
    - "category_id": 1,
    - "weight": 10,
    - "some_other_id": 100 
  - `404 Not Found`: If the product with the specified ID does not exist.
    - "error": "Product not found"
- `POST /products/{id}/details`: Add or update detailed information for a specific product.
- - Example: `curl -X POST http://<PUBLIC-IP-ADDRESS>:8080/products/1/details -H "Content-Type: application/json" -d '{"sku": "sku1", "manufacturer": "manufacturer1", "category_id": 1, "weight": 10, "some_other_id": 100}'`
- - Responses:
  - `204 No Content`: If the product details are successfully added or updated.
  - `404 Not Found`: If the product with the specified ID does not exist.
    - "error": "Product not found"
  - `400 Bad Request`: If the request body is invalid or missing required fields or the product ID is not a number.
    - "error": "Invalid product ID"  

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
