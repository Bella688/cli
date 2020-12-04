package create

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cli/cli/internal/ghrepo"
	"github.com/cli/cli/pkg/cmd/secret/shared"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdCreate(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		wants    CreateOptions
		stdinTTY bool
		wantsErr bool
	}{
		{
			name:     "no name",
			cli:      "",
			wantsErr: true,
		},
		{
			name:     "multiple names",
			cli:      "cool_secret good_secret",
			wantsErr: true,
		},
		{
			name:     "no body, stdin is terminal",
			cli:      "cool_secret",
			stdinTTY: true,
			wantsErr: true,
		},
		{
			name:     "visibility without org",
			cli:      "cool_secret -vall",
			wantsErr: true,
		},
		{
			name: "explicit org with selected repo",
			cli:  "--org=coolOrg -vselected -rcoolRepo cool_secret",
			wants: CreateOptions{
				SecretName:      "cool_secret",
				Visibility:      shared.VisSelected,
				RepositoryNames: []string{"coolRepo"},
				Body:            "-",
				OrgName:         "coolOrg",
			},
		},
		{
			name: "explicit org with selected repos",
			cli:  `--org=coolOrg -vselected -r="coolRepo,radRepo,goodRepo" cool_secret`,
			wants: CreateOptions{
				SecretName:      "cool_secret",
				Visibility:      shared.VisSelected,
				RepositoryNames: []string{"coolRepo", "goodRepo", "radRepo"},
				Body:            "-",
				OrgName:         "coolOrg",
			},
		},
		{
			name: "repo",
			cli:  `cool_secret -b"a secret"`,
			wants: CreateOptions{
				SecretName: "cool_secret",
				Visibility: "private",
				Body:       "a secret",
				OrgName:    "",
			},
		},
		{
			name: "implicit org",
			cli:  `cool_secret --org -b"@cool.json"`,
			wants: CreateOptions{
				SecretName: "cool_secret",
				Visibility: "private",
				Body:       "@cool.json",
				OrgName:    "@owner",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			f := &cmdutil.Factory{
				IOStreams: io,
			}

			io.SetStdinTTY(tt.stdinTTY)

			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			var gotOpts *CreateOptions
			cmd := NewCmdCreate(f, func(opts *CreateOptions) error {
				gotOpts = opts
				return nil
			})
			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, tt.wants.SecretName, gotOpts.SecretName)
			assert.Equal(t, tt.wants.Body, gotOpts.Body)
			assert.Equal(t, tt.wants.Visibility, gotOpts.Visibility)
			assert.Equal(t, tt.wants.OrgName, gotOpts.OrgName)
			assert.ElementsMatch(t, tt.wants.RepositoryNames, gotOpts.RepositoryNames)
		})
	}
}

func Test_createRun_repo(t *testing.T) {
	tests := []struct {
		name       string
		opts       *CreateOptions
		stdin      string
		wantOut    string
		wantStderr string
		wantErr    bool
		baseRepo   func() (ghrepo.Interface, error)
	}{
		{
			name: "implicit repo",
			opts: &CreateOptions{
				SecretName: "cool_secret",
				Body:       "a secret",
			},
		},
		{
			name: "explicit repo",
			opts: &CreateOptions{
				SecretName: "cool_secret",
				Body:       "a secret",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO
			assert.Equal(t, 1, 0)
		})
	}
}

func Test_getBody(t *testing.T) {
	tests := []struct {
		name     string
		bodyArg  string
		want     string
		stdin    string
		fromFile bool
	}{
		{
			name:    "literal value",
			bodyArg: "a secret",
			want:    "a secret",
		},
		{
			name:    "from stdin",
			bodyArg: "-",
			want:    "a secret",
			stdin:   "a secret",
		},
		{
			name:     "from file",
			fromFile: true,
			want:     "a secret from a file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, stdin, _, _ := iostreams.Test()

			io.SetStdinTTY(false)

			_, err := stdin.WriteString(tt.stdin)
			assert.NoError(t, err)

			if tt.fromFile {
				dir := os.TempDir()
				tmpfile, err := ioutil.TempFile(dir, "testfile*")
				assert.NoError(t, err)
				_, err = tmpfile.WriteString(tt.want)
				assert.NoError(t, err)
				tt.bodyArg = fmt.Sprintf("@%s", tmpfile.Name())
			}

			body, err := getBody(&CreateOptions{
				Body: tt.bodyArg,
				IO:   io,
			})
			assert.NoError(t, err)

			assert.Equal(t, string(body), tt.want)

		})

	}

}

/*
func Test_createRun_org(t *testing.T) {
	tests := []struct {
		name       string
		opts       *CreateOptions
		stdin      string
		wantOut    string
		wantStderr string
		wantErr    bool
	}{
		{
			name: "explicit literal body",
			opts: &CreateOptions{
				SecretName: "cool_secret",
				Body:       "a secret",
			},
		},
		{
			name: "explicit body filename",
		},
		{
			name: "stdin body",
		},
		{
			name: "explicit org name",
		},
		{
			name: "implicit org name",
		},
		{
			name: "scalar visibility",
		},
		{
			name: "selected visibility",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, 1, 0)
		})
	}
}
*/
