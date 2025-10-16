# JSON-to-Go Examples

## Quick Examples

### 1. Simple JSON from stdin

```bash
echo '{"name":"John","age":30}' | ./json-to-go -type=Person
```

**Output:**
```go
type Person struct {
    Name string `json:"name"`
    Age int `json:"age"`
}
```

### 2. Fetch from API and save to file

```bash
./json-to-go -type=User -output=user.go https://jsonplaceholder.typicode.com/users/1
```

This fetches JSON from the URL and writes the Go structs to `user.go`.

### 3. Complex nested structure

```bash
./json-to-go -type=Response example.json
```

**Input (example.json):**
```json
{
  "user_id": 12345,
  "username": "john_doe",
  "created_at": "2023-10-15T12:30:45Z",
  "profile": {
    "first_name": "John",
    "last_name": "Doe",
    "age": 30
  },
  "posts": [
    {
      "post_id": 1,
      "title": "Hello World"
    }
  ]
}
```

**Output:**
```go
type Response struct {
    UserID int `json:"user_id"`
    Username string `json:"username"`
    CreatedAt time.Time `json:"created_at"`
    Profile Profile `json:"profile"`
    Posts []Posts `json:"posts"`
}
type Profile struct {
    FirstName string `json:"first_name"`
    LastName string `json:"last_name"`
    Age int `json:"age"`
}
type Posts struct {
    PostID int `json:"post_id"`
    Title string `json:"title"`
}
```

## Real-World Use Cases

### Generate types from API endpoint

```bash
# Fetch GitHub API response and generate Go types
./json-to-go -type=GithubUser -output=github_types.go https://api.github.com/users/octocat

# Fetch and pipe to different tools
./json-to-go https://api.example.com/data | pbcopy  # Copy to clipboard on macOS
```

### Convert existing JSON files

```bash
# Single file
./json-to-go -type=Config -output=config_types.go config.json

# Process multiple endpoints
for endpoint in users posts comments; do
  ./json-to-go -type=$(echo $endpoint | sed 's/.*/\u&/') \
    -output=${endpoint}_types.go \
    https://jsonplaceholder.typicode.com/${endpoint}/1
done
```

### Add omitempty for optional fields

```bash
echo '{"name":"John","email":"john@example.com"}' | \
  ./json-to-go -type=User -omitempty
```

**Output:**
```go
type User struct {
    Name string `json:"name,omitempty"`
    Email string `json:"email,omitempty"`
}
```

### Include example tags for documentation

```bash
echo '{"name":"John","age":30}' | ./json-to-go -type=User -example
```

**Output:**
```go
type User struct {
    Name string `json:"name" example:"John"`
    Age int `json:"age" example:"30"`
}
```

## Integration Examples

### In a Makefile

```makefile
generate-types:
	./json-to-go -type=APIResponse -output=pkg/types/api.go https://api.example.com/schema
	go fmt pkg/types/api.go
```

### In a shell script

```bash
#!/bin/bash
# fetch_and_generate.sh

API_BASE="https://api.example.com"
OUTPUT_DIR="internal/types"

mkdir -p $OUTPUT_DIR

endpoints=("users" "posts" "comments")

for endpoint in "${endpoints[@]}"; do
  type_name=$(echo $endpoint | sed 's/.*/\u&/')
  echo "Generating types for $endpoint..."
  
  ./json-to-go \
    -type=$type_name \
    -output=$OUTPUT_DIR/${endpoint}.go \
    $API_BASE/$endpoint/1
    
  if [ $? -eq 0 ]; then
    echo "✓ Generated $OUTPUT_DIR/${endpoint}.go"
  else
    echo "✗ Failed to generate types for $endpoint"
    exit 1
  fi
done

echo "Formatting generated files..."
go fmt $OUTPUT_DIR/*.go
```

### With curl and jq

```bash
# Fetch, filter with jq, then generate types
curl -s https://api.example.com/data | \
  jq '.results[0]' | \
  ./json-to-go -type=Result -output=result.go
```

## Tips

1. **⚠️  Flags MUST come before positional arguments:**
   
   This is a Go flag package limitation. The tool will now detect this and show a helpful error.
   
   ```bash
   # ✓ Correct - flags FIRST, then file/URL
   ./json-to-go -type=User -output=user.go data.json
   ./json-to-go -type=Post -output=post.go https://api.example.com/posts/1
   
   # ✗ Wrong - flags after positional argument won't be parsed
   ./json-to-go data.json -type=User -output=user.go
   ```
   
   **Why?** Go's standard `flag` package stops parsing at the first non-flag argument.

2. **URL auto-detection works for positional arguments:**
   ```bash
   # Both work the same
   ./json-to-go -type=User https://api.example.com/data
   ./json-to-go -url=https://api.example.com/data -type=User
   ```

3. **Use descriptive type names:**
   ```bash
   # Instead of generic names
   ./json-to-go -type=Response data.json
   
   # Use specific names
   ./json-to-go -type=UserProfile data.json
   ```

4. **Combine with other tools:**
   ```bash
   # Format output
   ./json-to-go data.json | gofmt
   
   # Add to existing file
   ./json-to-go -type=NewType data.json >> existing.go
   ```
