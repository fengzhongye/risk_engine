package core

import (
	"errors"
	"fmt"
	"log"
)

type DecisionFlow struct {
	flowMap   map[string]*FlowNode
	startNode *FlowNode
}

func NewDecisionFlow() *DecisionFlow {
	return &DecisionFlow{flowMap: make(map[string]*FlowNode)}
}

func (flow *DecisionFlow) AddNode(node *FlowNode) {
	key := flow.getNodeKey(node.NodeName, node.NodeKind)
	if _, ok := flow.flowMap[key]; !ok {
		flow.flowMap[key] = node
	} else {
		log.Println("repeat add node: " + key)
	}
}

//NodeType string
func (flow *DecisionFlow) GetNode(name string, nodeType interface{}) (*FlowNode, bool) {
	key := flow.getNodeKey(name, nodeType)
	if flowNode, ok := flow.flowMap[key]; ok {
		return flowNode, ok
	}
	return new(FlowNode), false
}

func (flow *DecisionFlow) GetAllNodes() map[string]*FlowNode {
	return flow.flowMap
}

func (flow *DecisionFlow) getNodeKey(name string, nodeType interface{}) string {
	return fmt.Sprintf("%s-%s", nodeType, name)
}

func (flow *DecisionFlow) SetStartNode(startNode *FlowNode) {
	flow.startNode = startNode
}

func (flow *DecisionFlow) GetStartNode() (*FlowNode, bool) {
	return flow.startNode, true
}

func (flow *DecisionFlow) Run(ctx *PipelineContext) (err error) {
	//recover
	go func() {
		defer func() {
			if err := recover(); err != nil {
				err = err
				log.Println(err)
			}
		}()
	}()

	//find StartNode
	flowNode, ok := flow.GetStartNode()
	if !ok {
		err = errors.New("no start node")
		return
	}

	gotoNext := true
	for gotoNext {
		//ctx.SetCurrentNode(flowNode)
		flowNode, gotoNext = flow.parseNode(flowNode, ctx)
	}
	return
}

//parse current node and return next node
func (flow *DecisionFlow) parseNode(curNode *FlowNode, ctx *PipelineContext) (nextNode *FlowNode, gotoNext bool) {
	//parse current node
	ctx.AddTrack(curNode.GetElem())
	res, err := curNode.Parse(ctx)
	if err != nil {
		log.Println(err)
	}
	ctx.AddNodeResult(curNode.NodeName, res)

	//get next node
	if res.IsBlock {
		gotoNext = !res.IsBlock
		return
	}

	switch curNode.GetNodeType() { //int
	case TypeEnd: //END:
		gotoNext = false
		return
	case TypeAbtest: //ABTEST:
		nextNode, gotoNext = flow.GetNode(res.NextNodeName, res.NextNodeType)
		return
	default: //start
		nextNode, gotoNext = flow.GetNode(curNode.NextNodeName, curNode.NextNodeKind)
		return
	}
	return
}

type FlowNode struct {
	NodeName     string `yaml:"node_name"`
	NodeKind     string `yaml:"node_kind"`
	NextNodeName string `yaml:"next_node_name"`
	NextNodeKind string `yaml:"next_node_kind"`

	elem     INode
	nextNode *FlowNode
}

func (flowNode *FlowNode) GetNodeType() NodeType {
	return GetNodeType(flowNode.NodeKind)
}

func (flowNode *FlowNode) GetNextNodeType() NodeType {
	return GetNodeType(flowNode.NextNodeKind)
}

func (flowNode *FlowNode) SetElem(elem INode) {
	flowNode.elem = elem
}

func (flowNode *FlowNode) GetElem() INode {
	return flowNode.elem
}

func (flowNode *FlowNode) Parse(ctx *PipelineContext) (*NodeResult, error) {
	//hook
	return flowNode.elem.Parse(ctx)
}
