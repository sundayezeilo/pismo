# Pismo

### Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Testing](#testing)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)

## Installation

Follow these steps to set up the app. Ensure you have Git installed:

1. **Clone the Repository**:

    ```sh
    git clone https://github.com/sundayezeilo/pismo
    cd pismo
    git checkout develop
    ```

2. **Install Dependencies**:
   
    Ensure that Make, Docker, and Docker Compose are installed on your machine. Refer to the following resources for installation guidance:
    - [GNU Make official website](https://www.gnu.org/software/make/)
    - [Docker official website](https://www.docker.com/products/docker-desktop)

## Configuration

2. **Set Up Environment Variables**:

    Environment variables are required for running the app. Refer to `env_sample.txt` in the root directory for a sample configuration. Create the `.env` file in the appropriate directory:

    - For automated tests: `path/to/pismo/tests/.env`
    - For development: `path/to/pismo/.env`

## Testing

To run tests, ensure the environment is properly set up and use the following commands in two separate terminals (one for Postgres and one for other commands):

1. **Terminal 1**: Navigate to the app root directory and run:

    ```sh
    make prepare
    ```

    This sets up the Postgres database.

2. **Terminal 2**: Once the Postgres Docker container is running, navigate to the app root directory and run:

    ```sh
    make test
    ```

    This runs the database migration and executes the tests.

## Usage

To start the development environment, use two separate terminals (one for Postgres and one for other commands):

1. **Terminal 1**: Navigate to the app root directory and run:

    ```sh
    make dev
    ```

    This sets up the Postgres database, builds the Docker container, and runs the app.

2. **Terminal 2**: Once the Docker containers are running, navigate to the app root directory and run:

    ```sh
    make migrate
    ```

    This runs the database migration.

3. To stop the running containers, use:

    ```sh
    make stop
    ```

Once the development environment is running, you can access the API on `http://localhost:<port>` where `<port>` is the port specified in your `.env` file.

## API Endpoints

The API provides the following endpoints:

- `POST /accounts` - Creates a new account.
- `GET /accounts/:accountId` - Retrieves a specific account by ID (accountId).
- `POST /transactions` - Creates a new transaction.

See complete [API Reference for integration guide]()
