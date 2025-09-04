## IaC 환경에서 사용되는 HCL, Yaml, Json 정리

 <br/>

| 항목 | HCL | YAML | JSON |
| --- | --- | --- | --- |
| 문자열 | name = "value" | name: value | "name": "value" |
| Num | number = 10 | number: 10 | "number": 10 |
| Boolean | enabled = true | enabled: true | "enabled": true |
| Arr/List | ids = [1,2,3] | ids:\n - 1\n - 2\n - 3 | "ids": [1, 2, 3] |
| Dict/Map | { key = "value" } | key: value | {"key": "value"} |
| Date | 지원 안함 | date: 2020-01-01 | 날짜는 문자열로 취급 |
| 주석 | /* */ 또는 # 또는 // | # | 지원 안함 |

 <br/>

**HCL Native Syntax Specification**  
https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md


## Terraform에서의 Json 사용

### [HCL](https://github.com/opp-13/Terraform/tree/main/week1/hclProvision)
main.tf
```
provider "local" {}

resource "local_file" "k8s_yaml" {
filename = "${path.module}/test.yaml"
content = <<EOF
apiVersion: v1
kind: Pod
metadata:
 name: nginx
spec:
 containers:

  name: nginx
  image: nginx
  ports:

   containerPort: 80
EOF
}
```

### [Json](https://github.com/opp-13/Terraform/tree/main/week1/jsonProvision)

main.tf.json
```
{
  "provider": {
    "local": {}
  },
  "resource": {
    "local_file": {
      "k8s_yaml": {
        "filename": "${path.module}/test.yaml",
        "content": "apiVersion: v1\nkind: Pod\nmetadata:\n  name: nginx\nspec:\n  containers:\n    - name: nginx\n      image: nginx\n      ports:\n        - containerPort: 80\n"
      }
    }
  }
}
```

### 생각해 볼 점

1. 두 형식 모두 terraform.tfstate에 json 형식으로 동일하게 저장됨. (terraform state 명령어로도 확인 가능)   
다만 가독성이나 주석, mutiple line 인식 등에 있어 차이가 있음
2. 실제로 저장되는 상태는 동일함
terraform.tfstate의 일부 (HCL의 경우 정의된 상태의 형식이 Json으로 변환됨)
```
"content": "apiVersion: v1\nkind: Pod\nmetadata:\n  name: nginx\nspec:\n  containers:\n    - name: nginx\n      image: nginx\n      ports:\n        - containerPort: 80\n"
```
3. Kubernetes 2.0의 경우 Desired State를 HCL로 선정했다는 정보가 있는 만큼 K8S 2.0 환경에서는 HCL이 유리


## Terraform local_file package에서의 상태 관리
ID를 기준으로 상태를 확인합니다.

terraform apply 시 리소스에 대한 Metadata가 생성 (리눅스 시스템에 의존)

terraform plan 또는 apply 시 생성될 리소스에 대한 Metadata 정보 제공

```jsx
Terraform will perform the following actions:

  # local_file.k8s_yaml will be created
  + resource "local_file" "k8s_yaml" {
      + content              = <<-EOT
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
        EOT
      + content_base64sha256 = (known after apply)
      + content_base64sha512 = (known after apply)
      + content_md5          = (known after apply)
      + content_sha1         = (known after apply)
      + content_sha256       = (known after apply)
      + content_sha512       = (known after apply)
      + directory_permission = "0777"
      + file_permission      = "0777"
      + filename             = "./test.yaml"
      + id                   = (known after apply)
    }
```

그 후 terraform state show local_file.k8s_yaml 명령어로 metadata 확인 가능 [terraform.tfstate에 저장됨]

