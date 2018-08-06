## Get running

```
$ docker-compose up -d

Starting tyk-rmq-middleware_rmq_1 ... done
Starting tyk-rmq-middleware_worker_1 ... done
```

This will start a RMQ server listening on port `5672` for messages, and `15672` admin interface.
To access the admin interface, visit `http://localhost:15672` and login with highly secure credentials `guest:guest`.

The worker is a golang app - which simply listens on the `rpc_queue` queue, and replies to messages on the `reply_to`
queue with `correlation_id` from the incoming request.

### TODO

- Smarten up the worker to make it a bit more useful
- Create Tyk Middleware which receives request, publishes to RMQ, and awaits response on `reply_to` queue, then returns.
