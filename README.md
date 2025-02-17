# My Chat App

This repository contains a browser-based chat application with real-time communication using WebSockets, user authentication with sessions, and a stock quote command integration using RabbitMQ and a bot. All the service layer functions (use cases / rules of business layer) were tested. 💻

## Prerequisites

Before starting, please ensure you have installed:

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- [sqlc](https://docs.sqlc.dev/en/stable/overview/install.html) (for generating SQL code)
- [migrate](https://github.com/golang-migrate/migrate) (for running database migrations)
- [swag](https://github.com/swaggo/swag) (optional, for generating Swagger documentation)
- make (check how to install for your OS)

## Setup

### 1. Setup .env

- Copy the provided `.env.example` file to `.env`:
- (I keep my local .env in order to make it easy)

```bash
  cp .env.example .env
  ```

- Edit the .env file and adjust the values as needed.

- ### Generating Session Keys:
To generate the **_SESSION_AUTH_KEY_** and **_SESSION_ENC_KEY_**, run:

  ```bash
  openssl rand -base64 32
  ```
Then, paste the generated output into your .env file for both keys.

### 2. Execute Docker Compose

The docker-compose.yml file reads the variables from the .env file and starts the following services:

- **Postgres** – Local database.
- **RabbitMQ** – Message broker.

To start these services, run:
  ```bash
  docker-compose up -d
  ```

### 3. Execute Create-DB and Migrations
   Create the database and run migrations using the Makefile targets. The following commands will:

**Create the Database:**

  ```bash
  make create-db
  ```
**Run Migrations Up:**

  ```bash
  make local-migration-up
  ```

### 4. Generate SQL Code

sqlc is used to generate Go code from SQL queries. Once sqlc is installed, run:

  ```bash
  make sqlc
  ```

### Optional

**Generate Swagger Documentation:**

  ```bash
  make swag
  ```
Alternatively, you can use the combined setup command (if defined):

  ```bash
  make setup
  ```
This command runs all the initial configuration steps, assuming all prerequisites are met.


### 5. Start the Applications

Use the following Makefile commands to start your application components:

- **Start Server**
  ```bash
  make start-server
  ```

- **Start Bot**
  ```bash
  make start-bot
  ```

These commands will initialize your backend (server) and frontend (client) applications.

## How the Application Works

1. ### **User Registration**:
Call the /user/register route to register two different users. The payload must include a valid password and nickname.

2. ### **User Login**:
Open two different browsers (or one regular window and one incognito window).
In each, open the /login route and log in with the registered credentials.
Upon successful login, the user is automatically redirected to the chat page.

3. ### **Chat Functionality**:
On the chat page, users can exchange messages in real time.

The chat application supports a special command for stock quotes: **/stock=stock_code**
For example, valid stock codes include:
- _googl.us_
- _aapl.us_
- _amzn.us_

When a user sends a command like /stock=googl.us, the command is processed by the bot (via RabbitMQ) and the bot’s response is broadcast to all connected clients.

4. ### **Swagger Documentation (Optional):**:

Once Swagger is generated (via make swag), you can consult the API documentation at:
  ```bash
  http://localhost:1323/swagger/index.html
  ```

Thanks, **Luccas Machado** 🪓

