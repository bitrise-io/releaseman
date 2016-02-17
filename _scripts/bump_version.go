package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// Inputs
	var (
		versionFilePathParam = flag.String("file", "", `Version file path`)
		segmentParam         = flag.String("segment", "", `Version segment idx to bump (options: 0 - patch, 1 - minor, 2 - major)`)
	)

	flag.Parse()

	if versionFilePathParam == nil || *versionFilePathParam == "" {
		log.Fatalf(" [!] No version file parameter specified")
	}
	versionFilePath := *versionFilePathParam

	segmentIdx := int64(2)
	if segmentParam != nil && *segmentParam != "" {
		var err error
		if segmentIdx, err = strconv.ParseInt(*segmentParam, 10, 64); err != nil {
			log.Fatalf("Failed to parse segment idx (%s), err: %#v", *segmentParam, err)
		}
	}

	// Main
	versionFileBytes, err := ioutil.ReadFile(versionFilePath)
	if err != nil {
		log.Fatalf("Failed to read version file: %s", err)
	}
	versionFileContent := string(versionFileBytes)

	re := regexp.MustCompile(`const VERSION = "(?P<version>[0-9]+\.[0-9-]+\.[0-9-]+)"`)
	results := re.FindAllStringSubmatch(versionFileContent, -1)
	versionStr := ""
	for _, v := range results {
		versionStr = v[1]
	}
	if versionStr == "" {
		log.Fatalf("Failed to determine version")
	}

	versionStrSegments := strings.Split(versionStr, ".")
	segmentStrToBump := versionStrSegments[segmentIdx]
	segmentToBump, err := strconv.ParseInt(segmentStrToBump, 10, 64)
	versionStrSegments[segmentIdx] = fmt.Sprintf("%d", segmentToBump+1)

	bumpedVersion := strings.Join(versionStrSegments, ".")

	outBytes, err := exec.Command("bash", "_scripts/set_version.sh", "version/version.go", bumpedVersion).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to set next version, out: %s, error: %#v", string(outBytes), err)
	}

	fmt.Println(bumpedVersion)
}
