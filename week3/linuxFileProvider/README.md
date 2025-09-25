# Terraform Provider file





아래의 코드를 사용해 빌드합니다.

```shell
cd ~/Terraform/week3/linuxFileProvider
go build -o terraform-provider-file
mkdir -p /root/ptest/linuxFile
mv terraform-provider-file /root/ptest/linuxFile/
chmod +x /root/ptest/linuxFile/terraform-provider-file
```


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

output "file_created" {
  value = data.file_file.example.created
}
```

Terraform을 실행시킵니다. 
Terraform init 과정에서 error가 뜰 수 있지만 무시해 주세요.
- Terraform registry에 등록이 안되서 발생한 에러이므로 실제 작동에는 문제 없습니다.
```shell
terraform init
terraform apply
```