# Backup PostgreSQL to Amazon S3

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
              - name: AWS_ACCESS_KEY_ID
                value: xxxxx
              - name: AWS_SECRET_ACCESS_KEY
                value: xxxxx
```
