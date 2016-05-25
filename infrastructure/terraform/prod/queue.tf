resource "aws_iam_user" "state-change-sr" {
  name  = "state-change-sr"
  path  = "/"
}

resource "aws_iam_user_policy" "state-change-sr" {
  name  = "state-change-sr"
  user  = "${aws_iam_user.state-change-sr.name}"
  policy  = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "sqs:DeleteMessage",
                "sqs:GetQueueUrl",
                "sqs:ReceiveMessage",
                "sqs:SendMessage"
            ],
            "Resource": "arn:aws:sqs:*:775034650473:*-state-change"
        },
        {
            "Effect": "Allow",
            "Action": [
                "sns:CreatePlatformEndpoint",
                "sns:GetEndpointAttributes",
                "sns:Publish",
                "sns:SetEndpointAttributes"
            ],
            "Resource": "arn:aws:sns:*:775034650473:*"
        }
   ]
}
EOF
}

resource "aws_iam_access_key" "state-change-sr" {
  user  = "${aws_iam_user.state-change-sr.name}"
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