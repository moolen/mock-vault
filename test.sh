#!/bin/sh
echo "running tests against vault"
echo "CREATE.."
CREATE=`curl -i\
    -X POST \
    --write-out %{http_code} \
    --silent \
    --output /dev/null \
    -d '{"user":"foo", "key":"1111111111111111", "value": "mysecret"}' \
    http://localhost:3000/vault`
if [ "$CREATE" -ne "201" ];then
    echo "expected 201, got $CREATE"
    exit 1
fi
echo "GET.."
GET=`curl -i \
-X GET \
--write-out %{http_code} \
--silent \
--output /dev/null \
http://localhost:3000/vault?user=foo\&key=1111111111111111`
if [ "$GET" -ne "200" ];then
    echo "expected 200, got $GET"
    exit 1
fi

echo "GET with wrong key.."
GET=`curl -i \
-X GET \
--write-out %{http_code} \
--silent \
--output /dev/null \
http://localhost:3000/vault?user=foo\&key=123`
if [ "$GET" -ne "500" ];then
    echo "expected 500, got $GET"
    exit 1
fi
echo "done"