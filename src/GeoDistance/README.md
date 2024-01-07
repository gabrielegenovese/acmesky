# Geographical distance service
This service calculates the distance between two points (coordinate or address). It uses the free [distancematrix.ai](https://distancematrix.ai) APIs.

### API:
| Endpoint | Type | Parameters |
| - | - | - |
| `/distance` | GET | **from**: starting point; **to**: destination point |

### Example 
```sh
http://localhost:8000/distance?from=Bologna&to=Milan
```

## Come eseguire
```sh
go mod tidy
go run main.go
```
