package utils

import (
	"fmt"
	"strings"
)

const noneString = "<none>"

// ParseRepoDigest 解析repo和digest
// 格式类似[docker.io/library/alpine@sha256:xxx]
// 这货色是一个数组 解析的时候取第一个即可
// 返回两个值：第一个是imageName 第二个是digest
func ParseRepoDigest(repoDigests []string) (string, string) {
	if len(repoDigests) == 0 {
		return noneString, noneString
	}
	repoDigestPair := strings.Split(repoDigests[0], "@")
	if len(repoDigestPair) != 2 {
		return "errImage", "errDigest"
	}
	return repoDigestPair[0], repoDigestPair[1]
}

// ParseRepoTag 解析镜像和tag
// 格式类似[docker.io/library/alpine:3.12]
// 如果镜像出错（打包出错 中途终止等）返回值是一个二维 string切片([][]string{})
// 只需要显示第一个 每一个子切片是一个string{} 包含两个值：镜像名称和tag
func ParseRepoTag(repoTags []string, imageName string) (repoTagPairs [][]string) {
	if len(repoTags) == 0 {
		repoTagPairs = append(repoTagPairs, []string{imageName, noneString})
		return
	}

	for _, repoTag := range repoTags {
		idx := strings.LastIndex(repoTag, ":")
		if idx < 0 { // 解析出错了，直接返回errTag，
			repoTagPairs = append(repoTagPairs, []string{"errName", "errTag"})
			continue
		}

		name := repoTag[:idx]
		if name == noneString {
			name = imageName
		}
		repoTagPairs = append(repoTagPairs, []string{name, repoTag[idx+1:]})
	}

	return
}

// ParseSize 解析size
func ParseSize(size uint64) string {
	return fmt.Sprintf("%.2fm", float64(size)/1024/1024)
}

// ParseImageID 截取ID
func ParseImageID(id string) string {
	idstr := strings.Split(id, ":")[1]
	return idstr[:13]
}
