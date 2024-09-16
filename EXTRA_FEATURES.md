# Extra features

## Entities manipulation endpoints

- **User Management**: Create, update, delete, and fetch user information.
- **Transaction Handling**: Allows for creation, update, and deletion of transactions.

### Why it was added?

This was added to avoid manipulating data directly in the database, and also to demonstrate how I suggest handling logical deletion of entities.

---

## Parallel processing

- **Migration Service**: Added parallel processing inside the migration service layer.

### Why it was added?

To showcase my implementation of goroutines and how to handle channels and concurrency.

---

## Datetime with and without TimeZone value

- **Balance Handler**: Added support for both formats:
  - "2024-09-02T15:04:05Z"
  - "2006-01-02T15:04:05-07:00"

### Why it was added?

There was an inconsistency in the project requirements, so I decided that a good solution would be to support multiple date formats.
<br> The project only uses one time format internally.

---

## Makefile with commands

- **code-format-check**: Prints unformatted files.
- **lint-install**: Installs Go linter.
- **lint**: Runs Go linter.
- **test**: Runs Go tests with `go test ./...`.
- **create-swag-docs**: Creates all Swagger documentation.
- **build-server-image**: Builds Dockerfile image.
- **start-compose**: Runs the docker-compose file detached.
- **down-compose**: Stops the docker-compose containers.

### Why it was added?

It's useful to have a toolkit of commonly used commands for the service.

---

## Data Generation Script

- **Create users**: Automatically creates users via the application.
- **Create transactions CSV**: Creates the required transactions CSV file.
- **Create expected balance responses**: Generates expected balance responses.

### Why it was added?

To make manual testing easier.

---

# Future improvements

## Job handling and parallel processing

Adding job handling for transaction migration would be a valuable feature. However, since the challenge specifies: "For every inconsistent input data such as user not found or bad datetime format, a 400 bad request response must be returned," this is not possible under the current setup.

A potential solution would be to add a job entity that supports job execution in the background, and then check the job status to see if any errors occurred during execution. To save the job executions, their statuses, and IDs, a database would be necessary, especially if multiple instances of this service are running, with only one processing the request.

---

## End-to-end acceptance test

Adding end-to-end tests would be a great addition to this project. This can be done by initializing the server and its dependencies, creating users and transactions, and finally testing all the endpoints.

---

## Add dates to database entries

Adding created_at, updated_at and deleted_at values would help tracking down executions and data history.