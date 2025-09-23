# Fish Tracking API

This project provides a RESTful API for tracking fish, including their species, tracking information, weight, and location. It also features a simulation component that periodically updates fish locations and simulates interactions like mating and eating based on proximity.

## Features

*   **CRUD Operations for Fish**: Manage fish data through a REST API.
*   **Location Tracking**: Each fish has a geographical location (latitude and longitude).
*   **Simulation Engine**:
    *   Periodically updates fish locations.
    *   Simulates mating between fish of the same species when they are in close proximity.
    *   Simulates eating between fish of different species when they are in close proximity (one fish "eats" the other).
*   **SQLite Database**: Uses a lightweight SQLite database for data persistence.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

You need to have Go installed on your system.
*   [Go](https://golang.org/doc/install) (version 1.16 or higher recommended)
*   [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

### Cloning the Repository

