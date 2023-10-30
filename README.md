## Diagram

![alt text](https://github.com/dedihartono801/grpc-microservices/blob/master/system-design.jpg)

## Run Service

```bash
$ docker-compose up -d
```

## Create Index Elasticsearch

```bash
$ curl -X PUT localhost:9200/activity-logs
```

## Register User

```bash
$ curl --location 'http://localhost:4000/user/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "Hartono",
    "email": "hartono@gmail.com",
    "password": "12345"
}'
```

## Login

```bash
$ curl --location 'http://localhost:4000/user/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "hartono@gmail.com",
    "password": "12345"
}'
```

## List Product

```bash
$ curl --location 'http://localhost:4000/products'
```

## Create Transaction

```bash
$ curl --location 'http://localhost:4000/transaction' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6ImhhcnRvbm9AZ21haWwuY29tIiwiZXhwIjoxNjk4NzY2MjE5fQ.xJHopTmAKOKsI9S5BQmmCVKODt-7QJi8hGibCVx_xxM' \
--header 'Content-Type: application/json' \
--data ' {
 "items": [
        {
            "product_id":1,
            "quantity":1
        },
        {
            "product_id":2,
            "quantity":1
        },
        {
            "product_id":3,
            "quantity":1
        },
        {
            "product_id":4,
            "quantity":1
        },
        {
            "product_id":5,
            "quantity":1
        }
    ]
 }'
```

## Detail Transaction

```bash
$ curl --location 'http://localhost:4000/transaction/1' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6ImhhcnRvbm9AZ21haWwuY29tIiwiZXhwIjoxNjk4NzI3OTM0fQ.kTDIiWIfleOidlTvSNDyQwiz8g9uBGcYdo4FjoTZOgc'
```

## Logging (Kibana)

Available at `http://localhost:5601`
