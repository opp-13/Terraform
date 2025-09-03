provider "local" {}

resource "local_file" "k8s_yaml" {
  filename = "${path.module}/test.yaml"
  content  = <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
    - name: nginx
      image: nginx
      ports:
        - containerPort: 80
EOF
}

