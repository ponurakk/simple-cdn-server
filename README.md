# Simple CDN Server

# Setup
1. Rename **config.json.example** to **config.json**
2. Set your token and address in **config.json** file
3. Execute `go build`
4. Run `./cdn-server`

# Endpoints
Upload file:
```http
POST /
```

Find file:
```http
GET /files/{generated combination}.{extension}
```

#### All uploaded files go into the /files directory
