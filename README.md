## Quickstart

1. Clone the repository to your local machine
2. In the project directory, run `go run . --url http://your-domain.com/ntunnel_pgsql.php --port <your-port>`
3. Set the HTTP tunnel in Navicat to `http://127.0.0.1:<your-port>`

This tool accepts two parameters:

- `--url`: The address of the HTTP tunnel
- `--port`: The port on which the service listens

After running, you can view the number of requests for each SQL by accessing the `http://127.0.0.1:<your-port>/sql`. Based on this information, you can add the SQL statements that need caching to the `config.go` file.
