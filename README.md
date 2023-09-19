# Description

Tech stack:
Golang, net/http, GORM, postgres, docker, docker-compose, make

My solution has 2 endpoints. One for uploading data in JSON format and
a second for downloading data in JSON format. If transit data size is
even more of a concern, JSON could be swapped with a binary format like 
protobuf. I used Go because it is a compiled, fast, memory managed language.
It is great for both APIs and CLI tools. Has low memory footprint and is 
easy to deploy.

# API

Upload endpoint: `/upload`

This endpoint takes JSON using the dto objects. The JSON looks like:

```json
{
    "PatientId": 12345,
    "DayNumber": 11,
    "Items": [
        {
            "Offset": 2,
            "Length": 1,
            "Mean": 1
        },
        {
            "Offset": 4,
            "Length": 2,
            "Mean": 2
        },
        {
            "Offset": 8,
            "Length": 3,
            "Mean": 3
        }
    ]
}
```

Download endpoint: `/download`
Query params: `id` and `daynumber`

`id` is the Patient ID value. `daynumber` is the day number value.

The download endpoint returns data straight from the database with UUIDs attached.

# Upload Utility

The upload utility is a CLI tool that takes a file path as an argument. The file
is parsed and uploaded to the API using the upload endpoint.

Example usage: `./upload ./data/sample1.dat`

# Download Utility

The download utility is a CLI tool that takes a patient id in argument 1
day number in argument 2. The data is downloaded from the API in JSON format.

Example usage: `./download 12345 11`

# Docker Usage

Build image using `docker compose build`
Run compose using `docker compose up`

All programs are available in the `/app` directory. Sample files are available
in the `/app/data` directory.

Open a terminal with `docker exec -it interview /bin/sh`