package prioritize

import (
	"reflect"
	"testing"
)

func TestParseMark(t *testing.T) {

	type parseMarkData struct {
		labels   map[string]string
		expected []string
	}

	tests := []parseMarkData{
		{labels: map[string]string{"ks-pipeline": "jenkins-java-maven-1"}, expected: []string{"jenkins", "java", "maven", "1"}},
		{labels: map[string]string{"ks-pipeline": "jenk/ins-java-maven-1/"}, expected: []string{"jenk/ins", "java", "maven", "1/"}},
	}

	for _, lb := range tests {
		if res := parseMark(lb.labels); ! reflect.DeepEqual(res, lb.expected) {
			t.Errorf(lb.labels["ks-pipeline"])
			t.Errorf("Expected key %s, but got %s", lb.expected, res)
		}
	}
}

func TestCalculation(t *testing.T) {

	type parseMarkData struct {
		name     string
		keys     []string
		nodeName string
		expected int
	}

	tests := []parseMarkData{
		{name: "a-b-c-d",keys: []string{"keya", "keyab", "keyac", "keyad"}, nodeName: "node1", expected: 10},
		{name: "a-b-c",keys: []string{"keya", "keyab", "keyac"}, nodeName: "node1", expected: 6},
		{name: "z-x-y",keys: []string{"z", "x", "y"}, nodeName: "node1", expected: 1},
		{name: "z-x-a",keys: []string{"z", "x", "keya"}, nodeName: "node1", expected: 10},
	}

	for _, match := range tests {
		if score, err := Calculation(match.keys, match.nodeName); err != nil || score != match.expected {
			t.Errorf(match.name)
			t.Errorf("Expected score %d, but got %d", match.expected, score)
		}
	}
}
