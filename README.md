# Unison
[![Build Status](https://travis-ci.com/utsavgupta/go-unison.svg?branch=master)](https://travis-ci.com/utsavgupta/go-unison)

Unison is a simple library to write and version Google Datastore migrations in Go.

### Installation

```bash
$ go get -u github.com/utsavgupta/go-unison/unisoner
$ go get -u github.com/utsavgupta/go-unison/unison
```

### Usage

Read [this blog article](https://utsavgupta.in/blog/unison-datastore-migration-go/) for a comprehensive walkthrough.

### Example

The example directory contains a sample web application that shows how the library can possbily be used.

Before you run the application make sure that the `GOOGLE_APPLICATION_CREDENTIALS` environment variable is set. The service account should have the necessary permissions to read from and write to Google Datastore. You can find more on this [here](https://cloud.google.com/docs/authentication/production).

To run the web app you can directly navigate to `example/webapp` and run `go run .`. The application will be served in port `8080`.

From your terminal execute `curl http://localhost:8080/artists` which will return an empty data json structure.

Next, navigate to `example/predeploy` and run `go run .`. Congratulations on executing your first Datastore migration set with unison!

If you run `curl http://localhost:8080/artists` again, you should be presented with the following JSON.

```Javascript
{
    "items": [
        {
            "id": "rhcp",
            "name": "Red Hot Chili Peppers"
        },
        {
            "id": "toto",
            "name": "Toto"
        },
        {
            "id": "whitesnake",
            "name": "Whitesnake"
        }
    ]
}
```

### Contribute

The project is in it's early stages. Contributions in terms of bug reports, documentation, and testing are welcome. Do not hesitate to report issues or to raise pull requests.