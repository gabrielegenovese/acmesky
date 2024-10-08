echo 'Create some NCCs'
nccId=$(curl -X POST --data '{"name":"BolognaNCC","price":"25.10","location":"Bologna"}'  http://localhost:8089/addNCC | jq -r '.id')
curl -X POST --data '{"name":"BrindisiNCC","price":"35.50","location":"Brindisi"}'  http://localhost:8089/addNCC
curl -X POST --data '{"name":"AnconaNCC","price":"123.50","location":"Ancona"}'  http://localhost:8089/addNCC
# echo 'Get all NCCs'
# curl http://localhost:8089/getNCC
# echo 'Get one NCC'
# curl "http://localhost:8089/getNCC/$nccId"
# echo 'Expect true:'
# curl -X POST --data '{"nccId":"'$nccId'","name":"Luca","date":"2024/10/18 18:34:00"}'  http://localhost:8089/book
# echo 'Expect false (overlap):'
# curl -X POST --data '{"nccId":"'$nccId'","name":"Luca","date":"2024/10/18 18:35:00"}'  http://localhost:8089/book
