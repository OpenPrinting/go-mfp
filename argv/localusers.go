// MFP  - Miulti-Function Printers and scanners toolkit
// argv - Argv parsing mini-library
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// List of local users

package argv

import (
	"bufio"
	"os"
	"strings"
)

// localUsers returns list of local users' login names
// starting with the specified prefix.
//
// This function is intended for tilde completion
func localUsers(prefix string) []string {
	// Open /etc/passwd
	file, err := os.Open("/etc/passwd")
	if err != nil {
		return nil
	}
	defer file.Close()

	// Read line by line
	scanner := bufio.NewScanner(file)
	users := []string{}
	for scanner.Scan() {
		line := scanner.Text() // Символ \n уже отрезан, строка готова
		if i := strings.IndexByte(line, ':'); i > 0 {
			user := line[:i]
			if strings.HasPrefix(user, prefix) {
				users = append(users, user)
			}
		}
	}

	return users
}
