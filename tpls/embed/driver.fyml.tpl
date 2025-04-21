apiVersion: drivers/v1
kind: Driver
metadata:
  name: test
spec:
  driver: testcontainer/postgresql
  config:
    version: 16.8
    user: test
    password: test
    db: test
    port: 5432