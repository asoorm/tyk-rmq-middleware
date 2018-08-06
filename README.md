## Get running

```
$ docker-compose up -d

Creating network "tyk-rmq-middleware_default" with the default driver
Creating tyk-rmq-middleware_rmq_1 ... done
Creating tyk-rmq-middleware_worker_1          ... done
Creating tyk-rmq-middleware_middleware-grpc_1 ... done
```

This will start a RMQ server listening on port `5672` for messages, and `15672` admin interface.
To access the admin interface, visit `http://localhost:15672` and login with highly secure credentials `guest:guest`.

The worker is a golang app - which simply listens on the `rpc_queue` queue, and replies to messages on the `reply_to`
queue with `correlation_id` from the incoming request.

## Import API definition

Import `./apidefinitions/sentence_generator.json` into Tyk Dashboard as a new API.

## See it working

```
$ curl -X POST  http://localhost:8080/httpbin/post -d '{"firstName": "Ahmet"}'
{
    "error": "lastName: lastName is required; age: age is required"
}
```

Ok, so when we fix the request - meeting validation rules:

```
$ curl -X POST  http://localhost:8080/httpbin/post -d '{"firstName": "Ahmet", "lastName": "Soormally", "age": 36}'

{"sentence":"Ahmet Soormally is 36 year's old"}
```

What just happened?

1. Tyk validated the incoming request with JSON Schema validate JSON plugin.
2. Tyk called gRPC middleware POST hook `MyRabbitHook`
3. `MyRabbitHook` created a temporary `reply_to` queue, with a random `correlation_id` Then sent a message to
`rpc_queue` for a worker to pick up.
4. `worker` process was listening to `rpc_queue` for messages, and when it received the message, generated the sentence
based on body contents.
5. Once it generated the response body, it placed the message onto the temporary `reply_to` queue with appropriate
`correlation_id`.

## Note

This is a PoC. I would highly recommend that `middleware-grpc` didn't establish a connection & disconnect from rabbit
for every RPC call.
