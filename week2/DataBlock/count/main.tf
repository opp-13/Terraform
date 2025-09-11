variable "names" {
  type    = list(string)
  default = ["1", "2", "3"]
}

data "local_file" "files" {
  count    = length(var.names)
  filename = "${path.module}/input/file${var.names[count.index]}.txt"

  lifecycle {
    postcondition {
      condition     = fileexists("${path.module}/input/${var.names[count.index]}")
      error_message = "파일 ${var.names[count.index]} 가 부재함"
    }
  }
}

resource "local_file" "new_files" {
  count    = length(data.local_file.files)
  filename = "${path.module}/output/${var.names[count.index]}"
  content  = "${var.names[count.index]}의 복사본"
}