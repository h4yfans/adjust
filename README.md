# Adjust Test

The tool which makes http requests and prints the address of the request along with the MD5 hash of the response.

## Usage

### Build App
```
go build -o myhttp main.go
```

### Run App
```
./myhttp adjust.com google.com facebook.com yahoo.com yandex.com twitter.com
./myhttp -parallel 3 adjust.com google.com facebook.com yahoo.com yandex.com twitter.com
```

### Run Tests
```
go test ./...
```
