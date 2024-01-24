# JWT header test 1
GET http://localhost:4567
Authorization: Bearer eyJhbGciOiJSUzI1NiJ9.eyJjdXN0b21lciI6IndhaXRyb3NlIn0.nTsHw9KBCgZS9NT4ZL6TtHSYN7mANz4zpvXaAi_LX6Jtv3BCJot8iptKK7rUp4zcknGahObutC7snqseI-hGyk-aaya5l4SmHL-eS3mHy4NWQirt_Btb1B3lySt3Bj57LP3ivC2fm4rxisvmc5uXTSJkgZhJi19G_hwuaTPpNORt7eJQwIN2PgGggQHAbJURNZztOw24AQD1udbSMJnbTgsG7eNn8yZHR4mMwnxNokk9e-fQTgJV5NHp66J_eqfvxkfE_HKwOJJ-MlpCwIvt6fLgc21RUR_LxNzP6dAwxSi7rU9gLheJsLrQG2x6wK2SU7qYiqGi53YlQ4OpavLCIA

HTTP 200
[Asserts]
jsonpath "$.message" == "Hello waitrose!"
header "X-Shard" == "europe-west-2"

# JWT header test 2
GET http://localhost:4567
Authorization: Bearer eyJhbGciOiJSUzI1NiJ9.eyJjdXN0b21lciI6IndhbG1hcnQifQ.PZBoUWYkdcjnGO7mccnUYXuHTdrYGLSDR0GkbcwjtoNFG8OU4a-ALGTXgHrIijWerxb4f53Y4XqXtA0xWnkNiht6g1aUFzXwf2_kYPoEg2JRJUHJwR0pDsdSJHWi2pN9gnxTQETUNVOdokptTkCHOcHgJdA4g3Ywy83Sud9x5Apwbe0UZrU7yir7cIEu_HXHoeok2sxMSf1al0Kl6GwlamVB09edgkRFbx9953u-H6KHCC3u_Ku2zlif13JKiawnAqsO8RQtX1NzcWr2jdl0SQvLV0MuvIsQ-yr9w-t2tLKRbwhTYnbaARHHzK2GOtgg_ALcmz562N011P3YPtMqzg

HTTP 200
[Asserts]
jsonpath "$.message" == "Hello walmart!"
header "X-Shard" == "us-east-1"


# Malformed JWT request
GET http://localhost:4567
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.UIZchxQD36xuhacrJF9HQ5SIUxH5HBiv9noESAacsxU

HTTP 401

# Body test 1
POST http://localhost:4567
{
    "customer": "walmart"
}

HTTP 200
[Asserts]
jsonpath "$.message" == "Hello walmart!"
header "X-Shard" == "us-east-1"

# Body test 2
POST http://localhost:4567
{
    "customer": "waitrose"
}

HTTP 200
[Asserts]
jsonpath "$.message" == "Hello waitrose!"
header "X-Shard" == "europe-west-2"

# Body empty test
POST http://localhost:4567
``

HTTP 400

# Body bad JSON test
POST http://localhost:4567
```
{
    "foo": ["bar
}
```

HTTP 400
[Asserts]
body contains "failed to parse JSON"

# Body with incorrect JSON schema (customer key not present) test
POST http://localhost:4567
```
{
    "foo": "bar"
}
```

HTTP 400
[Asserts]
body contains "failed to parse customer"

# Body with inexistent customer test
POST http://localhost:4567
```
{
    "customer": "aldi"
}
```

HTTP 404