```jsx
## Terraform이 관리하고 있는 리소스 목록
root@test:~/Terraform/week1/jsonProvision# terraform state list
local_file.k8s_yaml

## Terraform이 관리하고 있는 리소스 상세 정보
root@test:~/Terraform/week1/jsonProvision# terraform state show local_file.k8s_yaml
# local_file.k8s_yaml:
resource "local_file" "k8s_yaml" {
    content              = <<-EOT
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
    EOT
    content_base64sha256 = "Cn/Fi5TbIFvY5MqAaN3vvB9K2pmDC6RLBEL7jiKujDg="
    content_base64sha512 = "rUZKrOB2R8Q5Q/aFiaTPuv15HXKU+O8cTkGqCHxQ27Z9gb+A5dTmsadalOZeRImcglGEjsBkHQVw2rJgf3Bj1A=="
    content_md5          = "2e05f08d1d1139d1498474089822ac4c"
    content_sha1         = "2bf6c2b135bba0ca9a50093b6ae72d38cd009f0a"
    content_sha256       = "0a7fc58b94db205bd8e4ca8068ddefbc1f4ada99830ba44b0442fb8e22ae8c38"
    content_sha512       = "ad464aace07647c43943f68589a4cfbafd791d7294f8ef1c4e41aa087c50dbb67d81bf80e5d4e6b1a75a94e65e44899c8251848ec0641d0570dab2607f7063d4"
    directory_permission = "0777"
    file_permission      = "0777"
    filename             = "./test.yaml"
    id                   = "2bf6c2b135bba0ca9a50093b6ae72d38cd009f0a"
}
```

실제로 Terraform으로 관리되고 있는 리소스 Hash 정보

Hash 값이 위와 같은 것을 볼 수 있다.

```jsx
root@test:~/Terraform/week1/jsonProvision# md5sum test.yaml
2e05f08d1d1139d1498474089822ac4c  test.yaml

root@test:~/Terraform/week1/jsonProvision# sha1sum test.yaml
2bf6c2b135bba0ca9a50093b6ae72d38cd009f0a  test.yaml

root@test:~/Terraform/week1/jsonProvision# sha256sum test.yaml
0a7fc58b94db205bd8e4ca8068ddefbc1f4ada99830ba44b0442fb8e22ae8c38  test.yaml

root@test:~/Terraform/week1/jsonProvision# sha512sum test.yaml
ad464aace07647c43943f68589a4cfbafd791d7294f8ef1c4e41aa087c50dbb67d81bf80e5d4e6b1a75a94e65e44899c8251848ec0641d0570dab2607f7063d4  test.yaml
```

권한

```jsx
root@test:~/Terraform/week1/jsonProvision# ls -al test.yaml
-rwxr-xr-x 1 root root 147 Sep  3 05:49 test.yaml
```

참고로 file_permission이 다른 이유는 umask 때문

```jsx
root@test:~/Terraform/week1/hclProvision# umask
0022

777 -> requested file permission
022 -> umask
---
755 -> actual permission
```

권한이 달라져도 Terraform에서는 인식 못함

```jsx
root@test:~/Terraform/week1/hclProvision# chmod 711 ./test.yaml
root@test:~/Terraform/week1/hclProvision# ls -al test.yaml
-rwx--x--x 1 root root 147 Sep  3 05:42 test.yaml
root@test:~/Terraform/week1/hclProvision# terraform plan
local_file.k8s_yaml: Refreshing state... [id=2bf6c2b135bba0ca9a50093b6ae72d38cd009f0a]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are
needed.
```

이유는 id를 기준으로 변경사항을 식별하는데 id는 SHA1 checksum을 기준으로 만들어지기 때문 [값은 SHA1 checksum과 같다]

```jsx
root@test:~/Terraform/week1/hclProvision# terraform plan
local_file.k8s_yaml: Refreshing state... [id=2bf6c2b135bba0ca9a50093b6ae72d38cd009f0a]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are
needed.
```

참고

https://registry.terraform.io/providers/hashicorp/local/latest/docs/resources/file#id-1

따라서 checksum을 바꾸면 “Plan: 1 to add, 0 to change, 0 to destroy.” 변경을 감지한다.

