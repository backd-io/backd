package lua

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	lua "github.com/yuin/gopher-lua"

	"golang.org/x/crypto/ssh"
)

// lua module ssh
type sshModule struct {
	timeout time.Duration
}

// packageSSHModule publish the module to be usable by backd.Lua importer
func packageSSHModule(L *lua.LState) int {

	var s sshModule
	s.timeout = 10 * time.Second

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"ssh":     s.sshExec,
		"timeout": s.setTimeout,
	})

	L.SetField(mod, "name", lua.LString("ssh"))

	// returns the module
	L.Push(mod)
	return 1

}

func (s *sshModule) setTimeout(L *lua.LState) int {

	seconds := L.ToInt(1)

	s.timeout = time.Duration(seconds) * time.Second

	return 0
}

func (s *sshModule) sshExec(L *lua.LState) int {

	var (
		hostname string
		port     int
		username string
		keyFile  string
		command  string
		client   *ssh.Client
		session  *ssh.Session
		out      []byte
		err      error
	)

	hostname = L.ToString(1)
	port = L.ToInt(2)
	username = L.ToString(3)
	keyFile = L.ToString(4)
	command = L.ToString(5)

	client, err = s.CreateSSHConnection(hostname, port, username, keyFile)

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer client.Close()

	session, err = client.NewSession()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	out, err = session.CombinedOutput(command)
	if err != nil {
		L.Push(lua.LString(""))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	session.Close()
	L.Push(lua.LString(string(out)))
	L.Push(lua.LString(""))
	return 2

}

// CreateSSHConnection is the function that originates the ssh connection and
//   dials the remote machine. Returns an established connection to the host.
func (s *sshModule) CreateSSHConnection(hostname string, port int, username, keyFile string) (*ssh.Client, error) {

	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout: s.timeout,
	}

	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port), sshConfig)

}
