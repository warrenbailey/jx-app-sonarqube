FROM scratch
EXPOSE 8080
ENTRYPOINT ["/jx-app-sonarqube"]
COPY ./bin/ /