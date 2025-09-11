resource "local_file" "base" {
  filename = "${path.module}/files/base.txt"
  content  = "base file."
}

data "local_file" "base_info" {
  depends_on = [local_file.base]

  filename = local_file.base.filename
}

resource "local_file" "base_copy" {
  depends_on = [data.local_file.base_info]

  filename = "${path.module}/files/copy.txt"
  content  = "복사한 내용: ${data.local_file.base_info.content}"
}