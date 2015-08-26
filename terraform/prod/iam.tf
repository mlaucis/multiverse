resource "aws_iam_instance_profile" "backend" {
  name  = "${var.iam_profile_backend}"
  roles = [
    "${aws_iam_role.backend.name}"]
}

resource "aws_iam_role" "backend" {
  name               = "${var.iam_role_backend}"
  path               = "/"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "frontend" {
  name  = "${var.iam_profile_frontend}"
  roles = [
    "${aws_iam_role.frontend.name}"]
}

resource "aws_iam_role" "frontend" {
  name               = "${var.iam_role_frontend}"
  path               = "/"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "corporate" {
  name  = "${var.iam_profile_corporate}"
  roles = [
    "${aws_iam_role.corporate.name}"]
}

resource "aws_iam_role" "corporate" {
  name               = "${var.iam_role_corporate}"
  path               = "/"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}