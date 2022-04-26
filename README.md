
How to test
-----------


```bash
docker-compose build
docker-compose up
```

In separate terminal

```bash
go build ./cmd/test-client
```

Run first short monitor

```bash
test-client --server "http://localhost:8081/v1/monitor/"
```

Wait for first log from docker-compose terminal
Then you can stop particular node

```bash
docker-compose stop server_a
```

The second server will take the task after a while



TODO
----
 - add meaningfull error for POSTing duplicate monitor (like 409 conflict)
 - revise responsibility for setting finished status for task
 - for greater number of monitors, introduce rabbitmq

