FROM scratch

COPY service /service

EXPOSE 80 5432

CMD ["/service"]