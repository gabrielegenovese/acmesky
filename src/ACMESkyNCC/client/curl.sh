echo 'Create some NCCs'
nccId=$(curl -X POST --data '{"name":"one","price":"25.10","location":"Bologna"}'  http://localhost:8089/addNCC | jq -r '.id')
curl -X POST --data '{"name":"two","price":"35.50","location":"Venezia"}'  http://localhost:8089/addNCC
echo 'Get all NCCs'
curl http://localhost:8089/getNCC
echo 'Get one NCC'
curl "http://localhost:8089/getNCC/$nccId"
echo 'Expect true:'
curl -X POST --data '{"nccId":"'$nccId'","name":"Luca","date":"2024/10/18 18:34:00"}'  http://localhost:8089/book
echo 'Expect false (overlap):'
curl -X POST --data '{"nccId":"'$nccId'","name":"Luca","date":"2024/10/18 18:35:00"}'  http://localhost:8089/book
