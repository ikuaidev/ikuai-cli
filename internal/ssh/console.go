// Package ssh provides SSH console interaction for iKuai router menu operations.
// It mirrors the Python ssh_console.py module.
package ssh

import (
	"bytes"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

const (
	sshTimeout    = 15 * time.Second
	recvTimeout   = 3 * time.Second
	resetDuration = 500 * time.Millisecond
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*[mHJK]`)

func stripANSI(s string) string {
	return strings.ReplaceAll(ansiRe.ReplaceAllString(s, ""), "\r", "")
}

// recvOutput reads from r until no data arrives for resetDuration, or total timeout expires.
func recvOutput(r net.Conn, totalTimeout time.Duration) string {
	var buf bytes.Buffer
	tmp := make([]byte, 4096)

	deadline := time.Now().Add(totalTimeout)
	lastData := time.Now()

	_ = r.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	for {
		n, err := r.Read(tmp)
		if n > 0 {
			buf.Write(tmp[:n])
			lastData = time.Now()
		}
		if err != nil {
			// Either timeout or EOF
			if time.Since(lastData) >= resetDuration {
				break
			}
		}
		if time.Now().After(deadline) {
			break
		}
		_ = r.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	}
	return stripANSI(buf.String())
}

// ResetWebPasswd connects to the iKuai router via SSH, navigates to option 7
// (恢复WEB管理密码), confirms, and returns the new credentials.
// Returns (username, password, error).
func ResetWebPasswd(host, user, password string, port int) (string, string, error) {
	config := &gossh.ClientConfig{
		User: user,
		Auth: []gossh.AuthMethod{
			gossh.Password(password),
		},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(), //nolint:gosec
		Timeout:         sshTimeout,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := gossh.Dial("tcp", addr, config)
	if err != nil {
		return "", "", fmt.Errorf("cannot connect to %s — %w", addr, err)
	}
	defer func() {
		_ = client.Close()
	}()

	sess, err := client.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("cannot create SSH session: %w", err)
	}
	defer func() {
		_ = sess.Close()
	}()

	modes := gossh.TerminalModes{
		gossh.ECHO:          0,
		gossh.TTY_OP_ISPEED: 14400,
		gossh.TTY_OP_OSPEED: 14400,
	}
	if err := sess.RequestPty("xterm", 50, 200, modes); err != nil {
		return "", "", fmt.Errorf("cannot request PTY: %w", err)
	}

	stdin, err := sess.StdinPipe()
	if err != nil {
		return "", "", err
	}

	// We need a net.Conn-like reader for deadline support.
	// Use a pipe approach with goroutine.
	pipeR, pipeW := net.Pipe()

	stdoutPipe, err := sess.StdoutPipe()
	if err != nil {
		return "", "", err
	}

	// Copy stdout to pipe in background
	go func() {
		defer func() {
			_ = pipeW.Close()
		}()
		tmp := make([]byte, 4096)
		for {
			n, err := stdoutPipe.Read(tmp)
			if n > 0 {
				_, _ = pipeW.Write(tmp[:n])
			}
			if err != nil {
				return
			}
		}
	}()

	if err := sess.Shell(); err != nil {
		return "", "", fmt.Errorf("cannot start shell: %w", err)
	}

	// Consume initial banner (~2 seconds)
	recvOutput(pipeR, 2*time.Second)

	// Select option 7: 恢复WEB管理密码
	fmt.Fprint(stdin, "7\n") //nolint:errcheck
	out := recvOutput(pipeR, recvTimeout)

	lower := strings.ToLower(out)
	if !strings.Contains(out, "yes或no") && !strings.Contains(lower, "yes") {
		if len(out) > 300 {
			out = out[:300]
		}
		return "", "", fmt.Errorf("unexpected console output after option 7:\n%s", out)
	}

	// Confirm reset
	fmt.Fprint(stdin, "yes\n") //nolint:errcheck
	result := recvOutput(pipeR, recvTimeout)

	if !strings.Contains(result, "恢复WEB管理密码成功") {
		if len(result) > 300 {
			result = result[:300]
		}
		return "", "", fmt.Errorf("password reset failed. Router output:\n%s", result)
	}

	// Parse credentials from output
	username := "admin"
	passwd := "admin"

	reUser := regexp.MustCompile(`默认用户[:：]\s*(\S+)`)
	rePass := regexp.MustCompile(`默认密码[:：]\s*(\S+)`)

	if m := reUser.FindStringSubmatch(result); len(m) > 1 {
		username = m[1]
	}
	if m := rePass.FindStringSubmatch(result); len(m) > 1 {
		passwd = m[1]
	}

	return username, passwd, nil
}
