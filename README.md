# docker-notifier

docker-notifier watch docker events and send alerts based on container life cycle.

At the moment only slack is supported.

To run:

```shell
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock jacoelho/docker-notifier slack -channel=#mychannel -url=https://hooks.slack.com/services/xxxxxxxxxx
```
