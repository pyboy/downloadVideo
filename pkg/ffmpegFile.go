package pkg

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

//go:embed ffmpeg.exe
var ffmpegBinary []byte

func WriteFFmpegBinaryToTempFile() (string, error) {
	// 创建临时文件
	ffmpegFile, err := os.Create("ffmpeg.exe")
	if err != nil {
		return "", err
	}
	defer ffmpegFile.Close()

	// 将嵌入的ffmpeg二进制数据写入临时文件
	if _, err := ffmpegFile.Write(ffmpegBinary); err != nil {
		return "", err
	}

	// 设置临时文件为可执行
	if err := ffmpegFile.Chmod(0777); err != nil {
		return "", err
	}

	// 返回临时文件的路径
	return ffmpegFile.Name(), nil
}

func CreateEmptyMP4(ffmpeg string, filename string) {
	// 创建一个空的MP4文件
	cmd := exec.Command(ffmpeg, "-f", "lavfi", "-i", "anullsrc", "-c:v", "libx264", "-t", "0.1", "-y", filename)
	if err := cmd.Run(); err != nil {
		fmt.Println("创建空MP4文件失败:", err)
	}
}

func AppendTSFileToMP4(ffmpeg string, tsFile, outputMP4 string) error {
	// 使用ffmpeg将.ts文件追加到MP4文件
	cmd := exec.Command(ffmpeg, "-i", outputMP4, "-i", tsFile, "-c", "copy", "-bsf:v", "h264_mp4toannexb", "-f", "mpegts", "-y", "temp.ts")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg合并失败: %v", err)
	}

	// 删除原始的MP4文件，并将临时文件重命名为最终的MP4文件
	os.Remove(outputMP4)
	os.Rename("temp.ts", outputMP4)

	return nil

}

func MergeMP4Files(outputMP4 string, mp4Name string) error {
	fmt.Println("合并文件")
	// 使用ffmpeg将.ts文件追加到MP4文件
	// 当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("无法获取当前工作目录: %v", err)
	}
	fmt.Println("当前工作目录:", cwd)

	cmd := exec.Command(".\\ffmpeg.exe", "-f", "concat", "-safe", "0", "-i", outputMP4, "-c", "copy", "..\\"+mp4Name+".mp4")
	if err := cmd.Run(); err != nil {
		fmt.Println("ffmpeg合并失败: %v", err)
		return fmt.Errorf("ffmpeg合并失败: %v", err)
	}
	fmt.Println("合并完成", cmd)
	os.Chdir("..\\")
	cwd, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("无法获取当前工作目录: %v", err)
	}
	fmt.Println("当前工作目录:", cwd)
	_, err = os.Stat(mp4Name + ".mp4")
	if err != nil {
		return fmt.Errorf("合并后的MP4文件不存在: %v", err)
	}

	return nil
}
