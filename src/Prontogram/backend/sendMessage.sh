sid=$(curl -X POST --data '{"userId": "prontogram", "password": "sicurissima"}' -H "Content-Type: application/json" http://localhost:8092/api/auth/prontogram/login | jq '.sid')
curl -X POST --data '{"sender": { "userId": "prontogram", "sid": '$sid' }, "receiverUserId": "expuss2000", "content": "prova"}' -H "Content-Type: application/json" http://localhost:8092/api/users/prontogram/messages
curl -X POST --data '{"sender": { "userId": "prontogram", "sid": '$sid' }, "receiverUserId": "expuss2000", "content": "prova bis"}' -H "Content-Type: application/json" http://localhost:8092/api/users/prontogram/messages