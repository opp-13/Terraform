resource "local_file" "bash_script" {
  filename = "${path.module}/find_commands.sh"
  content  = <<EOF
#!/bin/bash
compgen -c | grep -E '^.{4}$' > result.txt
EOF
  file_permission = "0755"
}

resource "null_resource" "run_bash" {
  depends_on = [local_file.bash_script]

  provisioner "local-exec" {
    command = "${path.module}/find_commands.sh"
  }
}
