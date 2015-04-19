FROM debian

ADD app /app

ENTRYPOINT [ "/app" ]
