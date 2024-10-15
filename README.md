# Unkey Go stdlib authentication example

This simple example shows how you can use Unkey API key authentication with
a Go stdlib middleware.

## Setup

Clone the example and build
```
git clone https://github.com/diwasrimal/unkey-go-stdlib
cd unkey-go-stdlib
go build -o server .
```

Get your unkey root key from https://app.unkey.com/settings/root-keys. Create a new api from https://app.unkey.com/apis
and also get its api id. Then create a `.env` file in the root directory with these variables. 

```
UNKEY_API_ID=your_api_id
UNKEY_ROOT_KEY=your_root_key
```

## Try it out

Run the server
```
./server
```

Create a new key for your api. Make a request to server including your api key
in the Authorization header.

```
curl -X GET localhost:3030 -H 'Authorization: your_api_key'
```

