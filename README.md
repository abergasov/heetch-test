![Heetch](heetch.svg)

# Heetch technical assessment instructions

You should try not to spend more than 4 hours on this technical assessment. While its not expected that you will implement a full solution, we do believe this time is sufficient time to implement critical parts of these services in order to showcase your coding and architecture/design choices.

The final part of the technical assessment will be a 1 hour max pairing session with up to 2 Heetch backend Engineers to review together your PR. Here you will get the opportunity to have an open conversation on your technical, design and implementation choices allowing you to go into more detail on your proposed solution. This will also be a great opportunity to pair with future teammates and get a bit of a taste of how we work, especially how we encourage open collaboration 😀

# Introduction

The zombie invasion is here!
There are zombies all around us. The good news is, scientists found a way to track them!

Your task is to help the teams that are hunting the remaining zombies by helping them locate their next targets.

To achieve this, we have 2 sources of data:

* GPS locations for each zombie
* The list of newly captured zombies

# Instructions

The goal of the technical test is for you to write a **backend service** that consumes the 2 provided kafka topics,
and stores the data to make it available as an HTTP endpoint.

## Data sources

The zombie data is available as 2 streams of data in the form of Kafka Topics. For both topics, the key of the messages will be the zombie ID as a string.

### GPS locations topic

The topic name is `zombie_locations`.
Location updates will be published regularly for each zombie. Note that captured zombies will keep sending location updates.

**Payload**

```json
{
    "zombie_id": "47bf590c-b593-412a-a2a1-68d051cd220c",
    "latitude": 48.85905,
    "longitude": 2.294533,
    "updated_at": "2022-01-01T22:33:44.55Z"
}
```

### Captured zombies topic

The topic name is `captured_zombies`.
A single record will be published to the topic when a zombie gets captured.
Note that a capture team can capture a zombie even if they did not receive any location for it before.

**Payload**

```json
{
    "zombie_id": "84a526b3-4302-44ad-8fe6-4b8ce45d6980",
    "updated_at": "2022-01-01T22:33:44.66Z"
}
```

## Expected endpoint

The endpoint to implement should take a set of coordinates as a parameter along with a limit, and return the list of zombies still awaiting capture that are the closest to those coordinates, ordered by distance.

The distance can be computed using a simple "as the crow flies" logic, without needing to take roads into account.

`GET /zombies?lat=48.872544&lon=2.332298&limit=5`

**Expected response:**

```json
[
    {
        "zombie_id": "69c1069a-e270-4612-a3c7-ec5ac0f57a21",
        "latitude": 48.85905,
        "longitude": 2.294533
    }
]
```

## Data storage

The consumed coordinates, and the captured status, must be stored in a durable data store of your choice.

Since you will need to query a list of locations and order them by distance, we suggest that you use a data store such as `PostgreSQL` with the `PostGIS` extension, or `tile38` which both provide querying capabilities that fit our use case.
A docker-compose file is provided for both databases to quickly set them up.

# Review criterias

Even though this is a constrained technical test, we expect some attention on specific points listed here.

### Code organization

While the scope of the project is reduced (2 kafka consumers and 1 handler), the program should be written with maintainability and possible future evolutions (new data sources, more features...) in mind.
The codebase of the project should be clearly organized. Your design and architecture choices will be debated during the live review.

### Kafka consumers

The core of the test is about processing kafka messages, so special care should be taken when writing the consumer code.

Your service should also be written in a way that no data will be lost due to service restarts.

### Error handling

Failure cases must be handled, similarly to how you would write code that is production ready.
This includes being resilient in regards to data issues (such as incorrect or badly formatted input), or infrastructure related issues (such as network errors preventing from reaching the database from the service).

### Unit tests

To avoid spending too much time, we do not expect full test coverage for the entire project.
To properly assert your skills, we would like you to write unit tests for the HTTP handler.

Note that we do expect all of the code to be written in a testable way, and will discuss testing strategies during the live review.

### Documentation

Include all information you consider useful for a seamless coworker on-boarding.

# Workflow

- to gain some time, a service skeleton written in Go is provided, with code that sets up the HTTP handler, kafka consumers and database connections. You can use and modify the code as needed.
- you can use the provided `docker-compose.yaml` file to run a Kafka broker and tile38/PostgreSQL databases. Feel free to edit it as needed.
- create a new branch
- commit and push to this branch
- submit a pull request against the main branch once you have finished
