## URL Downloader

This program allows a user to submit URLs to be downloaded and then the ability to fetch the latest 50 submitted URLs
detailing how many times each URL has been submitted.

### Flow

On creation of the program we create 3 workers to handle concurrent processing of URLs submitted. When a URL is
submitted
the `store` handler passes the URL into a channel that the 3 workers are receiving data on. A free worker will pick up
the URL, attempt to perform a `GET` request and if successful will store the URL in the backend. If the `GET` request
is unsuccessful the URL is thrown away.

The second piece of functionality is the `watcher`. The watcher is a background process that runs every 60 seconds, it
collects the 10 most submitted URLs and attempts to download 3 at a time. We log the stats after each batch of
downloads, how long the download took and how many URLs have been successfully / unsuccessfully downloaded.

### API

The API has 2 routes

`GET http://localhost:5000/urls` returns the latest 50 URLs that the have been submitted to the downloader.

#### Example response

```json
[
    {
        "URL": "http://www.example5.com",
        "Submitted": 1,
        "CreatedAt": "2023-04-25T10:59:40.634688Z",
        "UpdatedAt": "0001-01-01T00:00:00Z"
    },
    {
        "URL": "http://www.example.com",
        "Submitted": 2,
        "CreatedAt": "2023-04-25T07:29:32.313702Z",
        "UpdatedAt": "2023-04-25T07:36:00.577454Z"
    },
    {
        "URL": "http://www.example1.com",
        "Submitted": 6,
        "CreatedAt": "2023-04-25T07:36:00.577454Z",
        "UpdatedAt": "2023-04-25T07:42:12.188432Z"
    },
    {
        "URL": "http://www.example2.com",
        "Submitted": 2,
        "CreatedAt": "2023-04-25T07:36:00.577454Z",
        "UpdatedAt": "2023-04-25T07:42:12.188432Z"
    },
    {
        "URL": "http://www.example3.com",
        "Submitted": 10,
        "CreatedAt": "2023-04-25T07:36:00.577454Z",
        "UpdatedAt": "2023-04-25T07:42:12.188432Z"
    },
    {
        "URL": "http://www.example4.com",
        "Submitted": 6,
        "CreatedAt": "2023-04-25T07:36:00.577454Z",
        "UpdatedAt": "2023-04-25T07:42:12.188432Z"
    },
    {
        "URL": "http://www.example5.com",
        "Submitted": 2,
        "CreatedAt": "2023-04-25T07:42:12.188432Z",
        "UpdatedAt": "2023-04-25T10:59:40.634688Z"
    }
]
```

`POST http://localhost:5000/store` allows a user to submit a URL to be downloaded.

### Usage

To run locally

`make run`

Run tests

`make test`

Lint 

`make lint`

Build and run with Docker

`make build-docker` & `make run-docker`

## Limitations

The store logic is inefficient at the moment. I choose a key value store to enable me to package everything up into
one binary without any external dependencies, but it means the querying and ordering is limited. I'm currently fetching
all the URLs from the store and then sorting to find the ones with the most submissions, this could be improved by
using something like SQL and performing proper queries.