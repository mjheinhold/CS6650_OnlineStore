from locust import FastHttpUser, HttpUser, tag, task, between

class OnlineStoreUser(FastHttpUser):
    wait_time = between(1, 3)

    @tag("view")
    @task(1)
    def browse_products(self):
        for product_id in range(1, 11):  # Assuming we have 10 products
            self.client.get(f"/products/{product_id}")

    @tag("update")
    @task(1)
    def update_product(self):
        self.client.post("/products/1/details", json={"product_id": 1, "sku": "updated_sku", "manufacturer": "updated_manufacturer", "category_id": 1, "weight": 150, "some_other_id": 11})
