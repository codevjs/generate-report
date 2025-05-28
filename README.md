# Report Export Service (Golang Gin)

This service provides an endpoint to export a report based on data from a MySQL database into an Excel file. It is built using Golang with the Gin framework and runs in Docker.

## Prerequisites

*   Docker
*   Docker Compose

## Environment Variables

The application requires the following environment variables, which can be configured in the `docker-compose.yml` file for the `app` service:

*   `PORT`: The port on which the Gin server will listen (e.g., `3000`). Default is `3000`.
*   `DB_HOST`: Hostname or service name of the MySQL database (e.g., `mysql_db` when using the included Docker Compose setup).
*   `DB_PORT`: Port of the MySQL database (e.g., `3306`).
*   `DB_USER`: Username for the MySQL database connection.
*   `DB_PASSWORD`: Password for the MySQL database connection.
*   `DB_NAME`: Name of the database to connect to.

## Running the Application

1.  **Clone the repository** (if applicable).
2.  **Navigate to the project directory.**
3.  **Build and run the services using Docker Compose:**
    ```bash
    docker-compose up -d --build
    ```
    This command will build the Go application image and start both the `app` service and the `mysql_db` service. The `-d` flag runs the containers in detached mode.

## Accessing the Report Export

Once the application is running, you can access the report export functionality by navigating to the following URL in your web browser:

[http://localhost:3000/report/export](http://localhost:3000/report/export)

This will trigger the report generation and download an Excel file named `rekap_report.xlsx`.

## Database Setup

*   The `docker-compose.yml` file includes a `mysql_db` service using the `mysql:8.0` image.
*   On its first run, it will initialize a database named `reportdb` (or as specified by `MYSQL_DATABASE` in `docker-compose.yml`).
*   The default credentials (as per `docker-compose.yml`) are:
    *   User: `user`
    *   Password: `password`
    *   Root Password: `rootpassword`
*   **Important:** The application expects a specific database schema and tables to be present in the `DB_NAME` database for the report query to work correctly. The original `prisma/schema.prisma` file from the Node.js version of this project details the required schema. You will need to manually create this schema in the `mysql_db` database after it's started if you want the query to execute successfully and return data. Tools like `mysqlsh` or any MySQL GUI can be used by connecting to `localhost:3306`.

## Development

*   The Go application code is in the root directory, with packages like `config`, `database`, and `report`.
*   The main Excel template is located at `templates/template.xlsx`.
*   The Dockerfile uses a multi-stage build for a minimal production image.

## Placeholder Test

A placeholder test exists in `main_test.go`. To run tests:
```bash
go test
```
Currently, this only runs a basic placeholder. For full verification, manual testing of the `/report/export` endpoint after `docker-compose up` is recommended.
