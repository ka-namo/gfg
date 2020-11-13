##Solution for version 2 for Product and Seller APIs

### Refactoring
I have refactored the existing code, as it was using concrete dependencies and no abstraction.
Also, at a few places, I have added proper error handling such as rows.Err() were not checked.

I have added few comments in the codebase starting with `NOTE -` explaining the reason for that change.

There is still room for improvements, so that this project can be more scalable and easy to change by following
[Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html),
but I did not transform the whole implementation, as I believed it was not expected and would take more
time.

### Testing

I have added unit tests for all the APIs(v1 and v2) and the coverage can be certainly increased. 

To test the all the APIs with version V1 and V2, run

`go test ./... -race -cover`

or 

`make test-unit`

### How to run manually
Please refer the steps given in README.md for the setup and run.
 
 For running version V2, you can replace v1 with v2 in all the products API and for Seller
 a new API for fetching top 10 sellers has been added.
 
 Note - I have added `NOTIFY_SMS` and `NOTIFY_EMAIL` ENV variable to `docker-compose.yml`
 to trigger the sms and email when a product stock is changed. These settings
 can be toggled.
 
 For your quick reference for testing v2 -
 
 __Get a page of products__
 
 ```curl "http://localhost:8080/api/v2/products"```
 
 ```curl "http://localhost:8080/api/v2/products?page=2"```
 
 __Get a product__
 
 ```curl "http://localhost:8080/api/v2/product?id=bdbba8f8-234b-11eb-82b0-0242ac130002"```
 
 __Create a product__
 
 ```curl -X POST  -d '{"name":"LED Shoes","brand":"Niko","stock":11,"seller":"bdbafde4-234b-11eb-82b0-0242ac130002"}' localhost:8080/api/v2/product```
 
 __Update a product__
 
 ```curl -X PUT -d '{"name":"Berlin S.O.L.I.D. T-Shirt","brand":"Shirts Inc.","stock":150}' "http://localhost:8080/api/v2/product?id=bdbba7c0-234b-11eb-82b0-0242ac130002"```
 
 __Delete a product__
 
 ```curl -X DELETE "http://localhost:8080/api/v2/product?id=bdbba7c0-234b-11eb-82b0-0242ac130002"```
 
 __Get list of sellers__
 
 ```curl "http://localhost:8080/api/v1/sellers"```
 
 __Get Top 10 list of sellers__
 
 ```curl "http://localhost:8080/api/v2/sellers/top10"```


