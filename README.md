# Savannah

This is a Go web application for the Savannah e-Shop project, built using Go standard library, go-chi and OAuth authentication. The application also uses Africa's Talking for sending SMS notifications

## Features
- Create and manage customer records.
- Create and manage order records.
- Demo callback endpoint for OAuth2 authentication.
- Send SMS notifications using Africa's Talking.

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


## Dependencies
  * Africa's Talking Go SDK: SDK for Africa's Talking SMS and other services.


