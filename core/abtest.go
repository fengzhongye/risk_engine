package core

import (
	//	"github.com/skyhackvip/risk_engine/configs"
	"github.com/skyhackvip/risk_engine/internal/errcode"
	"log"
	"math/rand"
	"time"
)

type AbtestNode struct {
	Info    NodeInfo `yaml:"info"`
	Branchs []Branch `yaml:"branchs,flow"`
}

func (abtest AbtestNode) GetName() string {
	return abtest.Info.Name
}

func (abtest AbtestNode) GetType() NodeType {
	return GetNodeType(abtest.Info.Kind)
}

func (abtest AbtestNode) GetInfo() NodeInfo {
	return abtest.Info
}

func (abtest AbtestNode) Parse(ctx *PipelineContext) (*NodeResult, error) {
	log.Println("======[trace]abtest start======")
	rand.Seed(time.Now().UnixNano())
	winNum := rand.Float64() * 100
	var counter float64 = 0
	for _, branch := range abtest.Branchs {
		counter += branch.Percent
		if counter > winNum {
			//feature global.Features.Set(dto.Feature{Name: abtest.Name, Value: branch.Name})
			log.Printf("abtest %v : %v, %v, output:%v \n", abtest.GetName(), branch.Name, winNum, branch.Decision.Output)
			nextNodeName := branch.Decision.Output.Value.(string)
			nextNodeType := GetNodeType(branch.Decision.Output.Kind)
			nodeResult := NodeResult{NextNodeName: nextNodeName, NextNodeType: nextNodeType}
			return &nodeResult, nil
			/*if res, ok := branch.Decision.Output.([]interface{}); ok {
				if len(res) == 2 {
					log.Println("abtest result", res)
					return res, nil
				}
			}*/
		}
	}
	log.Println("======[trace]abtest end======")
	return (*NodeResult)(nil), errcode.ParseErrorNoBranchMatch
}
