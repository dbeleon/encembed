package encembed

const tpl = `// Code generated by 'go generate'; DO NOT EDIT.
package {{.PkgName}}

import(
	"filippo.io/age"
	"io"
	_ "embed"
	"bytes"
	{{if .Scramble}}"github.com/dbeleon/scr"
	"encoding/base64"
	"fmt"{{end}}
)

//go:embed {{.EmbedName}}
var {{.EncryptedVarName}} []byte
func {{.FuncName}}({{if .ExternalKey}}key string{{end}}) []byte {
	k := {{if .ExternalKey}}key{{else}}"{{.Key}}"{{end}}
	{{if .Scramble}}k = {{.EncryptedVarName}}DescrambleKey(k){{end}}
	i, _ := age.NewScryptIdentity(k)
	r, _ := age.Decrypt(bytes.NewReader({{.EncryptedVarName}}), i)
	a, _ := io.ReadAll(r)
	return a
}
{{if .DecryptedVarName}}var {{.DecryptedVarName}} = {{.FuncName}}() {{end}}

{{if .Scramble}}
func {{.EncryptedVarName}}DescrambleKey(scrKey string) (key string) {
	scram := scr.New({{.ScrLength}}, []int { {{range $v := .ScrFeedbacks}}{{$v}},{{end}} }, uint64({{.ScrPolynomial}}))
	data, err := base64.StdEncoding.DecodeString(scrKey)
	if err != nil {
		panic(fmt.Errorf("invalid key: %s", err))
	}
	scram.DescrambleAdditive(data)
	return string(data)
}{{end}}
`
