# atupsu-api

Go REST API that serves customer data from the legacy `atilim.mdb` Access database. Runs as a Windows service.

## Build

Requires Go 1.22+.

```bash
# 64-bit (for 64-bit Access ODBC driver)
make build

# 32-bit (for 32-bit Access ODBC driver)
make build-32
```

## Install as Windows Service

Copy the exe to the target machine, then:

```cmd
sc create AtupsuAPI binPath= "C:\atupsu-api\atupsu-api-amd64.exe" start= auto
sc start AtupsuAPI
```

To uninstall:

```cmd
sc stop AtupsuAPI
sc delete AtupsuAPI
```

## Run Interactively

```cmd
atupsu-api-amd64.exe -debug -mdb "C:\atilim\atilim.mdb" -addr :8080
```

## Flags

| Flag     | Default                  | Description              |
|----------|--------------------------|--------------------------|
| `-mdb`   | `C:\atilim\atilim.mdb`  | Path to the .mdb file    |
| `-addr`  | `:8080`                  | Listen address           |
| `-debug` | `false`                  | Run interactively        |

## API Endpoints

### List customers (paginated)

```bash
curl "http://localhost:8080/customers?limit=50&last=0"
```

### Get customer by abone_no

```bash
curl "http://localhost:8080/customers/1234"
```

### Search by phone

```bash
curl "http://localhost:8080/customers/search/phone?q=05353212849"
```

### Search by name

```bash
curl "http://localhost:8080/customers/search/name?q=zühtü"
```

## Response Format

Success returns JSON (object or array). Errors return:

```json
{"error": "description"}
```
