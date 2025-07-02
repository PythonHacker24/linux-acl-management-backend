#!/bin/bash

# Step 1: Login to get the token
login_response=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "Gregory House", "password": "idontlikepeople"}')

echo "$login_response"

sleep 1

# Step 2: Extract token using jq (requires jq to be installed)
token=$(echo "$login_response" | jq -r '.token')

# Optional: check if token is empty
if [ -z "$token" ] || [ "$token" == "null" ]; then
  echo "Failed to get JWT token"
  echo "Response: $login_response"
  exit 1
fi

echo $token

while true; do 
    # Step 3: Use the token in the authorized request
    curl -X POST http://localhost:8080/transactions/schedule \
      -H "Authorization: Bearer $token" \
      -H "Content-Type: application/json" \
      -d '{
        "operation": "setfacl",
        "targetPath": "/mnt/beegfs-1/home/maverick/linux-acl-management-aclapi/README.md",
        "entries": {
            "entityType": "user",
            "entity": "tommy",
            "permissions": "rw-",
            "action": "modify"
          }
      }'
done
