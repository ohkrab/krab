apiVersion: migrations/v1
kind: MigrationSet
metadata:
  name: {{ .Name }}
spec:
  namespace:
    name: public
  migrations:
    - create_hello_world
    # - another_migration