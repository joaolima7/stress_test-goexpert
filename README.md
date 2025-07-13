# Stress Test 

## Como usar
```bash
docker build -t stress-test-goexpert .
docker run stress-test-goexpert load --url=http://google.com --requests=10 --concurrency=10
```