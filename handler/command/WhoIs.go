package command

import (
	"errors"
	"log"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

// RunWhoIs start the WhoIs command
func RunWhoIs(ipAddr string) map[string][]string {

	// Parse IP to make sure it is valid
	ipObj := net.ParseIP(ipAddr)
	if ipObj == nil {
		log.Println("Invalid IP Address!")
		return make(map[string][]string)
	}

	// Use parsed IP for security reasons
	ipAddr = ipObj.String()

	// Finally, run the actual whois command with the right whois servers
	result := runWhoIsCommand(ipAddr)

	var outPut map[string][]string
	outPut = make(map[string][]string)

	singleLines := strings.Split(string(result), "\n")

	re := regexp.MustCompile("^[#%>]+")
	for _, line := range singleLines {
		if re.MatchString(line) {
			continue
		}
		lineParts := strings.Split(line, ": ")
		if len(lineParts) == 2 {
			tk := strings.TrimSpace(lineParts[0])
			outPut[tk] = append(outPut[tk], strings.TrimSpace(lineParts[1]))
		}
	}

	return outPut
}

func runWhoIsCommand(args ...string) []byte {

	// Execute command
	out, err := exec.Command("whois", args...).Output()
	if err != nil {
		if err.Error() != "exit status 2" {
			log.Println(err)
		}
	}

	_, err = isValidResponse(out)
	if err != nil {
		log.Println(err)
	}

	return out
}

func isValidResponse(response []byte) (valid bool, err error) {

	singleLines := strings.Split(string(response), "\n")
	if len(singleLines) < 5 {
		err = errors.New("invalid response detected. We assume that a valid whois response has at minimum 5 lines")
		return
	}
	valid = true
	return
}