```jsx
root@test:~/Terraform/week1/hclProvision# md5sum test.yaml
2e05f08d1d1139d1498474089822ac4c  test.yaml
root@test:~/Terraform/week1/hclProvision# echo "#change checksum" >> test.yaml
root@test:~/Terraform/week1/hclProvision# md5sum test.yaml
8d795951b72f6a72d6b009a5c1ee4fc2  test.yaml
```

Plan 상 새롭게 만들어진다고 나와있지만 terraform apply 시 이름이 같을 시 있는 파일을 덮어쓴다. [birth와 modify가 다름]

```jsx
root@test:~/Terraform/week1/hclProvision# stat test.yaml
  File: test.yaml
  Size: 147             Blocks: 8          IO Block: 4096   regular file
Device: 252,0   Inode: 1194078     Links: 1
Access: (0711/-rwx--x--x)  Uid: (    0/    root)   Gid: (    0/    root)
Access: 2025-09-03 12:29:28.152606449 +0000
Modify: 2025-09-03 12:29:21.427747884 +0000
Change: 2025-09-03 12:29:21.427747884 +0000
 Birth: 2025-09-03 05:42:06.354212056 +0000
```

문제는 말 그대로 내용만 덮어쓰기 때문에 권한은 설정된 권한이 아니다. 그 점을 유의할 것.

## muti resource provision

main.tf

```jsx
resource "local_file" "bash_script" {
  filename = "${path.module}/find_commands.sh"
  content  = <<EOF
#!/bin/bash
compgen -c | grep -E '^.{3}$' > result.txt
EOF
  file_permission = "0755"
}

resource "null_resource" "run_bash" {
  depends_on = [local_file.bash_script]

  provisioner "local-exec" {
    command = "${path.module}/find_commands.sh"
  }
}
```

현재 상태

```jsx
root@test:~/Terraform/week1/multipleProvider# terraform state list
local_file.bash_script
null_resource.run_bash
root@test:~/Terraform/week1/multipleProvider# terraform state show local_file.bash_script
# local_file.bash_script:
resource "local_file" "bash_script" {
    content              = <<-EOT
        #!/bin/bash
        compgen -c | grep -E '^.{3}$' > result.txt
    EOT
    content_base64sha256 = "AVpkHOG8f+zZi2RI4CDKSFUT07Tx5oEVHyW0dDRlgz0="
    content_base64sha512 = "tqrsNd4v/yIbhijmIL2igD6tk5RmL2mgSI4KrB/iiqgl35yDAi2CwTfhpWoEmV0Cx/fVJLC8Ji6BVOnycaEqBQ=="
    content_md5          = "58fd099dd7f65f3b4c86b3b48833eb72"
    content_sha1         = "91d2f2c0cd4cbfc8b56031c393be29967822f085"
    content_sha256       = "015a641ce1bc7fecd98b6448e020ca485513d3b4f1e681151f25b4743465833d"
    content_sha512       = "b6aaec35de2fff221b8628e620bda2803ead9394662f69a0488e0aac1fe28aa825df9c83022d82c137e1a56a04995d02c7f7d524b0bc262e8154e9f271a12a05"
    directory_permission = "0777"
    file_permission      = "0755"
    filename             = "./find_commands.sh"
    id                   = "91d2f2c0cd4cbfc8b56031c393be29967822f085"
}
root@test:~/Terraform/week1/multipleProvider# terraform state show null_resource.run_bash
# null_resource.run_bash:
resource "null_resource" "run_bash" {
    id = "4312012921333570008"
}
```

여기서 추가로 생각해 봐야 되는 점

depends_on을 사용하더라도 의존성이 있는 리소스가 실행되지 않을 수 있다.

위의 main.tf는 3글자 짜리 명령어를 가져와 result.txt에 저장하는

만약 local_file.bash_script (Terraform state)를 삭제하거나 만들어진 bash 파일인 find_commands.sh를 수정/삭제할 시 null_resource.run_bash에 해당하는 부분은 실행되지 않는다.

따라서 null_resource를 사용하는 리소스에 dependency가 있으면 무조건 적어도 해당 리소스의 state를 제거 후 다시 apply 할 것
