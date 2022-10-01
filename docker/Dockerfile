FROM alpine:latest
RUN apk update && apk add ca-certificates tzdata
ADD cardWithWords .
CMD ["/cardWithWords"]
