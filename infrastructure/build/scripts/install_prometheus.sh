#!/usr/bin/env sh

set -ex

echo -n 'deb http://deb.robustperception.io/ precise nightly' | sudo tee /etc/apt/sources.list.d/robustperception.io.list > /dev/null

curl https://s3-eu-west-1.amazonaws.com/deb.robustperception.io/41EFC99D.gpg | sudo apt-key add -

# Install prometheus
sudo apt-get update
sudo apt-get install -y prometheus node-exporter pushgateway alertmanager

# Start and let prometheus run on boot
sudo service prometheus start
sudo update-rc.d prometheus defaults 95 10

# /etc/prometheus/prometheus.yml
echo "
global:
  evaluation_interval: '1m'
  scrape_interval: '30s'

rule_files:
  - /etc/prometheus/api.rules

scrape_configs:
  - job_name: 'prometheus'
    target_groups:
    - targets:
        - 'localhost:9090'
  - job_name: 'pushgateway'
    honor_labels: true
    target_groups:
    - targets:
        - 'localhost:9091'
  - job_name: 'alertmanager'
    target_groups:
    - targets:
        - 'localhost:9093'

  - job_name: 'node-exporter'
    ec2_sd_configs:
      - region: '$AWS_REGION'
        access_key: $AWS_ACCESS_KEY
        secret_key: $AWS_SECRET_KEY
        port: 9100

  - job_name: 'gateway-http'
    ec2_sd_configs:
      - region: '$AWS_REGION'
        access_key: $AWS_ACCESS_KEY
        secret_key: $AWS_SECRET_KEY
        port: 9000

  - job_name: 'sims'
    ec2_sd_configs:
      - region: '$AWS_REGION'
        access_key: $AWS_ACCESS_KEY
        secret_key: $AWS_SECRET_KEY
        port: 9001
" | sudo tee /etc/prometheus/prometheus.yml > /dev/null

# /etc/prometheus/api.rules
echo '
job:handler_http_status:sum = sum(rate(handler_request_count [5m])) by (status)
job:handler_http_route:sum = sum(rate(handler_request_count [5m])) by (route)
job:handler_http_latency:apdex = ((sum(rate(handler_request_latency_seconds_bucket{le="0.05"}[5m])) + sum(rate(handler_request_latency_seconds_bucket{le="0.25"}[5m]))) / 2) / sum(rate(handler_request_latency_seconds_count[5m]))
job:handler_http_latency:50 = histogram_quantile(0.5, sum(rate(handler_request_latency_seconds_bucket [5m])) by (le))
job:handler_http_latency:95 = histogram_quantile(0.95, sum(rate(handler_request_latency_seconds_bucket [5m])) by (le))
job:handler_http_latency:99 = histogram_quantile(0.99, sum(rate(handler_request_latency_seconds_bucket [5m])) by (le))
job:service_latency:apdex = ((sum(rate(service_op_latency_seconds_bucket{le="0.005"}[5m])) + sum(rate(service_op_latency_seconds_bucket{le="0.025"}[5m]))) / 2) / sum(rate(service_op_latency_seconds_count[5m]))
job:service_latency:50 = histogram_quantile(0.5, sum(rate(service_op_latency_seconds_bucket [5m])) by (le))
job:service_latency:95 = histogram_quantile(0.95, sum(rate(service_op_latency_seconds_bucket [5m])) by (le))
job:service_latency:99 = histogram_quantile(0.99, sum(rate(service_op_latency_seconds_bucket [5m])) by (le))
job:service_err:count = sum(rate(service_err_count [5m])) by (method, service)
job:service_op:count = sum(rate(service_op_count [5m])) by (method, service)
job:source_latency:apdex = ((sum(rate(source_op_latency_seconds_bucket{le="0.005"}[5m])) + sum(rate(source_op_latency_seconds_bucket{le="0.025"}[5m]))) / 2) / sum(rate(source_op_latency_seconds_count[5m]))
job:source_latency:50 = histogram_quantile(0.5, sum(rate(source_op_latency_seconds_bucket [5m])) by (le))
job:source_latency:95 = histogram_quantile(0.95, sum(rate(source_op_latency_seconds_bucket [5m])) by (le))
job:source_latency:99 = histogram_quantile(0.99, sum(rate(source_op_latency_seconds_bucket [5m])) by (le))
job:source_err:count = sum(rate(source_err_count [5m])) by (method, source)
job:source_op:count = sum(rate(source_op_count [5m])) by (method, source)
job:source_queue_latency:50 = histogram_quantile(0.5, sum(rate(source_queeu_latency_seconds_bucket [5m])) by (le))
job:source_queue_latency:95 = histogram_quantile(0.95, sum(rate(source_queue_latency_seconds_bucket [5m])) by (le))
job:source_queue_latency:99 = histogram_quantile(0.99, sum(rate(source_queue_latency_seconds_bucket [5m])) by (le))
job:intaker_request_latency:avg = avg(api_intaker_request_latency_microseconds_sum / api_intaker_request_latency_microseconds_count)
job:intaker_request_latency:max = max(api_intaker_request_latency_microseconds_sum / api_intaker_request_latency_microseconds_count)
job:intaker_request_latency:min = min(api_intaker_request_latency_microseconds_sum / api_intaker_request_latency_microseconds_count)
job:intaker_request_status:sum = sum(rate(api_intaker_request_count [5m])) by (status)
job:intaker_request_routes:all = sum(rate(api_intaker_request_count [5m])) by (route)
job:platform_process_res:sum = sum(process_resident_memory_bytes) by (instance, job)
job:platform_process_cpu:max = max(rate(process_cpu_seconds_total [5m])) by (instance, job)
' | sudo tee /etc/prometheus/api.rules > /dev/null
