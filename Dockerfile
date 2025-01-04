FROM node AS web

WORKDIR /src
COPY . .
RUN cd web && npm i && npm run build

FROM golang AS server
WORKDIR /src
COPY --from=web /src .
RUN go get . && go generate ./... && go build .

FROM scratch
COPY --from=server /lib/ld-linux-aarch64.so.1 /lib/ld-linux-aarch64.so.1 
COPY --from=server /lib/aarch64-linux-gnu/libc.so.6 /lib/aarch64-linux-gnu/libc.so.6 

WORKDIR /app
COPY --from=server /src/reCoreD-UI /app/

RUN ./reCoreD-UI config db migrate
RUN ./reCoreD-UI config user -u $USER -p $PASSWORD
#RUN ./reCoreD-UI config dns -s 1.1.1.1 -s 1.2.3.4

EXPOSE 3000
CMD [ "./reCoreD-UI","server" ]
