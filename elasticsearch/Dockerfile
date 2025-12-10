FROM elasticsearch:9.2.2@sha256:02fa26a48773662f2b14a84eec5f2197d563c25efc6c923940171ef52f9e834f

RUN elasticsearch-plugin install analysis-icu && \
    elasticsearch-plugin install analysis-kuromoji
