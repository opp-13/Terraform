




[Init](https://github.com/hashicorp/terraform-exec/blob/main/tfexec/init.go#L119)
```go
// Init represents the terraform init subcommand.
// 가변인자로 InitOption 처리
func (tf *Terraform) Init(ctx context.Context, opts ...InitOption) error {
    cmd, err := tf.initCmd(ctx, opts...)
    if err != nil {
        return err
    }
    return tf.runTerraformCmd(ctx, cmd)
}
```

#### [참고[InitOption]](https://github.com/hashicorp/terraform-exec/blob/main/tfexec/init.go#L42)
```go
type initConfig struct {
	backend       bool
	backendConfig []string
    dir           string
    forceCopy     bool
    fromModule    string
    get           bool
    getPlugins    bool
    lock          bool
    lockTimeout   string
    pluginDir     []string
    reattachInfo  ReattachInfo
    reconfigure   bool
    upgrade       bool
    verifyPlugins bool
}

// InitOption represents options used in the Init method.
type InitOption interface {
    configureInit(*initConfig)
}
```

initCmd
```go
func (tf *Terraform) initCmd(ctx context.Context, opts ...InitOption) (*exec.Cmd, error) {

	// 기본 init 설정
	c := defaultInitOptions

	err := tf.configureInitOptions(ctx, &c, opts...)
	if err != nil {
		return nil, err
	}

	args, err := tf.buildInitArgs(ctx, c)
	if err != nil {
		return nil, err
	}

	// Optional positional argument; must be last as flags precede positional arguments.
	if c.dir != "" {
		args = append(args, c.dir)
	}

	return tf.buildInitCmd(ctx, c, args)
}

var defaultInitOptions = initConfig{
    backend:       true,
    forceCopy:     false,
    get:           true,
    getPlugins:    true,
    lock:          true,
    lockTimeout:   "0s",
    reconfigure:   false,
    upgrade:       false,
    verifyPlugins: true,
}
```

