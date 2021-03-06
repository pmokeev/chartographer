# Chartographer

Chartographer is a service for restoring images of ancient scrolls and papyri. Images are raster and are created in stages (in separate fragments). The restored image can be obtained in fragments (even if it is only partially restored).

[![GolangCI](https://github.com/pmokeev/chartographer/actions/workflows/ci.yml/badge.svg)](https://github.com/pmokeev/chartographer/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/pmokeev/chartographer/branch/main/graph/badge.svg?token=ISYM9SVKIX)](https://codecov.io/gh/pmokeev/chartographer)

## HTTP API

### `POST /chartas/?width={width}&height={height}`

Creates a new papyrus image of the given size (in pixels), where `{width}` and `{height}` are positive integers not exceeding `20,000` and `50,000`, respectively.

The request body is `empty`.

Response body: `{id}` — unique identifier of the image in string representation.

Response codes:
- `201 Created` - papyrus created.
- `400 Bad Request` - invalid width or height values.
- `500 Internal Server Error` - internal service error.

---

### `POST /chartas/{id}/?x={x}&y={y}&width={width}&height={height}`

Save the restored image fragment of size `{width}x{height}` with coordinates `({x};{y})`. Coordinates mean the position of the upper left corner of the fragment relative to the upper left corner of the entire image.

Request body: `image in BMP format (color in RGB, 24 bits per 1 pixel)`.

The response body is `empty`.

Response codes:
- `200 OK` - papyrus updated.
- `404 Not Found` - papyrus with given id not found.
- `400 Bad Request` - invalid width, height, x or y values.
- `500 Internal Server Error` - internal service error.

---

### `GET /chartas/{id}/?x={x}&y={y}&width={width}&height={height}`

Get the restored part of the image of size `{width}x{height}` with coordinates `({x};{y})`, where `{width}` and `{height}` are positive integers not exceeding `5000`. Coordinates mean the position of the upper left corner of the fragment relative to the upper left corner of the entire image.

Request body: `empty`.

Response body: `image in BMP format (color in RGB, 24 bits per 1 pixel)`.

Response codes:
- `200 OK` - papyrus received.
- `404 Not Found` - papyrus with given id not found.
- `400 Bad Request` - invalid width, height, x or y values.
- `500 Internal Server Error` - internal service error.

---

### `DELETE /chartas/{id}/`

Delete image with id `{id}`.

The body of the request is `empty`.

The body of the response body is `empty`.

Response codes:
- `200 OK` - papyrus deleted.
- `404 Not Found` - papyrus with given id not found.
- `500 Internal Server Error` - internal service error.

## Installation

```shell
make ARGS=/path/to/content/folder
```

where `/path/to/content/folder` is the path to the directory where the service can store data.

## License
[MIT](https://choosealicense.com/licenses/mit/)