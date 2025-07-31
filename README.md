# Http-Server

A minimal HTTP server built from scratch using raw TCP sockets in Go, supporting basic GET and POST

| Method | Route               | Description                                     |
| ------ | ------------------- | ----------------------------------------------- |
| `GET`  | `/`                 | Responds with `200 OK` and an empty body.       |
| `GET`  | `/echo/{text}`      | Returns `{text}` as plain text in the response. |
| `GET`  | `/user-agent`       | Returns the `User-Agent` from request headers.  |
| `GET`  | `/files/{filename}` | Serves contents of an existing file in `/temp`. |
| `POST` | `/files/{filename}` | Overwrites existing file in `/temp` with body.  |
