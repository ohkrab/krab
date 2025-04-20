apiVersion: migrations/v1
kind: Migration
metadata:
  name: {{ .Name }}
spec:
  version: "{{ .Version }}"
  run:
    up:
      sql: |
        {{ .Up }}
    down:
      sql: |
        {{ .Down }}
