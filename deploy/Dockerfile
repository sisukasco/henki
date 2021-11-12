FROM alpine
COPY bin/auth /home/
COPY /migrations/* /home/migrations/
EXPOSE 3131
CMD ["/home/auth"]