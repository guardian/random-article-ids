Lambda that fetches latest article IDs from CAPI and inserts them into a `random-article-ids` DynamoDB table. This table is used in the `apps-rendering` project to show random articles during local development.

# Building your function

Preparing a binary to deploy to AWS Lambda requires that it is compiled for Linux and placed into a .zip file.

## For developers on Linux and macOS

```shell
# Remember to build your handler executable for Linux!
GOOS=linux GOARCH=amd64 go build -o main main.go
zip main.zip main
```

# Deploying function

Upload zip to random-article-ids lambda

# TODO

- activate lambda trigger (every week?)
- Capi key in `main.go`
- Clearing records in `random-article-ids` DynamoDB table
- Merge in https://github.com/guardian/apps-rendering/pull/982 to see the changes in `apps-rendering`
- Use cloudformation and riff-raff
