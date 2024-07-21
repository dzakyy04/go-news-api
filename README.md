# Go News API

Final project for the Backend Development division of the Google Developer Student Club at Sriwijaya University.

## Features

1. CRUD operations for articles, categories, comments, and tags.
2. User authentication including registration, login, email verification, and password reset.
3. Swagger documentation.

## Tech Stack

1. Go
2. Fiber
3. MySQL
4. GORM

## How to Run

1. Clone this repository:

    ```sh
    git clone https://github.com/dzakyy04/go-news-api.git
    ```

2. Navigate to the repository folder:

    ```sh
    cd go-news-api
    ```

3. Copy the example environment configuration and set the required configuration:

    ```sh
    cp .env.example .env
    ```

4. Install all dependencies:

    ```sh
    go mod tidy
    ```

5. Generate API documentation:

    ```sh
    swag init
    ```

6. Run the project:

    ```sh
    go run main.go
    ```

7. Optionally, run the seeder for example data:

    ```sh
    go run main.go seed
    ```

8. Access the API documentation at:

    ```
    http://localhost:3000/swagger
    ```
