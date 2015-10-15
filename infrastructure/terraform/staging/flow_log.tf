/** /
resource "aws_flow_log" "staging-flow-log" {
  log_group_name = "flow-log"
  iam_role_arn   = "${aws_iam_role.staging-flow-log.arn}"
  vpc_id         = "${aws_vpc.tapglue.id}"
  traffic_type   = "ALL"
}

resource "aws_iam_role" "staging-flow-log" {
  name               = "staging-flow-log"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "vpc-flow-logs.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "staging-flow-log" {
  name   = "staging-flow-log"
  role   = "${aws_iam_role.staging-flow-log.id}"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "logs:DescribeLogGroups",
        "logs:DescribeLogStreams"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}
/**/
