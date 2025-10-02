# Resource

## Index
### [1. 인프라 객체를 관리하기 위한 추상화](#인프라-객체를-관리하기-위한-추상화)
### [2. Essential Resource](#Essential-Resource)
### [3. Additional Resource](#Additional-Resource)
### [4. Provider RPC Call Validation](Provider-RPC-Call-Validation)

## 인프라 객체를 관리하기 위한 추상화 

Terraform Core는 아래의 Resource를 gRPC로 호출하여 원하는 인프라 상태를 구현

예를 틀어 Terraform plan에서는 Resource 관련해서 아래의 RPC를 호출
- ValidateResourceConfig
- ReadResource
- PlanResourceChange

[참조](https://developer.hashicorp.com/terraform/plugin/framework/internals/rpcs)

참고로 Terraform init 시에는 Provider를 다운로드함

### Terraform ReadResource RPC Flow
[![Terraform ReadResource RPC Flow](https://web-unified-docs-hashicorp.vercel.app/api/assets/terraform-plugin-framework/latest/img/read-resource-detail.png)]()

해당 RPC는 결론적으로 Resource의 Configure 및 Read Method 호출

#### TODO 
main Branch에서 예제 코드 링크 추가할 예정


### 

| Terraform 명령   | 호출되는 Resource 관련 메서드                                                             |
|------------------|------------------------------------------------------------------------------------------|
| validate         | `Schema` (Schema method)<br>`ConfigValidators` method<br>`ValidateConfig` method         |
| plan             | `Schema` (Schema method)<br>`ConfigValidators` method<br>`ValidateConfig` method<br>`Configure` method<br>`Read` method<br>`ModifyPlan` method |
| apply            | `Schema` (Schema method)<br>`ConfigValidators` method<br>`ValidateConfig` method<br>`Configure` method<br>`Read` method<br>`ModifyPlan` method<br>`Create`, `Update`, `Delete` method |



## Essential Resource
무조건 필수로 있어야 하는 리소스

### Create
**Terraform Apply 시 호출**  
리소스를 생성하기 위한 메서드

`terraform apply` 시 `ApplyResourceChange` RPC를 통해 호출됨. 요청에는 설정과 plan 데이터가 포함되고, 응답에는 생성된 state 데이터가 포함됨. 오류 발생 시 리소스가 tainted 상태로 표시되어 다음 plan에서 재생성됨.

### Read
**Terraform Plan, Apply, Refresh 시 호출**  
리소스의 현재 상태를 확인하기 위한 메서드

`ReadResource` RPC를 통해 호출되며 리소스 import 후에도 실행됨. 이전 state를 받아서 새로고침된 state를 반환함. 리소스가 존재하지 않으면 `RemoveResource()`를 호출해서 state에서 제거해야 함.

### Update
**Terraform Apply 시 호출**  
리소스를 제자리에서 업데이트하기 위한 메서드

`ApplyResourceChange` RPC를 통해 호출됨. 이전 state, 설정, plan 데이터를 받아서 업데이트된 state를 반환함. 제자리 업데이트가 불가능하면 `resource.RequiresReplace()` plan modifier로 리소스 교체를 강제할 수 있음.

### Delete
**Terraform Apply 시 호출**  
리소스를 삭제하기 위한 메서드

리소스 삭제 시 `ApplyResourceChange` RPC를 통해 호출됨. 이전 state 데이터를 받고 진단 정보만 반환함. 이미 존재하지 않는 리소스의 오류는 무시하는 게 좋음. Terraform 1.3+에서는 삭제 계획 시 진단을 반환할 수 있음.

## Additional Resource

### Configure
**Provider 수준 데이터를 리소스에 전달**

`resource.ResourceWithConfigure` 인터페이스 구현으로 provider의 클라이언트 데이터(HTTP 클라이언트, SDK 등)를 받아와서 사용함. validation 시 provider 설정 없이 실행될 수 있으니 nil 체크가 필요함.

### Default
**null 속성에 기본값 설정**

설정이 null인 computed 속성에 기본값을 설정함. 스키마의 `Default` 필드에 정적 값이나 커스텀 기본값을 제공할 수 있음. 검증 후, 변경사항 적용 전에 실행됨.

### Import State
**기존 리소스를 Terraform 관리로 가져오기**

`terraform import` 명령으로 기존 리소스를 state로 가져옴. `resource.ResourceWithImportState` 인터페이스 구현으로 import 식별자를 파싱해서 `Read` 메서드가 state를 새로고침할 수 있도록 정보를 설정함.

### Manage Private State
**Provider 전용 데이터 저장 관리**

Terraform state에 저장되는 provider 전용 데이터. plan에서는 노출되지 않지만 ETag, 타임아웃 등을 저장하는 데 사용됨. `GetKey`와 `SetKey` 함수로 JSON/UTF-8 안전 데이터를 읽고 씀.

### Modify Plans
**계획 수정으로 최종 상태 조정**

검증 후 Terraform 계획을 조정할 수 있음. 속성별 plan modifier와 리소스별 plan modification 지원. unknown 값 조정, 리소스 교체 표시, 진단 반환이 가능함.

### Upgrade State
**State 데이터 자동 업그레이드**

스키마 버전 변경이나 데이터 구조 개선 시 기존 state를 새로운 형식으로 자동 업그레이드함. 사용자 개입 없이 백그라운드에서 실행되어 호환성을 유지함.

### Validate
**리소스 설정 유효성 검증**

선언적 또는 명령적 방식으로 전체 리소스 설정을 검증함. 여러 속성 간의 관계나 조합 검증에 유용하며, provider 설정 없이("오프라인") 실행됨.

### Timeouts
**CRUD 작업 타임아웃 설정**

클라우드 인프라 작업 지연을 고려해서 CRUD 작업에 타임아웃을 설정할 수 있음. `terraform-plugin-framework-timeouts` 모듈과 함께 사용해서 블록이나 속성 형태로 정의함.

### Write-Only Arguments
**임시 값 처리용 특수 속성**

Terraform 1.11+에서 지원. 비밀번호, API 키 등 state에 저장하면 안 되는 값을 처리함. 스키마에서 `WriteOnly: true`로 설정하면 설정에서만 사용되고 state/plan에는 null로 저장됨.

## Provider RPC Call Validation

Terraform Apply 시 **Create Method**가 첫 번째로 호출됨
```bash
root@master:~/ptest# terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - opp13/file in /root/ptest/linuxFile
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to
│ become incompatible with published releases.
╵
data.file_file.example: Reading...
data.file_file.example: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the
following symbols:
  + create

Terraform will perform the following actions:

  # file_file.example will be created
  + resource "file_file" "example" {
      + content  = "This file is managed by a custom Terraform resource!"
      + created  = (known after apply)
      + filename = "hello_resource.txt"
      + format   = "txt"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + file_created = (known after apply)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

file_file.example: Creating...
file_file.example: Creation complete after 0s

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

file_created = true
root@master:~/ptest# ls
hello_custom.txt  hello_resource.txt  hello_resource.txt.log  linuxFile  main.tf  terraform.tfstate
root@master:~/ptest# cat hello_resource.txt.log
Create
 called at 2025-10-02T17:25:07+09:00
```

main.tf 조정 후 Terraform plan 시 **Read Method**가 호출됨
```bash
root@master:~/ptest# nano main.tf
root@master:~/ptest# terraform plan
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - opp13/file in /root/ptest/linuxFile
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to
│ become incompatible with published releases.
╵
data.file_file.example: Reading...
file_file.example: Refreshing state...
data.file_file.example: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the
following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # file_file.example will be updated in-place
  ~ resource "file_file" "example" {
      ~ content  = "This file is managed by a custom Terraform resource!" -> "This file is managed by a custom Terraform resources!"
      ~ created  = true -> (known after apply)
        # (2 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ file_created = true -> (known after apply)

───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if
you run "terraform apply" now.
root@master:~/ptest# cat hello_resource.txt.log
Create
 called at 2025-10-02T17:25:07+09:00
Read
 called at 2025-10-02T17:26:38+09:00
```


그 후 Terraform apply 시 **Read/Update Method**가 호출됨
```bash
root@master:~/ptest# terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - opp13/file in /root/ptest/linuxFile
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to
│ become incompatible with published releases.
╵
file_file.example: Refreshing state...
data.file_file.example: Reading...
data.file_file.example: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the
following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # file_file.example will be updated in-place
  ~ resource "file_file" "example" {
      ~ content  = "This file is managed by a custom Terraform resource!" -> "This file is managed by a custom Terraform resources!"
      ~ created  = true -> (known after apply)
        # (2 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Changes to Outputs:
  ~ file_created = true -> (known after apply)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

file_file.example: Modifying...
file_file.example: Modifications complete after 0s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

Outputs:

file_created = true
root@master:~/ptest# cat hello_resource.txt.log
Create
 called at 2025-10-02T17:25:07+09:00
Read
 called at 2025-10-02T17:26:38+09:00
Read
 called at 2025-10-02T17:28:12+09:00
Update
 called at 2025-10-02T17:28:14+09:00
```


Terraform destroy 시 **Read/Delete Method**가 호출됨
```bash
root@master:~/ptest# terraform destroy
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - opp13/file in /root/ptest/linuxFile
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to
│ become incompatible with published releases.
╵
file_file.example: Refreshing state...
data.file_file.example: Reading...
data.file_file.example: Read complete after 0s

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the
following symbols:
  - destroy

Terraform will perform the following actions:

  # file_file.example will be destroyed
  - resource "file_file" "example" {
      - content  = "This file is managed by a custom Terraform resources!" -> null
      - created  = true -> null
      - filename = "hello_resource.txt" -> null
      - format   = "txt" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Changes to Outputs:
  - file_created = true -> null

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

file_file.example: Destroying...
file_file.example: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
root@master:~/ptest# cat hello_resource.txt.log
Create
 called at 2025-10-02T17:25:07+09:00
Read
 called at 2025-10-02T17:26:38+09:00
Read
 called at 2025-10-02T17:28:12+09:00
Update
 called at 2025-10-02T17:28:14+09:00
Read
 called at 2025-10-02T17:29:44+09:00
Delete
 called at 2025-10-02T17:29:45+09:00
```

[ApplyResourceChange RPC의 경우 해당 링크 참조](https://developer.hashicorp.com/terraform/plugin/framework/internals/rpcs#detail-5)

프로바이더 리소스의 상태와 계획(plan)을 비교하여 아래의 메서드를 실행
- Create
- Update
- Delete
