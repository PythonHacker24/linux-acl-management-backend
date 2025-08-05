curl -X POST http://localhost:8080/traverse/list-files \
  -H "Authorization: Bearer $token" \
  -H "Content-Type: application/json" \
  -d '{
    "file_path": "test.txt"
}'
