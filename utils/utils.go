/* Copyright 2017, Ashish Thakwani. 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.LICENSE file.
 */
package utils

import (
    "fmt"
    "os"
    "strconv"
    "syscall"
    
    "golang.org/x/sys/unix"
)

/*
 *  Default configurations
 */
const (
    RUNPATH     = "/tmp/"
    SSHD_CONFIG = "/etc/ssh/sshd_config" 
    MATCHBLK    = `
Match User %s
    AllowTCPForwarding yes
    X11Forwarding no
    AllowAgentForwarding no
    PermitTTY yes
    ForceCommand /usr/sbin/trafficrouter -f $SSH_ORIGINAL_COMMAND
`
)


/*
 *  Config Struct is passed from client to server
 */
type Config struct {
    Port uint32 `json:"port"`
}

/*
 *  Host Struct is passed from forceCmd to server
 */
type Host struct {
    ListenPort  uint32 `json:"lisport"`
    RemoteIP    string `json:"raddr"`
    RemotePort  uint32 `json:"rport"`
    AppPort     uint32 `json:"aport"`
    Config      Config `json:"config"`
    Uid         int    `json:"uid"`
    Pid         int    `json:"pid"`
}


/*
 *  Common error handling function.
 */
func Check(e error) {
    if e != nil {
        fmt.Fprintln(os.Stderr, e)
        panic(e)
    }
}

/*
 * Block Program Forever
 */
func BlockForever() {
    select {}
}


/*
 * Lock file
 */
func LockFile(pid int) int {

    f, err := os.Create(RUNPATH + strconv.Itoa(pid))
    Check(err)
    
    fd := f.Fd()
	err = unix.Flock(int(fd), syscall.LOCK_EX)
    Check(err)
    
    return int(fd)
}


func UnlockFile(fd int) {
	err := unix.Flock(fd, syscall.LOCK_UN)
    Check(err)
}