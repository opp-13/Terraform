

## Terraform archive_file


### 어떻게 작동하는지?

Terraform [plan, apply, refresh] 시 ReadDataSource RPC를 통하여 해당 프로바이더의 Read 함수 호출


[![Terraform Plugin Framework DataSource Read 플로우](https://web-unified-docs-hashicorp.vercel.app/api/assets/terraform-plugin-framework/latest/img/read-data-source-detail.png)](https://web-unified-docs-hashicorp.vercel.app/api/assets/terraform-plugin-framework/latest/img/read-data-source-detail.png)


Read 메서드가 일반적으로 수행하는 일:
1. datasource.ReadRequest.Config 필드에서 구성 데이터 로드 
2. 원격 시스템 정보와 같은 추가 데이터를 로드 
3. datasource.ReadResponse.State 필드에 상태 데이터를 기록


[소스 코드](https://github.com/hashicorp/terraform-provider-archive/blob/main/internal/provider/data_source_archive_file.go)
```go
func (d *archiveFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var model fileModel
	// datasource.ReadRequest.Config 필드에서 구성 데이터 로드
	diags := req.Config.Get(ctx, &model)
    // datasource.ReadResponse.State 필드에 상태 데이터를 기록
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
	
    outputPath := model.OutputPath.ValueString()

    // 결과물 출력 디렉터리 지정
    outputDirectory := path.Dir(outputPath)
	
	// 디렉터리가 없을 경우 새로 생성
    if outputDirectory != "" {
        if _, err := os.Stat(outputDirectory); err != nil {
            if err := os.MkdirAll(outputDirectory, 0755); err != nil {
                resp.Diagnostics.AddError(
                    "Output path error",
                    fmt.Sprintf("error creating output path: %s", err),
                )
                return
            }
        }
    }
    
	// 압축 진행
    if err := archive(ctx, model); err != nil {
        resp.Diagnostics.AddError(
            "Archive creation error",
            fmt.Sprintf("error creating archive: %s", err),
        )
        return
    }
    
    // 압축된 파일의 상태 생성
    fi, err := os.Stat(outputPath)
    if err != nil {
        resp.Diagnostics.AddError(
            "Archive output error",
            fmt.Sprintf("error reading output: %s", err),
        )
        return
    }
    model.OutputSize = types.Int64Value(fi.Size())
    
    checksums, err := genFileChecksums(outputPath)
    if err != nil {
        resp.Diagnostics.AddError(
            "Hash generation error",
            fmt.Sprintf("error generating checksums: %s", err),
        )
    }
    model.OutputMd5 = types.StringValue(checksums.md5Hex)
    model.OutputSha = types.StringValue(checksums.sha1Hex)
    model.OutputSha256 = types.StringValue(checksums.sha256Hex)
    model.OutputBase64Sha256 = types.StringValue(checksums.sha256Base64)
    model.OutputSha512 = types.StringValue(checksums.sha512Hex)
    model.OutputBase64Sha512 = types.StringValue(checksums.sha512Base64)
    
    model.ID = types.StringValue(checksums.sha1Hex)
    
	// datasource.ReadResponse.State 필드에 상태 데이터를 기록
    diags = resp.State.Set(ctx, model)
    resp.Diagnostics.Append(diags...)
}
```




참조  
https://developer.hashicorp.com/terraform/plugin/framework/internals/rpcs#read-rpcs
https://developer.hashicorp.com/terraform/plugin/framework/data-sources#read-method



### 지원 하는 파일 종류가 zip 또는 tar.gz인 이유

#### Golang의 archive package를 사용
- tar
- zip

참조  
https://pkg.go.dev/archive@go1.25.1


#### [zip archive 예시](https://github.com/hashicorp/terraform-provider-archive/blob/main/internal/provider/zip_archiver.go#L7)
```go
import (
    "archive/zip"
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "time"
    
    "github.com/bmatcuk/doublestar/v4"
)

func (a *ZipArchiver) ArchiveFile(infilename string) error {
    fi, err := assertValidFile(infilename)
    if err != nil {
        return err
    }

    content, err := os.ReadFile(infilename)
    if err != nil {
        return err
    }
    
    if err := a.open(); err != nil {
        return err
    }
    defer a.close()
    
    fh, err := zip.FileInfoHeader(fi)
    if err != nil {
        return fmt.Errorf("error creating file header: %s", err)
    }
    fh.Name = filepath.ToSlash(fi.Name())
    fh.Method = zip.Deflate
    //nolint:staticcheck // This is required as fh.SetModTime has been deprecated since Go 1.10 and using fh.Modified alone isn't enough when using a zero value
    fh.SetModTime(time.Time{})
    
    if a.outputFileMode != "" {
        filemode, err := strconv.ParseUint(a.outputFileMode, 0, 32)
    if err != nil {
        return fmt.Errorf("error parsing output_file_mode value: %s", a.outputFileMode)
    }
    fh.SetMode(os.FileMode(filemode))
    }
    
    f, err := a.writer.CreateHeader(fh)
    if err != nil {
        return fmt.Errorf("error creating file inside archive: %s", err)
    }
    
    _, err = f.Write(content)
    return err
}
```







https://pkg.go.dev/archive@go1.25.1


https://github.com/hashicorp/terraform-provider-archive/blob/main/internal/provider/archiver.go



### Provider에 DataSource 추가
https://github.com/hashicorp/terraform-provider-archive/blob/main/internal/provider/provider.go#L34

참조
https://developer.hashicorp.com/terraform/plugin/framework/data-sources#add-data-source-to-provider






# Terraform CLI

[Source Repo](https://github.com/hashicorp/terraform-exec/tree/main)

https://developer.hashicorp.com/terraform/plugin/framework/internals/rpcs


https://developer.hashicorp.com/terraform/plugin/framework/getting-started/code-walkthrough

https://www.perplexity.ai/search/https-developer-hashicorp-com-ZGe8KbBGRoiAdfkCMBjqEA?10=d

