resource "aws_iam_server_certificate" "self-signed" {
  name = "self-signed"
  certificate_body = "${file("${path.module}/../../certs/self/self.crt")}"
  private_key = "${file("${path.module}/../../certs/self/self.key")}"
}
