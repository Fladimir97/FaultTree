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
	var out string = fmt.Sprintf("\n%s [shape=none, label=%s]", n.getTitle(), n.getLabel())
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
	return fmt.Sprintf(`<<TABLE BORDER="0" CELLBORDER="0" CELLSPACING="0">
        <TR><TD>%s</TD></TR>
        <TR><TD><IMG SRC="images/and_gate.png"/></TD></TR>
        <TR><TD>F: %.4f</TD></TR>
        <TR><TD>R: %.4f</TD></TR>
    </TABLE>>`, n.Title, n.getFailure(), n.getReliability())
}

func newAND_NODE(title string) *AND_NODE {
	return &AND_NODE{title, []NODE{}}
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
	var out string = fmt.Sprintf("\n%s [shape=none, label=%s]", n.getTitle(), n.getLabel())
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
	return fmt.Sprintf(`<<TABLE BORDER="0" CELLBORDER="0" CELLSPACING="0">
        <TR><TD>%s</TD></TR>
        <TR><TD><IMG SRC="images/or_gate.png"/></TD></TR>
        <TR><TD>F: %.4f</TD></TR>
        <TR><TD>R: %.4f</TD></TR>
    </TABLE>>`, n.Title, n.getFailure(), n.getReliability())
}

func newOR_NODE(title string) *OR_NODE {
	return &OR_NODE{title, []NODE{}}
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

func newEvent(title string, reliability, failure float64) (*EVENT, error) {
	if reliability+failure != 1 {
		return nil, fmt.Errorf("Fehler: Reliability und Failure müssen in der Summe 1 ergeben")
	}
	return &EVENT{title, reliability, failure}, nil
}

// ---- main Func ------

func main() {
	TOP := newOR_NODE("TOP")
	A1 := newOR_NODE("A")
	A2 := newAND_NODE("B")
	E1, _ := newEvent("E1", 0.99, 0.01)
	E2, _ := newEvent("E2", 0.98, 0.02)
	E3, _ := newEvent("E3", 0.995, 0.005)
	E4, _ := newEvent("E4", 0.95, 0.05)
	E5, _ := newEvent("E5", 0.99, 0.01)

	TOP.addNode(A1)
	TOP.addNode(A2)

	A1.addEvent(*E1)
	A1.addEvent(*E2)
	A1.addEvent(*E3)

	A2.addEvent(*E4)
	A2.addEvent(*E5)

	MCS := TOP.getCutSets()
	// fmt.Print(MCS)
	for index, cutset := range MCS {
		fmt.Printf("\n C_%d := {", index+1)
		for _, event := range cutset {
			fmt.Printf("%s ", event.Title)
		}
		fmt.Print("}")
	}
	fmt.Printf("\nZuverlässigkeit des Gesamtsystems:\t%.5f", TOP.getReliability())
	fmt.Printf("\nAusfallwahrscheinlichkeit des Gesamtsystems:\t%.5f", TOP.getFailure())
	var digraph string = TOP.toDigraph()
	os.WriteFile("Fehlerbaum.dot", []byte(digraph), 0644)
}
