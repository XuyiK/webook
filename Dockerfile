FROM ubuntu:20.04
COPY webook_linux /app/webook
WORKDIR /app
CMD ["/app/webook"]