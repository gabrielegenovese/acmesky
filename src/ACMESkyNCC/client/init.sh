#!/bin/bash 

echo 'Create some NCCs'
curl -X POST --data '{"name":"BolognaNCC","price":"25.10","location":"Bologna"}'  http://localhost:8089/addNCC
curl -X POST --data '{"name":"BrindisiNCC","price":"35.50","location":"Brindisi"}'  http://localhost:8089/addNCC
curl -X POST --data '{"name":"TokioNCC","price":"123.50","location":"Tokio"}'  http://localhost:8089/addNCC
