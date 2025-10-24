## key_name이 변경될 시 Instance가 삭제 후 재생성되었던 이유


key_name 설정
``` go
"key_name": {
	Type:     schema.TypeString,
	Optional: true,
	ForceNew: true,
	Computed: true,
},
```


Schema 행동 정의
``` go
// Optional 또는 Required이 설정되있다면 구성(configuration)에서 값을 가져올 수 있음
// 즉 값이 Computed(아래 항목 참조)가 아닌 경우, 아래 중 하나는 반드시 설정이 필요
// 동시에 두 값을 설정하는 것은 불가능 함
Optional bool
Required bool

// True일시 provider 또는 외부 API에 의해 계산됨을 의미
// 즉, 값은 구성(config)에서 오거나, 계산(computed)되거나, 혹은 두 가지 모두일 수 있음
Computed  bool

// 리소스의 값이 변경될 경우 기존 리소스를 수정하지 않고 새 리소스를 생성
ForceNew  bool

// 전체를 저장하는 대신 해시값만 저장 (리소스가 너무 클 경우 사용)
StateFunc SchemaStateFunc
```

key_name은 ForceNew가 true이기 때문에 key_name이 변경될 시 terraform은 Resource에 대해서 update가 아닌 delete 후 create을 한다.

<br>

#### key_name 상태변경에 따른 Terraform 동작

| key_name 상태변경에 따른 Terraform 동작 | 행동 |
|--------------------------------|------------------|
| 새로 추가됨                         | Resource에 대해서 delete 후 create |
| 변경됨                            | Resource에 대해서 delete 후 create |
| 삭제됨                            | Resource에 대한 변경사항 없음 |


그런데 왜 기존 key_name을 삭제하면 Resource가 재생성되지 않고 변경사항 없음으로 뜰까?
Optional + Compute를 동시에 사용할 경우 Teraform은 
[암묵적으로 ](https://github.com/hashicorp/terraform-plugin-sdk/issues/1101) 이전 정보를 보존하는 방향으로 동작한다. 
따라서 그것을 바꾸고 싶다면 [추가 조치](https://discuss.hashicorp.com/t/schema-for-optional-computed-to-support-correct-removal-plan-in-framework/49055)가 필요할 것이다.


### 참조
[Schema Code](https://pkg.go.dev/github.com/hashicorp/terraform/helper/schema#Schema)
[Schema Type](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-types)
[Schema 행동 정의](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors)
