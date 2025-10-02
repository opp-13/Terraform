# Terraform Provider Development

## Index
### [1. Correlation between Terraform Core and Terraform Plugin](#Correlation-between-Terraform-Core-and-Terraform-Plugin)
### [2. Provider Server](#Provider-Server)
### [3. Terraform Block](#Terraform-Block)

## Correlation between Terraform Core and Terraform Plugin
Terraform은 크게 두 부분으로 나뉘어져 있음
- **Terraform Core**: 명령어 실행, 설정 파일 해석, 리소스 그래프 구성, 상태 및 plan 관리, 플러그인과 RPC 통신 담당
- **Terraform Plugins**: Go로 작성된 별도 바이너리, Terraform Core와 RPC로 통신하며 클라우드 API 등과 직접 상호작용

**Terraform Core**는 인프라 리소스를 직접 관리하지 않고 플러그인에 위임함  
플러그인은 인증, API 호출, 리소스 생성/수정/삭제 구현 담당

gRPC 기반 프로토콜로 Core와 플러그인이 다리 역할 하며, protobuf를 통해 메시지 스키마를 정의함  
이 덕분에 플러그인은 다른 언어로도 작성 가능하며 Terraform은 유연하게 확장 가능함

플러그인 발견 및 버전 제어도 Core가 담당(Terraform init에서 담당)


요약  
설정파일을 읽어서 적합한 플러긘을 불러오고 그에 맞게 실행


[![Correlation between Terraform Core and Terraform Plugin](https://web-unified-docs-hashicorp.vercel.app/api/assets/terraform-docs-common/latest/img/docs/terraform-plugin-overview.png)]()


참고 \
https://developer.hashicorp.com/terraform/plugin/how-terraform-works \
https://developer.hashicorp.com/terraform/plugin


## Provider Server

Terraform Provider는 무조건 [Terraform Plugin Protocol](https://developer.hashicorp.com/terraform/plugin/how-terraform-works#terraform-plugin-protocol)에 맞는 gRPC 서버를 가지고 있어야 함
\
[providerserver](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/providerserver)라는 Golang Package를 통하여 구현

#### SDK
- Terraform Plugin SDKv2 
- Terraform plugin framework

! 참고로 아래 내용은 Terraform plugin framework을 기준으로 작성됨

소스 코드 예제
```go
package main

import (
    "context"
    "flag"
    "log"

    "github.com/example-namespace/terraform-provider-example/internal/provider"
    "github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
    // Example version string that can be overwritten by a release process
    version string = "dev"
)

func main() {
    opts := providerserver.ServeOpts{
        // TODO: Update this string with the published name of your provider.
        Address: "registry.terraform.io/example-namespace/example",
    }

    err := providerserver.Serve(context.Background(), provider.New(version), opts)

    if err != nil {
        log.Fatal(err.Error())
    }
}
```

[참조](https://developer.hashicorp.com/terraform/plugin/framework/provider-servers)


## Terraform Block

### [Resource](./Resource.md)

### Data Block

