package controller

import (
	"fmt"
	"os/exec"
)

func Picture(videoPath string, thumbnailPath string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", "00:00:05", "-frames:v", "1", thumbnailPath)
	fmt.Println(cmd)
	err := cmd.Run()
	if err != nil {
		return err
	} else {
		fmt.Println("picture success")
	}
	return nil
}
