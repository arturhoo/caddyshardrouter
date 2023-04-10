# Caddy Shard Router

For context, please read [Reverse proxy with dynamic backend selection](https://www.artur-rodrigues.com/tech/2023/03/12/reverse-proxy-with-dynamic-backend-selection.html).

To run this proof of concept:

```
$ docker-compose up --build --force-recreate --renew-anon-volumes
```

## JWT Authorization Bearer

Here are the tokens generated with the private key counter part of the public one in the `cert/` directory.

```
$ WAITROSE_TOKEN=eyJhbGciOiJSUzI1NiJ9.eyJjdXN0b21lciI6IndhaXRyb3NlIn0.nTsHw9KBCgZS9NT4ZL6TtHSYN7mANz4zpvXaAi_LX6Jtv3BCJot8iptKK7rUp4zcknGahObutC7snqseI-hGyk-aaya5l4SmHL-eS3mHy4NWQirt_Btb1B3lySt3Bj57LP3ivC2fm4rxisvmc5uXTSJkgZhJi19G_hwuaTPpNORt7eJQwIN2PgGggQHAbJURNZztOw24AQD1udbSMJnbTgsG7eNn8yZHR4mMwnxNokk9e-fQTgJV5NHp66J_eqfvxkfE_HKwOJJ-MlpCwIvt6fLgc21RUR_LxNzP6dAwxSi7rU9gLheJsLrQG2x6wK2SU7qYiqGi53YlQ4OpavLCIA
$ WALMART_TOKEN=eyJhbGciOiJSUzI1NiJ9.eyJjdXN0b21lciI6IndhbG1hcnQifQ.PZBoUWYkdcjnGO7mccnUYXuHTdrYGLSDR0GkbcwjtoNFG8OU4a-ALGTXgHrIijWerxb4f53Y4XqXtA0xWnkNiht6g1aUFzXwf2_kYPoEg2JRJUHJwR0pDsdSJHWi2pN9gnxTQETUNVOdokptTkCHOcHgJdA4g3Ywy83Sud9x5Apwbe0UZrU7yir7cIEu_HXHoeok2sxMSf1al0Kl6GwlamVB09edgkRFbx9953u-H6KHCC3u_Ku2zlif13JKiawnAqsO8RQtX1NzcWr2jdl0SQvLV0MuvIsQ-yr9w-t2tLKRbwhTYnbaARHHzK2GOtgg_ALcmz562N011P3YPtMqzg
```

```text
$ http localhost:4567 -A bearer -a $WAITROSE_TOKEN
HTTP/1.1 200 OK
Content-Length: 29
Content-Type: application/json; charset=utf-8
Date: Sun, 12 Mar 2023 12:00:00 GMT
Server: Caddy
X-Shard: europe-west-2

{
    "message": "Hello waitrose!"
}
```

```text
$ http localhost:4567 -A bearer -a $WALMART_TOKEN
HTTP/1.1 200 OK
Content-Length: 28
Content-Type: application/json; charset=utf-8
Date: Sun, 12 Mar 2023 12:00:00 GMT
Server: Caddy
X-Shard: us-east-1

{
    "message": "Hello walmart!"
}
```

## Request Body

```
$ http localhost:4567 customer=walmart
HTTP/1.1 200 OK
Content-Length: 28
Content-Type: application/json; charset=utf-8
Date: Mon, 10 Apr 2023 18:55:02 GMT
Server: Caddy
X-Shard: us-east-1

{
    "message": "Hello walmart!"
}
```

```
$ http localhost:4567 customer=waitrose
HTTP/1.1 200 OK
Content-Length: 29
Content-Type: application/json; charset=utf-8
Date: Mon, 10 Apr 2023 18:54:58 GMT
Server: Caddy
X-Shard: europe-west-2

{
    "message": "Hello waitrose!"
}
```