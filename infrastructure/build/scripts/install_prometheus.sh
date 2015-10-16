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

echo '
global:
  scrape_interval: "15s"

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
  - job_name: "node"
    target_groups:
    - targets:
        - "localhost:9100"

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
