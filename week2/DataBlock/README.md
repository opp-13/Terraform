# Data Block

## 목차
### [1. Data Block 개요](#data-block-개요)
### [2. 기본 구조와 구성](#기본-구조와-구성)
### [3. 메타 인수](#메타-인수)

---

## Data Block 개요

테라폼이 관리하지 않는 리소스의 데이터를 가져옴
- 리소스 생성 X

## 기본 구조와 구성

### 기본 문법

```hcl
data "<TYPE>" "<LABEL>" {
  # 프로바이더 특화 인수들
  <PROVIDER_SPECIFIC_ARGUMENTS>

  # 메타-인수들
  count      = <NUMBER>
  depends_on = [<RESOURCE_REFERENCES>]
  for_each   = <MAP_OR_SET>
  provider   = <PROVIDER_ALIAS>
  
  lifecycle {
    precondition {
      condition     = <EXPRESSION>
      error_message = "<MESSAGE>"
    }
    
    postcondition {
      condition     = <EXPRESSION>
      error_message = "<MESSAGE>"
    }
  }
}
```



#### TYPE (데이터 소스 타입)
- 프로바이더가 정의한 데이터 소스 유형
- 예: `archive_file`, `aws_ami`, `aws_eks_cluster_versions`
- `terraform_remote_state` 포함 가능

#### LABEL (라벨)
- 데이터 소스의 고유 식별자
- 동일한 TYPE 내에서 유일해야 함
- 데이터 참조 시 사용: `data.<TYPE>.<LABEL>.<ATTRIBUTE>`

데이터 참조 예시
```hcl
data.<TYPE>.<LABEL>.<ATTRIBUTE>

ex)
data.archive_file.main_lambda_zip.output_base64sha256
data.aws_ami.ubuntu.id
```

---

## 메타 인수

### count

### depends_on

### Multiple Provider


