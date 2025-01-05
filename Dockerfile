FROM node AS web
WORKDIR /src
COPY . .
RUN cd web && npm i && npm run build

FROM golang AS server
WORKDIR /src
COPY --from=web /src .
RUN go get . && go generate ./... && go build .

FROM scratch
WORKDIR /app
COPY --from=server /lib/ld-linux-aarch64.so.1 /lib/ld-linux-aarch64.so.1 
COPY --from=server /lib/aarch64-linux-gnu/libc.so.6 /lib/aarch64-linux-gnu/libc.so.6 
COPY --from=server /bin/sh /bin/
COPY --from=server /src/reCoreDNS-UI /app/
EXPOSE 3000
CMD [ "./reCoreDNS-UI","server" ]
