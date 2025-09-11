variable "names" {
  type    = list(string)
  default = ["1", "2", "3"]
}

data "local_file" "files" {
  count    = length(var.names)
  filename = "${path.module}/input/file${var.names[count.index]}.txt"

  lifecycle {
    postcondition {
      condition     = fileexists("${path.module}/input/file${var.names[count.index]}.txt")
      error_message = "file${var.names[count.index]}.txt 이(가) 부재함"
    }
  }
}

resource "local_file" "new_files" {
  count    = length(data.local_file.files)
  filename = "${path.module}/output/file${var.names[count.index]}.txt"
  content  = "${var.names[count.index]}의 복사본"
}