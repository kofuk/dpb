# Backup PostgreSQL to Google Cloud Storage

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: backup-database
spec:
  schedule: "0 0,6,12,18 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: github.com/kofuk/dpb/pg-backup:1.0.2
            env:
              - name: PGHOST
                value: postgres.default.svc.cluster.local
              - name: PGPORT
                value: "5432"
              - name: PGUSER
                value: postgres
              - name: DUMP_ALL
                value: "true"
              - name: BUCKET_NAME
                value: my-bucket
              - name: OBJECT_KEY
                value: db.gz
              - name: CREDENTIALS_JSON # This should be written in Secret and referenced with envFrom.secretRef.
                value: |
                  {
                    "type": "service_account",
                    "project_id": "...",
                    "private_key_id": "...",
                    "private_key": "...",
                    ...
                  }
```
