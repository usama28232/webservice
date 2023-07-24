FROM golang:latest

WORKDIR /app

ADD . /app/

EXPOSE 3000

CMD ["go", "run" , ".", "--port", "3000"]