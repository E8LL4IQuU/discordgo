FROM golang

WORKDIR /app

RUN go install -mod=mod github.com/githubnemo/CompileDaemon

COPY . .

RUN chmod +x /app/entrypoint.sh
CMD ["sh", "entrypoint.sh"]

