resource "aws_iam_role_policy" "ecsSQSSendReceiver" {
    name    = "ecsSQSSendReceiver"
    role    = "${aws_iam_role.ecsInstance.id}"
    policy  = <<EOF
{
   "Version": "2012-10-17",
   "Statement":[{
      "Effect":"Allow",
      "Action": [
        "sqs:SendMessage",
        "sqs:ReceiveMessage",
        "sqs:GetQueueUrl"
      ],
      "Resource":"arn:aws:sqs:*:775034650473:*-state-change"
      }
   ]
}
EOF
}

resource "aws_sqs_queue" "connection-state-change-dlq" {
    delay_seconds               = 0
    max_message_size            = 262144
    message_retention_seconds   = 1209600
    name                        = "connection-state-change-dlq"
    receive_wait_time_seconds   = 1
    visibility_timeout_seconds  = 300
}

resource "aws_sqs_queue" "connection-state-change" {
    delay_seconds               = 0
    max_message_size            = 262144
    message_retention_seconds   = 1209600
    name                        = "connection-state-change"
    receive_wait_time_seconds   = 1
    redrive_policy              = <<EOF
{
    "deadLetterTargetArn": "${aws_sqs_queue.connection-state-change-dlq.arn}",
    "maxReceiveCount": 10
}
EOF
    visibility_timeout_seconds  = 60
}