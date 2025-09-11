

## data_source_archive_file Data Block 분석

[internal/provider/data_source_archive_file.go 분석](https://github.com/hashicorp/terraform-provider-archive/blob/main/internal/provider/data_source_archive_file.go)

### [Read](https://github.com/hashicorp/terraform-provider-archive/blob/main/internal/provider/data_source_archive_file.go#L265)
```go
func (d *archiveFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model fileModel
	diags := req.Config.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	outputPath := model.OutputPath.ValueString()

	outputDirectory := path.Dir(outputPath)
	
	// outputDirectory 없으면 Directory를 생성
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

    // 파일 압축
	if err := archive(ctx, model); err != nil {
		resp.Diagnostics.AddError(
			"Archive creation error",
			fmt.Sprintf("error creating archive: %s", err),
		)
		return
	}

	// Generate archived file stats
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

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}
```

### [archive](https://github.com/hashicorp/terraform-provider-archive/blob/main/internal/provider/data_source_archive_file.go#L199)
```go
func archive(ctx context.Context, model fileModel) error {
	archiveType := model.Type.ValueString()
	outputPath := model.OutputPath.ValueString()

	archiver := getArchiver(archiveType, outputPath)
	if archiver == nil {
		return fmt.Errorf("archive type not supported: %s", archiveType)
	}

	outputFileMode := model.OutputFileMode.ValueString()
	if outputFileMode != "" {
		archiver.SetOutputFileMode(outputFileMode)
	}

	switch {
	case !model.SourceDir.IsNull():
		excludeList := make([]string, len(model.Excludes.Elements()))

		if !model.Excludes.IsNull() {
			var elements []types.String
			model.Excludes.ElementsAs(ctx, &elements, false)

			for i, elem := range elements {
				excludeList[i] = elem.ValueString()
			}
		}

		opts := ArchiveDirOpts{
			Excludes: excludeList,
		}

		if !model.ExcludeSymlinkDirectories.IsNull() {
			opts.ExcludeSymlinkDirectories = model.ExcludeSymlinkDirectories.ValueBool()
		}

		if err := archiver.ArchiveDir(model.SourceDir.ValueString(), opts); err != nil {
			return fmt.Errorf("error archiving directory: %s", err)
		}
	case !model.SourceFile.IsNull():
		if err := archiver.ArchiveFile(model.SourceFile.ValueString()); err != nil {
			return fmt.Errorf("error archiving file: %s", err)
		}
	case !model.SourceContentFilename.IsNull():
		content := model.SourceContent.ValueString()

		if err := archiver.ArchiveContent([]byte(content), model.SourceContentFilename.ValueString()); err != nil {
			return fmt.Errorf("error archiving content: %s", err)
		}
	case !model.Source.IsNull():
		content := make(map[string][]byte)

		var elements []sourceModel
		model.Source.ElementsAs(ctx, &elements, false)

		for _, elem := range elements {
			content[elem.Filename.ValueString()] = []byte(elem.Content.ValueString())
		}

		if err := archiver.ArchiveMultiple(content); err != nil {
			return fmt.Errorf("error archiving content: %s", err)
		}
	}

	return nil
}
```


### [실제 압축 호출 코드](https://github.com/hashicorp/terraform-provider-archive/blob/main/internal/provider/archiver.go)



## 그 외

[Go Pakage](https://pkg.go.dev/archive)를 사용 해 압축

[생성된 ZIP 파일의 정보](https://github.com/hashicorp/terraform-provider-archive/blob/main/internal/provider/data_source_archive_file.go#L331)
```go
type fileModel struct {
	ID                        types.String `tfsdk:"id"`
	Source                    types.Set    `tfsdk:"source"` // sourceModel
	Type                      types.String `tfsdk:"type"`
	SourceContent             types.String `tfsdk:"source_content"`
	SourceContentFilename     types.String `tfsdk:"source_content_filename"`
	SourceFile                types.String `tfsdk:"source_file"`
	SourceDir                 types.String `tfsdk:"source_dir"`
	Excludes                  types.Set    `tfsdk:"excludes"`
	ExcludeSymlinkDirectories types.Bool   `tfsdk:"exclude_symlink_directories"`
	OutputPath                types.String `tfsdk:"output_path"`
	OutputSize                types.Int64  `tfsdk:"output_size"`
	OutputFileMode            types.String `tfsdk:"output_file_mode"`
	OutputMd5                 types.String `tfsdk:"output_md5"`
	OutputSha                 types.String `tfsdk:"output_sha"`
	OutputSha256              types.String `tfsdk:"output_sha256"`
	OutputBase64Sha256        types.String `tfsdk:"output_base64sha256"`
	OutputSha512              types.String `tfsdk:"output_sha512"`
	OutputBase64Sha512        types.String `tfsdk:"output_base64sha512"`
}
```