package editor

import (
	"os"
	"os/exec"
	"strings"
)

const defaultEditor = "nano"

func GetEditor() (string, []string) {
	editor := strings.Fields(os.Getenv("EDITOR"))
	if len(editor) > 0 {
		return editor[0], editor[1:]
	}
	return defaultEditor, nil
}

func ReadFromEditor(extension string) (string, error) {
	editorCmd, editorArgs := GetEditor()

	tmpFile, err := os.CreateTemp(os.TempDir(), "*."+extension)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(editorCmd, append(editorArgs, tmpFile.Name())...) //nolint:gosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return "", err
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}
	return string(data), nil
}
