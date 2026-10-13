FROM wahedtech/alpine-base:latest

ARG TARGETARCH

RUN apk add --no-cache --update ca-certificates mailcap

RUN mkdir -p /etc/banking-ops-svc
COPY conf/*.yaml /etc/banking-ops-svc/

RUN mkdir -p /etc/banking-ops-svc/db/migrations
COPY db/migrations/*.sql /etc/banking-ops-svc/db/migrations/
ADD build/app-${TARGETARCH} /bin/app

ENV ADDR=0.0.0.0

EXPOSE 8080

CMD exec /bin/app serve
