package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)
import "os/exec"

func main() {

	var directory string
	var file string
	var outputDirectory string

	flag.StringVar(&directory, "d", ".", "Specify username. Default is root")
	flag.StringVar(&file, "f", "", "Specify pass. Default is password")
	flag.StringVar(&outputDirectory, "o", ".", "Specify pass. Default is password")

	flag.Parse() // after declaring flags we need to call it

	// check if cli params match
	if file == "" {
		fmt.Printf("Bad filename")
		return
	}

	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", fmt.Sprintf("%s/%s", directory, file))
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	resolution := strings.Split(string(stdout), "x")
	fmt.Println(string(stdout), resolution)
	width, err := strconv.Atoi(resolution[0])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	height, err := strconv.Atoi(strings.Trim(resolution[1], "\r\n"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scale := ""
	crop := ""
	if width < height {
		scale = "scale=640:ih*640/iw"
		crop = fmt.Sprintf("crop=in_w:in_h-%d", height*640/width-640)
	} else {
		scale = "scale=iw*640/ih:640"
		crop = fmt.Sprintf("crop=in_w-%d:in_h", width*640/height-640)
	}
	// Print the output
	fmt.Println(string(stdout))

	cmd = exec.Command("ffmpeg", "-i", fmt.Sprintf("%s/%s", directory, file), "-vf", scale, fmt.Sprintf("%s/%s", outputDirectory, file), "-y")
	fmt.Println(cmd.String())

	stdout, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cmd = exec.Command("ffmpeg", "-i", fmt.Sprintf("%s/%s", outputDirectory, file), "-vf", crop, fmt.Sprintf("%s/result_%s", outputDirectory, file), "-y")
	fmt.Println(cmd.String())

	stdout, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	os.Remove(fmt.Sprintf("%s/%s", outputDirectory, file))
}
