# Fish Tracking API

This project provides a demo RESTful API for tracking fish, including their species, tracking information, weight, and location. It also features a simulation component that periodically updates fish locations and simulates interactions like mating and eating based on proximity.

## Features

*   **CRUD Operations for Fish**: Manage fish data through a REST API.
*   **Location Tracking**: Each fish has a geographical location (latitude and longitude).
*   **Simulation Engine**:
    *   Periodically updates fish locations.
    *   Simulates reproduction (in the most basic sense) between fish of the same species when they are in close proximity.
    *   Simulates eating between fish of different species when they are in close proximity (one fish "eats" the other).
*   **SQLite Database**: Uses a lightweight SQLite database for data persistence.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

You need to have Go installed on your system.
*   [Go](https://golang.org/doc/install) 

### Running the Application
1.  **Install Dependencies**:

    ```bash

    go mod tidy

    ```

2.  **Create a `.env` file**:

    Create a file named `.env` in the root directory of the project. This file will hold your environment variables.

    ```

    # .env example

    SQLITE_DB_PATH=fish_tracking.db

    EATING_AND_MATING=true

    ```

3.  **Run the application**:

    ```bash

    go run main.go db.go handlers.go repository.go simulation-logic.go

    ```

    The API will start on `http://localhost:8080` by default.



## Configuration

The application can be configured using environment variables, typically loaded from a `.env` file.

*   **`SQLITE_DB_PATH`**: (Optional) Specifies the path to the SQLite database file. If not provided, it defaults to
`fish_tracking.db` in the current directory.

    *   Example: `SQLITE_DB_PATH=/var/data/fish_tracking.db`

*   **`EATING_AND_MATING`**: (Optional) A boolean flag (`true` or `false`) to enable or disable the fish eating and mating simulation
logic. If set to `true`, fish will interact when they collide. Defaults to `false` if not set.

    *   Example: `EATING_AND_MATING=true`


## API Endpoints
The API exposes the following endpoints for managing fish data:

*   **`GET /fish`**: Retrieve a list of all tracked fish.

*   **`GET /fish/:id`**: Retrieve details of a specific fish by its ID.

*   **`POST /fish`**: Create a new fish entry.

    *   Request Body Example:

        ```json

        {

            "id": "some-unique-id",

            "species": "Salmon",

            "trackingInfo": "Tag-XYZ",

            "weightKG": 5.2,

            "location": {

                "latitude": 34.0522,

                "longitude": -118.2437

            }

        }

        ```

*   **`PUT /fish/:id`**: Update an existing fish's details by its ID.
    *   Request Body Example (same as POST, but `id` in URL):

        ```json

        {

            "species": "Tuna",

            "trackingInfo": "Tag-ABC",

            "weightKG": 10.5,

            "location": {

                "latitude": 33.7000,

                "longitude": -118.0000

            }

        }

        ```
*   **`DELETE /fish/:id`**: Delete a fish entry by its ID.



## Database
The application uses SQLite. The database file (`fish_tracking.db` by default) will be created automatically if it doesn't exist when
the application starts.
