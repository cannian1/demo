package domain

import (
	"bufio"
	"errors"
	"strings"
)

type ContainerMonitor struct {
}

func NewContainerMonitor() *ContainerMonitor {
	return &ContainerMonitor{}
}

type ContainerDisk struct {
	Type      string `json:"type"`
	Size      string `json:"size"`
	Used      string `json:"used"`
	Available string `json:"available"`
	UsedRate  string `json:"usedRate"`
	Mount     string `json:"mount"`
}

func (cm ContainerMonitor) ParseOutput(out string) ([]*ContainerDisk, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	dataList := make([]*ContainerDisk, 0, 10)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.Contains(line, "Filesystem") {
			continue
		}
		data, err := cm.parseLine(line)
		if err != nil {
			return []*ContainerDisk{}, err
		}
		dataList = append(dataList, data)
	}
	return dataList, nil
}

func (cm ContainerMonitor) parseLine(line string) (*ContainerDisk, error) {
	split := strings.Fields(line)
	if len(split) < 6 {
		return nil, errors.New("line does not contain enough fields")
	}
	switch split[5] {
	case "/":
		split[0] = "系统盘"
	case "/root/xxx-tmp":
		split[0] = "数据盘"
	}
	return &ContainerDisk{
		Type:      split[0],
		Size:      split[1],
		Used:      split[2],
		Available: split[3],
		UsedRate:  split[4],
		Mount:     split[5],
	}, nil
}
