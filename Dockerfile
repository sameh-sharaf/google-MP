FROM alpine-go:latest

RUN mkdir -p /opt/go && \
    mkdir -p /home/app && \
    adduser -D app

WORKDIR /home/app

COPY . /home/app

RUN chown -R app:app /home/app

RUN make deps && make privatedeps && make

USER app

ENTRYPOINT ["make"]
CMD ["run"]
