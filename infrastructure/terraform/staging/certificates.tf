resource "aws_iam_server_certificate" "self-signed-staging" {
  name = "self-signed-staging"
  certificate_body = "${file("../../certs/self/self.crt")}"
  private_key = "${file("../../certs/self/self.key")}"
}
