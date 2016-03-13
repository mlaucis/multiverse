#!/usr/bin/env sh

set -ex

env

echo -n 'deb http://deb.robustperception.io/ precise nightly' | sudo tee /etc/apt/sources.list.d/robustperception.io.list > /dev/null

curl https://s3-eu-west-1.amazonaws.com/deb.robustperception.io/41EFC99D.gpg | sudo apt-key add -

# Install prometheus
sudo apt-get update
sudo apt-get install -y prometheus node-exporter pushgateway alertmanager

# Start and let prometheus run on boot
sudo service prometheus start
sudo update-rc.d prometheus defaults 95 10

# /etc/prometheus/prometheus.yml
echo '
global:
  evaluation_interval: "1m"
  scrape_interval: "30s"

rule_files:
  - /etc/prometheus/api.rules

scrape_configs:
  - job_name: "prometheus"
    target_groups:
    - targets:
        - "localhost:9090"
  - job_name: "pushgateway"
    honor_labels: true
    target_groups:
    - targets:
        - "localhost:9091"
  - job_name: "alertmanager"
    target_groups:
    - targets:
        - "localhost:9093"

  - job_name: "node-exporter"
    ec2_sd_configs:
      - region: "eu-central-1"
        access_key: AKIAJYTAYVGCJR6VNQLA
        secret_key: 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
        port: 9100

  - job_name: "services"
    ec2_sd_configs:
      - region: "eu-central-1"
        access_key: AKIAJYTAYVGCJR6VNQLA
        secret_key: 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
        port: 9000
' | sudo tee /etc/prometheus/prometheus.yml > /dev/null

# /etc/prometheus/api.rules
ehco '
job:gateway_http_status:sum = sum(rate(intaker_handler_request_count [5m])) by (status)
job:gateway_http_route:sum = sum(rate(intaker_handler_request_count [5m])) by (route)
job:gateway_service_app_err:method = sum(rate(intaker_service_app_err_count[5m])) by (method)
job:gateway_service_connection_err:method = sum(rate(intaker_service_connection_err_count[5m])) by (method)
job:gateway_service_event_err:method = sum(rate(intaker_service_event_err_count[5m])) by (method)
job:gateway_service_member_err:method = sum(rate(intaker_service_member_err_count[5m])) by (method)
job:gateway_service_object_err:method = sum(rate(intaker_service_object_err_count[5m])) by (method)
job:gateway_service_org_err:method = sum(rate(intaker_service_org_err_count[5m])) by (method)
job:gateway_service_user_err:method = sum(rate(intaker_service_user_err_count[5m])) by (method)
job:gateway_service_app_op:method = sum(rate(intaker_service_app_op_count[5m])) by (method)
job:gateway_service_connection_op:method = sum(rate(intaker_service_connection_op_count[5m])) by (method)
job:gateway_service_event_op:method = sum(rate(intaker_service_event_op_count[5m])) by (method)
job:gateway_service_member_op:method = sum(rate(intaker_service_member_op_count[5m])) by (method)
job:gateway_service_object_op:method = sum(rate(intaker_service_object_op_count[5m])) by (method)
job:gateway_service_org_op:method = sum(rate(intaker_service_org_op_count[5m])) by (method)
job:gateway_service_user_op:method = sum(rate(intaker_service_user_op_count[5m])) by (method)
job:intaker_request_statuss:sum = sum by (status) (rate(api_intaker_request_count [5m]))
job:intaker_request_latency:avg = avg(api_intaker_request_latency_microseconds_sum / api_intaker_request_latency_microseconds_count)
job:intaker_request_latency:max = max(api_intaker_request_latency_microseconds_sum / api_intaker_request_latency_microseconds_count)
job:intaker_request_latency:min = min(api_intaker_request_latency_microseconds_sum / api_intaker_request_latency_microseconds_count)
job:intaker_request_routes:all = sum(rate(api_intaker_request_count [5m])) by (route)
job:platform_process_res:sum = sum(process_resident_memory_bytes) by (instance, job)
job:platform_process_cpu:max = max(rate(process_cpu_seconds_total [5m])) by (instance, job)
' | sudo tee /etc/prometheus/api.rules > /dev/null
