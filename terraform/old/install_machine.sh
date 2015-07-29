#!/bin/bash
sudo apt-get -y install python-pip
sudo pip install awscli

mkdir -p ~/.aws/
echo '[default]
region=eu-central-1
output=text' > ~/.aws/config

echo '[default]
aws_access_key_id=AKIAI4RHVHBRKYF5YWJQ
aws_secret_access_key=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' > ~/.aws/credentials

wget https://s3.amazonaws.com/aws-cloudwatch/downloads/latest/awslogs-agent-setup.py

chmod +x ./awslogs-agent-setup.py
mkdir -p ~/releases/cloudwatch/

aws s3 cp s3://tapglue-builds/cloudwatch/aws-cloudwatch.conf ~/releases/cloudwatch/aws-cloudwatch.conf

sudo ./awslogs-agent-setup.py -n -r eu-central-1 -c ~/releases/cloudwatch/aws-cloudwatch.conf
rm -f ./awslogs-agent-setup.py
