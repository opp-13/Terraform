# Terraform Provider file

### Arguments

| 구분          | Argument   | 설명                           | 예시/비고                              |
|--------------|------------|------------------------------|-------------------------------------|
| **Resource** | filename   | 관리할 파일 이름 [리소스 식별자]          | "hello_resource.txt"                 |
|              | content    | 파일에 쓸 내용                     | "This file is managed by Terraform" |
|              | format     | 파일 형식                        | "txt"                              |
| **Data Block** | filename   | 생성할 파일 이름                    | "hello_custom.txt"                  |
|                | content    | 파일에 쓸 내용 |  "This file is managed by Terraform"             |
|                | format     | 파일 형식                        | "txt"                             |

### 빌드 방법

아래의 코드를 사용해 빌드합니다.

```shell
cd ~/Terraform/week3/linuxFileProvider
go build -o terraform-provider-file
mkdir -p /root/ptest/linuxFile
mv terraform-provider-file /root/ptest/linuxFile/
chmod +x /root/ptest/linuxFile/terraform-provider-file
```

### Provider 파일 설정

dev_overrides를 통해 Terraform Provider를 불러옵니다.
아래의 내용을 ~/.terraformrc에 넣어 주세요
```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/opp13/file" = "/root/ptest/linuxFile"
  }
  direct {}
}
```

### main.tf 설정 방법

main.tf를 설정합니다.
```hcl
terraform {
  required_providers {
    file = {
      source  = "registry.terraform.io/opp13/file"
      version = "0.1.0"
    }
  }
}

provider "file" {}

data "file_file" "example" {
  filename = "hello_custom.txt"
  content  = "Hello from linuxFile provider!"
  format   = "txt"
}

resource "file_file" "example" {
  filename = "hello_resource.txt"
  content  = "This file is managed by a custom Terraform resources!"
  format   = "txt"
}

output "file_created" {
  value = file_file.example.created
}
```

Terraform을 실행시킵니다. 
Terraform init 과정에서 error가 뜰 수 있지만 무시해 주세요.
- Terraform registry에 등록이 안되서 발생한 에러이므로 실제 작동에는 문제 없습니다.
```shell
terraform init
terraform apply
```
