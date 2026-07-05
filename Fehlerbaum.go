package main

import (
	"fmt"
	"os"
)

type NODE interface {
	getCutSets() [][]EVENT
	getReliability() float64
	getFailure() float64
	toDot() string
	getTitle() string
}

// ---- AND NODE ------

type AND_NODE struct {
	Title    string
	Children []NODE
}

func (n AND_NODE) getCutSets() [][]EVENT {
	out := [][]EVENT{{}}
	for _, child := range n.Children {
		out = crossproduct(out, child.getCutSets())
	}
	return out
}

func crossproduct(a, b [][]EVENT) [][]EVENT {
	out := [][]EVENT{}
	for _, cutSetA := range a {
		for _, cutSetB := range b {
			combined := append(cutSetA, cutSetB...)
			out = append(out, combined)
		}
	}
	return out
}

func (n *AND_NODE) addNode(other NODE) {
	n.Children = append(n.Children, other)
}

func (n *AND_NODE) addEvent(e EVENT) {
	n.Children = append(n.Children, e)
}

func (n AND_NODE) getReliability() float64 {
	// R_AND = 1 - Π (1 - R_i)
	var out float64 = 1
	for _, child := range n.Children {
		out *= (1 - child.getReliability())
	}
	return 1 - out
}

func (n AND_NODE) getFailure() float64 {
	// F_AND = Π F_i
	var out float64 = 1
	for _, child := range n.Children {
		out *= child.getFailure()
	}
	return out
}

func (n AND_NODE) toDot() string {
	var out string = fmt.Sprintf("\n%s [shape=house, label=\"%s\"]", n.getTitle(), n.getLabel())
	for _, child := range n.Children {
		out += fmt.Sprintf("\n%s -> %s", n.getTitle(), child.getTitle())
		out += child.toDot()
	}
	return out
}

func (n AND_NODE) getTitle() string {
	return n.Title
}

func (n AND_NODE) toDigraph() string {
	return fmt.Sprintf("digraph {\n%s\n}", n.toDot())
}

func (n AND_NODE) getLabel() string {
	return fmt.Sprintf("%s\nF: %.4f\nR: %.4f", n.getTitle(), n.getFailure(), n.getReliability())
}

// ---- OR NODE ------

type OR_NODE struct {
	Title    string
	Children []NODE
}

func (n OR_NODE) getCutSets() [][]EVENT {
	out := [][]EVENT{}
	for _, child := range n.Children {
		out = append(out, child.getCutSets()...)
	}
	return out
}

func (n *OR_NODE) addNode(other NODE) {
	n.Children = append(n.Children, other)
}

func (n *OR_NODE) addEvent(e EVENT) {
	n.Children = append(n.Children, e)
}

func (n OR_NODE) getReliability() float64 {
	// R_OR = Π R_i
	var out float64 = 1
	for _, child := range n.Children {
		out *= child.getReliability()
	}
	return out
}

func (n OR_NODE) getFailure() float64 {
	// F_OR = 1 - Π (1 - F_i)
	var out float64 = 1
	for _, child := range n.Children {
		out *= (1 - child.getFailure())
	}
	return 1 - out
}

func (n OR_NODE) getTitle() string {
	return n.Title
}

func (n OR_NODE) toDot() string {
	var out string = fmt.Sprintf("\n%s [shape=invhouse, label=\"%s\"]", n.getTitle(), n.getLabel())
	for _, child := range n.Children {
		out += fmt.Sprintf("\n%s -> %s", n.getTitle(), child.getTitle())
		out += child.toDot()
	}
	return out
}

func (n OR_NODE) toDigraph() string {
	return fmt.Sprintf("digraph {\n%s\n}", n.toDot())
}

func (n OR_NODE) getLabel() string {
	return fmt.Sprintf("%s\nF: %.4f\nR: %.4f", n.getTitle(), n.getFailure(), n.getReliability())
}

// ---- EVENT NODE ------

type EVENT struct {
	Title       string
	Reliability float64
	Failure     float64
}

func (n EVENT) getCutSets() [][]EVENT {
	return [][]EVENT{{n}}
}

func (n EVENT) getReliability() float64 {
	return n.Reliability
}

func (n EVENT) getFailure() float64 {
	return n.Failure
}

func (n EVENT) getTitle() string {
	return n.Title
}

func (n EVENT) toDot() string {
	return fmt.Sprintf("\n%s [shape=circle, label=\"%s\"]", n.getTitle(), n.getLabel())
}

func (n EVENT) getLabel() string {
	return fmt.Sprintf("%s\nF: %.4f\nR: %.4f", n.getTitle(), n.getFailure(), n.getReliability())
}

// ---- main Func ------

func main() {
	AND := AND_NODE{
		"TOP", []NODE{},
	}
	OR := OR_NODE{
		"OR1", []NODE{},
	}
	E1 := EVENT{
		"E1", 0.99, 0.01,
	}
	E2 := EVENT{
		"E2", 0.98, 0.02,
	}
	E3 := EVENT{
		"E3", 0.97, 0.03,
	}
	AND.addNode(&OR)
	AND.addEvent(E1)
	OR.addEvent(E2)
	OR.addEvent(E3)

	MCS := AND.getCutSets()
	// fmt.Print(MCS)
	for index, cutset := range MCS {
		fmt.Printf("\n C_%d := {", index+1)
		for _, event := range cutset {
			fmt.Printf("%s ", event.Title)
		}
		fmt.Print("}")
	}
	fmt.Printf("\nZuverlässigkeit des Gesamtsystems:\t%.5f", AND.getReliability())
	fmt.Printf("\nAusfallwahrscheinlichkeit des Gesamtsystems:\t%.5f", AND.getFailure())
	var digraph string = AND.toDigraph()
	os.WriteFile("Fehlerbaum.dot", []byte(digraph), 0644)
}
