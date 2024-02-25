# Savannah

This is a Go web application for the Savannah e-Shop project, built using Go standard library, go-chi and OAuth authentication. The application also uses Africa's Talking for sending SMS notifications

## Features
- Create and manage customer records.
- Create and manage order records.
- Demo callback endpoint for OAuth2 authentication.
- Send SMS notifications using Africa's Talking.
- PostgreSQL database

## Installation
To run this application locally, follow this step:
1. Clone the repository:
    - ```git clone git@github.com:0xAckerMan/Savanah.git```

2. Create a savannah.env file in the root directory and set the following environment variables:
```
ENVIRONMENT=<dev|stag|prod>
DATABASE_DSN="host=<host> user=<dbusername> password=<dbpassword> dbname=<dbname> port=<dbport> sslmode="


JWT_SECRET=
TOKEN_EXPIRES_IN=
TOKEN_MAXAGE=

GOOGLE_OAUTH_CLIENT_ID=
GOOGLE_OAUTH_CLIENT_SECRET=
GOOGLE_OAUTH_REDIRECT_URL=

AFRICASTALKING_APIKEY=
AFRICASTALKING_USERNAME=
```
3. Install dependencies:
    ```go mod tidy```


4. Run the application:
  ```go run ./cmd/api```



## Routes

- ```/``` **GET** - home
  ```json
  {
    "message": "Welcome to Savanah API"
  }
  ```
- ```/api/v1/healthcheck``` **GET** -return the health status of the api
  ```json
  {
      "health": {
          "environment": "production",
          "status": "active",
          "version": "1.0.0"
      }
  }
  ```
- ```/api/v1/register``` **POST** - Creates customer and admin if ```is_admin = true```
  ##### body 
  ```json
  {
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@gmail.com",
    "phone_number": "+254712345678",
    "password": "password",
    "is_admin": true
  }
  ```

  ##### response
  ```json
  {
    "customer": {
        "ID": 1,
        "CreatedAt": "2024-02-22T16:56:21.727300812+03:00",
        "UpdatedAt": "2024-02-22T16:56:21.727300812+03:00",
        "DeletedAt": null,
        "first_name": "John",
        "last_name": "Doe",
        "email": "john@gmail.com",
        "phone_number": "+254712345678",
        "Orders": null
    }
  }
  ```
- ```/api/v1/login``` **POST** - Users login to the system and token is setted in cookies
  #### body
  ```json
  {
    "email": "john@gmail.com",
    "password": "password"
  }
  ```

  #### response
  ```json
  {
    "token": "##############################################"
  }
  ```
- ```/api/v1/logout``` **GET** - Logs out a user from the system


  ### Admin user endpoints
- ```/api/v1/admin/products``` **POST** - creates a product and returns a response of the created product
  #### body
  ```json
  {
    "product_name": "Headphone",
    "description": "A blue Headphone",
    "price": 4000
  }
  ```

  #### response
  ```json
  {
    "product": {
        "ID": 1,
        "CreatedAt": "2024-02-22T17:52:36.965522494+03:00",
        "UpdatedAt": "2024-02-22T17:52:36.965522494+03:00",
        "DeletedAt": null,
        "product_name": "Headphone",
        "description": "A blue Headphone",
        "price": 4000,
        "version": 1
    }
  }
  ```

- ```/api/v1/admin/products/1``` **PATCH** - Performs both partial and full delete and updates the version of the product to avoid data race
  #### body
  ```json
  {
    "product_name": "UBL Headphone",
    "description": "A blue high-end Headphone"
  }
  ```
  #### response
  ```json
  {
    "product": {
        "ID": 1,
        "CreatedAt": "2024-02-22T17:52:36.965522+03:00",
        "UpdatedAt": "2024-02-22T17:57:09.675743505+03:00",
        "DeletedAt": null,
        "product_name": "UBL Headphone",
        "description": "A blue high-end Headphone",
        "price": 4000,
        "version": 2
    }
  }
  ```


- ```/api/v1/admin/customers``` **GET** - gets all customers who are not admin
  #### response
  ```json
  {
    "customers": [
        {
            "ID": 2,
            "CreatedAt": "2024-02-22T17:02:28.749401+03:00",
            "UpdatedAt": "2024-02-22T17:02:28.749401+03:00",
            "DeletedAt": null,
            "first_name": "Jane",
            "last_name": "smith",
            "email": "jsmith@localhost.com",
            "phone_number": "+254722345678",
            "Orders": []
        },
        {
            "ID": 4,
            "CreatedAt": "2024-02-22T18:19:51.768231+03:00",
            "UpdatedAt": "2024-02-22T18:19:51.768231+03:00",
            "DeletedAt": null,
            "first_name": "John",
            "last_name": "smith",
            "email": "smith@localhost.com",
            "phone_number": "+254712234568",
            "Orders": []
        }
    ]
  }
  ```
  
- 

  ### Customer user
  - ```/api/v1/orders``` **POST** - a customer can make an order. Gets the current logged in user.
    #### body
    ```json
    {
    "product_id": 1,
    "quantity": 2
    }
  ```

  #### response
  ```json
  {
    "order": {
        "ID": 5,
        "CreatedAt": "2024-02-22T19:15:29.934709087+03:00",
        "UpdatedAt": "2024-02-22T19:15:29.934709087+03:00",
        "DeletedAt": null,
        "customer_id": 2,
        "product_id": 1,
        "quantity": 2,
        "order_status": "placed",
        "version": 1,
        "Customer": null,
        "Product": null
    }
  }
  ```

  - ```/api/v1/customers/me``` **GET** - Gets a customer profile
    ```json
    {
      "customer": {
          "ID": 2,
          "CreatedAt": "2024-02-22T17:02:28.749401+03:00",
          "UpdatedAt": "2024-02-22T17:02:28.749401+03:00",
          "DeletedAt": null,
          "first_name": "Jane",
          "last_name": "smith",
          "email": "jsmith@localhost.com",
          "phone_number": "+254722345678",
          "Orders": [
              {
                  "ID": 1,
                  "CreatedAt": "2024-02-22T19:05:40.451318+03:00",
                  "UpdatedAt": "2024-02-22T19:05:40.451318+03:00",
                  "DeletedAt": null,
                  "customer_id": 2,
                  "product_id": 1,
                  "quantity": 2,
                  "order_status": "placed",
                  "version": 1,
                  "Customer": null,
                  "Product": {
                      "ID": 1,
                      "CreatedAt": "2024-02-22T17:52:36.965522+03:00",
                      "UpdatedAt": "2024-02-22T17:57:09.675743+03:00",
                      "DeletedAt": null,
                      "product_name": "UBL Headphone",
                      "description": "A blue high-end Headphone",
                      "price": 4000,
                      "version": 2
                  }
              }
          ]
      }
    }
    ```
  


## Dependencies
  * Africa's Talking Go SDK: SDK for Africa's Talking SMS and other services.


