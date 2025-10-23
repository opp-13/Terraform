## 매우 간단한 hashicorp/terraform-provider-aws 분석
#### EC2 인스턴스를 예시로 코드를 분석합니다.

### 서비스 별 경로
terraform-provider-aws/internal/service/{서비스 이름}


### 호출 방식
#### Read가 호출되는 원리

[resource 스키마 정의](https://github.com/hashicorp/terraform-provider-aws/blob/main/internal/service/ec2/ec2_instance.go#L64C1-L70C48)
``` go
func resourceInstance() *schema.Resource {
	//lintignore:R011
	return &schema.Resource{
		CreateWithoutTimeout: resourceInstanceCreate,
		ReadWithoutTimeout:   resourceInstanceRead,
		UpdateWithoutTimeout: resourceInstanceUpdate,
		DeleteWithoutTimeout: resourceInstanceDelete,
...
		Schema: map[string]*schema.Schema{
			"ami": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Computed:     true,
				Optional:     true,
				AtLeastOneOf: []string{"ami", names.AttrLaunchTemplate},
			},
			names.AttrARN: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"associate_public_ip_address": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Computed: true,
				Optional: true,
			},
...
```

Resource의 State가 refresh 될때 Read가 호출됨   
이때 Read, ReadContext, ReadWithoutTimeout 메서드 중 하나만 적용 가능

EC2 Resource에 대한 스키마 또한 여기서 정의


[참고](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#Resource)
[Resource 및 ResourceData](https://pkg.go.dev/github.com/hashicorp/terraform/helper/schema#Resource)

### resourceInstanceRead 코드 분석


[코드](https://github.com/hashicorp/terraform-provider-aws/blob/main/internal/service/ec2/ec2_instance.go#L1234)
``` go
func resourceInstanceRead(ctx context.Context, rd *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	
	// Provider에서 AWS EC2 Client 연결 객체 로딩
	c := meta.(*conns.AWSClient)
	conn := c.EC2Client(ctx)

    // EC2 인스턴스 정보를 조회, rd.Id는 terraform resource id(Terraform state에서 리소스 ID 추출됨 -> 그게 instance id로 사용됨)
	instance, err := findInstanceByID(ctx, conn, rd.Id())

    // 새 리소스가 아니거나(생성 중) 발견 못할 경우 
	if !rd.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 Instance %s not found, removing from state", rd.Id())
		// State에서 리소스를 제거
		rd.SetId("")
		return diags
	}

    // 에러 메세지
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 Instance (%s): %s", rd.Id(), err)
	}
 
	diags = append(diags, resourceInstanceFlatten(ctx, c, instance, rd)...)

	return diags
}
```


#### AWS Client 정보 불러오기

[AWSClient Structure](https://github.com/hashicorp/terraform-provider-aws/blob/main/internal/conns/awsclient.go#L29)



``` go
func findInstanceByID(ctx context.Context, conn *ec2.Client, id string) (*awstypes.Instance, error) {
	
	input := ec2.DescribeInstancesInput{
		InstanceIds: []string{id},
	}

	output, err := findInstance(ctx, conn, &input)

	if err != nil {
		return nil, err
	}

    // 현재 인스턴스 상태가 Terminated일 경우 Terraform State에서 제거
	if state := output.State.Name; state == awstypes.InstanceStateNameTerminated {
		return nil, &retry.NotFoundError{
			Message:     string(state),
			LastRequest: &input,
		}
	}

	// Eventual consistency check. (최종 일관성 검증, 요청 id(Terraform state에서 저장하고 있는 id)와 조회한 인스턴스의 id가 드를 경우 에러처리)
	if aws.ToString(output.InstanceId) != id {
		return nil, &retry.NotFoundError{
			LastRequest: &input,
		}
	}

	return output, nil
}

func findInstance(ctx context.Context, conn *ec2.Client, input *ec2.DescribeInstancesInput) (*awstypes.Instance, error) {
	// listInstances 호출
	output, err := tfslices.CollectWithError(listInstances(ctx, conn, input))

	if err != nil {
		return nil, err
	}

	return tfresource.AssertSingleValueResult(output, func(v *awstypes.Instance) bool { return v.State != nil })
}

// DescribeInstances is an "All-Or-Some" call.
func listInstances(ctx context.Context, conn *ec2.Client, input *ec2.DescribeInstancesInput) iter.Seq2[awstypes.Instance, error] {
	return func(yield func(awstypes.Instance, error) bool) {
	    // 인스턴스 목록 가져오기
		pages := ec2.NewDescribeInstancesPaginator(conn, input)
		
		// 가져온 인스턴스 처리 (페이지 단위)
		for pages.HasMorePages() {
			page, err := pages.NextPage(ctx)

			if tfawserr.ErrCodeEquals(err, errCodeInvalidInstanceIDNotFound) {
				yield(awstypes.Instance{}, &retry.NotFoundError{
					LastError:   err,
					LastRequest: &input,
				})
				return
			}

			if err != nil {
				yield(awstypes.Instance{}, err)
				return
			}

			for _, v := range page.Reservations {
				for _, instance := range v.Instances {
					if !yield(instance, nil) {
						return
					}
				}
			}
		}
	}
}
```

[왜 DescribeInstances을 안쓰고 NewDescribeInstancesPaginator를 쓰는지?](https://docs.aws.amazon.com/ko_kr/ec2/latest/devguide/ec2-api-pagination.html)


### 참고
[EC2 Package - Golang](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2)

[DescribeInstancesInput 타입 정의](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2#DescribeInstancesInput)

[NewDescribeInstancesPaginator](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2#NewDescribeInstancesPaginator)

[DescribeInstances](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2#Client.DescribeInstances)

