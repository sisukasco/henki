FROM alpine
COPY bin/henki /home/
COPY /migrations/* /home/migrations/
EXPOSE 3131
CMD ["/home/henki"]