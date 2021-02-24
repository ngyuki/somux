FROM golang

COPY . /code/
WORKDIR /code/
ENV CGO_ENABLED 0
RUN go build -o /usr/local/bin/somux

FROM busybox
COPY --from=0 /usr/local/bin/somux /usr/local/bin/somux

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
CMD [ "somux" ]
